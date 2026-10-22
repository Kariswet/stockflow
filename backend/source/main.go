package main

import (
	"log"

	"stockflow/controller"
	"stockflow/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           stockflow
// @version         1.0
// @description     stockflow documentation
func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(cors.Default())

	basePath := "/api/"
	docs.SwaggerInfo.BasePath = basePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/rapidoc.html", "docs/rapidoc.html")

	api := router.Group(basePath)
	api.GET("/ping", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "hello"}) })

	// service.NewStockService()
	controller.NewProductController(api)
	controller.NewOrderController(api)

	router.Run(":8000")
}
