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

type UserService struct {
	conn    *grpc.ClientConn
	user    apis.UserServiceClient
	address apis.AddressServiceClient
}

func NewUserService(cfg *configs.Service, l *zap.Logger) (*UserService, error) {
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
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.Name), cfg.Host, cfg.Port,
	)
	return &UserService{
		conn:    conn,
		user:    apis.NewUserServiceClient(conn),
		address: apis.NewAddressServiceClient(conn),
	}, nil
}

func (s *UserService) UserClient() apis.UserServiceClient {
	return s.user
}

func (s *UserService) AddressClient() apis.AddressServiceClient {
	return s.address
}

func (s *UserService) Close() error {
	return s.conn.Close()
}
