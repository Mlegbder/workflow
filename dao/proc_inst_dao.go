package dao

import (
	"gorm.io/gorm"
	"workflow/model"
)

type ProcInstDao struct {
}

// Add 创建流程实例
func (ProcInstDao) Add(db *gorm.DB, procInst *model.ProcInst) (err error) {
	err = db.Create(&procInst).Error
	return
}

// GetByProcInstById 根据ID获取流程实例
func (ProcInstDao) GetByProcInstById(db *gorm.DB, procInstId int) (def model.ProcInst, err error) {
	err = db.Where("id", procInstId).Find(&def).Error
	return
}

// Update 更新实例
func (ProcInstDao) Update(db *gorm.DB, inst model.ProcInst) (err error) {
	err = db.Updates(&inst).Error
	return
}
