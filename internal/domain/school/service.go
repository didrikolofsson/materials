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

func (s *ServiceDomainSchool) CreateMaterialWithInitialVersion(
	ctx context.Context, m *models.CreateMaterialRequest,
) (models.MaterialID, models.MaterialVersionID, error) {
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

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return materialID, materialVersionID, nil
}

func (s *ServiceDomainSchool) UpdateCurrentVersion(ctx context.Context, m, v int64) error {
	err := s.r.Materials.UpdateCurrentVersion(ctx, s.db, m, v)
	return err
}

func (s *ServiceDomainSchool) GetMaterialByID(ctx context.Context, id int64) (*models.Material, error) {
	m, err := s.r.Materials.GetByID(ctx, s.db, id)
	return m, err
}
