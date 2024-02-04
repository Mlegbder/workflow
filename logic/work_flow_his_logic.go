package workflow

import (
	"gorm.io/gorm"
	"workflow/dao"
	"workflow/types"
)

// WorkFlowHis 流程历史
func WorkFlowHis(db *gorm.DB, req *types.ProcReq) (resp *[]types.ProcHisResp, err error) {
	resp, err = dao.ProcHisDao{}.ListByProcInstId(db, req.Id)
	return
}
