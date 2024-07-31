package main

import (
	"github.com/stringintech/task-broker/services/notification"
	"github.com/stringintech/task-broker/types"
	"log"
)

func main() {
	c := notification.ExternalServiceConfig{
		ConnectionUrl: "amqp://guest:guest@localhost:5672/",
		QueueName:     "task-queue",
	}
	var notifService notification.Service
	var err error
	if notifService, err = notification.NewExternalService(&c); err != nil {
		panic(err)
	}
	if err = notifService.Start(); err != nil {
		panic(err)
	}
	defer notifService.Close()

	if err = notifService.OnTaskCreated(types.Task{
		Id:    "dummy-identifier",
		Title: "dummy-title",
	}); err != nil {
		panic(err)
	}
	log.Println("sent on task created notification")
}
