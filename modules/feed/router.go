package feed

import (
	"github.com/gin-gonic/gin"
)

func (e *Feed) InitRouter(app *gin.Engine) {
	group := app.Group("feed")

	group.GET("/list", e.GetFeedList)
}
