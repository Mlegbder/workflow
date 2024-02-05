package dao

import (
	"github.com/Mlegbder/workflow/model"
	"gorm.io/gorm"
)

type ProcDefDao struct {
}

// GetByProcDefById 根据ID获取流程定义
func (ProcDefDao) GetByProcDefById(db *gorm.DB, procDefId int) (def model.ProcDef, err error) {
	err = db.Where("id", procDefId).Find(&def).Error
	return
}
