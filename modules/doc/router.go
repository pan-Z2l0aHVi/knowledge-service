package doc

import "github.com/gin-gonic/gin"

func (e *Doc) InitRouter(app *gin.Engine) {
	group := app.Group("doc")

	group.GET("/info", e.GetInfo)
	group.POST("/create", e.Create)
	group.POST("/update", e.Update)
}
