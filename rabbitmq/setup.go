package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var AmqpChannel *amqp.Channel

func RabbitmqInit() {
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")

	HandleError(err, "Can't connect to AMQP")
	amqpChannel, err := conn.Channel()
	HandleError(err, "Can't create a amqpChannel")
	AmqpChannel = amqpChannel
	err = AmqpChannel.Qos(1, 0, false)
	HandleError(err, "Could not configure QoS")

	_, err = AmqpChannel.QueueDeclare("favoritePost", true, false, false, false, nil)
	HandleError(err, "Could not declare `favoritePost` queue")
	_, err = AmqpChannel.QueueDeclare("addPost", true, false, false, false, nil)
	HandleError(err, "Could not declare `addPost` queue")
	_, err = AmqpChannel.QueueDeclare("deletePost", true, false, false, false, nil)
	HandleError(err, "Could not declare `deletePost` queue")
	_, err = AmqpChannel.QueueDeclare("updatePost", true, false, false, false, nil)
	HandleError(err, "Could not declare `updatePost` queue")
}
