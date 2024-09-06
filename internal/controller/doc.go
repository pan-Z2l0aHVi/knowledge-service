package controller

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type DocController struct{}

// @Summary 获取文档详情
// @Description
// @Produce json
// @Param doc_id query string true "文档ID"
// @Success 200 {object} entity.GetDocInfoResp "ok" "文档详情"
// @Router /doc/info [get]
func (e *DocController) GetInfo(ctx *gin.Context) {
	var query entity.GetDocInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	doc, err := docD.Find(ctx, query.DocID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	docS := service.DocService{}
	docInfo, err := docS.FormatDoc(ctx, doc)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetDocInfoResp{
		DocInfo: docInfo,
	}
	tools.RespSuccess(ctx, res)
}

// @Summary 新建文档
// @Description 校验登录态
// @Produce json
// @Param request body entity.CreateDocPayload true "空间ID、标题、封面、内容"
// @Success 200 {object} model.Doc "ok" "文档详情"
// @Router /doc/create [post]
func (e *DocController) Create(ctx *gin.Context) {
	var payload entity.CreateDocPayload
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

// @Summary 更新文档
// @Description 校验登录态
// @Produce json
// @Param request body entity.UpdateDocPayload true "文档ID、空间ID、标题、封面、内容、摘要"
// @Success 200 {object} model.Doc "ok" "更新后的文档详情"
// @Router /doc/update [post]
func (e *DocController) Update(ctx *gin.Context) {
	var payload entity.UpdateDocPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var authorID string
	if userID, exist := ctx.Get("uid"); exist {
		authorID = userID.(string)
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
	if payload.Public != nil {
		feedS := service.FeedService{}
		if *payload.Public {
			_, err := feedS.SyncFeed(ctx, authorID, docInfo.AuthorID, docInfo.ID.Hex(), consts.DocFeed)
			if err != nil {
				tools.RespFail(ctx, consts.Fail, err.Error(), nil)
				return
			}
		} else {
			feedD := dao.FeedDAO{}
			subjectIDs := append([]string{}, docInfo.ID.Hex())
			delErr := feedD.DeleteManyBySubject(ctx, subjectIDs, consts.DocFeed)
			if delErr != nil {
				tools.RespFail(ctx, consts.Fail, delErr.Error(), nil)
				return
			}
		}
		err := feedS.RemoveAllFeedListCache()
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
	}
	tools.RespSuccess(ctx, docInfo)
}

// @Summary 删除文档
// @Description 校验登录态
// @Produce json
// @Param request body entity.DeleteDocPayload true "文档ID列表"
// @Success 200 "ok"
// @Router /doc/delete [post]
func (e *DocController) Delete(ctx *gin.Context) {
	var payload entity.DeleteDocPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docS := service.DocService{}
	err := docS.DeleteDocs(ctx, payload.DocIDs)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

// @Summary 搜索文档
// @Description 使用但不校验登录态
// @Produce json
// @Param query query entity.SearchDocsQuery true "query参数"
// @Success 200 {object} entity.GetDocsResp "ok" "文档列表"
// @Router /doc/docs [get]
func (e *DocController) SearchDocs(ctx *gin.Context) {
	var query entity.SearchDocsQuery
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
	if query.SortType == "" {
		query.SortType = "desc"
	}
	var asc int
	if query.SortType == "desc" {
		asc = -1
	} else if query.SortType == "asc" {
		asc = 1
	}
	if query.SortBy == "" {
		query.SortBy = "update_time"
	}
	docD := dao.DocDAO{}
	docs, total, err := docD.FindListWithTotal(ctx,
		query.Page,
		query.PageSize,
		userID,
		query.SpaceID,
		query.Keywords,
		query.SortBy,
		asc,
		nil,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	docS := service.DocService{}
	docList := []entity.DocInfo{}
	for _, doc := range docs {
		docInfo, err := docS.FormatDoc(ctx, doc)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		docList = append(docList, docInfo)
	}
	res := entity.GetDocsResp{
		Total: total,
		List:  docList,
	}
	tools.RespSuccess(ctx, res)
}

// @Summary 获取文档草稿
// @Description 校验登录态
// @Produce json
// @Param query query entity.GetDraftQuery true "文档ID和分页参数"
// @Success 200 {array} model.Draft "ok" "草稿列表"
// @Router /doc/drafts [get]
func (e *DocController) GetDrafts(ctx *gin.Context) {
	var query entity.GetDraftQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	drafts, err := docD.FindDrafts(ctx, query.DocID, query.Page, query.PageSize)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, drafts)
}

// @Summary 更新文档草稿
// @Description 校验登录态
// @Produce json
// @Param request body entity.UpdateDraftPayload true "文档ID和草稿内容"
// @Success 200 {object} model.Draft "ok" "更新后的草稿详情"
// @Router /doc/update_drafts [post]
func (e *DocController) UpdateDrafts(ctx *gin.Context) {
	var payload entity.UpdateDraftPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	docD := dao.DocDAO{}
	draft, err := docD.CreateDraft(ctx, payload.DocID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, draft)
}
