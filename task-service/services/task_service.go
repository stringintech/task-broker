package services

import (
	"github.com/stringintech/task-broker/services/notification"
	"github.com/stringintech/task-broker/services/storage"
	"github.com/stringintech/task-broker/types"
)

type TaskService struct {
	NotificationService notification.Service
	StorageService      storage.Service
}

func (s *TaskService) CreateTask(task *types.Task) error { //TODO return typed error in order to handle task error and notification error separately
	err := s.StorageService.CreateTask(task)
	if err != nil {
		return err
	}
	return s.NotificationService.OnTaskCreated(*task)
}
