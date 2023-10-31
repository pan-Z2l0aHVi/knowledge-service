package feed

import (
	"knowledge-base-service/consts"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func (e *Feed) GetDetail(ctx *gin.Context) {
	var query GetFeedDetailQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	dao := FeedDao{}
	feed, err := dao.FindFeed(ctx, query.FeedID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	author, err := dao.FindAuthorByID(ctx, feed.AuthorID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	res := GetFeedDetailResp{
		FeedListItem: FeedListItem{
			Feed:       feed,
			AuthorInfo: author,
		},
	}
	tools.RespSuccess(ctx, res)
}

func (e *Feed) SearchFeedList(ctx *gin.Context) {
	var query SearchFeedsListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var asc int
	if query.SortType == "desc" {
		asc = -1
	} else if query.SortType == "asc" {
		asc = 1
	}
	dao := FeedDao{}
	feedList, err := dao.FindFeedList(
		ctx,
		query.Page,
		query.PageSize,
		query.Keywords,
		query.SortBy,
		asc,
		query.AuthorID,
	)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	list := make([]FeedListItem, len(feedList))
	for i, item := range feedList {
		author, err := dao.FindAuthorByID(ctx, item.AuthorID)
		if err != nil {
			tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
			return
		}
		list[i] = FeedListItem{
			Feed:       item,
			AuthorInfo: author,
		}
	}
	res := GetFeedListResp{
		Total: len(list),
		List:  list,
	}
	tools.RespSuccess(ctx, res)
}

func (e *Feed) PraiseFeed(ctx *gin.Context) {
	var payload PraiseFeedPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.FailCode, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	} else {
		tools.RespFail(ctx, consts.FailCode, "当前用户不存在", nil)
		return
	}
	dao := FeedDao{}
	liked, err := dao.CheckLiked(ctx, userID, payload.FeedID)
	if err != nil {
		tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
		return
	}
	switch payload.Event {
	case "like":
		if !liked {
			err := dao.Like(ctx, userID, payload.FeedID)
			if err != nil {
				tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
				return
			}
		}
		tools.RespSuccess(ctx, nil)
		return

	case "unlike":
		if liked {
			err := dao.UnLike(ctx, userID, payload.FeedID)
			if err != nil {
				tools.RespFail(ctx, consts.FailCode, err.Error(), nil)
				return
			}
		}
		tools.RespSuccess(ctx, nil)
		return

	default:
		tools.RespFail(ctx, consts.FailCode, "参数错误，event:"+payload.Event, nil)
	}
}
