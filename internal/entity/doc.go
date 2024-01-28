package entity

import (
	"knowledge-service/internal/model"
)

type GetDocInfoQuery struct {
	DocID string `form:"doc_id" binding:"required"`
}

type GetDocInfoResp struct {
	DocInfo
}

type DocInfo struct {
	model.Doc
	Author Author `json:"author"`
}

type CreateDocPayload struct {
	SpaceID string `json:"space_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Cover   string `json:"cover"`
}

type UpdateDocPayload struct {
	DocID   string  `json:"doc_id" binding:"required"`
	Content *string `json:"content"`
	Summary *string `json:"summary"`
	Title   *string `json:"title"`
	Cover   *string `json:"cover"`
	Public  *bool   `json:"public"`
}

type DeleteDocPayload struct {
	DocIDs []string `json:"doc_ids" binding:"required"`
}

type SearchDocsQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	AuthorID string `form:"author_id"`
	SpaceID  string `form:"space_id"`
	SortBy   string `form:"sort_by"`
	SortType string `form:"sort_type"`
	Keywords string `form:"keywords"`
}

type GetDocsResp struct {
	Total int64     `json:"total"`
	List  []DocInfo `json:"list"`
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

type Author struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
