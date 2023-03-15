package stream

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TutorialEdge/ctxlog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// A client represents a single websocket connection
// in our system
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// A pool represents all websocket clients on our service
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Channels   map[string][]*Client
	Clients    []*Client
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Service) SetupRoutes(c *gin.Engine) {
	c.GET("/api/v1/stream", s.stream)
}

func (s *Service) stream(c *gin.Context) {
	ctx := c.Request.Context()
	channelID := c.Query("channel")
	ctx = ctxlog.WithFields(ctx, ctxlog.Fields{
		"channel-id": channelID,
	})
	s.log.Info(ctx, "new websocket connection")
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Info(
			ctx,
			fmt.Sprintf("Failed to set websocket upgrade: %+v", err.Error()),
		)
		return
	}
	client := &Client{
		Conn: conn,
	}
	defer conn.Close()
	s.log.Info(ctx, "registering client for channel")
	if _, ok := s.Pool.Channels[channelID]; ok {
		s.Pool.Channels[channelID] = append(s.Pool.Clients, client)
	} else {
		s.Pool.Channels[channelID] = []*Client{client}
	}
	s.keepAlive(ctx, client)
}

func (s *Service) keepAlive(ctx context.Context, client *Client) {
	defer func() {
		s.log.Info(ctx, "connection closed")
		s.Pool.Unregister <- client
	}()
	s.log.Info(ctx, "keepalive started")

	for {
		time.Sleep(15 * time.Second)
		s.log.Info(ctx, "sending keepalive")
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte("I'm alive"))
		if err != nil {
			s.log.Error(ctx, "failed to send message on connection")
			return
		}
	}
}
