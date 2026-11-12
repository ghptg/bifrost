package logstore

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSQLitePragmasApplied guards the connection settings: the pure-Go driver
// silently ignores connection-string keys it does not interpret, so a mistyped
// DSN would otherwise leave SQLite on its defaults without any error.
func TestSQLitePragmasApplied(t *testing.T) {
	store, err := newSqliteLogStore(context.Background(), &SQLiteConfig{
		Path: filepath.Join(t.TempDir(), "pragma.db"),
	}, testLogger{})
	require.NoError(t, err)
	sqlDB, err := store.db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	want := []struct {
		pragma string
		value  string
	}{
		{"foreign_keys", "1"},
		{"journal_mode", "wal"},
		{"busy_timeout", "60000"},
		{"cache_size", "10000"},
		{"synchronous", "1"},
		{"wal_autocheckpoint", "1000"},
	}
	for _, w := range want {
		var got string
		require.NoError(t, store.db.Raw("PRAGMA "+w.pragma).Row().Scan(&got), w.pragma)
		require.Equal(t, w.value, got, w.pragma)
	}
}
