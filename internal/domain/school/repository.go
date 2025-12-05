package school

import (
	"context"
	"database/sql"
)

// TxOrDB interface allows methods to work with either *sql.DB or *sql.Tx
type TxOrDB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type SubjectsRepository interface {
	// List accepts optional transaction - pass nil for direct DB query
	List(ctx context.Context, tx *sql.Tx) ([]Subject, error)
}

type MySQLSubjectsRepository struct {
	db *sql.DB
}

func NewMySQLSubjectsRepository(db *sql.DB) *MySQLSubjectsRepository {
	return &MySQLSubjectsRepository{db: db}
}

func (r *MySQLSubjectsRepository) List(ctx context.Context, tx *sql.Tx) ([]Subject, error) {
	const query = `
		SELECT id, name, created_at
		FROM subjects
	`

	// Use transaction if provided, otherwise use repository's db connection
	var db TxOrDB = r.db
	if tx != nil {
		db = tx
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var subject Subject
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.CreatedAt); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subjects, nil
}
