package material

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Material) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.Find(ctx, query.MaterialID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Material) Upload(ctx *gin.Context) {
	var payload UploadPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.Create(ctx, payload.Type, payload.URL)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := UploadResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Material) Search(ctx *gin.Context) {
	var query MaterialSearchQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialList, err := dao.Search(ctx, query.Type, query.Keywords, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	total, err := dao.GetCount(ctx, query.Type, query.Keywords)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if len(materialList) == 0 {
		materialList = []Material{}
	}
	res := MaterialSearchResp{
		Total: total,
		List:  materialList,
	}
	tools.RespSuccess(ctx, res)
}
