package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	TableNameProcTask = "proc_task"
	TableNameProcInst = "proc_inst"
	TableNameProcHis  = "proc_his"
	TableNameProcDef  = "proc_def"
)

type (
	// ProcDef 流程定义表
	ProcDef struct {
		Id        int            `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
		ProcName  string         `gorm:"column:proc_name;type:varchar(255);comment:流程名称;NOT NULL" json:"procName"`
		Version   int            `gorm:"column:version;type:int(11);default:1;comment:版本 默认1;NOT NULL" json:"version"`
		Remark    string         `gorm:"column:remark;type:varchar(255);comment:备注;NOT NULL" json:"remark"`
		Resource  string         `gorm:"column:resource;type:varchar(500);comment:流程JSON;NOT NULL" json:"resource"`
		CreatedAt time.Time      `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"createdAt"`
		UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`
		DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deletedAt"`
	}
	// ProcHis 流程历史记录表
	ProcHis struct {
		Id             int       `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
		Assignee       string    `gorm:"column:assignee;type:varchar(50);comment:任务处理人;NOT NULL" json:"assignee"`
		ProcInstId     int       `gorm:"column:proc_inst_id;type:bigint(20);comment:流程实例ID;NOT NULL" json:"procInstId"`
		NodeId         int       `gorm:"column:node_id;type:bigint(20);comment:节点id;NOT NULL" json:"nodeId"`
		ApprovalStatus int       `gorm:"column:approval_status;type:tinyint(1);comment:审批状态: 1.审批通过 2.审批驳回 3.发起流程 4.中止流程 5.结束流程;NOT NULL" json:"approvalStatus"`
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"createdAt"`
	}

	// ProcInst 流程实例表
	ProcInst struct {
		Id             int       `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
		ProcDefId      int       `gorm:"column:proc_def_id;type:bigint(20);comment:流程定义ID;NOT NULL" json:"proc_def_id"`
		ProcDefVersion int       `gorm:"column:proc_def_version;type:int(11);comment:流程版本号;NOT NULL" json:"proc_def_version"`
		NodeInfo       string    `gorm:"column:node_info;type:varchar(500);comment:流程节点;NOT NULL" json:"node_info"`
		NodeId         int       `gorm:"column:node_id;type:bigint(20);comment:当前节点;NOT NULL" json:"node_id"`
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"createdAt"`
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`
		IsComplete     int       `gorm:"column:is_complete;type:tinyint(1);default:1;comment:是否完成 1:进行中 2已完成;NOT NULL" json:"is_complete"`
		Assignee       string    `gorm:"column:assignee;type:varchar(50);comment:发起人;NOT NULL" json:"assignee"`
	}
	// ProcTask 流程任务表
	ProcTask struct {
		Id            int       `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
		ProcInstId    int       `gorm:"column:proc_inst_id;type:bigint(20);comment:流程实例ID;NOT NULL" json:"proc_inst_id"`
		NodeId        int       `gorm:"column:node_id;type:bigint(20);comment:节点ID;NOT NULL" json:"node_id"`
		Assignee      string    `gorm:"column:assignee;type:varchar(50);comment:任务处理人;NOT NULL" json:"assignee"`
		MemberCount   int       `gorm:"column:member_count;type:int(11);default:1;comment:表示当前任务需要多少人审批之后才能结束，默认是 1;NOT NULL" json:"member_count"`
		UnCompleteNum int       `gorm:"column:un_complete_num;type:int(11);default:1;comment:表示还有多少人没有审批，默认是1;NOT NULL" json:"un_complete_num"`
		AgreeNum      int       `gorm:"column:agree_num;type:int(11);default:0;comment:表示通过的人数;NOT NULL" json:"agree_num"`
		ActType       string    `gorm:"column:act_type;type:varchar(10);" json:"act_type"`
		CreatedAt     time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;NOT NULL" json:"createdAt"`
		UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updatedAt"`
		IsComplete    int       `gorm:"column:is_complete;type:tinyint(1);default:1;comment:是否完成 1:进行中 2已完成;NOT NULL" json:"is_complete"`
	}
)

func (m *ProcTask) TableName() string {
	return TableNameProcTask
}

func (m *ProcInst) TableName() string {
	return TableNameProcInst
}

func (m *ProcHis) TableName() string {
	return TableNameProcHis
}

func (m *ProcDef) TableName() string {
	return TableNameProcDef
}
