package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type (
	AuthContext struct {
		AuthID          uint64
		Role            string
		IsEmailVerified bool
		PharmacyID      *uuid.UUID
	}

	TransportContext struct {
		IPAddress string
		UserAgent string
	}
)

// Get Auth Information from Context
func CtxAuth(ctx context.Context) *AuthContext {
	if v, ok := ctx.Value(constants.CtxKeyAuth).(*AuthContext); ok {
		return v
	}
	return nil
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

// Inject Metadata into Context
func CtxWithMetadata(ctx context.Context, kv ...string) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		append(
			kv,
			constants.MDKeyRequestID,
			CtxRequestID(ctx),
		)...,
	)
}
