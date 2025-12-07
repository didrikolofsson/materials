// Package school provides service layer for business logic.
package school

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/didrikolofsson/materials/internal/models"
)

type ServiceDomainSchool struct {
	r  RepositoryDomainSchool
	db *sql.DB
}

func NewServiceDomainSchool(db *sql.DB) *ServiceDomainSchool {
	return &ServiceDomainSchool{r: NewRepositoryDomainSchool(), db: db}
}

func (s *ServiceDomainSchool) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.r.Subjects.List(ctx, s.db)
}

func (s *ServiceDomainSchool) ListTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.r.Teachers.List(ctx, s.db)
}

func (s *ServiceDomainSchool) ListAllMaterialsByTeacherID(ctx context.Context, teacherID models.GenericID) ([]models.Material, error) {
	return s.r.Materials.ListAllByTeacherID(ctx, s.db, teacherID)
}

func (s *ServiceDomainSchool) CreateMaterialWithInitialVersion(
	ctx context.Context, m *models.CreateMaterialRequest,
) (models.GenericID, models.GenericID, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	materialID, err := s.r.Materials.Create(ctx, tx, &models.Material{
		TeacherID: m.Params.TeacherID,
		SubjectID: &m.Params.SubjectID,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create material: %w", err)
	}

	_, err = s.r.Materials.UpdateOriginalMaterialID(ctx, tx, materialID, materialID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update original material id: %w", err)
	}

	materialVersion := &models.MaterialVersion{
		MaterialID:    materialID,
		Title:         m.Body.Title,
		Summary:       m.Body.Summary,
		Description:   m.Body.Description,
		Content:       m.Body.Content,
		VersionNumber: 1,
		IsMain:        true,
	}

	materialVersionID, err := s.r.MaterialVersions.Create(ctx, tx, materialVersion)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create material version: %w", err)
	}

	materialVersionID, err = s.r.Materials.UpdateCurrentVersion(ctx, tx, materialID, materialVersionID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update current version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return materialID, materialVersionID, nil
}

// ListMaterialVersionsByMaterialID lists all material versions for a given material ID
func (s *ServiceDomainSchool) ListMaterialVersionsByMaterialID(ctx context.Context, materialID models.GenericID) ([]models.MaterialVersion, error) {
	materialVersions, err := s.r.MaterialVersions.ListAllByMaterialID(ctx, s.db, materialID)
	if err != nil {
		return nil, fmt.Errorf("failed to list material versions by material id: %w", err)
	}
	return materialVersions, nil
}

func (s *ServiceDomainSchool) UpdateCurrentVersion(ctx context.Context, m, v models.GenericID) (models.GenericID, error) {
	materialVersionID, err := s.r.Materials.UpdateCurrentVersion(ctx, s.db, m, v)
	if err != nil {
		return 0, fmt.Errorf("failed to update current version: %w", err)
	}
	return materialVersionID, nil
}

func (s *ServiceDomainSchool) GetMaterialByID(ctx context.Context, id models.GenericID) (*models.Material, error) {
	return s.r.Materials.GetByID(ctx, s.db, id)
}

func (s *ServiceDomainSchool) UpdateMainVersionForMaterialByVersionID(ctx context.Context, r *models.UpdateMainVersionForMaterialRequest) error {
	err := s.r.MaterialVersions.SetMainForMaterialVersion(ctx, s.db, r.Params.MaterialID, r.Body.MaterialVersionID)
	if err != nil {
		return fmt.Errorf("failed to set main version for material by version id: %w", err)
	}
	return nil
}
