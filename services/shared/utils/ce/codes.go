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
	CodeAlreadyExists          errCode = "ERR_ALREADY_EXISTS"
	CodeBCryptHashingFailed    errCode = "ERR_BCRYPT_HASHING_FAILED"
	CodeCacheCommandExec       errCode = "ERR_CACHE_COMMAND_EXECUTION"
	CodeCacheScriptExec        errCode = "ERR_CACHE_SCRIPT_EXECUTION"
	CodeDBQueryExec            errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTx                   errCode = "ERR_DB_TX"
	CodeEmailNotAvailable      errCode = "ERR_EMAIL_NOT_AVAILABLE"
	CodeEventPublishingFailed  errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeInternal               errCode = "ERR_INTERNAL"
	CodeInvalidPayload         errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeJWTGenerationFailed    errCode = "ERR_JWT_GENERATION_FAILED"
	CodeMissingContextValue    errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata        errCode = "ERR_MISSING_METADATA"
	CodeUnknown                errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed   errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidPayload         string = "Invalid payload"
)
