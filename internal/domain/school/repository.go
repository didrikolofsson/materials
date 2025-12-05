package school

import (
	"github.com/didrikolofsson/materials/internal/repositories"
)

type RepositoryDomainSchool struct {
	Subjects         repositories.SubjectsRepository
	Teachers         repositories.TeachersRepository
	Materials        repositories.MaterialsRepository
	MaterialVersions repositories.MaterialVersionsRepository
}

func NewRepositoryDomainSchool() RepositoryDomainSchool {
	return RepositoryDomainSchool{
		Subjects:         repositories.NewMySQLSubjectsRepository(),
		Teachers:         repositories.NewMySQLTeachersRepository(),
		Materials:        repositories.NewMySQLMaterialsRepository(),
		MaterialVersions: repositories.NewMySQLMaterialVersionsRepository(),
	}
}
