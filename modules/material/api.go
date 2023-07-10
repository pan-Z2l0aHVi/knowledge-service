package material

type UploadPayload struct {
	URL string `form:"url" binding:"required"`
}

type UploadResp struct {
	URL string `json:"url"`
}

type GetInfoQuery struct {
	MaterialID string `form:"material_id" binding:"required"`
}

type GetInfoResp struct {
	Material
}

type MaterialSearchQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Type     int    `form:"type" binding:"required"`
	Keywords string `form:"keywords"`
}

type MaterialSearchResp struct {
	Data []Material `json:"data"`
}
