package main

import (
	_ "knowledge-service/docs"
	"knowledge-service/internal/router"
	"knowledge-service/middleware"
	"knowledge-service/pkg/tools"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := tools.ParseConfigure()
	if err != nil {
		panic(err)
	}
	app := gin.Default()
	app.Use(middleware.CORS())
	app.Use(gzip.Gzip(gzip.DefaultCompression))

	var mongo *tools.Mongo
	mongo.InitDB()
	var redis *tools.Redis
	redis.InitRedis()

	initSwagger(app)
	registerRoutes(app)

	addr := cfg.Host + ":" + cfg.Port
	if err := app.Run(addr); err != nil {
		panic(err)
	}
}

func initSwagger(app *gin.Engine) {
	if mode := os.Getenv("GIN_MODE"); mode != "release" {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func registerRoutes(app *gin.Engine) {
	router.InitCommonRouter(app)
	router.InitUserRouter(app)
	router.InitDocRouter(app)
	router.InitFeedRouter(app)
	router.InitSpaceRouter(app)
	router.InitWallpaperRouter(app)
}
