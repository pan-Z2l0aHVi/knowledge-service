package feed

type GetFeedDetailQuery struct {
	FeedID string `form:"feed_id"`
}

type GetFeedDetailResp struct {
	FeedListItem
}

type SearchFeedsListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keywords string `form:"keywords"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
}

type FeedListItem struct {
	Feed
	AuthorInfo Author `json:"author_info"`
}

type GetFeedListResp struct {
	Total int            `json:"total"`
	List  []FeedListItem `json:"list"`
}

type PraiseFeedPayload struct {
	Event  string `json:"event"`
	FeedID string `json:"feed_id"`
}
