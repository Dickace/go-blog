package main

import (
	"awesomeProject/consumer"
	"awesomeProject/controllers"
	"awesomeProject/models"
	"awesomeProject/rabbitmq"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {

	rabbitmq.RabbitmqInit()

	consumer.AddPost()
	consumer.AddPost()
	consumer.AddPost()
	consumer.DeletePost()
	consumer.FavoritePost()
	consumer.UpdatePost()

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
	r.Run(":8081")
	//r.Run(":82")
}
