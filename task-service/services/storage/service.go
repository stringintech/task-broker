package storage

import types "github.com/stringintech/task-broker/types/base"

type Service interface {
	Start() error
	Close() error
	CreateTask(*types.Task) error
}
