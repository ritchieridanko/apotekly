package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.Storage, l *zap.Logger) (*cloudinary.Cloudinary, error) {
	c, err := cloudinary.NewFromParams(cfg.Cloud, cfg.Key, cfg.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.Admin.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err)
	}

	l.Sugar().Infof("[STORAGE] initialized (provider=%s, cloud=%s)", cfg.Provider, cfg.Cloud)
	return c, nil
}
