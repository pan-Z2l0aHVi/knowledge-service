package router

import (
	"knowledge-service/internal/controller"
	"knowledge-service/middleware"

	"github.com/gin-gonic/gin"
)

func InitWallpaperRouter(app *gin.Engine) {
	group := app.Group("wallpaper")
	wallpaperC := controller.WallpaperController{}
	group.GET("/search", middleware.UseToken(), wallpaperC.Search)
	group.GET("/info", middleware.UseToken(), wallpaperC.GetInfo)
}
