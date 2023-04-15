package user

type User struct {
	UserID   string `json:"_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
