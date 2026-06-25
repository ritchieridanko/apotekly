package services

import (
	"fmt"
	"strings"

	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/ritchieridanko/apotekly/services/shared/contract/apis/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	conn   *grpc.ClientConn
	client apis.AuthServiceClient
}

func NewAuthService(cfg *configs.Service, l *zap.Logger) (*AuthService, error) {
	conn, err := grpc.NewClient(
		cfg.Addr,
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.Name), cfg.Host, cfg.Port,
	)
	return &AuthService{
		conn:   conn,
		client: apis.NewAuthServiceClient(conn),
	}, nil
}

func (s *AuthService) Client() apis.AuthServiceClient {
	return s.client
}

func (s *AuthService) Close() error {
	return s.conn.Close()
}
