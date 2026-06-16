package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Init(env string) (*zap.Logger, error) {
	var conf zap.Config
	if env == "prod" {
		conf = zap.NewProductionConfig()
		conf.DisableStacktrace = true
		conf.DisableCaller = false
	} else {
		conf = zap.NewDevelopmentConfig()
		conf.DisableStacktrace = false
		conf.DisableCaller = false
	}

	l, err := conf.Build(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	l.Sugar().Infof("[LOGGER] initialized (env=%s, level=%s)", env, l.Level().String())
	return l, nil
}
