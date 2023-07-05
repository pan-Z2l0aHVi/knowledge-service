package material

import "github.com/gin-gonic/gin"

func (e *Material) InitRouter(app *gin.Engine) {
	group := app.Group("material")

	group.GET("/info", e.getInfo)
	group.POST("/create", e.create)
	group.POST("/update", e.update)
	group.PUT("upload", e.upload)
}
