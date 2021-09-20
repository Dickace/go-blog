package main

import (
	"awesomeProject/controllers"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
)



func main() {

	r:= gin.Default();

	models.OpenConnection()

	r.GET("/api/v1/posts", controllers.GetPosts)
	r.POST("/api/v1/posts", controllers.AddPost)
	r.GET("/api/v1/posts/:id", controllers.GetPost)
	r.PATCH("/api/v1/posts/:id", controllers.EditPost)
	r.DELETE("/api/v1/posts/:id", controllers.DeletePost)

	r.Run();
}
