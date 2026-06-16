package constants

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "x-request-id"
	CtxKeyTransport ctxKey = "x-transport"
)
