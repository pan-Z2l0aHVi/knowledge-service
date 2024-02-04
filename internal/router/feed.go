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

	group.GET("/comment_list", feedC.GetCommentList)
	group.POST("/comment", middleware.VerifyToken(), feedC.Comment)
	group.POST("/reply", middleware.VerifyToken(), feedC.Reply)
	group.POST("/comment_update", middleware.VerifyToken(), feedC.UpdateComment)
	group.POST("/comment_delete", middleware.VerifyToken(), feedC.DeleteComment)

	group.GET("/related_feeds", middleware.UseToken(), feedC.GetRelatedFeeds)
}
