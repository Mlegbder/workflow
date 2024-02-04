package workflow

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"strconv"
	"workflow/consts"
	"workflow/dao"
	"workflow/model"
	"workflow/types"
)

// WorkFlowJump 流程跳跃(特殊业务)
func WorkFlowJump(tx *gorm.DB, req *types.ProcReq) (resp types.ProcInstResp, err error) {
	var (
		inst         model.ProcInst
		jumpProcNode types.ProcNode
		nextProcNode types.ProcNode
		procTask     model.ProcTask
	)
	//查询当前流程实例
	inst, err = dao.ProcInstDao{}.GetByProcInstById(tx, req.Id)
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
	jumpProcNodeStr := gjson.Get(resource, "resource.#(type="+consts.NodeTypeJump+")").String()
	if len(jumpProcNodeStr) == 0 {
		err = errors.New("the process is not configured with a jump node, and cannot be jumped")
		return
	}
	err = json.Unmarshal([]byte(jumpProcNodeStr), &jumpProcNode)
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
	//更新历史任务已完成
	procTask.UnCompleteNum = procTask.UnCompleteNum - 1
	procTask.AgreeNum = procTask.AgreeNum + 1
	procTask.IsComplete = consts.IsCompleteYes
	err = dao.ProcTaskDao{}.Update(tx, procTask)

	//获取跳跃流程
	nextProcNodeStr := gjson.Get(resource, "resource.#(nodeId="+strconv.Itoa(jumpProcNode.NextNode)+")").String()
	if len(nextProcNodeStr) == 0 {
		err = errors.New("process configuration jump node does not exist")
		return
	}
	err = json.Unmarshal([]byte(nextProcNodeStr), &nextProcNode)
	if err != nil {
		return
	}

	//创建新的任务
	task := model.ProcTask{
		ProcInstId:    inst.Id,
		NodeId:        nextProcNode.NodeId,
		Assignee:      req.UserName,
		MemberCount:   nextProcNode.MemberCount,
		UnCompleteNum: nextProcNode.UnCompleteNum,
		AgreeNum:      nextProcNode.AgreeNum,
		ActType:       nextProcNode.ActType,
		IsComplete:    consts.IsCompleteNo,
	}
	err = dao.ProcTaskDao{}.Add(tx, task)
	if err != nil {
		return
	}

	//更新流程实例
	inst.NodeId = nextProcNode.NodeId
	err = dao.ProcInstDao{}.Update(tx, inst)
	if err != nil {
		return
	}

	//写入审批历史
	his := model.ProcHis{
		Assignee:       req.UserName,
		ProcInstId:     inst.Id,
		NodeId:         inst.NodeId,
		ApprovalStatus: consts.ApprovalStatusJump,
	}
	err = dao.ProcHisDao{}.Add(tx, his)
	if err != nil {
		return
	}

	resp.ProcInstId = inst.Id
	resp.NodeName = nextProcNode.NodeName
	resp.NodeId = nextProcNode.NodeId
	return
}
