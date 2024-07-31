package notification

import "github.com/stringintech/task-broker/types"

type Service interface {
	Start() error
	Close() error
	OnTaskCreated(task types.Task) error //TODO? pass by reference?
}
