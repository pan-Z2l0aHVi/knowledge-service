package user

import "github.com/gin-gonic/gin"

func (e *User) InitRouter(app *gin.Engine) {
	group := app.Group("user")

	group.GET("/profile", e.getProfile)
	group.POST("/profile", e.updateProfile)

	group.POST("/sign_up", e.signUp)
	group.POST("/sign_in", e.signIn)
}
