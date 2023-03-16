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
		router:        gin.Default(),
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	handler.router.Use(cors.New(config))

	// basic health check endpoint
	handler.router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world!"})
	})

	return handler
}

// A handy-dandy way for us to allow components to register new routes
func (h *Handler) RegisterRoutes(service Service) {
	service.SetupRoutes(h.router)
}

// Serve - kicks off our server, no graceful handling of shutdowns yet
func (h *Handler) Serve() error {
	if err := h.router.Run(":8080"); err != nil {
		return err
	}
	return nil
}
