// Package school provides service layer for business logic.
package school

import (
	"context"

	"github.com/didrikolofsson/materials/internal/models"
	"github.com/didrikolofsson/materials/internal/repositories"
)

type SubjectService struct {
	r repositories.SubjectsRepository
}

func NewSubjectService(r repositories.SubjectsRepository) *SubjectService {
	return &SubjectService{r: r}
}

func (s *SubjectService) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	return s.r.List(ctx)
}
