package controller

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type SpaceController struct{}

// @Summary 获取空间详情
// @Description
// @Produce json
// @Param space_id query string true "空间ID"
// @Success 200 {object} entity.GetSpaceInfoResp "ok" "空间详情"
// @Router /space/info [get]
func (e *SpaceController) GetInfo(ctx *gin.Context) {
	var query entity.GetSpaceInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	spaceD := dao.SpaceDAO{}
	spaceInfo, err := spaceD.Find(ctx, query.SpaceID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetSpaceInfoResp{
		Space: spaceInfo,
	}
	tools.RespSuccess(ctx, res)
}

// @Summary 新建空间
// @Description 校验登录态
// @Produce json
// @Param request body entity.CreateSpacePayload true "名称和描述"
// @Success 200 {object} model.Space "ok" "空间详情"
// @Router /space/create [post]
func (e *SpaceController) Create(ctx *gin.Context) {
	var payload entity.CreateSpacePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var ownerID string
	if userID, exist := ctx.Get("uid"); exist {
		ownerID = userID.(string)
	}
	spaceD := dao.SpaceDAO{}
	spaceInfo, err := spaceD.Create(
		ctx,
		ownerID,
		payload.Name,
		payload.Desc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, spaceInfo)
}

// @Summary 更新空间
// @Description 校验登录态
// @Produce json
// @Param request body entity.UpdateSpacePayload true "空间ID、名称和描述"
// @Success 200 {object} model.Space "ok" "更新后的空间详情"
// @Router /space/update [post]
func (e *SpaceController) Update(ctx *gin.Context) {
	var payload entity.UpdateSpacePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	spaceD := dao.SpaceDAO{}
	spaceInfo, err := spaceD.Update(
		ctx,
		payload.SpaceID,
		payload.Name,
		payload.Desc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, spaceInfo)
}

// @Summary 删除空间
// @Description 校验登录态
// @Produce json
// @Param request body entity.DeleteSpacePayload true "空间ID列表"
// @Success 200 "ok"
// @Router /space/delete [post]
func (e *SpaceController) Delete(ctx *gin.Context) {
	var payload entity.DeleteSpacePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	spaceD := dao.SpaceDAO{}
	if err := spaceD.DeleteMany(ctx, payload.SpaceIDs); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	for _, spaceID := range payload.SpaceIDs {
		if err := deleteSpaceContent(ctx, spaceID); err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
	}
	tools.RespSuccess(ctx, nil)
}

// @Summary 搜索空间
// @Description 校验登录态
// @Produce json
// @Param query query entity.SearchSpacesQuery true "query参数"
// @Success 200 {object} entity.GetSpacesResp "ok" "空间列表"
// @Router /space/search [get]
func (e *SpaceController) SearchSpaces(ctx *gin.Context) {
	var query entity.SearchSpacesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.OwnerID != "" {
		userID = query.OwnerID
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
	spaceD := dao.SpaceDAO{}
	spaces, total, err := spaceD.FindListWithTotal(ctx,
		query.Page,
		query.PageSize,
		userID,
		query.Keywords,
		query.SortBy,
		asc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetSpacesResp{
		Total: total,
		List:  spaces,
	}
	tools.RespSuccess(ctx, res)
}

func deleteSpaceContent(ctx *gin.Context, spaceID string) error {
	docD := dao.DocDAO{}
	docs, err := docD.FindManyBySpace(ctx, spaceID)
	if err != nil {
		return err
	}
	docIDs := []string{}
	for _, doc := range docs {
		docIDs = append(docIDs, doc.ID.Hex())
	}
	docS := service.DocService{}
	if err := docS.DeleteDocs(ctx, docIDs); err != nil {
		return err
	}
	return nil
}
