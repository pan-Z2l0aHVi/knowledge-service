package doc

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Doc) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.Find(ctx, query.DocID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Doc: docInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Doc) Create(ctx *gin.Context) {
	var payload CreatePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var authorID string
	if userID, exist := ctx.Get("uid"); exist {
		authorID = userID.(string)
	}
	dao := DocDAO{}
	docInfo, err := dao.Create(
		ctx,
		authorID,
		payload.Title,
		payload.Content,
		payload.Cover,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *Doc) Update(ctx *gin.Context) {
	var payload UpdatePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	docInfo, err := dao.Update(
		ctx,
		payload.DocID,
		payload.Title,
		payload.Content,
		payload.Cover,
		payload.Public,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, docInfo)
}

func (e *Doc) Delete(ctx *gin.Context) {
	var payload DeletePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	err := dao.Delete(ctx, payload.DocIDs)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *Doc) GetDocs(ctx *gin.Context) {
	var query GetDocsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.AuthorID != "" {
		userID = query.AuthorID
	} else if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	dao := DocDAO{}
	docs, err := dao.FindDocsByAuthor(ctx, userID, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	count, err := dao.FindCountByAuthor(ctx, userID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetDocsResp{
		Total: count,
		List:  docs,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Doc) GetDrafts(ctx *gin.Context) {
	var query GetDraftQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	drafts, err := dao.FindDraftsByDoc(ctx, query.DocID, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, drafts)
}

func (e *Doc) UpdateDraft(ctx *gin.Context) {
	var payload UpdateDraftPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := DocDAO{}
	draft, err := dao.UpdateDraft(ctx, payload.DocID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, draft)
}
