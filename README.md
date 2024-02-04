## WOMATA 2.0

## 说明
> 极简工作流,以json方式配置, 支持条件解析
> * 首次使用请执行init.sql初始化数据表
> * 配置流程到proc_dec
> * 流程支持修改,不影响产生的流程, 只对新创建的流程生效

## 流程配置类型
> * start (起始节点)
> * end（结束节点）
> * condition（条件节点）
> * process（流程节点）
> * suspend (中止)
> * force_end (强制结束节点)
> * jump (跳跃)

## 流程示例
> {"procId":1,"procName":"退款流程","version":1,"resource":[{"nodeId":10,"nodeName":"待退款","type":"start","condition":"","nextNode":102,"prevNode":40},{"nodeId":102,"nodeName":"","type":"condition","condition":"orderChannel == 3","nextNode":30,"prevNode":101},{"nodeId":101,"nodeName":"","type":"condition","condition":"refund - refundable >= 0 && refundable > 0 && refund > 0","nextNode":30,"prevNode":20},{"nodeId":20,"nodeName":"待核对","type":"process","condition":"","nextNode":30,"prevNode":0},{"nodeId":30,"nodeName":"已完成","type":"end","condition":"","nextNode":0,"prevNode":0},{"nodeId":40,"nodeName":"已取消","type":"end","condition":"","nextNode":0,"prevNode":0}]}

## 创建流程
```
// 创建工作流
flow := workflow.WorkFlowAddLogic{
	Logger: l.Logger,
	Ctx:    l.ctx,
	SvcCtx: l.svcCtx,
}
flowParam := types.ProcInstCreateReq{
	Id: consts.FlowActivityCreate, //数据库配置的流程id
}
flowResp, err = flow.WorkFlowAdd(tx, &flowParam)
```

## 审批流程
```
var(
    flow types.ProcInstResp                //调用工作流返回的结果
	flowCond = make(map[string]interface{}, 0) //创建流程时带的条件
)
//带条件审批
flowCond["refundStationFee"] = req.WarehouseRefund.RefundPayPickupFee //判断是否需要支付退货驿站费用
flowCond["afterType"] = req.WarehouseRefund.AfterType                 //判断是否仅退款
flowParam := types.ProcProcApprovalReq{
	ProcInstId:     warehouseRefund.FlowId,
	ApprovalStatus: req.ApprovalStatus,
	Parameters:     flowCond,
}
flow, err = workflow.NewWorkFlowApprovalLogic(l.ctx, l.svcCtx).WorkFlowApproval(tx, &flowParam)
```


## 强制结束/废弃流程 force_end
```
flowParam := types.ProcReq{
	Id: womataPurchase.FlowId,
}
flow := workflow.WorkFlowForceEndLogic{
	Logger: l.Logger,
	Ctx:    l.ctx,
	SvcCtx: l.svcCtx,
}
flowResp, err = flow.WorkFlowForceEnd(tx, &flowParam)
if err != nil {
	return
}
```

## 跳流程,特殊业务. 直接跳至jump的节点
```
flowParam := types.ProcReq{
	Id: warehouseRefund.FlowId,
}
flow, err = workflow.NewWorkFlowJumpLogic(l.ctx, l.svcCtx).WorkFlowJump(tx, &flowParam)
```


## 中止流程 直接跳至 suspend 节点
```
//调用工作流
flowParam := types.ProcReq{
	Id: activity.FlowId,
}
flow := workflow.WorkFlowSuspendLogic{
	Logger: l.Logger,
	Ctx:    l.ctx,
	SvcCtx: l.svcCtx,
}
flowResp, err = flow.WorkFlowSuspend(tx, &flowParam)
if err != nil {
	return
}
activity.ActivityStatus = flowResp.NodeId
err = tx.Updates(&activity).Error
if err != nil {
	return
}
```