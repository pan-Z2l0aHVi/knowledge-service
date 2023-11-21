package controller

import (
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type MaterialController struct{}

func (e *MaterialController) GetInfo(ctx *gin.Context) {
	var query api.GetMaterialInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialInfo, err := materialD.Find(ctx, query.MaterialID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := api.GetMaterialInfoResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *MaterialController) Upload(ctx *gin.Context) {
	var payload api.UploadMaterialPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialInfo, err := materialD.Create(ctx, payload.Type, payload.URL)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := api.UploadMaterialResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *MaterialController) Search(ctx *gin.Context) {
	var query api.SearchMaterialQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialList, err := materialD.Search(ctx, query.Type, query.Keywords, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	if len(materialList) == 0 {
		materialList = []model.Material{}
	}
	res := api.SearchMaterialResp{
		Total: len(materialList),
		List:  materialList,
	}
	tools.RespSuccess(ctx, res)
}
