package ce

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Internal Errors
var (
	ErrCacheNoResult   error = redis.Nil
	ErrCookieNotFound  error = http.ErrNoCookie
	ErrDBAffectNoRows  error = errors.New("no rows affected")
	ErrDBQueryNoRows   error = pgx.ErrNoRows
	ErrInvalidJWTClaim error = jwt.ErrTokenInvalidClaims
	ErrJWTExpired      error = jwt.ErrTokenExpired
	ErrJWTMalformed    error = jwt.ErrTokenMalformed
)

// Internal Error Codes
const (
	CodeAlreadyExists           errCode = "ERR_ALREADY_EXISTS"
	CodeAuthNotFound            errCode = "ERR_AUTH_NOT_FOUND"
	CodeAuthNotRegistered       errCode = "ERR_AUTH_NOT_REGISTERED"
	CodeBCryptHashingFailed     errCode = "ERR_BCRYPT_HASHING_FAILED"
	CodeCacheCommandExec        errCode = "ERR_CACHE_COMMAND_EXECUTION"
	CodeCacheScriptExec         errCode = "ERR_CACHE_SCRIPT_EXECUTION"
	CodeDBQueryExec             errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTx                    errCode = "ERR_DB_TX"
	CodeEmailAlreadyVerified    errCode = "ERR_EMAIL_ALREADY_VERIFIED"
	CodeEmailDeliveryFailed     errCode = "ERR_EMAIL_DELIVERY_FAILED"
	CodeEmailNotAvailable       errCode = "ERR_EMAIL_NOT_AVAILABLE"
	CodeEmailNotRegistered      errCode = "ERR_EMAIL_NOT_REGISTERED"
	CodeEmailNotVerified        errCode = "ERR_EMAIL_NOT_VERIFIED"
	CodeEmailTemplatingFailed   errCode = "ERR_EMAIL_TEMPLATING_FAILED"
	CodeEventCommittingFailed   errCode = "ERR_EVENT_COMMITTING_FAILED"
	CodeEventFetchingFailed     errCode = "ERR_EVENT_FETCHING_FAILED"
	CodeEventPublishingFailed   errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeEventTopicNotRegistered errCode = "ERR_EVENT_TOPIC_NOT_REGISTERED"
	CodeFailedPrecondition      errCode = "ERR_FAILED_PRECONDITION"
	CodeInternal                errCode = "ERR_INTERNAL"
	CodeInvalidParams           errCode = "ERR_INVALID_PARAMS"
	CodeInvalidPayload          errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata  errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeInvalidToken            errCode = "ERR_INVALID_TOKEN"
	CodeJSONRawEncodingFailed   errCode = "ERR_JSON_RAW_ENCODING_FAILED"
	CodeJSONUnmarshallingFailed errCode = "ERR_JSON_UNMARSHALLING_FAILED"
	CodeJWTGenerationFailed     errCode = "ERR_JWT_GENERATION_FAILED"
	CodeMissingContextValue     errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata         errCode = "ERR_MISSING_METADATA"
	CodeNoPendingEventInbox     errCode = "ERR_NO_PENDING_EVENT_INBOX"
	CodeNotFound                errCode = "ERR_NOT_FOUND"
	CodeOAuthRegularSignIn      errCode = "ERR_OAUTH_REGULAR_SIGN_IN"
	CodeOrphanedEventInbox      errCode = "ERR_ORPHANED_EVENT_INBOX"
	CodePanicOccurred           errCode = "ERR_PANIC_OCCURRED"
	CodePermissionDenied        errCode = "ERR_PERMISSION_DENIED"
	CodeProtobufParsingFailed   errCode = "ERR_PROTOBUF_PARSING_FAILED"
	CodeRefreshTokenNotFound    errCode = "ERR_REFRESH_TOKEN_NOT_FOUND"
	CodeRoleNotAuthorized       errCode = "ERR_ROLE_NOT_AUTHORIZED"
	CodeSessionExpired          errCode = "ERR_SESSION_EXPIRED"
	CodeSessionNotFound         errCode = "ERR_SESSION_NOT_FOUND"
	CodeSessionNotOwned         errCode = "ERR_SESSION_NOT_OWNED"
	CodeTokenNotOwned           errCode = "ERR_TOKEN_NOT_OWNED"
	CodeTypeConversionFailed    errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeUnauthenticated         errCode = "ERR_UNAUTHENTICATED"
	CodeUnknown                 errCode = "ERR_UNKNOWN"
	CodeURLGenerationFailed     errCode = "ERR_URL_GENERATION_FAILED"
	CodeUUIDGenerationFailed    errCode = "ERR_UUID_GENERATION_FAILED"
	CodeWrongPassword           errCode = "ERR_WRONG_PASSWORD"
)

// External Error Messages
const (
	MsgAuthNotFound           string = "Auth not found"
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgEmailAlreadyVerified   string = "Email is already verified"
	MsgEmailNotVerified       string = "Email is not verified"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidCredentials     string = "Invalid credentials"
	MsgInvalidParams          string = "Invalid params"
	MsgInvalidPayload         string = "Invalid payload"
	MsgInvalidSession         string = "Invalid session"
	MsgInvalidToken           string = "Invalid token"
	MsgNoPendingEventInbox    string = "No pending event inbox"
	MsgOrphanedEventInbox     string = "Orphaned event inbox"
	MsgSessionExpired         string = "Session expired"
	MsgSessionNotFound        string = "Session not found"
	MsgUnauthenticated        string = "Unauthenticated"
	MsgUnauthorized           string = "Unauthorized"
)
