package doc

type Doc struct {
	ID      string      `json:"_id"`
	Content string      `json:"content"`
	Author  interface{} `json:"author"`
}
