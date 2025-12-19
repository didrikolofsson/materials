// Package services provides services for the application.
package services

import (
	"context"
	"database/sql"

	"github.com/didrikolofsson/materials/generated/queries"
	customerrors "github.com/didrikolofsson/materials/internal/errors"
	"github.com/didrikolofsson/materials/internal/models"
	"github.com/didrikolofsson/materials/internal/repositories"
)

type Services struct {
	db    *sql.DB
	repos *repositories.MySQLRepository
}

func New(db *sql.DB, repos *repositories.MySQLRepository) *Services {
	return &Services{
		db:    db,
		repos: repos,
	}
}

func (s *Services) ListTeachers(ctx context.Context) ([]models.Teacher, error) {
	res, err := s.repos.ListTeachers(ctx)
	if err != nil {
		return nil, err
	}
	teachers := make([]models.Teacher, len(res))
	for i, teacher := range res {
		teachers[i] = models.Teacher{
			ID:        teacher.ID,
			Name:      teacher.Name,
			CreatedAt: teacher.CreatedAt,
		}
	}
	return teachers, nil
}

func (s *Services) ListSubjects(ctx context.Context) ([]models.Subject, error) {
	res, err := s.repos.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}
	subjects := make([]models.Subject, len(res))
	for i, subject := range res {
		subjects[i] = models.Subject{
			ID:        subject.ID,
			Name:      subject.Name,
			CreatedAt: subject.CreatedAt,
		}
	}
	return subjects, nil
}

func (s *Services) ListMaterials(ctx context.Context) ([]models.Material, error) {
	res, err := s.repos.ListMaterials(ctx)
	if err != nil {
		return nil, err
	}
	materials := make([]models.Material, len(res))
	for i, material := range res {
		materials[i] = models.Material{
			ID:          material.ID,
			TeacherName: material.TeacherName,
			SubjectName: material.SubjectName,
			CreatedAt:   material.CreatedAt,
			Title:       material.Title,
			Description: nullStringToPointer(material.Description),
			Summary:     nullStringToPointer(material.Summary),
		}
	}
	return materials, nil
}

func (s *Services) ListMaterialVersionsByMaterialID(ctx context.Context, materialID int64) ([]models.MaterialVersion, error) {
	res, err := s.repos.ListMaterialVersionsByMaterialID(ctx, materialID)
	if err != nil {
		return nil, err
	}
	materialVersions := make([]models.MaterialVersion, len(res))
	for i, materialVersion := range res {
		materialVersions[i] = models.MaterialVersion{
			ID:            materialVersion.ID,
			Title:         materialVersion.Title,
			Description:   nullStringToPointer(materialVersion.Description),
			Summary:       nullStringToPointer(materialVersion.Summary),
			Content:       materialVersion.Content,
			VersionNumber: int(materialVersion.VersionNumber),
			IsMain:        materialVersion.IsMain,
			CreatedAt:     materialVersion.CreatedAt,
		}
	}
	return materialVersions, nil
}

func (s *Services) GetTeacherByID(ctx context.Context, id int64) (models.Teacher, error) {
	teacher, err := s.repos.GetTeacherByID(ctx, id)
	if err != nil {
		return models.Teacher{}, err
	}
	return models.Teacher{
		ID:        teacher.ID,
		Name:      teacher.Name,
		CreatedAt: teacher.CreatedAt,
	}, nil
}

func (s *Services) GetTeacherMaterials(ctx context.Context, teacherID int64) ([]models.Material, error) {
	res, err := s.repos.GetTeacherMaterials(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	materials := make([]models.Material, len(res))
	for i, material := range res {
		materials[i] = models.Material{
			ID:          material.ID,
			TeacherName: material.TeacherName,
			SubjectName: material.SubjectName,
			CreatedAt:   material.CreatedAt,
			Title:       material.Title,
			Description: nullStringToPointer(material.Description),
			Summary:     nullStringToPointer(material.Summary),
		}
	}
	return materials, nil
}

func (s *Services) GetTeacherMaterialByID(ctx context.Context, teacherID, materialID int64) (models.Material, error) {
	material, err := s.repos.GetTeacherMaterialByID(ctx, teacherID, materialID)
	if err != nil {
		return models.Material{}, err
	}
	return models.Material{
		ID:          material.ID,
		TeacherName: material.TeacherName,
		SubjectName: material.SubjectName,
		CreatedAt:   material.CreatedAt,
		Title:       material.Title,
		Description: nullStringToPointer(material.Description),
		Summary:     nullStringToPointer(material.Summary),
	}, nil
}

func (s *Services) CreateInitialTeacherMaterial(ctx context.Context, teacherID int64, req models.CreateMaterialRequest) (int64, error) {
	// Create the material
	subjectID := int64PtrToValue(req.SubjectID)

	materialID, err := s.repos.CreateInitialTeacherMaterial(
		ctx,
		teacherID,
		subjectID,
		req,
	)
	if err != nil {
		return 0, err
	}

	return materialID, nil
}

func (s *Services) UpdateTeacherMaterialByID(ctx context.Context, teacherID, materialID int64, req models.UpdateMaterialRequest) (models.Material, error) {
	_, err := s.repos.GetTeacherMaterialByID(ctx, teacherID, materialID)
	if err != nil {
		return models.Material{}, err
	}

	// Get the current version to use as defaults
	versions, err := s.repos.ListMaterialVersionsByMaterialID(ctx, materialID)
	if err != nil {
		return models.Material{}, err
	}
	if len(versions) == 0 {
		return models.Material{}, customerrors.ErrBadRequest
	}

	// Find the main version
	var mainVersion *queries.ListMaterialVersionsByMaterialIDRow
	for i := range versions {
		if versions[i].IsMain {
			mainVersion = &versions[i]
			break
		}
	}
	if mainVersion == nil {
		mainVersion = &versions[0] // Use first version if no main version found
	}

	// Use provided values or defaults from main version
	title := mainVersion.Title
	if req.Title != nil {
		title = *req.Title
	}
	summary := nullStringToPointer(mainVersion.Summary)
	if req.Summary != nil {
		summary = req.Summary
	}
	description := nullStringToPointer(mainVersion.Description)
	if req.Description != nil {
		description = req.Description
	}
	content := mainVersion.Content
	if req.Content != nil {
		content = *req.Content
	}

	// Create a new version with updated content
	versionID, err := s.repos.CreateMaterialVersion(
		ctx, materialID, title, summary, description, content, true,
	)
	if err != nil {
		return models.Material{}, err
	}

	// Update all versions to set is_main = false, then set the new one to true
	err = s.repos.UpdateMaterialVersionMain(ctx, materialID, versionID)
	if err != nil {
		return models.Material{}, err
	}

	// Update the material's current_version_id
	err = s.repos.UpdateMaterialCurrentVersion(ctx, materialID, versionID)
	if err != nil {
		return models.Material{}, err
	}

	// Return the updated material
	return s.GetTeacherMaterialByID(ctx, teacherID, materialID)
}

func (s *Services) DeleteTeacherMaterialByID(ctx context.Context, teacherID, materialID int64) error {
	// Verify the material belongs to the teacher
	_, err := s.repos.GetTeacherMaterialByID(ctx, teacherID, materialID)
	if err != nil {
		return err
	}

	// Delete the material (cascade should handle versions)
	err = s.repos.DeleteMaterial(ctx, materialID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Services) UpdateMaterialVersionMain(ctx context.Context, materialID, versionID int64) error {
	err := s.repos.UpdateMaterialVersionMain(ctx, materialID, versionID)
	if err != nil {
		return err
	}

	// Also update the material's current_version_id
	err = s.repos.UpdateMaterialCurrentVersion(ctx, materialID, versionID)
	if err != nil {
		return err
	}
	return nil
}
