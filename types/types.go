package types

type (
	// ProcInstCreateReq 创建流程
	ProcInstCreateReq struct {
		Id         int                    `json:"id"`
		Parameters map[string]interface{} `json:"parameters,optional"`
		UserName   string                 `json:"userName"`
	}

	//ProcInstResp 实例
	ProcInstResp struct {
		ProcInstId int    `json:"procInstId"`
		NodeId     int    `json:"nodeId"`
		NodeName   string `json:"nodeName"`
	}

	// ProcProcApprovalReq 审批
	ProcProcApprovalReq struct {
		ProcInstId     int                    `json:"procInstId"`
		ApprovalStatus int                    `json:"approvalStatus"`
		Parameters     map[string]interface{} `json:"parameters"`
		UserName       string                 `json:"userName"`
	}

	// ProcNode 节点
	ProcNode struct {
		NodeId        int    `json:"nodeId"`
		NodeName      string `json:"nodeName"`
		Type          string `json:"type"`
		Condition     string `json:"condition"`
		NextNode      int    `json:"nextNode"`
		PrevNode      int    `json:"prevNode"`
		Assignee      string `json:"assignee"`
		MemberCount   int    `json:"memberCount"`
		UnCompleteNum int    `json:"unCompleteNum"`
		AgreeNum      int    `json:"agreeNum"`
		ActType       string `json:"actType"`
	}

	// ProcReq 历史记录
	ProcReq struct {
		Id       int    `json:"id"`
		UserName string `json:"userName"`
	}
	ProcHisResp struct {
		Id             int    `json:"id"`
		Assignee       string `json:"assignee"`
		ProcInstId     int    `json:"procInstId"`
		NodeId         int    `json:"nodeId"`
		ApprovalStatus int    `json:"approvalStatus"`
	}
)
