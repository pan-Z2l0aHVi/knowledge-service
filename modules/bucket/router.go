package bucket

import "github.com/gin-gonic/gin"

func (e *Bucket) InitRouter(app *gin.Engine) {
	group := app.Group("bucket")

	group.GET("/token", e.getToken)
}
