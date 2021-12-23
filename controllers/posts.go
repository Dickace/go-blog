package controllers

import (
	"awesomeProject/models"
	"awesomeProject/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type PostId struct {
	Id string `json:"ID"`
}
type EditPostData struct {
	Id    string            `json:"ID"`
	Input models.UpdatePost `json:"UpdatePost"`
}
type CreatePostData struct {
	Id    string            `json:"ID"`
	Input models.CreatePost `json:"UpdatePost"`
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	models.DB.Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}
func GetPost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}
func EditPost(c *gin.Context) {

	var input models.UpdatePost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	editPost := EditPostData{Id: c.Param("id"), Input: input}

	body, err := json.Marshal(editPost)
	if err != nil {
		rabbitmq.HandleError(err, "Error encoding JSON")
	}

	err = rabbitmq.AmqpChannel.Publish("", "updatePost", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": editPost})

}
func DeletePost(c *gin.Context) {

	post := PostId{Id: c.Param("id")}

	body, err := json.Marshal(post)
	if err != nil {
		rabbitmq.HandleError(err, "Error encoding JSON")
	}

	err = rabbitmq.AmqpChannel.Publish("", "deletePost", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
func AddPost(c *gin.Context) {
	var input models.CreatePost

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, err := json.Marshal(input)
	if err != nil {
		rabbitmq.HandleError(err, "Error encoding JSON")
	}
	fmt.Print(body)
	err = rabbitmq.AmqpChannel.Publish("", "addPost", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("addPost send")
	c.JSON(http.StatusOK, gin.H{"data": input})
}
func Favorite(c *gin.Context) {

	postFavorite := PostId{Id: c.Param("id")}
	body, err := json.Marshal(postFavorite)
	if err != nil {
		rabbitmq.HandleError(err, "Error encoding JSON")
	}

	err = rabbitmq.AmqpChannel.Publish("", "favoritePost", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{"data": postFavorite})
}
