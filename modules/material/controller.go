package material

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Material) getInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.find(ctx, query.MaterialID)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Material) upload(ctx *gin.Context) {
	var payload UploadPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.create(ctx, payload.Type, payload.URL)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	res := UploadResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Material) search(ctx *gin.Context) {
	var query MaterialSearchQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialList, err := dao.search(ctx, query.Type, query.Keywords, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	total, err := dao.getCount(ctx, query.Type, query.Keywords)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
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
