package repositories

import (
	"context"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialsRepository interface {
	Create(ctx context.Context, tx TxOrDB, m *models.Material) (models.GenericID, error)
	UpdateOriginalMaterialID(ctx context.Context, tx TxOrDB, m, o models.GenericID) (models.GenericID, error)
	ListAllByTeacherID(ctx context.Context, tx TxOrDB, teacherID models.GenericID) ([]models.Material, error)
	UpdateCurrentVersion(ctx context.Context, tx TxOrDB, m, v models.GenericID) (models.GenericID, error)
	GetByID(ctx context.Context, tx TxOrDB, id models.GenericID) (*models.Material, error)
}

type MySQLMaterialsRepository struct{}

func NewMySQLMaterialsRepository() *MySQLMaterialsRepository {
	return &MySQLMaterialsRepository{}
}

func (r *MySQLMaterialsRepository) Create(
	ctx context.Context, tx TxOrDB, m *models.Material,
) (models.GenericID, error) {
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
	return models.GenericID(id), nil
}

func (r *MySQLMaterialsRepository) UpdateOriginalMaterialID(ctx context.Context, tx TxOrDB, m, o models.GenericID) (models.GenericID, error) {
	query := `
		UPDATE materials
		SET original_material_id = ?
		WHERE id = ?;
	`
	result, err := tx.ExecContext(ctx, query, o, m)
	if err != nil {
		return 0, fmt.Errorf("failed to update original material id: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to update original material id: %w", err)
	}
	return models.GenericID(id), nil
}

func (r *MySQLMaterialsRepository) ListAllByTeacherID(ctx context.Context, tx TxOrDB, teacherID models.GenericID) ([]models.Material, error) {
	query := `
		SELECT id, teacher_id, subject_id, original_material_id, current_version_id, created_at
		FROM materials
		WHERE teacher_id = ?;
	`
	rows, err := tx.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to list all materials by teacher id: %w", err)
	}
	defer rows.Close()

	var materials []models.Material
	for rows.Next() {
		var m models.Material
		if err := rows.Scan(&m.ID, &m.TeacherID, &m.SubjectID, &m.OriginalMaterialID, &m.CurrentVersionID, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan material: %w", err)
		}
		materials = append(materials, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list all materials by teacher id: %w", err)
	}
	return materials, nil
}

func (r *MySQLMaterialsRepository) UpdateCurrentVersion(ctx context.Context, tx TxOrDB, m, v models.GenericID) (models.GenericID, error) {
	// First verify that the version exists and belongs to this material
	var versionMaterialID models.GenericID
	checkQuery := `
		SELECT material_id
		FROM material_versions
		WHERE id = ?;
	`
	if err := tx.QueryRowContext(ctx, checkQuery, v).Scan(&versionMaterialID); err != nil {
		return 0, fmt.Errorf("material version %d does not exist: %w", v, err)
	}
	if versionMaterialID != m {
		return 0, fmt.Errorf("material version %d does not belong to material %d", v, m)
	}

	query := `
		UPDATE materials
		SET current_version_id = ?
		WHERE id = ?;
	`

	_, err := tx.ExecContext(ctx, query, v, m)
	if err != nil {
		return 0, fmt.Errorf("failed to update current version: %w", err)
	}
	return v, nil
}

func (r *MySQLMaterialsRepository) GetByID(ctx context.Context, tx TxOrDB, id models.GenericID) (*models.Material, error) {
	query := `
		SELECT id, teacher_id, subject_id, original_material_id, current_version_id, created_at
		FROM materials
		WHERE id = ?;
	`

	var m models.Material
	if err := tx.QueryRowContext(
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
