package logstore

import (
	"errors"

	sqlite "github.com/glebarez/go-sqlite"
)

// SQLite primary result codes (https://www.sqlite.org/rescode.html).
const (
	sqliteResultBusy   = 5
	sqliteResultLocked = 6
)

// isSQLiteTransientError reports whether err is SQLITE_BUSY or SQLITE_LOCKED
// (any extended variant, e.g. SQLITE_BUSY_SNAPSHOT) from the SQLite driver.
// Both mean another connection holds the lock the write needs; they say nothing
// about the rows, so the batch can be retried whole once the holder commits.
// Falling back to per-row inserts instead makes every row wait out the
// connection's busy_timeout again.
func isSQLiteTransientError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		// An extended code carries its primary code in the low byte.
		base := sqliteErr.Code() & 0xff
		return base == sqliteResultBusy || base == sqliteResultLocked
	}
	return false
}
