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

type PharmacyService struct {
	conn     *grpc.ClientConn
	pharmacy apis.PharmacyServiceClient
}

func NewPharmacyService(cfg *configs.Service, l *zap.Logger) (*PharmacyService, error) {
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
		return nil, fmt.Errorf("failed to connect to pharmacy service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.Name), cfg.Host, cfg.Port,
	)
	return &PharmacyService{
		conn:     conn,
		pharmacy: apis.NewPharmacyServiceClient(conn),
	}, nil
}

func (s *PharmacyService) PharmacyClient() apis.PharmacyServiceClient {
	return s.pharmacy
}

func (s *PharmacyService) Close() error {
	return s.conn.Close()
}
