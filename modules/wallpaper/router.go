package wallpaper

import "github.com/gin-gonic/gin"

func (e *Wallpaper) InitRouter(app *gin.Engine) {
	group := app.Group("wallpaper")

	group.GET("/search", e.Search)
	group.GET("/info", e.GetInfo)
}
