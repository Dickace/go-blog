package main

import (
	"awesomeProject/controllers"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {

	r := gin.Default()

	r.Use(cors.Default())

	models.OpenConnection()

	r.GET("/api/v1/posts", controllers.GetPosts)
	r.POST("/api/v1/posts", controllers.AddPost)
	r.GET("/api/v1/posts/:id", controllers.GetPost)
	r.PATCH("/api/v1/posts/:id", controllers.EditPost)
	r.DELETE("/api/v1/posts/:id", controllers.DeletePost)
	r.PUT("/api/v1/posts/:id/favorite", controllers.Favorite)
	r.Run()
	//r.Run(":8001")
}
