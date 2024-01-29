package controller

import (
	"knowledge-service/internal/dao"
	"knowledge-service/internal/entity"
	"knowledge-service/internal/model"
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
	tools.RespSuccess(ctx, entity.GetFeedInfoResp{
		FeedInfo: feedInfo,
	})
}

func (e *FeedController) SearchFeedList(ctx *gin.Context) {
	var query entity.SearchFeedsListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
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
	feedD := dao.FeedDao{}
	feeds := []model.Feed{}
	var feedsTotal int64 = 0
	if query.Keywords != "" {
		trueBool := true
		docD := dao.DocDAO{}
		docs, total, err := docD.FindListWithTotal(
			ctx,
			query.Page,
			query.PageSize,
			query.AuthorID,
			"",
			query.Keywords,
			query.SortBy,
			asc,
			&trueBool,
		)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		feedsTotal = total
		for _, doc := range docs {
			feed, err := feedD.FindBySubject(ctx, doc.ID.Hex(), "doc")
			if err != nil {
				tools.RespFail(ctx, consts.Fail, err.Error(), nil)
				return
			}
			feeds = append(feeds, feed)
		}
	} else {
		list, total, err := feedD.FindListWithTotal(
			ctx,
			query.Page,
			query.PageSize,
			query.SortBy,
			asc,
			query.AuthorID,
		)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		feedsTotal = total
		feeds = list
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
		Total: feedsTotal,
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

func (e *FeedController) Comment(ctx *gin.Context) {
	var payload entity.CommentPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	feedD := dao.FeedDao{}
	comment, err := feedD.CreateComment(ctx, payload.FeedID, userID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feedS := service.FeedService{}
	commentInfo, err := feedS.FormatComment(ctx, comment)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, entity.CommentResp{
		CommentInfo: commentInfo,
	})
}

func (e *FeedController) GetCommentList(ctx *gin.Context) {
	var query entity.GetCommentListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
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
	comments, total, err := feedD.FindCommentListWithTotal(ctx, query.FeedID, query.Page, query.PageSize, query.SortBy, asc)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feedS := service.FeedService{}
	commentList, err := feedS.FormatComments(ctx, comments)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetCommentListResp{
		Total: total,
		List:  commentList,
	}
	tools.RespSuccess(ctx, res)
}

func (e *FeedController) Reply(ctx *gin.Context) {
	var payload entity.ReplyPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	feedD := dao.FeedDao{}
	subComment, err := feedD.ReplyComment(ctx, payload.FeedID, payload.CommentID, userID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feedS := service.FeedService{}
	replyInfo, err := feedS.FormatSubComment(ctx, subComment)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, replyInfo)
}

func (e *FeedController) DeleteComment(ctx *gin.Context) {
	var payload entity.DeleteCommentPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
	if err := feedD.DeleteComment(ctx, payload.FeedID, payload.CommentID, payload.SubCommentID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

func (e *FeedController) UpdateComment(ctx *gin.Context) {
	var payload entity.UpdateCommentPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDao{}
	feedS := service.FeedService{}
	if payload.SubCommentID != "" {
		subComment, err := feedD.UpdateSubComment(ctx, payload.FeedID, payload.CommentID, payload.SubCommentID, payload.Content)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		subCommentInfo, err := feedS.FormatSubComment(ctx, subComment)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		tools.RespSuccess(ctx, subCommentInfo)
		return
	}
	comment, err := feedD.UpdateComment(ctx, payload.FeedID, payload.CommentID, payload.Content)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	commentInfo, err := feedS.FormatComment(ctx, comment)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, commentInfo)
}
