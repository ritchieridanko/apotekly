package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
)

// Internal Error Codes
const (
	CodeBCryptHashingFailed   errCode = "ERR_BCRYPT_HASHING_FAILED"
	CodeCacheCommandExec      errCode = "ERR_CACHE_COMMAND_EXECUTION"
	CodeCacheScriptExec       errCode = "ERR_CACHE_SCRIPT_EXECUTION"
	CodeDBQueryExec           errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTx                  errCode = "ERR_DB_TX"
	CodeEmailNotAvailable     errCode = "ERR_EMAIL_NOT_AVAILABLE"
	CodeEventPublishingFailed errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeInvalidPayload        errCode = "ERR_INVALID_PAYLOAD"
	CodeJWTGenerationFailed   errCode = "ERR_JWT_GENERATION_FAILED"
	CodeMissingContextValue   errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata       errCode = "ERR_MISSING_METADATA"
	CodeUnknown               errCode = "ERR_UNKNOWN"
)

// External Error Messages
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidPayload         string = "Invalid payload"
)
