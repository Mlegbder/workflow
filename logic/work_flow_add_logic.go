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

// WorkFlowAdd 创建工作流
func WorkFlowAdd(tx *gorm.DB, req *types.ProcInstCreateReq) (resp types.ProcInstResp, err error) {
	var (
		nextProcNode types.ProcNode
		def          model.ProcDef
		expr         *govaluate.EvaluableExpression
		result       interface{} //表达式计算结果
		flag         bool        //是否满足表达式
	)
	isCompleteNo := consts.IsCompleteNo
	//根据ID获取流程定义
	def, err = dao.ProcDefDao{}.GetByProcDefById(tx, req.Id)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to get process definition, reason for failure: %s", err.Error()))
		return
	}
	if def.Id == 0 {
		err = errors.New("get process definition does not exist")
		return
	}
	//解析流程JSON
	if len(def.Resource) > 0 {
		if !gjson.Valid(def.Resource) {
			err = errors.New("incorrect process definition Json format")
			return
		}
	} else {
		err = errors.New("process definition Json is empty, can not generate workflow")
		return
	}
	//获取节点
	nextProcNodeStr := gjson.Get(def.Resource, `resource.#(type="start")`).String()
	if utils.IsEmpty(nextProcNodeStr) {
		err = errors.New("please check the process configuration if you can't find the corresponding configured node")
		return
	}
	err = json.Unmarshal([]byte(nextProcNodeStr), &nextProcNode)
	if err != nil {
		return
	}
	//取流程定义信息
	nodeId := gjson.Get(def.Resource, `resource.#(type="start").nodeId`).Num
	memberCount := gjson.Get(def.Resource, `resource.#(type="start").memberCount`).Num
	agreeNum := gjson.Get(def.Resource, `resource.#(type="start").agreeNum`).Num
	actType := gjson.Get(def.Resource, `resource.#(type="start").actType`).String()
	nodeName := gjson.Get(def.Resource, `resource.#(type="start").nodeName`).String()
	cond := gjson.Get(def.Resource, `resource.#(type="start").condition`).String()
	retNode := int(nodeId)
	retMemberCount := int(memberCount)
	retAgreeNum := int(agreeNum)
	//起始节点存在条件
	if !utils.IsEmpty(cond) {
		expr, err = govaluate.NewEvaluableExpression(cond)
		if err != nil {
			err = errors.New("conditional Expression Configuration Exception")
			return
		}
		result, err = expr.Evaluate(req.Parameters)
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
		cdProcNodeStr := gjson.Get(def.Resource, "resource.#(nodeId="+strconv.Itoa(retNode)+")").String()
		if len(cdProcNodeStr) == 0 {
			err = errors.New("can't find the next node")
			return
		}
		err = json.Unmarshal([]byte(cdProcNodeStr), &nextProcNode)
		if err != nil {
			return
		}
		if nextProcNode.Type == consts.NodeTypeCondition {
			nextProcNode, err = Process(nextProcNode.NodeId, req.Parameters, def.Resource)
			if err != nil {
				return
			}
			retNode = nextProcNode.NodeId
			retMemberCount = nextProcNode.MemberCount
			retAgreeNum = nextProcNode.AgreeNum
		}
		//判断走条件后是否直接已完成
		if nextProcNode.Type == consts.NodeTypeEnd {
			isCompleteNo = consts.IsCompleteYes
		}
	}
	//创建流程
	procInst := model.ProcInst{
		ProcDefId:      def.Id,
		ProcDefVersion: def.Version,
		NodeInfo:       def.Resource,
		NodeId:         retNode,
		Assignee:       req.UserName,
		IsComplete:     isCompleteNo,
	}
	err = dao.ProcInstDao{}.Add(tx, &procInst)
	if err != nil {
		return
	}
	//创建流程任务
	task := model.ProcTask{
		ProcInstId:    procInst.Id,
		NodeId:        retNode,
		Assignee:      req.UserName,
		MemberCount:   retMemberCount,
		UnCompleteNum: retMemberCount,
		AgreeNum:      retAgreeNum,
		ActType:       actType,
		IsComplete:    isCompleteNo,
	}
	err = dao.ProcTaskDao{}.Add(tx, task)
	if err != nil {
		return
	}
	//创建流程历史
	his := model.ProcHis{
		Assignee:       req.UserName,
		ProcInstId:     procInst.Id,
		NodeId:         retNode,
		ApprovalStatus: consts.ApprovalStatusLaunch,
	}
	err = dao.ProcHisDao{}.Add(tx, his)
	if err != nil {
		return
	}
	//返回
	resp.ProcInstId = procInst.Id
	resp.NodeName = nodeName
	resp.NodeId = retNode
	return
}
