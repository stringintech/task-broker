package rabbit_mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stringintech/task-broker/types/event"
	"google.golang.org/protobuf/proto"
)

type ServiceConfig struct {
	ConnectionUrl string
	QueueName     string
}

type Service struct {
	config     *ServiceConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func NewService(config *ServiceConfig) (*Service, error) {
	return &Service{config: config}, nil
}

func (s *Service) Start() error {
	var err error
	s.connection, err = amqp.Dial(s.config.ConnectionUrl)
	if err != nil {
		return err
	}

	s.channel, err = s.connection.Channel()
	if err != nil {
		return err
	}

	q, err := s.channel.QueueDeclare(
		s.config.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = s.channel.Close()
		_ = s.connection.Close()
		return err
	}
	s.queue = &q

	return nil
}

func (s *Service) Close() error {
	if err := s.ensureStarted(); err != nil {
		return err
	}
	_ = s.channel.Close() //TODO? handle/log errors
	_ = s.connection.Close()
	return nil
}

func (s *Service) OnTaskCreated(e event.TaskCreated) error {
	if err := s.ensureStarted(); err != nil {
		return err
	}

	body, err := proto.Marshal(&e)
	if err != nil {
		return err
	}

	return s.channel.Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (s *Service) ensureStarted() error {
	if s.queue == nil {
		return fmt.Errorf("service not started yet")
	}
	return nil
}
