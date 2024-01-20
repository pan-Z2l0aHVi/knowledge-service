package entity

import "knowledge-service/internal/model"

type UploadMaterialPayload struct {
	Type int    `json:"type" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type UploadMaterialResp struct {
	model.Material
}

type GetMaterialInfoQuery struct {
	MaterialID string `form:"material_id" binding:"required"`
}

type GetMaterialInfoResp struct {
	model.Material
}

type SearchMaterialQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Type     int    `form:"type" binding:"required"`
	Keywords string `form:"keywords"`
}

type SearchMaterialResp struct {
	List  []model.Material `json:"list"`
	Total int              `json:"total"`
}
