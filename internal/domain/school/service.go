// Package school provides service layer for business logic.
package school

import (
	"context"
)

type SubjectService struct {
	r SubjectsRepository
}

func NewSubjectService(r SubjectsRepository) *SubjectService {
	return &SubjectService{r: r}
}

func (s *SubjectService) ListSubjects(ctx context.Context) ([]Subject, error) {
	return s.r.List(ctx, nil)
}
