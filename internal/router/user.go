package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(app *gin.Engine) {
	group := app.Group("user")
	userC := controller.UserController{}
	group.GET("/profile", middleware.VerifyToken(), userC.GetProfile)
	group.POST("/profile", middleware.VerifyToken(), userC.UpdateProfile)

	group.POST("/sign_in", userC.Login)

	group.GET("/yd_qrcode", userC.GetYDQRCode)
	group.GET("/yd_login_status", userC.GetYDLoginStatus)
	group.POST("/yd_callback", userC.YDCallback)
}
