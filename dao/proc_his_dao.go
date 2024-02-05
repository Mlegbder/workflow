package dao

import (
	"github.com/Mlegbder/workflow/model"
	"github.com/Mlegbder/workflow/types"
	"gorm.io/gorm"
)

type ProcHisDao struct {
}

// Add 创建流程历史记录
func (ProcHisDao) Add(db *gorm.DB, his model.ProcHis) (err error) {
	err = db.Create(&his).Error
	return
}

// ListByProcInstId 流程历史记录
func (ProcHisDao) ListByProcInstId(db *gorm.DB, procInstId int) (resp *[]types.ProcHisResp, err error) {
	err = db.Table(model.TableNameProcHis).
		Where("proc_inst_id", procInstId).
		Select("`id`, `assignee`, `proc_inst_id`, `node_id`, `approval_status`, `created_at` created_at_str").
		Find(&resp).Error
	return
}
