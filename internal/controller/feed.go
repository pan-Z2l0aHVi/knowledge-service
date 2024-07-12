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

// @Summary 获取动态详情
// @Description 使用但不校验登录态
// @Produce json
// @Param feed_id query string true "动态ID"
// @Success 200 {object} entity.GetFeedInfoResp "ok" "动态详情"
// @Router /feed/info [get]
func (e *FeedController) GetInfo(ctx *gin.Context) {
	var query entity.GetFeedInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDAO{}
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

// @Summary 搜索动态
// @Description 使用但不校验登录态
// @Produce json
// @Param query query entity.SearchFeedsListQuery true "query参数"
// @Success 200 {object} entity.GetFeedListResp "ok" "动态列表"
// @Router /feed/list [get]
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
	feedD := dao.FeedDAO{}
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
			feed, err := feedD.FindBySubject(ctx, doc.ID.Hex(), consts.DocFeed)
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

// @Summary 点赞|取消点赞 动态
// @Description 校验登录态
// @Produce json
// @Param request body entity.LikeFeedPayload true "动态ID、点赞或取消点赞"
// @Success 200 "ok"
// @Router /feed/like [post]
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
	feedD := dao.FeedDAO{}
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

// @Summary 评论动态
// @Description 校验登录态
// @Produce json
// @Param request body entity.CommentPayload true "动态ID、评论内容"
// @Success 200 {object} entity.CommentResp "ok" "评论详情"
// @Router /feed/comment [post]
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
	feedD := dao.FeedDAO{}
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

// @Summary 获取评论列表
// @Description
// @Produce json
// @Param query query entity.GetCommentListQuery true "query参数"
// @Success 200 {object} entity.GetCommentListResp "ok" "评论列表"
// @Router /feed/comment_list [get]
func (e *FeedController) GetCommentList(ctx *gin.Context) {
	var query entity.GetCommentListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDAO{}
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

// @Summary 回复评论
// @Description 校验登录态
// @Produce json
// @Param request body entity.ReplyPayload true "动态ID、评论ID、目标用户ID、回复内容"
// @Success 200 {object} entity.ReplyInfo "ok" "回复评论详情"
// @Router /feed/reply [post]
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
	feedD := dao.FeedDAO{}
	subComment, err := feedD.ReplyComment(
		ctx,
		payload.FeedID,
		payload.CommentID,
		payload.ReplyUserID,
		userID,
		payload.Content,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feedS := service.FeedService{}
	replyInfo, err := feedS.FormatSubComment(ctx, payload.FeedID, payload.CommentID, subComment)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, replyInfo)
}

// @Summary 删除评论
// @Description 校验登录态
// @Produce json
// @Param request body entity.DeleteCommentPayload true "动态ID、主评论ID、子评论ID"
// @Success 200 "ok"
// @Router /feed/comment_delete [post]
func (e *FeedController) DeleteComment(ctx *gin.Context) {
	var payload entity.DeleteCommentPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDAO{}
	if err := feedD.DeleteComment(ctx, payload.FeedID, payload.CommentID, payload.SubCommentID); err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, nil)
}

// @Summary 更新评论
// @Description 校验登录态
// @Produce json
// @Param request body entity.UpdateCommentPayload true "动态ID、主评论ID、子评论ID、评论内容"
// @Success 200 {object} entity.CommentInfo "ok" "更新后的评论详情"
// @Router /feed/comment_update [post]
func (e *FeedController) UpdateComment(ctx *gin.Context) {
	var payload entity.UpdateCommentPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	feedD := dao.FeedDAO{}
	feedS := service.FeedService{}
	if payload.SubCommentID != "" {
		subComment, err := feedD.UpdateSubComment(ctx, payload.FeedID, payload.CommentID, payload.SubCommentID, payload.Content)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		subCommentInfo, err := feedS.FormatSubComment(ctx, payload.FeedID, payload.CommentID, subComment)
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

// @Summary 获取同一空间下相关动态列表
// @Description 使用但不校验登录态。
// @Produce json
// @Param query query entity.GetRelatedFeedsQuery true "query参数"
// @Success 200 {object} entity.GetRelatedFeedsResp "ok" "动态列表"
// @Router /feed/related_feeds [get]
func (e *FeedController) GetRelatedFeeds(ctx *gin.Context) {
	var query entity.GetRelatedFeedsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, consts.Fail, "参数错误:"+err.Error(), nil)
		return
	}
	var userID string
	if uid, exist := ctx.Get("uid"); exist {
		userID = uid.(string)
	}
	trueBool := true
	docD := dao.DocDAO{}
	feedD := dao.FeedDAO{}
	docs, total, err := docD.FindListWithTotal(
		ctx,
		query.Page,
		query.PageSize,
		"",
		query.SpaceID,
		"",
		"update_time",
		-1,
		&trueBool,
	)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	feeds := []model.Feed{}
	for _, doc := range docs {
		feed, err := feedD.FindBySubject(ctx, doc.ID.Hex(), consts.DocFeed)
		if err != nil {
			tools.RespFail(ctx, consts.Fail, err.Error(), nil)
			return
		}
		feeds = append(feeds, feed)
	}
	feedS := service.FeedService{}
	feedList, err := feedS.FormatFeedList(ctx, feeds, userID)
	if err != nil {
		tools.RespFail(ctx, consts.Fail, err.Error(), nil)
		return
	}
	res := entity.GetRelatedFeedsResp{
		Total: total,
		List:  feedList,
	}
	tools.RespSuccess(ctx, res)
}
