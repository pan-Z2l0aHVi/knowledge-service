package material

type GetInfoParams struct {
	MaterialID string `form:"material_id" binding:"required"`
}

type GetInfoResp struct {
	Material
	MaterialID string `json:"material_id"`
}

func (*GetInfoResp) _ID() {
	return
}
