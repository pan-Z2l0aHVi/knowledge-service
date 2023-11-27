package controller

import (
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type DocController struct{}

func (e *DocController) GetInfo(ctx *gin.Context) {
	var query api.GetDocInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	docInfo, err := docD.Find(ctx, query.DocID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := api.GetDocInfoResp{
		Doc: docInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *DocController) Create(ctx *gin.Context) {
	var payload api.CreateDocPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var authorID string
	if userID, exist := ctx.Get("uid"); exist {
		authorID = userID.(string)
	}
	docD := dao.DocDAO{}
	docInfo, err := docD.Create(
		ctx,
		authorID,
		payload.SpaceID,
		payload.Title,
		payload.Content,
		payload.Cover,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *DocController) Update(ctx *gin.Context) {
	var payload api.UpdateDocPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	docInfo, err := docD.Update(
		ctx,
		payload.DocID,
		payload.Title,
		payload.Content,
		payload.Summary,
		payload.Cover,
		payload.Public,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *DocController) Delete(ctx *gin.Context) {
	var payload api.DeleteDocPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	err := docD.Delete(ctx, payload.DocIDs)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *DocController) SearchDocs(ctx *gin.Context) {
	var query api.SearchDocsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.AuthorID != "" {
		userID = query.AuthorID
	} else if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	var asc int
	if query.SortType == "desc" {
		asc = -1
	} else if query.SortType == "asc" {
		asc = 1
	}
	docD := dao.DocDAO{}
	docs, err := docD.FindDocs(ctx,
		query.Page,
		query.PageSize,
		userID,
		query.SpaceID,
		query.Keywords,
		query.SortBy,
		asc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := api.GetDocsResp{
		Total: len(docs),
		List:  docs,
	}
	tools.RespSuccess(ctx, res)
}

func (e *DocController) GetDrafts(ctx *gin.Context) {
	var query api.GetDraftQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	drafts, err := docD.FindDraftsByDoc(ctx, query.DocID, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, drafts)
}

func (e *DocController) UpdateDraft(ctx *gin.Context) {
	var payload api.UpdateDraftPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	draft, err := docD.UpdateDraft(ctx, payload.DocID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, draft)
}
