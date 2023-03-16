package main

import (
	"context"
	"time"

	"github.com/TutorialEdge/ctxlog"
	"github.com/TutorialEdge/go-websocket-course/internal/api"
	"github.com/TutorialEdge/go-websocket-course/internal/queue"
	"github.com/TutorialEdge/go-websocket-course/internal/stream"
)

func Run() error {
	// start up our logger
	ctx := context.Background()
	log := ctxlog.New(
		ctxlog.WithJSONFormat(),
	)
	log.Info(ctx, "starting stream service")

	// we need a way of connecting in to our RabbitMQ queue
	// the queue package implements both publishing and consuming
	eventQueue, err := queue.New(log)
	if err != nil {
		log.Error(ctx, err.Error())
		return err
	}

	// For now, we want to simulate new events being published to our queue
	// for a test-channel. This will be consumed by the service and sent to
	// the respective clients
	go func() {
		for {
			time.Sleep(5 * time.Second)
			eventQueue.Publish("{\"channel_id\": \"test-channel\"}")
		}
	}()

	// The stream service handles the websocket connections into
	// our system and sending new events to the appropriate clients
	streamService := stream.New(eventQueue, log)
	go streamService.Start(ctx)

	// We're exposing these websockets as endpoints on a server.
	// Here is where we instantiate the server which is based on gin and
	// handle route registration for components.
	myAPI := api.New(streamService)
	myAPI.RegisterRoutes(streamService)

	// Let's kick off the API and wait for incoming WebSocket connections!
	if err := myAPI.Serve(); err != nil {
		log.Info(ctx, "failed to start API")
		return err
	}
	return nil
}

// The entrypoint for our super-duper cool websocket system
func main() {
	if err := Run(); err != nil {
		// oh no
		panic(err.Error())
	}
}
