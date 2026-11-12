package logstore

import (
	"context"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/maximhq/bifrost/framework/storehooks"
	"gorm.io/gorm"
)

// SQLiteConfig represents the configuration for a SQLite database.
type SQLiteConfig struct {
	Path string `json:"path"`
}

// newSqliteLogStore creates a new SQLite log store.
func newSqliteLogStore(ctx context.Context, config *SQLiteConfig, logger schemas.Logger) (*RDBLogStore, error) {
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		// Create DB file
		f, err := os.Create(config.Path)
		if err != nil {
			return nil, err
		}
		_ = f.Close()
	}
	// Configure SQLite with proper settings to handle concurrent access
	// The pure-Go driver (modernc) silently ignores connection-string keys it does not
	// interpret, so every setting is spelled as an explicit _pragma. busy_timeout comes
	// first so the later pragmas (journal_mode takes a lock) wait instead of failing.
	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(60000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=cache_size(10000)&_pragma=wal_autocheckpoint(1000)&_pragma=foreign_keys(1)", config.Path)
	logger.Debug("opening DB with dsn: %s", dsn)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newGormLogger(logger),
	})

	if err != nil {
		return nil, err
	}
	logger.Debug("db opened for logstore")
	storehooks.PoolOpened(storehooks.RoleLogs, db)

	s := &RDBLogStore{db: db, logger: logger}
	// Run migrations
	if err := triggerMigrations(ctx, db, logger); err != nil {
		return nil, err
	}

	return s, nil
}
