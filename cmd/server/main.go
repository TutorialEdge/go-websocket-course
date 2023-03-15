package main

import (
	"context"

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

	eventQueue := queue.New()
	streamService := stream.New(eventQueue, log)
	go streamService.Start()

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
