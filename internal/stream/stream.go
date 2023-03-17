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

// Streamer - is comprised of an event consumer
// a pool which maps all the connections to channels
// and a logger
type Streamer struct {
	// Consumer - used for consuming events from
	// a queue and then sending them off to our connected clients
	consumer EventConsumer

	// Channels used for registering
	// and unregistering clients to our streamer
	Register   chan *Client
	Unregister chan *Client

	// Channels - keeps a map of channel-ids
	// and all the clients subscribed to that channel
	Channels map[string][]string

	// Clients - keeps a map of all clients
	// and their IDs current subscribed to the streamer
	Clients map[string]*Client
	log     *ctxlog.CtxLogger
}

// New - returns a new streaming service
func New(
	consumer EventConsumer,
	log *ctxlog.CtxLogger,
) *Streamer {
	streamer := &Streamer{
		consumer:   consumer,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Channels:   make(map[string][]string),
		Clients:    make(map[string]*Client),
		log:        log,
	}
	return streamer
}

// Start - starts the streaming service
// This will continually be fed events from the consumer
// which it will then route to the appropriate channel and clients
// connected to that channel.
func (s *Streamer) Start(ctx context.Context) {
	eventChannel := make(chan internal.Event)
	go s.consumer.Consume(eventChannel)

	for {
		select {
		case e := <-eventChannel:
			// For example, show received message in a console.
			s.log.Info(ctx, "event consumed")

			for _, clientID := range s.Channels[e.ChannelID] {
				err := s.Clients[clientID].Conn.WriteJSON(e)
				if err != nil {
					s.log.Error(
						ctx,
						fmt.Sprintf("failed to send event: %s", err.Error()),
					)
				}
			}
		case client := <-s.Register:
			s.log.Info(ctx, "new client registered")
			s.Clients[client.ID] = client
			s.Channels[client.ChannelID] = append(s.Channels[client.ChannelID], client.ID)
			s.log.Info(ctx, fmt.Sprintf("Client Length: %d", len(s.Clients)))
		case client := <-s.Unregister:
			s.log.Info(ctx, "client unregistered")
			s.log.Info(ctx, fmt.Sprintf("%+v\n", client))
			delete(s.Clients, client.ID)
			s.log.Info(ctx, fmt.Sprintf("Client Length: %d", len(s.Clients)))

			for index, clientID := range s.Channels[client.ChannelID] {
				if clientID == client.ID {
					s.Channels[client.ChannelID] = append(s.Channels[client.ChannelID][:index], s.Channels[client.ChannelID][index+1:]...)
				}
			}
			s.log.Info(ctx, fmt.Sprintf("Channel Length: %d", len(s.Channels[client.ChannelID])))
		}
	}
}
