package controller

import (
	"knowledge-service/internal/api"
	"knowledge-service/internal/dao"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type FeedController struct{}

func (e *FeedController) GetDetail(ctx *gin.Context) {
	var query api.GetFeedDetailQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
	feed, err := feedD.FindFeed(ctx, query.FeedID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	author, err := feedD.FindAuthorByID(ctx, feed.AuthorID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := api.GetFeedDetailResp{
		FeedItem: api.FeedItem{
			Feed:       feed,
			AuthorInfo: author,
		},
	}
	tools.RespSuccess(ctx, res)
}

func (e *FeedController) SearchFeedList(ctx *gin.Context) {
	var query api.SearchFeedsListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var asc int
	if query.SortType == "desc" {
		asc = -1
	} else if query.SortType == "asc" {
		asc = 1
	}
	feedD := dao.FeedDao{}
	feedList, err := feedD.FindFeedList(
		ctx,
		query.Page,
		query.PageSize,
		query.Keywords,
		query.SortBy,
		asc,
		query.AuthorID,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	list := make([]api.FeedItem, len(feedList))
	for i, item := range feedList {
		author, err := feedD.FindAuthorByID(ctx, item.AuthorID)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		list[i] = api.FeedItem{
			Feed:       item,
			AuthorInfo: author,
		}
	}
	res := api.GetFeedListResp{
		Total: len(list),
		List:  list,
	}
	tools.RespSuccess(ctx, res)
}

func (e *FeedController) PraiseFeed(ctx *gin.Context) {
	var payload api.PraiseFeedPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.Fail, "当前用户不存在", nil)
		return
	}
	feedD := dao.FeedDao{}
	liked, err := feedD.CheckLiked(ctx, userID, payload.FeedID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	switch payload.Event {
	case "like":
		if !liked {
			err := feedD.Like(ctx, userID, payload.FeedID)
			if err != nil {
				tools.RespFail(ctx, consts.Fail, err.Error(), nil)
				return
			}
		}
		tools.RespSuccess(ctx, nil)
		return

	case "unlike":
		if liked {
			err := feedD.UnLike(ctx, userID, payload.FeedID)
			if err != nil {
				tools.RespFail(ctx, consts.Fail, err.Error(), nil)
				return
			}
		}
		tools.RespSuccess(ctx, nil)
		return

	default:
		tools.RespFail(ctx, consts.Fail, "参数错误，event:"+payload.Event, nil)
	}
}
