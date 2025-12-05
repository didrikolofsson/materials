package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialsRepository interface {
	Create(ctx context.Context, m *models.Material) error
	UpdateCurrentVersion(ctx context.Context, m, v string) error
	GetByID(ctx context.Context, id string) (*models.Material, error)
}

type MySQLMaterialsRepository struct {
	db *sql.DB
}

func NewMySQLMaterialsRepository(db *sql.DB) *MySQLMaterialsRepository {
	return &MySQLMaterialsRepository{db: db}
}

func (r *MySQLMaterialsRepository) Create(ctx context.Context, m *models.Material) error {
	query := `
		INSERT INTO materials (id, teacher_id, subject_id, original_material_id)
		VALUES (?, ?, ?, ?);
	`

	_, err := r.db.ExecContext(ctx, query, m.ID, m.TeacherID, m.SubjectID, m.OriginalMaterialID)
	if err != nil {
		return fmt.Errorf("failed to create material: %w", err)
	}
	return nil
}

func (r *MySQLMaterialsRepository) UpdateCurrentVersion(ctx context.Context, m, v string) error {
	query := `
		UPDATE materials
		SET current_version_id = ?
		WHERE id = ?;
	`

	_, err := r.db.ExecContext(ctx, query, v, m)
	if err != nil {
		return fmt.Errorf("failed to update current version: %w", err)
	}
	return nil
}

func (r *MySQLMaterialsRepository) GetByID(ctx context.Context, id string) (*models.Material, error) {
	query := `
		SELECT id, teacher_id, subject_id, original_material_id, current_version_id, created_at
		FROM materials
		WHERE id = ?;
	`

	var m models.Material
	if err := r.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&m.ID,
		&m.TeacherID,
		&m.SubjectID,
		&m.OriginalMaterialID,
		&m.CurrentVersionID,
		&m.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get material by id: %w", err)
	}
	return &m, nil
}
