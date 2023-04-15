package user

type GetProfileParams struct {
	UserID string `form:"user_id" binding:"required"`
}
