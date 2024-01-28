package entity

import (
	"knowledge-service/internal/model"
)

type GetFeedInfoQuery struct {
	FeedID string `form:"feed_id" binding:"required"`
}

type GetFeedInfoResp struct {
	FeedInfo
}

type SearchFeedsListQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Keywords string `form:"keywords"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
	AuthorID string `form:"author_id"`
}

type LikeInfo struct {
	model.Like
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type FeedInfo struct {
	model.Feed
	Likes     []LikeInfo `json:"likes"`
	Creator   Creator    `json:"creator"`
	Subject   DocInfo    `json:"subject"`
	Collected bool       `json:"collected"`
}

type Creator struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type GetFeedListResp struct {
	Total int64      `json:"total"`
	List  []FeedInfo `json:"list"`
}

type LikeFeedPayload struct {
	Event  string `json:"event" binding:"required"`
	FeedID string `json:"feed_id" binding:"required"`
}

type Commentator struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ReplyInfo struct {
	model.SubComment
	Commentator Commentator `json:"commentator"`
}

type CommentInfo struct {
	model.Comment
	Commentator Commentator `json:"commentator"`
	SubComments []ReplyInfo `json:"sub_comments"`
}

type GetCommentListQuery struct {
	FeedID   string `form:"feed_id" binding:"required"`
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
}

type CommentPayload struct {
	FeedID  string `json:"feed_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ReplyPayload struct {
	FeedID    string `json:"feed_id" binding:"required"`
	CommentID string `json:"comment_id"`
	Content   string `json:"content" binding:"required"`
}

type DeleteCommentPayload struct {
	FeedID       string `json:"feed_id" binding:"required"`
	CommentID    string `json:"comment_id"`
	SubCommentID string `json:"sub_comment_id"`
}

type UpdateCommentPayload struct {
	FeedID       string `json:"feed_id" binding:"required"`
	CommentID    string `json:"comment_id"`
	SubCommentID string `json:"sub_comment_id"`
	Content      string `json:"content"`
}

type GetCommentListResp struct {
	Total int           `json:"total"`
	List  []CommentInfo `json:"list"`
}

type CommentResp struct {
	CommentInfo
}

type UpdateCommentResp struct {
	CommentInfo
	ReplyCommentID string `json:"reply_comment_id" bson:"reply_comment_id"`
}
