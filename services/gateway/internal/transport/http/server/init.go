package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/router"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
)

type Server struct {
	config *configs.HTTPServer
	server *http.Server
	router *router.Router
	logger *logger.Logger
}

func Init(cfg *configs.HTTPServer, r *router.Router, l *logger.Logger) *Server {
	return &Server{
		config: cfg,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r.Router(),
			ReadTimeout:  cfg.Timeout.Read,
			WriteTimeout: cfg.Timeout.Write,
		},
		router: r,
		logger: l,
	}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.logger.Log("[SERVER] is running (host=%s, port=%d)", s.config.Host, s.config.Port)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan error, 1)
	go func() {
		stopped <- s.server.Shutdown(ctx)
	}()

	select {
	case <-ctx.Done():
		if err := s.server.Close(); err != nil {
			return fmt.Errorf("failed to shutdown server: %v: %w", err, ctx.Err())
		}
		return fmt.Errorf("failed to shutdown server: %w", ctx.Err())
	case err := <-stopped:
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	}
}
