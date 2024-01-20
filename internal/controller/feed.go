package controller

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/service"
	"knowledge-service/pkg/consts"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

type FeedController struct{}

func (e *FeedController) GetInfo(ctx *gin.Context) {
	var query entity.GetFeedInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
	feed, err := feedD.Find(ctx, query.FeedID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	feedS := service.FeedService{}
	feedInfo, err := feedS.FormatFeed(ctx, feed, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, feedInfo)
}

func (e *FeedController) SearchFeedList(ctx *gin.Context) {
	var query entity.SearchFeedsListQuery
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
	feeds, err := feedD.FindList(
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
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	feedS := service.FeedService{}
	feedList, err := feedS.FormatFeedList(ctx, feeds, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetFeedListResp{
		Total: len(feedList),
		List:  feedList,
	}
	tools.RespSuccess(ctx, res)
}

func (e *FeedController) LikeFeed(ctx *gin.Context) {
	var payload entity.LikeFeedPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
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
