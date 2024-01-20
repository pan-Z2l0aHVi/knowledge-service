package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitFeedRouter(app *gin.Engine) {
	group := app.Group("feed")
	feedC := controller.FeedController{}
	group.GET("/info", middleware.UseToken(), feedC.GetInfo)
	group.GET("/list", middleware.UseToken(), feedC.SearchFeedList)
	group.POST("/like", middleware.VerifyToken(), feedC.LikeFeed)
}
