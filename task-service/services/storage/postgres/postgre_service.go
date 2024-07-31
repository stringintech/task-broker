package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringintech/task-broker/types"
	"strconv"
)

type ServiceConfig struct {
	ConnectionUri string
}

type Service struct {
	config *ServiceConfig
	pool   *pgxpool.Pool
}

func NewService(config *ServiceConfig) (*Service, error) {
	return &Service{config: config}, nil
}

func (s *Service) Start() error {
	var err error
	s.pool, err = pgxpool.New(context.Background(), s.config.ConnectionUri)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Close() error {
	if err := s.ensureStarted(); err != nil {
		return err
	}
	s.pool.Close()
	return nil
}

func (s *Service) CreateTask(task *types.Task) error {
	if err := s.ensureStarted(); err != nil {
		return err
	}
	var id uint64
	if err := s.pool.QueryRow(context.Background(), `insert into task (title) values ($1) returning id`, task.Title).
		Scan(&id); err != nil {
		return err
	}
	task.Id = strconv.FormatUint(id, 10)
	return nil
}

func (s *Service) ensureStarted() error {
	if s.pool == nil {
		return fmt.Errorf("service not started yet")
	}
	return nil
}
