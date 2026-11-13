//go:build lite

package objectstore

import (
	"bytes"
	"errors"
)

// The parquet encoder is not compiled into the lite build, so payload objects are
// stored as JSON. These stand-ins keep the package API intact: IsParquet still
// recognizes a parquet object (it only inspects the magic bytes) and the codec
// functions fail with a clear error instead of silently mis-handling the data.

var errParquetNotBuilt = errors.New("parquet payload encoding is not included in this build (built with the lite tag); use the json object format")

// IsParquet reports whether data starts with the parquet file magic.
func IsParquet(data []byte) bool {
	return len(data) >= 8 && bytes.Equal(data[:4], []byte("PAR1"))
}

// EncodeJSONObjectToParquet is unavailable in the lite build.
func EncodeJSONObjectToParquet(_ []byte, _ string) ([]byte, error) {
	return nil, errParquetNotBuilt
}

// DecodeParquetToJSONObject is unavailable in the lite build.
func DecodeParquetToJSONObject(_ []byte) ([]byte, error) {
	return nil, errParquetNotBuilt
}
