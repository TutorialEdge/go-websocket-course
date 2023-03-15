package stream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TutorialEdge/ctxlog"
	"github.com/TutorialEdge/go-websocket-course/internal"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type EventConsumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

type Service struct {
	consumer EventConsumer
	Pool     *Pool
	log      *ctxlog.CtxLogger
}

func New(
	consumer EventConsumer,
	log *ctxlog.CtxLogger,
) *Service {
	return &Service{
		consumer: consumer,
		Pool: &Pool{
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
			Channels:   make(map[string][]*Client),
		},
		log: log,
	}
}

func (s *Service) Start() {
	s.log.Info(context.Background(), "start")
	ctx := context.Background()
	messages, err := s.consumer.Consume()
	if err != nil {
		s.log.Error(ctx, err.Error())
		return
	}

	forever := make(chan bool)
	go func() {
		for e := range messages {
			// For example, show received message in a console.
			s.log.Info(context.TODO(), "event consumed")
			var event internal.Event
			err := json.Unmarshal(e.Body, &event)
			if err != nil {
				return
			}

			for _, c := range s.Pool.Channels[event.ChannelID] {
				err := c.Conn.WriteMessage(websocket.TextMessage, e.Body)
				if err != nil {
					s.log.Error(
						ctx,
						fmt.Sprintf("failed to send event: %s", err.Error()),
					)
					// s.Pool.Channels[event.ChannelID][c.ID]
				}
			}
		}
	}()
	<-forever

}
