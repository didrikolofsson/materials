// Package school provides route registration for the school domain.
package school

import (
	"github.com/didrikolofsson/materials/internal/repositories"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all routes for the school domain
// Repository should be created outside (in server.go) and passed in
func RegisterRoutes(r chi.Router, repo repositories.SubjectsRepository) {
	service := NewSubjectService(repo)
	handler := NewSchoolHandler(service)

	r.Get("/subjects", handler.handleListSubjects)
}
