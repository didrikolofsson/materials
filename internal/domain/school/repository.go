// Package school provides repository interfaces for the school domain.
package school

import (
	"context"
	"database/sql"
)

type SubjectsRepository interface {
	List(ctx context.Context) ([]Subject, error)
}

type MySQLSubjectsRepository struct {
	db *sql.DB
}

func NewMySQLSubjectsRepository(db *sql.DB) *MySQLSubjectsRepository {
	return &MySQLSubjectsRepository{db: db}
}

func (r *MySQLSubjectsRepository) List(ctx context.Context) ([]Subject, error) {
	const query = `
		SELECT id, name, created_at
		FROM subjects
	`
	rows, err := r.db.QueryContext(ctx, query)
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
