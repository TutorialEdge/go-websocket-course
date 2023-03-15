package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Service - if we want to expose a
type Service interface {
	SetupRoutes(*gin.Engine)
}

type Handler struct {
	router        *gin.Engine
	StreamService Service
}

func New(stream Service) *Handler {
	handler := &Handler{
		StreamService: stream,
	}
	handler.setupRoutes()
	return handler
}

func (h *Handler) RegisterRoutes(service Service) {
	service.SetupRoutes(h.router)
}

func (h *Handler) setupRoutes() {
	h.router = gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	h.router.Use(cors.New(config))
	// basic health check endpoint
	h.router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world!"})
	})
}

func (h *Handler) Serve() error {
	if err := h.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}
