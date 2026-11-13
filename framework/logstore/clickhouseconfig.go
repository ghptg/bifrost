package logstore

import "github.com/maximhq/bifrost/core/schemas"

// Kept outside the clickhouse*.go files, which are excluded by the lite build tag,
// because the config types and the store factory reference it in every build.

// ClickHouseConfig represents the configuration for a ClickHouse log store.
//
// ClickHouse is an append-only columnar OLAP store. The backend uses
// ReplacingMergeTree tables with a connection-level `final = 1` setting so
// reads transparently see the latest version of each row (see clickhousestore.go
// for the mutation strategy).
type ClickHouseConfig struct {
	Host     *schemas.SecretVar `json:"host"`
	Port     *schemas.SecretVar `json:"port"`
	Database *schemas.SecretVar `json:"database"`
	Username *schemas.SecretVar `json:"username"`
	Password *schemas.SecretVar `json:"password"`
	// Protocol selects the ClickHouse wire protocol: "native" (default, port
	// 9000/9440) or "http" (port 8123/8443). clickhouse-go derives the protocol
	// from the DSN scheme, so this maps to clickhouse:// vs http(s)://.
	Protocol string `json:"protocol,omitempty"`
	// Secure enables TLS (native: secure=true; http: switches to https).
	Secure bool `json:"secure,omitempty"`
	// DialTimeout is the connection dial timeout in milliseconds (JSON config
	// duration fields are integer milliseconds). 0 means the 10s default.
	DialTimeout int `json:"dial_timeout,omitempty"`
	// MaxQuerySize is the ClickHouse max_query_size setting in bytes. 0 means the 16 MiB default.
	MaxQuerySize int `json:"max_query_size,omitempty"`
	// Cluster, when set, makes DDL run as `ON CLUSTER <name>` against
	// ReplicatedReplacingMergeTree engines. Empty means single-node.
	Cluster string `json:"cluster,omitempty"`
}
