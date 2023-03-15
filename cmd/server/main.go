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
	ctx := context.Background()
	log := ctxlog.New(
		ctxlog.WithJSONFormat(),
	)
	log.Info(ctx, "starting stream service")

	eventQueue, err := queue.New(log)
	if err != nil {
		log.Error(ctx, err.Error())
		return err
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			eventQueue.Publish("{\"channel_id\": \"test-channel\"}")
		}
	}()

	streamService := stream.New(eventQueue, log)
	go streamService.Start(ctx)

	myAPI := api.New(streamService)
	myAPI.RegisterRoutes(streamService)

	if err := myAPI.Serve(); err != nil {
		log.Info(ctx, "failed to start API")
		return err
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		// oh no
		panic(err.Error())
	}
}
