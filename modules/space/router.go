package space

import (
	"knowledge-base-service/middlewares"

	"github.com/gin-gonic/gin"
)

func (e *Space) InitRouter(app *gin.Engine) {
	group := app.Group("space")

	group.GET("/info", e.GetInfo)
	group.GET("/search", middlewares.UseToken(), e.SearchSpaces)
	group.POST("/create", middlewares.VerifyToken(), e.Create)
	group.POST("/update", middlewares.VerifyToken(), e.Update)
	group.POST("/delete", middlewares.VerifyToken(), e.Delete)
}
