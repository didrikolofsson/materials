// Package school provides route registration for the school domain.
package school

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, db *sql.DB) {
	repo := NewMySQLSubjectsRepository(db)
	service := NewSubjectService(repo)
	handler := NewSchoolHandler(service)

	r.Get("/subjects", handler.handleListSubjects)
}
