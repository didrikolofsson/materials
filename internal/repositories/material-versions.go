package repositories

import (
	"context"
	"database/sql"

	"github.com/didrikolofsson/materials/internal/models"
)

type MaterialVersionsRepository interface {
	Create(ctx context.Context, v *models.MaterialVersion) error
	GetAllByMaterialID(ctx context.Context, materialID string) ([]models.MaterialVersion, error)
	GetByID(ctx context.Context, materialID, versionID string) (*models.MaterialVersion, error)
}

type MySQLMaterialVersionsRepository struct {
	db *sql.DB
}

func NewMySQLMaterialVersionsRepository(db *sql.DB) *MySQLMaterialVersionsRepository {
	return &MySQLMaterialVersionsRepository{db: db}
}
