package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(app *gin.Engine) {
	group := app.Group("common")
	commonC := controller.CommonController{}
	group.GET("/qiniu_token", commonC.GetQiniuToken)
	group.GET("/r2_signed_url", commonC.GetR2SignedURL)

	group.POST("/report", middleware.UseToken(), commonC.Report)
	group.GET("/statics", commonC.GetStatics)
}
