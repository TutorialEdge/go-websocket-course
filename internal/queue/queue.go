package queue

import (
	"fmt"

	"github.com/TutorialEdge/ctxlog"
	"github.com/streadway/amqp"
)

type Service struct {
	QueueName string
	log       *ctxlog.CtxLogger
	Conn      *amqp.Connection
	Channel   *amqp.Channel
}

func New(log *ctxlog.CtxLogger) (*Service, error) {
	service := &Service{
		QueueName: "events",
		log:       log,
	}
	err := service.Connect()
	if err != nil {
		return nil, err
	}
	return service, nil
}

// Connect - establishes a connection to our RabbitMQ instance
// and declares the queue we are going to be using
func (s *Service) Connect() error {
	var err error
	s.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	fmt.Println("Successfully Connected to RabbitMQ")

	// We need to open a channel over our AMQP connection
	// This will allow us to declare queues and subsequently consume/publish
	// messages
	s.Channel, err = s.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Here we declare our new queue that we want to publish to and consume
	// from:
	_, err = s.Channel.QueueDeclare(
		s.QueueName, // Queue Name
		false,       // durable
		false,       // Delete when not used
		false,       // exclusive
		false,       // no wait
		nil,         // additional args
	)
	return nil
}

func (s *Service) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := s.Channel.Consume(
		s.QueueName,
		"Event-Consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// Publish - publishes a message to the queue
func (s *Service) Publish(message string) error {
	// attempt to publish a message to the queue!
	err := s.Channel.Publish(
		"",
		s.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully Published Message to Queue")
	return nil
}
