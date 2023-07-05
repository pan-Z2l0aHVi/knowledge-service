package material

type Material struct {
	ID       string      `json:"_id"`
	Type     string      `json:"type"`
	URL      string      `json:"content"`
	Uploader interface{} `json:"uploader"`
}
