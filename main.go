package main

import (
	"knowledge-base-service/middlewares"
	"knowledge-base-service/modules/doc"
	"knowledge-base-service/modules/material"
	"knowledge-base-service/modules/user"
	"knowledge-base-service/modules/wallpaper"
	"knowledge-base-service/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := tools.ParseConfigure("./conf/app.json")
	if err != nil {
		panic(err)
	}
	app := gin.Default()
	app.Use(middlewares.CORS())

	var mongo *tools.Mongo
	mongo.InitDB()
	registerRouter(app)

	addr := cfg.Host + ":" + cfg.Port
	if err := app.Run(addr); err != nil {
		panic(err)
	}
}

func registerRouter(app *gin.Engine) {
	new(user.User).InitRouter(app)
	new(doc.Doc).InitRouter(app)
	new(wallpaper.Wallpaper).InitRouter((app))
	new(material.Material).InitRouter((app))
}
