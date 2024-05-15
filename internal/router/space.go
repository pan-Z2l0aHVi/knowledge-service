package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitSpaceRouter(app *gin.Engine) {
	group := app.Group("space")
	spaceC := controller.SpaceController{}
	group.GET("/info", spaceC.GetInfo)
	group.GET("/search", middleware.VerifyToken(), spaceC.SearchSpaces)
	group.POST("/create", middleware.VerifyToken(), spaceC.Create)
	group.POST("/update", middleware.VerifyToken(), spaceC.Update)
	group.POST("/delete", middleware.VerifyToken(), spaceC.Delete)
}
