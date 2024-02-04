package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitRDocRouter(app *gin.Engine) {
	group := app.Group("doc")
	docC := controller.DocController{}
	group.GET("/info", docC.GetInfo)
	group.GET("/docs", middleware.UseToken(), docC.SearchDocs)
	group.POST("/create", middleware.VerifyToken(), docC.Create)
	group.POST("/update", middleware.VerifyToken(), docC.Update)
	group.POST("/delete", middleware.VerifyToken(), docC.Delete)

	group.GET("/drafts", middleware.VerifyToken(), docC.GetDrafts)
	group.PUT("/update_drafts", middleware.VerifyToken(), docC.UpdateDrafts)
}
