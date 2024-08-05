package main

import (
	"github.com/stringintech/task-broker/services"
	"github.com/stringintech/task-broker/services/notification"
	rabbitMqNotification "github.com/stringintech/task-broker/services/notification/rabbit_mq"
	"github.com/stringintech/task-broker/services/storage"
	postgresStorage "github.com/stringintech/task-broker/services/storage/postgres"
	types "github.com/stringintech/task-broker/types/base"
	"log"
)

func main() {
	dbConfig := postgresStorage.ServiceConfig{
		ConnectionUri: "postgres://postgres:postgres@localhost:6432/task-broker",
	}
	var storageService storage.Service
	var err error
	if storageService, err = postgresStorage.NewService(&dbConfig); err != nil {
		panic(err)
	}
	if err = storageService.Start(); err != nil {
		panic(err)
	}
	defer storageService.Close()

	c := rabbitMqNotification.ServiceConfig{
		ConnectionUrl: "amqp://guest:guest@localhost:5672/",
		QueueName:     "task-queue",
	}
	var notifService notification.Service
	if notifService, err = rabbitMqNotification.NewService(&c); err != nil {
		panic(err)
	}
	if err = notifService.Start(); err != nil {
		panic(err)
	}
	defer notifService.Close()

	taskService := services.TaskService{
		NotificationService: notifService,
		StorageService:      storageService,
	}

	if err = taskService.CreateTask(&types.Task{
		Id:    "dummy-identifier",
		Title: "dummy-title",
	}); err != nil {
		panic(err)
	}
	log.Println("task created")
}
