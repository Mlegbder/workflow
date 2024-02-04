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

// WorkFlowForceEnd 强制结束流程
func WorkFlowForceEnd(tx *gorm.DB, req *types.ProcReq) (resp types.ProcInstResp, err error) {
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
	suspendProcNodeStr := gjson.Get(resource, "resource.#(type="+consts.NodeTypeForceEnd+")").String()
	if len(suspendProcNodeStr) == 0 {
		err = errors.New("the process is not configured with a forced end node")
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

	//更新实例
	inst.NodeId = suspendProcNode.NodeId
	inst.IsComplete = consts.IsCompleteYes
	err = dao.ProcInstDao{}.Update(tx, inst)

	//写入审批历史
	his := model.ProcHis{
		Assignee:       req.UserName,
		ProcInstId:     inst.Id,
		NodeId:         inst.NodeId,
		ApprovalStatus: consts.ApprovalStatusForceEnd,
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
