//go:build lite

package logstore

import (
	"context"
	"errors"

	"github.com/maximhq/bifrost/core/schemas"
)

// newClickHouseLogStore is the lite-build stand-in: the ClickHouse driver is not
// compiled in, so selecting a ClickHouse log store fails at startup with a clear error.
func newClickHouseLogStore(_ context.Context, _ *ClickHouseConfig, _ int, _ schemas.Logger) (LogStore, error) {
	return nil, errors.New("the ClickHouse log store is not included in this build (built with the lite tag)")
}
