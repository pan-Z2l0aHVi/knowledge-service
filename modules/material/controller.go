package material

import (
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Material) getInfo(ctx *gin.Context) {
	var params GetInfoQuery
	if err := ctx.ShouldBindQuery(&params); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialInfo, err := dao.find(ctx, params.MaterialID)
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
	materialInfo, err := dao.create(ctx)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
	}
	tools.RespSuccess(ctx, materialInfo)
}

func (e *Material) search(ctx *gin.Context) {
	var params MaterialSearchQuery
	if err := ctx.ShouldBindQuery(&params); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	dao := MaterialDAO{}
	materialList, err := dao.search(ctx, params.Type, params.Keywords)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
	}
	materialLen := len(materialList)
	if materialLen > 0 {
		tools.RespSuccess(ctx, materialList)
	} else {
		tools.RespSuccess(ctx, []Material{})
	}
}
