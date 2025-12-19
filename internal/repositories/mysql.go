// Package repositories provides a MySQL repository.
package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/didrikolofsson/materials/generated/queries"
	customerrors "github.com/didrikolofsson/materials/internal/errors"
	"github.com/didrikolofsson/materials/internal/models"
)

type MySQLRepository struct {
	q  *queries.Queries
	db *sql.DB
}

func New(db *sql.DB) *MySQLRepository {
	q := queries.New(db)
	return &MySQLRepository{
		q:  q,
		db: db,
	}
}

func (r *MySQLRepository) ListTeachers(ctx context.Context) ([]queries.Teacher, error) {
	teachers, err := r.q.ListTeachers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternal
	}
	return teachers, nil
}

func (r *MySQLRepository) ListSubjects(ctx context.Context) ([]queries.Subject, error) {
	subjects, err := r.q.ListSubjects(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternal
	}
	return subjects, nil
}

func (r *MySQLRepository) ListAllMaterials(ctx context.Context) ([]queries.ListAllMaterialsRow, error) {
	materials, err := r.q.ListAllMaterials(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternal
	}
	return materials, nil
}

func (r *MySQLRepository) ListMaterialVersionsByMaterialID(ctx context.Context, materialID int64) ([]queries.ListMaterialVersionsByMaterialIDRow, error) {
	materialVersions, err := r.q.ListMaterialVersionsByMaterialID(ctx, materialID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternal
	}
	return materialVersions, nil
}

func (r *MySQLRepository) GetTeacherByID(ctx context.Context, id int64) (queries.Teacher, error) {
	teacher, err := r.q.GetTeacherByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queries.Teacher{}, customerrors.ErrNotFound
		}
		return queries.Teacher{}, customerrors.ErrInternal
	}
	return teacher, nil
}

func (r *MySQLRepository) GetTeacherMaterials(ctx context.Context, teacherID int64) ([]queries.GetTeacherMaterialsRow, error) {
	materials, err := r.q.GetTeacherMaterials(ctx, teacherID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerrors.ErrNotFound
		}
		return nil, customerrors.ErrInternal
	}
	return materials, nil
}

func (r *MySQLRepository) GetTeacherMaterialByID(ctx context.Context, teacherID, materialID int64) (queries.GetTeacherMaterialByIDRow, error) {
	material, err := r.q.GetTeacherMaterialByID(ctx, queries.GetTeacherMaterialByIDParams{
		TeacherID: teacherID,
		ID:        materialID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return queries.GetTeacherMaterialByIDRow{}, customerrors.ErrNotFound
		}
		return queries.GetTeacherMaterialByIDRow{}, customerrors.ErrInternal
	}
	return material, nil
}

func (r *MySQLRepository) CreateInitialTeacherMaterial(
	ctx context.Context,
	teacherID int64,
	subjectID int64,
	req models.CreateMaterialRequest,
) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Create material
	res, err := qtx.CreateMaterial(ctx, queries.CreateMaterialParams{
		TeacherID:          teacherID,
		SubjectID:          toNullInt64(&subjectID),
		OriginalMaterialID: toNullInt64(nil),
	})
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	materialID, err := res.LastInsertId()
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	// Create version
	res, err = qtx.CreateMaterialVersion(ctx, queries.CreateMaterialVersionParams{
		MaterialID:  materialID,
		Title:       req.Title,
		Summary:     toNullString(req.Summary),
		Description: toNullString(req.Description),
		Content:     req.Content,
		IsMain:      true,
	})
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	versionID, err := res.LastInsertId()
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	// Update material current version
	err = qtx.UpdateMaterialCurrentVersion(ctx, queries.UpdateMaterialCurrentVersionParams{
		CurrentVersionID: toNullInt64(&versionID),
		ID:               materialID,
	})
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	err = tx.Commit()
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	return materialID, nil
}

func (r *MySQLRepository) GetMaxVersionNumber(ctx context.Context, materialID int64) (int32, error) {
	maxVersion, err := r.q.GetMaxVersionNumber(ctx, materialID)
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	return int32(maxVersion), nil
}

func (r *MySQLRepository) CreateMaterialVersion(ctx context.Context, materialID int64, title string, summary, description *string, content string, isMain bool) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)
	maxVersion, err := qtx.GetMaxVersionNumber(ctx, materialID)
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	result, err := qtx.CreateMaterialVersion(ctx, queries.CreateMaterialVersionParams{
		MaterialID:    materialID,
		Title:         title,
		Summary:       toNullString(summary),
		Description:   toNullString(description),
		Content:       content,
		IsMain:        isMain,
		VersionNumber: int32(maxVersion + 1),
	})
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	err = tx.Commit()
	if err != nil {
		return 0, customerrors.ErrInternal
	}

	versionID, err := result.LastInsertId()
	if err != nil {
		return 0, customerrors.ErrInternal
	}
	return versionID, nil
}

func (r *MySQLRepository) UpdateMaterialCurrentVersion(ctx context.Context, materialID, versionID int64) error {
	err := r.q.UpdateMaterialCurrentVersion(ctx, queries.UpdateMaterialCurrentVersionParams{
		CurrentVersionID: toNullInt64(&versionID),
		ID:               materialID,
	})
	if err != nil {
		return customerrors.ErrInternal
	}
	return nil
}

func (r *MySQLRepository) UpdateMaterialVersionMain(ctx context.Context, materialID, versionID int64) error {
	err := r.q.UpdateMaterialVersionMain(ctx, queries.UpdateMaterialVersionMainParams{
		ID:         versionID,
		MaterialID: materialID,
	})
	if err != nil {
		return customerrors.ErrInternal
	}
	return nil
}

func (r *MySQLRepository) DeleteMaterial(ctx context.Context, id int64) error {
	err := r.q.DeleteMaterial(ctx, id)
	if err != nil {
		return customerrors.ErrInternal
	}
	return nil
}
