package stream

import (
	"context"
	"fmt"

	"github.com/TutorialEdge/ctxlog"
	"github.com/TutorialEdge/go-websocket-course/internal"
)

type EventConsumer interface {
	Consume(chan internal.Event)
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
	ctx := context.Background()
	eventChannel := make(chan internal.Event)
	go s.consumer.Consume(eventChannel)

	for {
		select {
		case e := <-eventChannel:
			// For example, show received message in a console.
			s.log.Info(context.TODO(), "event consumed")

			for _, c := range s.Pool.Channels[e.ChannelID] {
				err := c.Conn.WriteJSON(e)
				if err != nil {
					s.log.Error(
						ctx,
						fmt.Sprintf("failed to send event: %s", err.Error()),
					)
				}
			}
		}
	}
}
