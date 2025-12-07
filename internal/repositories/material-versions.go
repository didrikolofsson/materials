package repositories

import (
	"context"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialVersionsRepository interface {
	Create(ctx context.Context, tx TxOrDB, v *models.MaterialVersion) (models.GenericID, error)
	SetMainForMaterialVersion(ctx context.Context, tx TxOrDB, materialID, materialVersionID models.GenericID) error
	ListAllByMaterialID(ctx context.Context, tx TxOrDB, materialID models.GenericID) ([]models.MaterialVersion, error)
	// GetByID(ctx context.Context, tx TxOrDB, id models.GenericID) (*models.MaterialVersion, error)
}

type MySQLMaterialVersionsRepository struct{}

func NewMySQLMaterialVersionsRepository() *MySQLMaterialVersionsRepository {
	return &MySQLMaterialVersionsRepository{}
}

func (r *MySQLMaterialVersionsRepository) Create(
	ctx context.Context, tx TxOrDB, v *models.MaterialVersion,
) (models.GenericID, error) {
	query := `
		INSERT INTO material_versions (material_id, title, summary, description, content, version_number, is_main)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	res, err := tx.ExecContext(
		ctx, query, v.MaterialID, v.Title, v.Summary, v.Description, v.Content, v.VersionNumber, v.IsMain,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create material version: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to create material version: %w", err)
	}
	return models.GenericID(id), nil
}

// SetMainForMaterialVersion sets the material version as main for a given material, only one material version can be main.
func (r *MySQLMaterialVersionsRepository) SetMainForMaterialVersion(
	ctx context.Context, tx TxOrDB, materialID, materialVersionID models.GenericID,
) error {
	// Verify that both the material and material version exist
	var materialExists, materialVersionExists bool
	checkQuery := `
		SELECT 
			EXISTS(SELECT 1 FROM materials WHERE id = ?) as material_exists,
			EXISTS(SELECT 1 FROM material_versions WHERE id = ?) as material_version_exists
	`
	if err := tx.QueryRowContext(
		ctx, checkQuery, materialID, materialVersionID,
	).Scan(&materialExists, &materialVersionExists); err != nil {
		return fmt.Errorf("failed to check if material and material version exist: %w", err)
	}
	if !materialExists {
		return fmt.Errorf("material %d does not exist", materialID)
	}
	if !materialVersionExists {
		return fmt.Errorf("material version %d does not exist", materialVersionID)
	}

	// Verify that the material version belongs to the material
	var materialVersionBelongsToMaterial bool
	checkMaterialVersionBelongsToMaterialQuery := `
		SELECT 
			EXISTS(SELECT 1 FROM material_versions WHERE id = ? AND material_id = ?) as material_version_belongs_to_material
	`
	if err := tx.QueryRowContext(
		ctx, checkMaterialVersionBelongsToMaterialQuery, materialVersionID, materialID,
	).Scan(&materialVersionBelongsToMaterial); err != nil {
		return fmt.Errorf("failed to check if material version belongs to material: %w", err)
	}
	if !materialVersionBelongsToMaterial {
		return fmt.Errorf("material version %d does not belong to material %d", materialVersionID, materialID)
	}

	// Set the material version as main
	query := `
		UPDATE material_versions
		SET is_main = (id = ?)
		WHERE material_id = ?;
	`
	_, err := tx.ExecContext(
		ctx, query, materialVersionID, materialID,
	)
	if err != nil {
		return fmt.Errorf("failed to set material version as main: %w", err)
	}
	return nil
}

func (r *MySQLMaterialVersionsRepository) ListAllByMaterialID(
	ctx context.Context, tx TxOrDB, materialID models.GenericID,
) ([]models.MaterialVersion, error) {
	// Verify that the material exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM materials WHERE id = ?)`
	if err := tx.QueryRowContext(ctx, checkQuery, materialID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check if material exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("material %d does not exist", materialID)
	}

	query := `
		SELECT id, material_id, title, summary, description, content, version_number, is_main, created_at
		FROM material_versions
		WHERE material_id = ?;
	`
	rows, err := tx.QueryContext(ctx, query, materialID)
	if err != nil {
		return nil, fmt.Errorf("failed to list all material versions by material id: %w", err)
	}
	defer rows.Close()

	var materialVersions []models.MaterialVersion
	for rows.Next() {
		var materialVersion models.MaterialVersion
		if err := rows.Scan(
			&materialVersion.ID,
			&materialVersion.MaterialID,
			&materialVersion.Title,
			&materialVersion.Summary,
			&materialVersion.Description,
			&materialVersion.Content,
			&materialVersion.VersionNumber,
			&materialVersion.IsMain,
			&materialVersion.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan material version: %w", err)
		}
		materialVersions = append(materialVersions, materialVersion)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list all material versions by material id: %w", err)
	}
	return materialVersions, nil
}
