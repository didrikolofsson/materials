// Package repositories provides data access for the school domain.
package repositories

import (
	"context"
	"database/sql"
)

// TxOrDB interface allows methods to work with either *sql.DB or *sql.Tx
type TxOrDB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
