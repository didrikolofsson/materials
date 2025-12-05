package repositories

import (
	"context"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialVersionsRepository interface {
	Create(ctx context.Context, tx TxOrDB, v *models.MaterialVersion) (models.MaterialVersionID, error)
	// GetAllByMaterialID(ctx context.Context, materialID string) ([]models.MaterialVersion, error)
	// GetByID(ctx context.Context, materialID, versionID string) (*models.MaterialVersion, error)
}

type MySQLMaterialVersionsRepository struct{}

func NewMySQLMaterialVersionsRepository() *MySQLMaterialVersionsRepository {
	return &MySQLMaterialVersionsRepository{}
}

func (r *MySQLMaterialVersionsRepository) Create(
	ctx context.Context, tx TxOrDB, v *models.MaterialVersion,
) (models.MaterialVersionID, error) {
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
	return models.MaterialVersionID(id), nil
}
