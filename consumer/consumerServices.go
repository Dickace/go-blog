package consumer

import (
	"awesomeProject/controllers"
	"awesomeProject/models"
	"awesomeProject/rabbitmq"
	"encoding/json"
	"log"
	"os"
)

func AddPost() {
	messageChanel, err := rabbitmq.AmqpChannel.Consume(
		"addPost",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.HandleError(err, "Could not register consumer")
	go func() {
		for d := range messageChanel {
			addPost := &models.CreatePost{}
			err := json.Unmarshal(d.Body, addPost)

			post := models.Post{
				Title:    addPost.Title,
				PostDate: addPost.PostDate,
			}
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			log.Print(post)

			models.DB.Create(&post)
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}

	}()

}

func FavoritePost() {

	messageChanel, err := rabbitmq.AmqpChannel.Consume(
		"favoritePost",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.HandleError(err, "Could not register consumer")

	go func() {
		for d := range messageChanel {
			post := &models.Post{}
			favoritePost := &controllers.PostId{}
			err := json.Unmarshal(d.Body, favoritePost)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			log.Print(favoritePost)

			if err := models.DB.Where("id = ?", favoritePost.Id).First(&post).Error; err != nil {
				log.Printf("Error acknowledging message : %s", err)
			}
			if !post.IsFavorite {
				models.DB.Model(&post).Updates(map[string]interface{}{
					"IsFavorite": true,
				})
			} else {
				models.DB.Model(&post).Updates(map[string]interface{}{
					"IsFavorite": false,
				})
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}

	}()

}

func DeletePost() {

	messageChanel, err := rabbitmq.AmqpChannel.Consume(
		"deletePost",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.HandleError(err, "Could not register consumer")

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChanel {
			deletePost := &controllers.PostId{}
			err := json.Unmarshal(d.Body, deletePost)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			log.Print(deletePost)

			var post models.Post
			if err := models.DB.Where("id = ?", deletePost.Id).First(&post).Error; err != nil {
				log.Printf("Error acknowledging message : %s", err)
			}
			models.DB.Delete(&post)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}

	}()
}

func UpdatePost() {

	messageChanel, err := rabbitmq.AmqpChannel.Consume(
		"updatePost",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	rabbitmq.HandleError(err, "Could not register consumer")

	go func() {
		for d := range messageChanel {
			updatePost := &controllers.EditPostData{}
			err := json.Unmarshal(d.Body, updatePost)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			log.Print(updatePost)

			var post models.Post
			if err := models.DB.Where("id = ?", updatePost.Id).First(&post).Error; err != nil {
				log.Printf("Error acknowledging message : %s", err)
				return
			}
			models.DB.Model(&post).Updates(updatePost.Input)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}

	}()

}
