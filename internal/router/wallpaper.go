package router

import (
	"knowledge-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func InitWallpaperRouter(app *gin.Engine) {
	group := app.Group("wallpaper")
	wallpaperC := controller.WallpaperController{}
	group.GET("/search", wallpaperC.Search)
	group.GET("/info", wallpaperC.GetInfo)
}
