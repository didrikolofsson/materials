package repositories

import (
	"context"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

// TeachersRepository defines the interface for teacher data access
type TeachersRepository interface {
	// List returns all teachers
	List(ctx context.Context, tx TxOrDB) ([]models.Teacher, error)
}

type MySQLTeachersRepository struct{}

func NewMySQLTeachersRepository() *MySQLTeachersRepository {
	return &MySQLTeachersRepository{}
}

func (r *MySQLTeachersRepository) List(ctx context.Context, tx TxOrDB) ([]models.Teacher, error) {
	query := `
		SELECT id, name, created_at
		FROM teachers;
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teachers, nil
}
