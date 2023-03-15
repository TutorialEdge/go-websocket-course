package event

import (
	"context"
	"encoding/json"
	"fmt"

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

func (s *Service) ListenForEvents() error {
	return nil
}

// PublishEvent - when we
func (s *Service) PublishEvent(ctx context.Context, e internal.Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// attempt to publish a message to the queue!
	err = s.Channel.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Successfully Published Message to Queue")
	return nil
}
