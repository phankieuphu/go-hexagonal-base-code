package http

import (
	"account-service/config"
	"account-service/internal/adapters/http/handler"
	"account-service/internal/domain/ports"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
}

func NewServer(cfg config.API, accountService ports.AccountService) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	v1 := engine.Group("/api/v1")
	handler.NewAccountHandler(accountService).RegisterRoutes(v1)

	return &Server{
		engine: engine,
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      engine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Start() {
	log.Printf("HTTP server listening on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
