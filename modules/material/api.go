package material

type UploadPayload struct {
	Type int    `form:"type" binding:"required"`
	URL  string `form:"url" binding:"required"`
}

type UploadResp struct {
	Material
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
	List  []Material `json:"list"`
	Total int64      `json:"total"`
}
