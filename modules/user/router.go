package user

import (
	"knowledge-base-service/middlewares"

	"github.com/gin-gonic/gin"
)

func (e *User) InitRouter(app *gin.Engine) {
	group := app.Group("user")

	group.GET("/profile", middlewares.TokenAuth(), e.GetProfile)
	group.POST("/profile", e.UpdateProfile)

	group.POST("/sign_in", e.SignIn)
}
