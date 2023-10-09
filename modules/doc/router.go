package doc

import (
	"knowledge-base-service/middlewares"

	"github.com/gin-gonic/gin"
)

func (e *Doc) InitRouter(app *gin.Engine) {
	group := app.Group("doc")

	group.GET("/info", e.GetInfo)
	group.GET("/docs", middlewares.VerifyToken(), e.GetDocs)
	group.POST("/create", middlewares.VerifyToken(), e.Create)
	group.POST("/update", middlewares.VerifyToken(), e.Update)
	group.POST("/delete", middlewares.VerifyToken(), e.Delete)
	group.GET("/drafts", middlewares.VerifyToken(), e.GetDrafts)
	group.PUT("/update_draft", middlewares.VerifyToken(), e.UpdateDraft)
}
