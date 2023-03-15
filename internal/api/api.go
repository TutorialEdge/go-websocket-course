package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Stream interface {
	Stream(context.Context) error
}

type Handler struct {
	router *gin.Engine
	Stream Stream
}

func New() *Handler {
	var handler *Handler
	handler.setupRoutes()
	return handler
}

func (h *Handler) setupRoutes() {
	h.router = gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:1313", "https://tutorialedge.net"}

	h.router.Use(cors.New(config))

	h.router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world!"})
	})

	// /api/v1/channels - used on initial web page load
	// we can list all the channels and the client can select which
	// channel they'd like to stream in real-time
	h.router.GET("/api/v1/channels", h.getChannels)

	// /api/v1/stream - websocket endpoint that folks connect into
	// whenever new events are consumed from RabbitMQ, they're then sent to all clients
	// connected into this channel.
	h.router.GET("/api/v1/stream", h.stream)
}

func (h *Handler) getChannels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"channels": "a list of channels"})
}

func (h *Handler) stream(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {

	}
}
