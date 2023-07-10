package material

import "github.com/gin-gonic/gin"

func (e *Material) InitRouter(app *gin.Engine) {
	group := app.Group("material")

	group.GET("/search", e.search)
	group.GET("/info", e.getInfo)
	group.POST("/upload", e.upload)
}
