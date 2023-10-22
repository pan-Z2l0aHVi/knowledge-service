package feed

import (
	"knowledge-base-service/middlewares"

	"github.com/gin-gonic/gin"
)

func (e *Feed) InitRouter(app *gin.Engine) {
	group := app.Group("feed")

	group.GET("/detail", e.GetDetail)
	group.GET("/list", e.SearchFeedList)
	group.POST("/praise", middlewares.VerifyToken(), e.PraiseFeed)
}
