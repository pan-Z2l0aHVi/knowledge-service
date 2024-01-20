package entity

import "knowledge-service/internal/model"

type GetSpaceInfoQuery struct {
	SpaceID string `form:"space_id" binding:"required"`
}

type GetSpaceInfoResp struct {
	model.Space
}

type CreateSpacePayload struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type UpdateSpacePayload struct {
	SpaceID string  `json:"space_id" binding:"required"`
	Name    *string `json:"name"`
	Desc    *string `json:"desc"`
}

type DeleteSpacePayload struct {
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
	Total int           `json:"total"`
	List  []model.Space `json:"list"`
}
