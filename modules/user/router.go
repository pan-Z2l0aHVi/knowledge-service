package user

import (
	"knowledge-base-service/middlewares"

	"github.com/gin-gonic/gin"
)

func (e *User) InitRouter(app *gin.Engine) {
	group := app.Group("user")

	group.GET("/profile", middlewares.VerifyToken(), e.GetProfile)
	group.POST("/profile", middlewares.VerifyToken(), e.UpdateProfile)

	group.POST("/sign_in", e.Login)

	group.GET("/yd_qrcode", e.GetYDQRCode)
	group.GET("/yd_login_status", e.GetYDLoginStatus)
	group.POST("/yd_callback", e.YDCallback)
}
