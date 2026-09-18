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

type ProductService struct {
	conn    *grpc.ClientConn
	product apis.ProductServiceClient
}

func NewProductService(cfg *configs.Service, l *zap.Logger) (*ProductService, error) {
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
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.Name), cfg.Host, cfg.Port,
	)
	return &ProductService{
		conn:    conn,
		product: apis.NewProductServiceClient(conn),
	}, nil
}

func (s *ProductService) ProductClient() apis.ProductServiceClient {
	return s.product
}

func (s *ProductService) Close() error {
	return s.conn.Close()
}
