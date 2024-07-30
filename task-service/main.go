package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stringintech/task-broker/types"
	"google.golang.org/protobuf/proto"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare a queue")

	msg := &types.TaskMessage{
		Content: "Hello World",
	}
	body, err := proto.Marshal(msg)
	failOnError(err, "failed to encode message")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "failed to publish a message")
	log.Printf("sent message: %s", msg.Content)
}
