package feed

type GetFeedListQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetFeedListResp struct {
	Total int64  `json:"total"`
	List  []Feed `json:"list"`
}
