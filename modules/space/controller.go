package space

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Space) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := SpaceDAO{}
	spaceInfo, err := dao.Find(ctx, query.SpaceID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetInfoResp{
		Space: spaceInfo,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Space) Create(ctx *gin.Context) {
	var payload CreatePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var ownerID string
	if userID, exist := ctx.Get("uid"); exist {
		ownerID = userID.(string)
	}
	dao := SpaceDAO{}
	spaceInfo, err := dao.Create(
		ctx,
		ownerID,
		payload.Name,
		payload.Desc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, spaceInfo)
}

func (e *Space) Update(ctx *gin.Context) {
	var payload UpdatePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := SpaceDAO{}
	spaceInfo, err := dao.Update(
		ctx,
		payload.SpaceID,
		payload.Name,
		payload.Desc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, spaceInfo)
}

func (e *Space) Delete(ctx *gin.Context) {
	var payload DeletePayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := SpaceDAO{}
	err := dao.Delete(ctx, payload.SpaceIDs)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *Space) SearchSpaces(ctx *gin.Context) {
	var query SearchSpacesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if query.OwnerID != "" {
		userID = query.OwnerID
	} else if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.FailCode, "当前用户不存在", nil)
		return
	}
	var asc int
	if query.SortType == "desc" {
		asc = -1
	} else if query.SortType == "asc" {
		asc = 1
	}
	dao := SpaceDAO{}
	spaces, err := dao.FindSpaces(ctx,
		query.Page,
		query.PageSize,
		userID,
		query.Keywords,
		query.SortBy,
		asc,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetSpacesResp{
		Total: len(spaces),
		List:  spaces,
	}
	tools.RespSuccess(ctx, res)
}
