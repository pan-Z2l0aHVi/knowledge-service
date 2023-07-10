package user

type GetProfileQuery struct {
	UserID string `form:"user_id" binding:"required"`
}
