package repositories

import (
	"context"

	"github.com/didrikolofsson/materials/internal/models"
)

// SubjectsRepository defines the interface for subject data access
type SubjectsRepository interface {
	// List returns all subjects
	List(ctx context.Context, tx TxOrDB) ([]models.Subject, error)
}

type MySQLSubjectsRepository struct{}

func NewMySQLSubjectsRepository() *MySQLSubjectsRepository {
	return &MySQLSubjectsRepository{}
}

func (r *MySQLSubjectsRepository) List(ctx context.Context, tx TxOrDB) ([]models.Subject, error) {
	const query = `
		SELECT id, name, created_at
		FROM subjects
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
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
