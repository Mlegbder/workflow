package dao

import (
	"gorm.io/gorm"
	"workflow/model"
)

type ProcTaskDao struct {
}

// Add 创建流程任务
func (ProcTaskDao) Add(db *gorm.DB, task model.ProcTask) (err error) {
	err = db.Create(&task).Error
	return
}

// GetTaskByInstIdAndNodeId 获取当前节点任务
func (ProcTaskDao) GetTaskByInstIdAndNodeId(db *gorm.DB, procInstId, nodeId int) (task model.ProcTask, err error) {
	err = db.Where("proc_inst_id", procInstId).Where("node_id", nodeId).Find(&task).Error
	return
}

// Update 更新任务
func (ProcTaskDao) Update(db *gorm.DB, task model.ProcTask) (err error) {
	err = db.Updates(&task).Error
	return
}
