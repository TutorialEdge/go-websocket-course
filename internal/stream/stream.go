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

// Service - is comprised of an event consumer
// a pool which maps all the connections to channels
// and a logger
type Service struct {
	consumer EventConsumer
	Pool     *Pool
	log      *ctxlog.CtxLogger
}

// New - returns a new streaming service
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

// Start - starts the streaming service
// This will continually be fed events from the consumer
// which it will then route to the appropriate channel and clients
// connected to that channel.
func (s *Service) Start(ctx context.Context) {
	eventChannel := make(chan internal.Event)
	go s.consumer.Consume(eventChannel)

	forever := make(chan bool)
	// For every event that is consumed and sent to our
	// shared eventChannel, we want to send it to all the clients
	// connected to a given channel.
	for e := range eventChannel {
		// For example, show received message in a console.
		s.log.Info(ctx, "event consumed")

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
	<-forever
}
