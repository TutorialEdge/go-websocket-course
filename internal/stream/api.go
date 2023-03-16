package stream

import (
	"context"
	"net/http"
	"time"

	"github.com/TutorialEdge/ctxlog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// A client represents a single websocket connection
// in our system
type Client struct {
	ID        string
	ChannelID string
	Conn      *websocket.Conn
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SetupRoutes - maps our websocket endpoint
func (s *Streamer) SetupRoutes(c *gin.Engine) {
	c.GET("/api/v1/stream", s.stream)
}

// stream -
func (s *Streamer) stream(c *gin.Context) {
	ctx := c.Request.Context()
	channelID := c.Query("channel")
	ctx = ctxlog.WithFields(ctx, ctxlog.Fields{
		"channel_id": channelID,
	})
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.log.Info(ctx, "Failed to set websocket upgrade")
		return
	}
	client := &Client{
		ID:        uuid.New().String(),
		Conn:      conn,
		ChannelID: channelID,
	}
	defer conn.Close()
	ctx = ctxlog.WithFields(ctx, ctxlog.Fields{
		"client_id": client.ID,
	})

	s.Register <- client
	s.keepAlive(ctx, client)
}

// keepalive - a simple keepalive that sends a websocket
// event every 15 seconds.
func (s *Streamer) keepAlive(ctx context.Context, client *Client) {
	defer func() {
		s.log.Info(ctx, "connection closed")
		s.Unregister <- client
	}()
	s.log.Info(ctx, "keepalive started")

	for {
		time.Sleep(15 * time.Second)
		s.log.Info(ctx, "sending keepalive")
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte("{\"message\": \"I'm alive\"}"))
		if err != nil {
			s.log.Error(ctx, "failed to send message on connection")
			return
		}
	}
}
