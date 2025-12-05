// Package school provides route registration for the school domain.
package school

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, db *sql.DB) {
	repository := NewRepositoryDomainSchool(db)
	service := NewServiceDomainSchool(repository)
	handler := NewHandlerDomainSchool(service)

	r.Get("/subjects", handler.handleListSubjects)
	r.Get("/teachers", handler.handleListTeachers)
}
