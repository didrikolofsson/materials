// Package seed provides seeding functionality for the database.
package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/didrikolofsson/materials/generated/queries"
	"github.com/go-faker/faker/v4"
)

func seedTeachers(ctx context.Context, qtx *queries.Queries) error {
	nTeachers := 5

	teachers, err := qtx.ListTeachers(ctx)
	if err != nil {
		return fmt.Errorf("failed to query teachers: %w", err)
	}
	if len(teachers) > 0 {
		log.Println("Teachers already seeded..")
		return nil
	}

	log.Println("Seeding teachers...")
	for range nTeachers {
		d := queries.Teacher{
			Name: faker.Name(),
		}

		if err := qtx.SeedTeachers(ctx, d.Name); err != nil {
			return fmt.Errorf("failed to insert teacher: %w", err)
		}
	}
	return nil
}

func seedSubjects(ctx context.Context, qtx *queries.Queries) error {
	nSubjects := 5

	subjects, err := qtx.ListSubjects(ctx)
	if err != nil {
		return fmt.Errorf("failed to query subjects: %w", err)
	}
	if len(subjects) > 0 {
		log.Println("Subjects already seeded..")
		return nil
	}

	log.Println("Seeding subjects...")
	for range nSubjects {
		d := queries.Subject{
			Name: faker.Word(),
		}
		if err := qtx.SeedSubjects(ctx, d.Name); err != nil {
			return fmt.Errorf("failed to insert subject: %w", err)
		}
	}
	return nil
}

func seedMaterials(ctx context.Context, qtx *queries.Queries) error {
	nMaterials := 3

	materials, err := qtx.ListMaterials(ctx)
	if err != nil {
		return fmt.Errorf("failed to query materials: %w", err)
	}
	if len(materials) > 0 {
		log.Println("Materials already seeded..")
		return nil
	}

	log.Println("Seeding materials...")
	teachers, err := qtx.ListTeachers(ctx)
	if err != nil {
		return fmt.Errorf("failed to query teachers: %w", err)
	}

	subjects, err := qtx.ListSubjects(ctx)
	if err != nil {
		return fmt.Errorf("failed to query subjects: %w", err)
	}

	for _, teacher := range teachers {
		shuffled := make([]queries.Subject, len(subjects))
		copy(shuffled, subjects)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		for _, subject := range shuffled[:nMaterials] {
			result, err := qtx.CreateMaterial(ctx, queries.CreateMaterialParams{
				TeacherID: teacher.ID,
				SubjectID: sql.NullInt64{Int64: subject.ID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert material: %w", err)
			}

			// Get the generated ID
			materialID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get material ID: %w", err)
			}

			// Update original_material_id to reference itself
			if err := qtx.UpdateMaterialOriginalID(ctx, queries.UpdateMaterialOriginalIDParams{
				OriginalMaterialID: sql.NullInt64{Int64: materialID, Valid: true},
				ID:                 materialID,
			}); err != nil {
				return fmt.Errorf("failed to update original_material_id: %w", err)
			}
		}
	}
	return nil
}

func seedMaterialVersions(ctx context.Context, qtx *queries.Queries) error {
	nVersions := 3
	versionNumber := 0

	materialVersions, err := qtx.ListAllMaterialVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to query material versions: %w", err)
	}
	if len(materialVersions) > 0 {
		log.Println("Material versions already seeded..")
		return nil
	}
	log.Println("Seeding material versions...")

	// Get all materials
	materials, err := qtx.ListMaterials(ctx)
	if err != nil {
		return fmt.Errorf("failed to query materials: %w", err)
	}

	for _, material := range materials {
		for range nVersions {
			versionNumber++
			result, err := qtx.CreateMaterialVersion(ctx, queries.CreateMaterialVersionParams{
				MaterialID:    material.ID,
				Title:         faker.Word(),
				Summary:       sql.NullString{String: faker.Sentence(), Valid: true},
				Description:   sql.NullString{String: faker.Sentence(), Valid: true},
				Content:       faker.Paragraph(),
				VersionNumber: int32(versionNumber),
				IsMain:        versionNumber == 1,
			})
			if err != nil {
				return fmt.Errorf("failed to insert material version: %w", err)
			}
			if versionNumber == 1 {
				versionID, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("failed to get material version ID: %w", err)
				}
				if err := qtx.UpdateMaterialCurrentVersion(ctx, queries.UpdateMaterialCurrentVersionParams{
					CurrentVersionID: sql.NullInt64{Int64: versionID, Valid: true},
					ID:               material.ID,
				}); err != nil {
					return fmt.Errorf("failed to update material current version: %w", err)
				}
			}
		}
		versionNumber = 0
	}
	return nil
}

func Run(ctx context.Context, db *sql.DB) error {
	q := queries.New(db)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)

	if err := seedTeachers(ctx, qtx); err != nil {
		return fmt.Errorf("failed to seed teachers: %w", err)
	}
	if err := seedSubjects(ctx, qtx); err != nil {
		return fmt.Errorf("failed to seed subjects: %w", err)
	}
	if err := seedMaterials(ctx, qtx); err != nil {
		return fmt.Errorf("failed to seed materials: %w", err)
	}
	if err := seedMaterialVersions(ctx, qtx); err != nil {
		return fmt.Errorf("failed to seed material versions: %w", err)
	}

	return tx.Commit()
}
