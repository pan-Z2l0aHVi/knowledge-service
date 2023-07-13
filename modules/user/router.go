package user

import "github.com/gin-gonic/gin"

func (e *User) InitRouter(app *gin.Engine) {
	group := app.Group("user")

	group.GET("/profile", e.GetProfile)
	group.POST("/profile", e.UpdateProfile)

	group.POST("/sign_up", e.SignUp)
	group.POST("/sign_in", e.SignIn)
}
