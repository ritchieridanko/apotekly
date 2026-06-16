package logger

import (
	"context"
	"time"

	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(l *zap.Logger) *Logger {
	return &Logger{logger: l}
}

func (l *Logger) Log(message string, args ...any) {
	l.logger.Sugar().Infof(message, args...)
}

func (l *Logger) Info(ctx context.Context, message string, fields ...Field) {
	l.logger.Info(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) Warn(ctx context.Context, message string, fields ...Field) {
	l.logger.Warn(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...Field) {
	l.logger.Error(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) toFields(ctx context.Context, fields ...Field) []zap.Field {
	defaultSize := 3
	zf := append(
		make([]zap.Field, 0, len(fields)+defaultSize),
		zap.Time("timestamp", time.Now().UTC()),
		zap.String("request_id", utils.CtxRequestID(ctx)),
		zap.String("trace_id", utils.CtxTraceID(ctx)),
	)

	for _, field := range fields {
		zf = append(zf, zap.Any(field.key, field.value))
	}
	return zf
}
