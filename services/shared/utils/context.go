package utils

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"go.opentelemetry.io/otel/trace"
)

type TransportContext struct {
	IPAddress string
	UserAgent string
}

// Get Request ID from Context
func CtxRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// Get Trace ID from Context
func CtxTraceID(ctx context.Context) string {
	if sp := trace.SpanFromContext(ctx); sp.SpanContext().HasTraceID() {
		return sp.SpanContext().TraceID().String()
	}
	return ""
}

// Get Transport Information from Context
func CtxTransport(ctx context.Context) *TransportContext {
	if v, ok := ctx.Value(constants.CtxKeyTransport).(*TransportContext); ok {
		return v
	}
	return nil
}
