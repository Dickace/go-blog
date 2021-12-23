package main

import (
	"awesomeProject/models"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("rabbitmq://admin:password@localhost:5672/")

	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	queue, err := amqpChannel.QueueDeclare("addPost", true, false, false, false, nil)
	handleError(err, "Could not declare `addPost` queue")
	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	addPost := models.Post{Title: "New Post from rabbit", PostDate: time.Now()}

	body, err := json.Marshal(addPost)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddPost send")
}
