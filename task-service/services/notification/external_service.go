package notification

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stringintech/task-broker/types"
	"google.golang.org/protobuf/proto"
)

type ExternalServiceConfig struct {
	ConnectionUrl string
	QueueName     string
}

type ExternalService struct {
	config     *ExternalServiceConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func NewExternalService(config *ExternalServiceConfig) (*ExternalService, error) {
	return &ExternalService{config: config}, nil
}

func (s *ExternalService) Start() error {
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

func (s *ExternalService) Close() error {
	_ = s.channel.Close() //TODO? handle/log errors
	_ = s.connection.Close()
	return nil
}

func (s *ExternalService) OnTaskCreated(task types.Task) error {
	if err := s.ensureStarted(); err != nil {
		return err
	}

	body, err := proto.Marshal(&task)
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

func (s *ExternalService) ensureStarted() error {
	if s.queue == nil {
		return fmt.Errorf("service not started yet")
	}
	return nil
}
