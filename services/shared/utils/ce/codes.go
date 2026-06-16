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
	CodeDBQueryExec errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTx        errCode = "ERR_DB_TX"
)

// External Error Messages
const (
	MsgInternalServer string = "Internal server error"
)
