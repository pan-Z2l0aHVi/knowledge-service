package doc

type GetInfoQuery struct {
	DocID string `form:"doc_id" binding:"required"`
}

type GetInfoResp struct {
	Doc
}

type CreatePayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Cover   string `json:"cover"`
}

type UpdatePayload struct {
	DocID   string `json:"doc_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Title   string `json:"title"`
	Cover   string `json:"cover"`
	Public  *bool  `json:"public"`
}

type DeletePayload struct {
	DocIDs []string `json:"doc_ids" binding:"required"`
}

type GetDocsQuery struct {
	AuthorID string `form:"author_id"`
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
}

type GetDocsResp struct {
	Total int64 `json:"total"`
	List  []Doc `json:"list"`
}

type GetDraftQuery struct {
	DocID    string `form:"doc_id" binding:"required"`
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
}

type UpdateDraftPayload struct {
	DocID   string `json:"doc_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}
