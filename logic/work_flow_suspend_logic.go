package workflow

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"workflow/consts"
	"workflow/dao"
	"workflow/model"
	"workflow/types"
)

// WorkFlowSuspend 中止流程
func WorkFlowSuspend(tx *gorm.DB, req *types.ProcReq) (resp types.ProcInstResp, err error) {
	var (
		inst            model.ProcInst
		suspendProcNode types.ProcNode
		procTask        model.ProcTask
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
	suspendProcNodeStr := gjson.Get(resource, "resource.#(type="+consts.NodeTypeSuspend+")").String()
	if len(suspendProcNodeStr) == 0 {
		err = errors.New("the process is not configured with an abort node and cannot be aborted")
		return
	}
	err = json.Unmarshal([]byte(suspendProcNodeStr), &suspendProcNode)
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

	//创建新的任务
	task := model.ProcTask{
		ProcInstId:    inst.Id,
		NodeId:        suspendProcNode.NodeId,
		Assignee:      req.UserName,
		MemberCount:   suspendProcNode.MemberCount,
		UnCompleteNum: suspendProcNode.UnCompleteNum,
		AgreeNum:      suspendProcNode.AgreeNum,
		ActType:       suspendProcNode.ActType,
		IsComplete:    consts.IsCompleteNo,
	}
	err = dao.ProcTaskDao{}.Add(tx, task)
	if err != nil {
		return
	}

	//更新流程实例
	inst.NodeId = suspendProcNode.NodeId
	err = dao.ProcInstDao{}.Update(tx, inst)
	if err != nil {
		return
	}

	//写入审批历史
	his := model.ProcHis{
		Assignee:       req.UserName,
		ProcInstId:     inst.Id,
		NodeId:         inst.NodeId,
		ApprovalStatus: consts.ApprovalStatusSuspend,
	}
	err = dao.ProcHisDao{}.Add(tx, his)
	if err != nil {
		return
	}

	resp.ProcInstId = inst.Id
	resp.NodeName = suspendProcNode.NodeName
	resp.NodeId = suspendProcNode.NodeId
	return
}
