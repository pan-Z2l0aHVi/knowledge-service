package controller

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type MaterialController struct{}

func (e *MaterialController) GetInfo(ctx *gin.Context) {
	var query entity.GetMaterialInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialInfo, err := materialD.Find(ctx, query.MaterialID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetMaterialInfoResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *MaterialController) Upload(ctx *gin.Context) {
	var payload entity.UploadMaterialPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialInfo, err := materialD.Create(ctx, payload.Type, payload.URL)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.UploadMaterialResp{
		Material: materialInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *MaterialController) Search(ctx *gin.Context) {
	var query entity.SearchMaterialQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	materialD := dao.MaterialDAO{}
	materialList, err := materialD.FindList(ctx, query.Type, query.Keywords, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	if len(materialList) == 0 {
		materialList = []model.Material{}
	}
	res := entity.SearchMaterialResp{
		Total: len(materialList),
		List:  materialList,
	}
	tools.RespSuccess(ctx, res)
}
