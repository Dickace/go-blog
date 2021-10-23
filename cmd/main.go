package main

import (
	"awesomeProject/controllers"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
	}))

	models.OpenConnection()

	r.GET("/api/v1/posts", controllers.GetPosts)
	r.POST("/api/v1/posts", controllers.AddPost)
	r.GET("/api/v1/posts/:id", controllers.GetPost)
	r.PATCH("/api/v1/posts/:id", controllers.EditPost)
	r.DELETE("/api/v1/posts/:id", controllers.DeletePost)
	r.PUT("/api/v1/posts/:id/favorite", controllers.Favorite)
	r.POST("/api/v1/telemetry", controllers.AddTelemtry)
	r.GET("/api/v1/getTelemetry", controllers.GetTelemetries)
	r.Run()
	//r.Run(":8001")
}
