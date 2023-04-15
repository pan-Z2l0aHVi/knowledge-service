package doc

type GetInfoParams struct {
	DocID string `form:"doc_id" binding:"required"`
}

type GetInfoResp struct {
	Doc
	DocID string `json:"doc_id"`
}

func (*GetInfoResp) _ID() {
	return
}
