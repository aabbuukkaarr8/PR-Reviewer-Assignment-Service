package store

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrDBCheckNilConn = errors.New("db connection nil")
var ErrUpdateFailed = errors.New("update failed")
var ErrInsertFailed = errors.New("insert failed")
var ErrGetFailed = errors.New("get failed")
var ErrFindFailed = errors.New("find failed")
var ErrNotFound = errors.New("not found")
var ErrNoRecords = errors.New("no records found")
var ErrDuplicateKey = errors.New("duplicate key")
var ErrForeignKeyViolated = errors.New("foreign key violated")
var ErrConflict = errors.New("conflict")

func DBErrToErr(err error) (res error, ok bool) {
	var dbErr *pq.Error
	if !errors.As(err, &dbErr) {
		return nil, false
	}

	switch {
	case dbErr.Code == "23503":
		return fmt.Errorf("%w: %s: %s", ErrForeignKeyViolated, dbErr.Message, dbErr.Detail), true
	case dbErr.Code == "23505":
		return fmt.Errorf("%w: %s: %s", ErrDuplicateKey, dbErr.Message, dbErr.Detail), true
	default:
		return nil, false
	}
}
