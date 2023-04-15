package doc

import "github.com/gin-gonic/gin"

func (e *Doc) InitRouter(app *gin.Engine) {
	group := app.Group("doc")

	group.GET("/info", e.getInfo)
	group.POST("/create", e.create)
	group.POST("/update", e.update)
}
