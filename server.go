package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hadyanadam/golang-tutorial-gin/controller"
	"github.com/hadyanadam/golang-tutorial-gin/middlewares"
	"github.com/hadyanadam/golang-tutorial-gin/service"
)


var(
	videoService service.VideoService = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}


func main() {
	setupLogOutput()

	server := gin.Default()

	server.Static("/css", "./templates/css")

	server.LoadHTMLGlob("templates/*.html")

	apiRoutes := server.Group("/api")
	apiRoutes.Use(middlewares.BasicAuth())
	{
		apiRoutes.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST("/videos", func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video input is valid"})
			}
		})
	}

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}

	port := os.Getenv("PORT")

	server.Run(":" + port)
}