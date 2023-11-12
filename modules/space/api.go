package space

type GetInfoQuery struct {
	SpaceID string `form:"space_id" binding:"required"`
}

type GetInfoResp struct {
	Space
}

type CreatePayload struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type UpdatePayload struct {
	SpaceID string  `json:"space_id" binding:"required"`
	Name    *string `json:"name"`
	Desc    *string `json:"desc"`
}

type DeletePayload struct {
	SpaceIDs []string `json:"space_ids" binding:"required"`
}

type SearchSpacesQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	OwnerID  string `form:"owner_id"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
	Keywords string `form:"keywords"`
}

type GetSpacesResp struct {
	Total int     `json:"total"`
	List  []Space `json:"list"`
}
