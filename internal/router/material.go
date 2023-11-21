package router

import (
	"knowledge-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitMaterialRouter(app *gin.Engine) {
	group := app.Group("material")
	materialC := controller.MaterialController{}
	group.GET("/search", materialC.Search)
	group.GET("/info", materialC.GetInfo)
	group.POST("/upload", materialC.Upload)
}
