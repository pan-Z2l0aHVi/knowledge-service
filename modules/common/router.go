package common

import (
	"github.com/gin-gonic/gin"
)

func (e *Common) InitRouter(app *gin.Engine) {
	group := app.Group("common")

	group.GET("/qiniu_token", e.GetQiniuToken)
	group.POST("/report", e.Report)
}
