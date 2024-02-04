package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"strconv"
	"workflow/consts"
	"workflow/dao"
	"workflow/model"
	"workflow/types"
	"workflow/utils"
)

// WorkFlowApproval 审批
func WorkFlowApproval(tx *gorm.DB, req *types.ProcProcApprovalReq) (resp types.ProcInstResp, err error) {
	var (
		inst           model.ProcInst
		nodeId         int
		nodeName       string
		procNode       types.ProcNode
		procTask       model.ProcTask
		isTaskComplete bool
	)
	//查询当前流程实例
	inst, err = dao.ProcInstDao{}.GetByProcInstById(tx, req.ProcInstId)
	if err != nil {
		return
	}
	if inst.Id == 0 {
		err = errors.New("process does not exist")
		return
	}
	if inst.IsComplete == consts.IsCompleteYes { //校验是否已完成
		err = errors.New("the process has ended")
		return
	}
	//获取当前节点信息
	resource := inst.NodeInfo
	currProcNodeStr := gjson.Get(resource, "resource.#(nodeId="+strconv.Itoa(inst.NodeId)+")").String()
	var currProcNode types.ProcNode
	err = json.Unmarshal([]byte(currProcNodeStr), &currProcNode)
	if err != nil {
		return
	}
	//当前节点Task
	procTask, err = dao.ProcTaskDao{}.GetTaskByInstIdAndNodeId(tx, inst.Id, inst.NodeId)
	if err != nil {
		return
	}
	if procTask.Id == 0 {
		err = errors.New("process task does not exist")
		return
	}
	if req.ApprovalStatus == consts.ApprovalStatusAdopt {
		//是否多次审核
		if currProcNode.ActType != consts.NodeActTypeAnd {
			//直接下一步
			nodeId = currProcNode.NextNode
			isTaskComplete = true
		} else {
			//多次审核
			if procTask.MemberCount > (procTask.AgreeNum + 1) {
				nodeId = procTask.NodeId
				//审批次数不够
				isTaskComplete = false
			} else {
				nodeId = currProcNode.NextNode
				//审批完下一步
				isTaskComplete = true
			}
		}

	} else {
		nodeId = currProcNode.PrevNode
		isTaskComplete = true
	}

	if isTaskComplete {
		//当前审批任务已完成,进入下一步
		procNode, err = Process(nodeId, req.Parameters, resource)
		if err != nil {
			return
		}
		nodeName = procNode.NodeName
		nodeId = procNode.NodeId
		//更新历史任务已完成
		procTask.UnCompleteNum = procTask.UnCompleteNum - 1
		procTask.AgreeNum = procTask.AgreeNum + 1
		procTask.IsComplete = consts.IsCompleteYes
		err = dao.ProcTaskDao{}.Update(tx, procTask)
		if err != nil {
			return
		}
		if procNode.Type == consts.NodeTypeEnd {
			//节点走完,更新实例
			inst.NodeId = procNode.NodeId
			inst.IsComplete = consts.IsCompleteYes
			err = dao.ProcInstDao{}.Update(tx, inst)
			if err != nil {
				return
			}
		} else {
			//创建新的任务
			task := model.ProcTask{
				ProcInstId:    inst.Id,
				NodeId:        procNode.NodeId,
				Assignee:      req.UserName,
				MemberCount:   procNode.MemberCount,
				UnCompleteNum: procNode.UnCompleteNum,
				AgreeNum:      procNode.AgreeNum,
				ActType:       procNode.ActType,
				IsComplete:    consts.IsCompleteNo,
			}
			err = dao.ProcTaskDao{}.Add(tx, task)
			if err != nil {
				return
			}

			//更新流程实例
			inst.NodeId = procNode.NodeId
			err = dao.ProcInstDao{}.Update(tx, inst)
			if err != nil {
				return
			}
		}
	} else {
		nodeName = currProcNode.NodeName
		//审批任务未完成,继续当前节点
		procTask.UnCompleteNum = procTask.UnCompleteNum - 1
		procTask.AgreeNum = procTask.AgreeNum + 1
		err = dao.ProcTaskDao{}.Update(tx, procTask)
		if err != nil {
			return
		}
		//更新流程实例
		err = dao.ProcInstDao{}.Update(tx, inst)
		if err != nil {
			return
		}
	}

	//写入审批历史
	his := model.ProcHis{
		Assignee:       req.UserName,
		ProcInstId:     inst.Id,
		NodeId:         inst.NodeId,
		ApprovalStatus: req.ApprovalStatus,
	}
	err = dao.ProcHisDao{}.Add(tx, his)
	if err != nil {
		return
	}

	resp.ProcInstId = inst.Id
	resp.NodeName = nodeName
	resp.NodeId = nodeId
	return
}

// Process 流程节点解析,返回下一个节点
func Process(nodeId int, parameters map[string]interface{}, resource string) (procNode types.ProcNode, err error) {
	var (
		nextProcNode types.ProcNode
		expr         *govaluate.EvaluableExpression
		result       interface{} //表达式计算结果
		flag         bool        //是否满足表达式
		retNode      int         //最终走向
	)
	//获取节点
	nextProcNodeStr := gjson.Get(resource, "resource.#(nodeId="+strconv.Itoa(nodeId)+")").String()
	if utils.IsEmpty(nextProcNodeStr) {
		err = errors.New("please check the process configuration if you can't find the corresponding configured node")
		return
	}
	err = json.Unmarshal([]byte(nextProcNodeStr), &nextProcNode)
	if err != nil {
		return
	}
	if len(nextProcNode.Type) > 0 {
		if nextProcNode.Type == consts.NodeTypeCondition {
			//下一步是条件判断
			if len(nextProcNode.Condition) > 0 {
				expr, err = govaluate.NewEvaluableExpression(nextProcNode.Condition)
				if err != nil {
					err = errors.New("conditional Expression Configuration Exception")
					return
				}
				result, err = expr.Evaluate(parameters)
				if err != nil {
					err = errors.New("conditional Expression Parsing Exception")
					return
				}
				flag, err = strconv.ParseBool(fmt.Sprint(result))
				if err != nil {
					err = errors.New("expression Result Getting Exception")
					return
				}
				if flag {
					//满足 下一步
					retNode = nextProcNode.NextNode
				} else {
					//不满足 上一步
					retNode = nextProcNode.PrevNode
				}
			} else {
				err = errors.New("conditional node expression is empty")
				return
			}

			cdProcNodeStr := gjson.Get(resource, "resource.#(nodeId="+strconv.Itoa(retNode)+")").String()
			if len(cdProcNodeStr) == 0 {
				err = errors.New("can't find the next node")
				return
			}
			err = json.Unmarshal([]byte(cdProcNodeStr), &nextProcNode)
			if err != nil {
				return
			}
		}

		//如果下一个点类型还是条件 则再次解析
		//场景 例如: amount > 100 < 1000 经理审批  / amount > 1000 < 10000 主管审批  / amount > 10000  总监审批 可以通过condition 无限扩展分支
		if nextProcNode.Type == consts.NodeTypeCondition {
			nextProcNode, err = Process(nextProcNode.NodeId, parameters, resource)
			if err != nil {
				return
			}
		}
		procNode = nextProcNode
	} else {
		err = errors.New("node type configuration exception")
		return
	}
	return
}
