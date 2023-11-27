package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(app *gin.Engine) {
	group := app.Group("user")
	userC := controller.UserController{}
	group.GET("/profile", middleware.VerifyToken(), userC.GetProfile)
	group.POST("/profile", middleware.VerifyToken(), userC.UpdateProfile)

	group.GET("/user_info", middleware.UseToken(), userC.GetUserInfo)

	group.POST("/sign_in", userC.Login)

	group.GET("/yd_qrcode", userC.GetYDQRCode)
	group.GET("/yd_login_status", userC.GetYDLoginStatus)
	group.POST("/yd_callback", userC.YDCallback)

	group.GET("/collected_feeds", middleware.VerifyToken(), userC.GetCollectedFeeds)
	group.POST("/collect_feed", middleware.VerifyToken(), userC.CollectFeed)
	group.POST("/cancel_collect_feed", middleware.VerifyToken(), userC.CancelCollectFeed)

	group.GET("/followed_users", middleware.VerifyToken(), userC.GetFollowedUsers)
	group.POST("/follow_user", middleware.VerifyToken(), userC.FollowUser)
	group.POST("/unfollow_user", middleware.VerifyToken(), userC.UnfollowUser)

	group.GET("/collected_wallpapers", middleware.VerifyToken(), userC.GetCollectedWallpapers)
	group.POST("/collect_wallpaper", middleware.VerifyToken(), userC.CollectWallpaper)
	group.POST("/cancel_collect_wallpaper", middleware.VerifyToken(), userC.CancelCollectWallpaper)
}
