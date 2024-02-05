package workflow

import (
	"github.com/Mlegbder/workflow/dao"
	"github.com/Mlegbder/workflow/types"
	"gorm.io/gorm"
)

// WorkFlowHis 流程历史
func WorkFlowHis(db *gorm.DB, req *types.ProcReq) (resp *[]types.ProcHisResp, err error) {
	resp, err = dao.ProcHisDao{}.ListByProcInstId(db, req.Id)
	return
}
