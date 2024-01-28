package main

import (
	"knowledge-service/internal/router"
	"knowledge-service/middleware"
	"knowledge-service/pkg/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := tools.ParseConfigure()
	if err != nil {
		panic(err)
	}
	app := gin.Default()
	app.Use(middleware.CORS())

	var mongo *tools.Mongo
	mongo.InitDB()
	var redis *tools.Redis
	redis.InitRedis()

	registerRoutes(app)

	addr := cfg.Host + ":" + cfg.Port
	if err := app.Run(addr); err != nil {
		panic(err)
	}
}

func registerRoutes(app *gin.Engine) {
	router.InitCommonRouter(app)
	router.InitUserRouter(app)
	router.InitRDocRouter(app)
	router.InitFeedRouter(app)
	router.InitSpaceRouter(app)
	router.InitWallpaperRouter(app)
}
