package school

import (
	"database/sql"

	"github.com/didrikolofsson/materials/internal/repositories"
)

type RepositoryDomainSchool struct {
	Subjects repositories.SubjectsRepository
	Teachers repositories.TeachersRepository
}

func NewRepositoryDomainSchool(db *sql.DB) RepositoryDomainSchool {
	return RepositoryDomainSchool{
		Subjects: repositories.NewMySQLSubjectsRepository(db),
		Teachers: repositories.NewMySQLTeachersRepository(db),
	}
}
