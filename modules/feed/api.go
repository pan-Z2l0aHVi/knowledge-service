package feed

type GetFeedListQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type FeedListItem struct {
	Feed
	AuthorInfo Author `json:"author_info"`
}

type GetFeedListResp struct {
	Total int64          `json:"total"`
	List  []FeedListItem `json:"list"`
}
