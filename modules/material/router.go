package material

import "github.com/gin-gonic/gin"

func (e *Material) InitRouter(app *gin.Engine) {
	group := app.Group("material")

	group.GET("/search", e.Search)
	group.GET("/info", e.GetInfo)
	group.POST("/upload", e.Upload)
}
