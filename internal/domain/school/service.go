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
