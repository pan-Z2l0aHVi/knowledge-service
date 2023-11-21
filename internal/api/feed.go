package api

import "knowledge-service/internal/model"

type GetFeedDetailQuery struct {
	FeedID string `form:"feed_id"`
}

type GetFeedDetailResp struct {
	FeedItem
}

type SearchFeedsListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keywords string `form:"keywords"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
	AuthorID string `form:"author_id"`
}

type FeedItem struct {
	model.Feed
	AuthorInfo model.Author `json:"author_info"`
}

type GetFeedListResp struct {
	Total int        `json:"total"`
	List  []FeedItem `json:"list"`
}

type PraiseFeedPayload struct {
	Event  string `json:"event"`
	FeedID string `json:"feed_id"`
}
