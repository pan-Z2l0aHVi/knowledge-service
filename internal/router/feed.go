package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitFeedRouter(app *gin.Engine) {
	group := app.Group("feed")
	feedC := controller.FeedController{}
	group.GET("/detail", feedC.GetDetail)
	group.GET("/list", feedC.SearchFeedList)
	group.POST("/praise", middleware.VerifyToken(), feedC.PraiseFeed)
}
