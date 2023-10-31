package common

import (
	"github.com/gin-gonic/gin"
)

func (e *Common) InitRouter(app *gin.Engine) {
	group := app.Group("common")

	group.GET("/qiniu_token", e.GetQiniuToken)
	group.GET("/r2_signed_url", e.GetR2SignedURL)
	group.POST("/report", e.Report)
}
