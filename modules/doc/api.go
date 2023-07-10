package doc

type GetInfoQuery struct {
	DocID string `form:"doc_id" binding:"required"`
}

type GetInfoResp struct {
	Doc
}
