// Package school provides service layer for business logic.
package school

import (
	"context"

	"github.com/didrikolofsson/materials/internal/models"
)

type ServiceDomainSchool struct {
	r RepositoryDomainSchool
}

func NewServiceDomainSchool(r RepositoryDomainSchool) *ServiceDomainSchool {
	return &ServiceDomainSchool{r: r}
}

func (s *ServiceDomainSchool) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.r.Subjects.List(ctx)
}

func (s *ServiceDomainSchool) ListTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.r.Teachers.List(ctx)
}

func (s *ServiceDomainSchool) CreateMaterial(ctx context.Context, m *models.Material) error {
	_, err := s.r.Materials.Create(ctx, m)
	return err
}

func (s *ServiceDomainSchool) UpdateCurrentVersion(ctx context.Context, m, v int64) error {
	err := s.r.Materials.UpdateCurrentVersion(ctx, m, v)
	return err
}

func (s *ServiceDomainSchool) GetMaterialByID(ctx context.Context, id int64) (*models.Material, error) {
	m, err := s.r.Materials.GetByID(ctx, id)
	return m, err
}
