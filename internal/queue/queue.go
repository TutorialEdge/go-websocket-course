package queue

import (
	"github.com/TutorialEdge/go-websocket-course/internal"
	"github.com/streadway/amqp"
)

type Service struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func New() *Service {
	return &Service{}
}

func (s *Service) Consume() internal.Event {

	// Make a channel to receive messages into infinite loop.
	// forever := make(chan bool)

	// <-forever
	return internal.Event{}
}
