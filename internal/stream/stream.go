package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TutorialEdge/ctxlog"
	"github.com/TutorialEdge/go-websocket-course/internal"
	"github.com/gorilla/websocket"
)

type EventConsumer interface {
	Consume() internal.Event
}

type Service struct {
	consumer EventConsumer
	Pool     *Pool
	log      *ctxlog.CtxLogger
}

func New(consumer EventConsumer, log *ctxlog.CtxLogger) *Service {
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
	for {
		time.Sleep(1 * time.Second)
		ctx := context.Background()
		event := s.consumer.Consume()
		event.ChannelID = "test-channel"

		data, err := json.Marshal(event)
		if err != nil {
			s.log.Error(ctx, "failed to marshal event")
		}

		for _, c := range s.Pool.Channels[event.ChannelID] {
			err := c.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				s.log.Error(
					ctx,
					fmt.Sprintf("failed to send event: %s", err.Error()),
				)
				// s.Pool.Channels[event.ChannelID][c.ID]
			}
		}
	}
}
