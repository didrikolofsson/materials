// Package mysql provides a MySQL database connection.
package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func New(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to open database: %w", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}
	return db
}
