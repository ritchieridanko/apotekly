package ce

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode string

type Error struct {
	code    errCode
	message string
	err     error
	fields  []logger.Field
}

func NewError(ec errCode, message string, err error, fields ...logger.Field) *Error {
	return &Error{
		code:    ec,
		message: message,
		err:     err,
		fields:  fields,
	}
}

func (e *Error) Code() errCode {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *Error) Fields() []logger.Field {
	return e.fields
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Append(fields ...logger.Field) *Error {
	e.fields = append(e.fields, fields...)
	return e
}

func (e *Error) Bind(ctx *gin.Context) {
	ctx.Error(e)
}

func (e *Error) ToGRPCErr() error {
	switch e.code {
	case
		CodeInvalidParams,
		CodeInvalidPayload,
		CodeInvalidRequestMetadata,
		CodeInvalidToken,
		CodeTokenNotOwned:
		return status.Error(codes.InvalidArgument, e.message)
	case
		CodeAddressNotFound,
		CodeAuthNotFound,
		CodeNotFound,
		CodeSessionNotFound,
		CodeUserNotFound:
		return status.Error(codes.NotFound, e.message)
	case
		CodeAlreadyExists,
		CodeEmailNotAvailable,
		CodePharmacyAlreadyExists,
		CodeUserAlreadyExists:
		return status.Error(codes.AlreadyExists, e.message)
	case
		CodePermissionDenied,
		CodeRoleNotAuthorized:
		return status.Error(codes.PermissionDenied, e.message)
	case
		CodeEmailAlreadyVerified,
		CodeEmailNotVerified,
		CodeFailedPrecondition,
		CodeOAuthEmailChange,
		CodeOAuthPasswordChange:
		return status.Error(codes.FailedPrecondition, e.message)
	case
		CodeAuthNotRegistered,
		CodeEmailNotRegistered,
		CodeInvalidSession,
		CodeOAuthRegularSignIn,
		CodeSessionExpired,
		CodeSessionNotOwned,
		CodeUnauthenticated,
		CodeWrongPassword:
		return status.Error(codes.Unauthenticated, e.message)
	case
		CodeBCryptHashingFailed,
		CodeCacheCommandExec,
		CodeCacheScriptExec,
		CodeDBQueryExec,
		CodeDBTx,
		CodeEmailDeliveryFailed,
		CodeEmailTemplatingFailed,
		CodeEventCommittingFailed,
		CodeEventFetchingFailed,
		CodeEventPublishingFailed,
		CodeEventTopicNotRegistered,
		CodeInternal,
		CodeJSONRawEncodingFailed,
		CodeJSONUnmarshallingFailed,
		CodeJWTGenerationFailed,
		CodeMissingContextValue,
		CodeMissingMetadata,
		CodeNoPendingEventInbox,
		CodeOrphanedEventInbox,
		CodePanicOccurred,
		CodeProtobufParsingFailed,
		CodeTypeAssertionFailed,
		CodeTypeConversionFailed,
		CodeURLGenerationFailed,
		CodeUUIDGenerationFailed:
		return status.Error(codes.Internal, e.message)
	case CodeUnknown:
		return status.Error(codes.Unknown, e.message)
	default:
		return status.Error(codes.Unknown, e.message)
	}
}

func (e *Error) ToHTTPErr() int {
	switch e.code {
	case
		CodeInvalidParams,
		CodeInvalidPayload,
		CodeInvalidRequestMetadata:
		return http.StatusBadRequest
	case
		CodeRefreshTokenNotFound,
		CodeUnauthenticated:
		return http.StatusUnauthorized
	case
		CodeFailedPrecondition,
		CodePermissionDenied:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists:
		return http.StatusConflict
	case
		CodeInternal,
		CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
