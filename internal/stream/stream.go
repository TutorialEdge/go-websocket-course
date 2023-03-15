package stream

import (
	"context"

	"golang.org/x/net/websocket"
)

type EventConsumer interface {
	Listen()
}

type Service struct {
	consumer EventConsumer
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func New(consumer EventConsumer) *Service {
	return &Service{
		consumer: consumer,
	}
}

func (s *Service) Stream(ctx context.Context) error {
	for {

	}
	return nil
}
