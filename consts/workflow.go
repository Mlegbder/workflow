package consts

const (
	//审批状态:  1.审批通过 2.审批驳回  3.发起流程 4.中止流程 5.结束流程 6.跳跃流程
	ApprovalStatusAdopt    = 1
	ApprovalStatusReject   = 2
	ApprovalStatusLaunch   = 3
	ApprovalStatusSuspend  = 4
	ApprovalStatusForceEnd = 5
	ApprovalStatusJump     = 6

	//节点类型
	NodeTypeStart     = "start"     //开始
	NodeTypeEnd       = "end"       //结束
	NodeTypeCondition = "condition" //条件
	NodeTypeProcess   = "process"   //流程
	NodeTypeSuspend   = "suspend"   //中止
	NodeTypeForceEnd  = "force_end" //强制结束
	NodeTypeJump      = "jump"      //跳跃

	//"or"表示或签，即一个人通过或者驳回就结束，"and"表示会签，要所有人通过就流转到下一步，如果有一个人驳回那么就跳转到上一步
	NodeActTypeAnd = "and"
	NodeActTypeOr  = "or"

	//是否完成 1:进行中 2已完成
	IsCompleteYes = 2
	IsCompleteNo  = 1
)

var (
	ApprovalStatusList = []int{
		ApprovalStatusAdopt,
		ApprovalStatusReject,
		ApprovalStatusLaunch,
		ApprovalStatusSuspend,
		ApprovalStatusForceEnd,
	}
)
