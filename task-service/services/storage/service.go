package storage

import "github.com/stringintech/task-broker/types"

type Service interface {
	Start() error
	Close() error
	CreateTask(task *types.Task) error
}
