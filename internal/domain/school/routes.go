// Package school provides route registration for the school domain.
package school

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(r chi.Router, db *sql.DB, validate *validator.Validate) {
	service := NewServiceDomainSchool(db)
	handler := NewHandlerDomainSchool(service, validate)

	r.Get("/subjects", handler.handleListSubjects)
	r.Get("/teachers", handler.handleListTeachers)

	r.Get("/materials/{id}", handler.handleGetMaterialByID)
	r.Post("/teachers/{teacher_id}/materials/{subject_id}", handler.handleCreateMaterial)

}
