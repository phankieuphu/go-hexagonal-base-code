package http

import (
	"account-service/config"
	"account-service/internal/adapters/http/handler"
	"account-service/internal/domain/ports"
	"account-service/pkg/logger"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	health "github.com/phankieuphu/go-health-check"
	"github.com/phankieuphu/go-health-check/ginhealth"
)

type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
}

func NewServer(cfg config.API, accountService ports.AccountService, h *health.Handler) *Server {
	engine := gin.New()
	// Probes are registered before engine.Use, so they skip the middleware:
	// kubelet hits them every few seconds, which would flood the access log.
	ginhealth.RegisterRoutes(engine, h)
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

// Start blocks serving requests until Shutdown is called (returning nil)
// or the listener fails (returning the error).
func (s *Server) Start() error {
	logger.Info("HTTP server listening", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown stops accepting new connections and waits for in-flight requests
// to finish, or for ctx to expire.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
