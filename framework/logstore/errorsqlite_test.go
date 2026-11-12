package logstore

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestIsTransientWriteErrorRealSQLiteBusy drives a real SQLITE_BUSY through the
// pure-Go driver: one connection holds the write lock with BEGIN IMMEDIATE while a
// second (busy_timeout=0, so it fails at once) tries to insert. The error the logs
// writer sees for a contended writer lock must classify as transient, so the batch
// is retried whole instead of fanning out into per-row inserts.
func TestIsTransientWriteErrorRealSQLiteBusy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "busy.db")
	dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(0)"
	holder, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, holder.Exec("CREATE TABLE t (id TEXT PRIMARY KEY)").Error)
	holderSQL, err := holder.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = holderSQL.Close() })

	conn, err := holderSQL.Conn(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	_, err = conn.ExecContext(context.Background(), "BEGIN IMMEDIATE")
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = conn.ExecContext(context.Background(), "ROLLBACK") })

	writer, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	writerSQL, err := writer.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = writerSQL.Close() })

	err = writer.Exec("INSERT INTO t (id) VALUES ('a')").Error
	require.Error(t, err)
	require.True(t, isSQLiteTransientError(err), "real SQLITE_BUSY must classify as transient: %T %v", err, err)
	require.True(t, IsTransientWriteError(err), "real SQLITE_BUSY must be transient: %v", err)
	require.False(t, IsPostgresStatementTimeoutError(err))
}

// TestIsTransientWriteErrorRealSQLiteConstraint checks a row-level failure is not
// classified as transient.
func TestIsTransientWriteErrorRealSQLiteConstraint(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "c.db")), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	require.NoError(t, db.Exec("CREATE TABLE t (id TEXT PRIMARY KEY)").Error)
	require.NoError(t, db.Exec("INSERT INTO t (id) VALUES ('a')").Error)
	err = db.Exec("INSERT INTO t (id) VALUES ('a')").Error
	require.Error(t, err)
	require.False(t, isSQLiteTransientError(err), "constraint error must not be transient: %T %v", err, err)
	require.False(t, IsTransientWriteError(err))
}
