package notification

import "github.com/stringintech/task-broker/types/event"

type Service interface {
	Start() error
	Close() error
	OnTaskCreated(event.TaskCreated) error //TODO? pass by reference?
}
