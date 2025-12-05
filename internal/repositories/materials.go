package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialsRepository interface {
	Create(ctx context.Context, tx TxOrDB, m *models.Material) (models.MaterialID, error)
	UpdateCurrentVersion(ctx context.Context, m, v int64) error
	GetByID(ctx context.Context, id int64) (*models.Material, error)
}

type MySQLMaterialsRepository struct {
	db *sql.DB
}

func NewMySQLMaterialsRepository(db *sql.DB) *MySQLMaterialsRepository {
	return &MySQLMaterialsRepository{db: db}
}

func (r *MySQLMaterialsRepository) Create(ctx context.Context, tx TxOrDB, m *models.Material) (models.MaterialID, error) {
	query := `
		INSERT INTO materials (teacher_id, subject_id)
		VALUES (?, ?);
	`

	result, err := tx.ExecContext(ctx, query, m.TeacherID, m.SubjectID)
	if err != nil {
		return 0, fmt.Errorf("failed to create material: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to create material: %w", err)
	}
	return models.MaterialID(id), nil
}

func (r *MySQLMaterialsRepository) UpdateCurrentVersion(ctx context.Context, m, v int64) error {
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

func (r *MySQLMaterialsRepository) GetByID(ctx context.Context, id int64) (*models.Material, error) {
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
