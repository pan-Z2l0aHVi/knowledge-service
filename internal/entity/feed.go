package entity

import (
	"knowledge-service/internal/model"
)

type GetFeedInfoQuery struct {
	FeedID string `form:"feed_id"`
}

type GetFeedInfoResp struct {
	FeedInfo
}

type SearchFeedsListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
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
	Total int        `json:"total"`
	List  []FeedInfo `json:"list"`
}

type LikeFeedPayload struct {
	Event  string `json:"event"`
	FeedID string `json:"feed_id"`
}
