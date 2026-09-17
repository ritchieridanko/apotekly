package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ritchieridanko/apotekly/services/product/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/apotekly/services/product/internal/transport/rpc/interceptors"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	config *configs.GRPCServer
	server *grpc.Server
	logger *logger.Logger
	ph     *handlers.ProductHandler
}

func Init(cfg *configs.GRPCServer, l *logger.Logger, ph *handlers.ProductHandler) *Server {
	srv := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
		grpc.ChainUnaryInterceptor(
			interceptors.Request(l),
			interceptors.Recovery(l),
			interceptors.Logging(l),
			interceptors.Auth(),
		),
	)

	apis.RegisterProductServiceServer(srv, ph)

	return &Server{
		config: cfg,
		server: srv,
		logger: l,
		ph:     ph,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to build listener: %w", err)
	}
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.logger.Log("[SERVER] is running (host=%s, port=%d)", s.config.Host, s.config.Port)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return fmt.Errorf("failed to shutdown server: %w", ctx.Err())
	case <-stopped:
		return nil
	}
}
