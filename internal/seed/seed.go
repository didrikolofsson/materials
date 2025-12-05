// Package seed provides seeding functionality for the database.
package seed

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/didrikolofsson/materials/internal/models"
	"github.com/google/uuid"
)

func seedTeachers(db *sql.DB) error {
	log.Println("Seeding teachers...")
	teachers := []string{
		"John Doe",
		"Jane Doe",
		"John Smith",
		"Jane Smith",
	}
	for _, teacher := range teachers {
		id := uuid.NewString()
		if _, err := db.Exec("INSERT INTO teachers (id, name) VALUES (?, ?)", id, teacher); err != nil {
			return fmt.Errorf("failed to insert teacher: %w", err)
		}
	}
	return nil
}

func seedSubjects(db *sql.DB) error {
	log.Println("Seeding subjects...")
	subjects := []string{
		"Mathematics",
		"Science",
		"History",
		"English",
		"Art",
		"Music",
	}
	for _, subject := range subjects {
		id := uuid.NewString()
		if _, err := db.Exec("INSERT INTO subjects (id, name) VALUES (?, ?)", id, subject); err != nil {
			return fmt.Errorf("failed to insert subject: %w", err)
		}
	}
	return nil
}

func seedMaterials(db *sql.DB) error {
	log.Println("Seeding materials...")
	// Get all teachers and subjects
	teachers, err := db.Query("SELECT id, name FROM teachers")
	if err != nil {
		return fmt.Errorf("failed to query teachers: %w", err)
	}
	defer teachers.Close()

	var subjects []models.Subject
	rows, err := db.Query("SELECT id, name FROM subjects")
	if err != nil {
		return fmt.Errorf("failed to query subjects: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return fmt.Errorf("failed to scan subject: %w", err)
		}
		subjects = append(subjects, subject)
	}

	for teachers.Next() {
		var teacher models.Teacher
		if err := teachers.Scan(&teacher.ID, &teacher.Name); err != nil {
			return fmt.Errorf("failed to scan teacher: %w", err)
		}
		shuffled := make([]models.Subject, len(subjects))
		copy(shuffled, subjects)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		for _, subject := range shuffled[:3] {
			newID := uuid.NewString()
			if _, err := db.Exec("INSERT INTO materials (id, teacher_id, subject_id, original_material_id) VALUES (?, ?, ?, ?)", newID, teacher.ID, subject.ID, newID); err != nil {
				return fmt.Errorf("failed to insert material: %w", err)
			}
		}
	}
	return nil
}

func seedMaterialVersions(db *sql.DB) error {
	log.Println("Seeding material versions...")
	// Get all materials
	materials, err := db.Query("SELECT id, teacher_id, subject_id, original_material_id FROM materials")
	if err != nil {
		return fmt.Errorf("failed to query materials: %w", err)
	}
	defer materials.Close()

	for materials.Next() {
		var material models.Material
		if err := materials.Scan(&material.ID, &material.TeacherID, &material.SubjectID, &material.OriginalMaterialID); err != nil {
			return fmt.Errorf("failed to scan material: %w", err)
		}
		nVersions := 3
		for versionNumber := range nVersions {
			versionNumber++ // range starts at 0, so increment to start at 1
			newID := uuid.NewString()
			if _, err := db.Exec("INSERT INTO material_versions (id, material_id, title, summary, description, content, version_number, is_main) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", newID, material.ID, "Title", "Summary", "Description", "Content", versionNumber, versionNumber == 1); err != nil {
				return fmt.Errorf("failed to insert material version: %w", err)
			}
		}
	}
	return nil
}

func Run(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM teachers").Scan(&count); err != nil {
		return fmt.Errorf("failed to count teachers: %w", err)
	}
	if count == 0 {
		if err := seedTeachers(db); err != nil {
			return fmt.Errorf("failed to seed teachers: %w", err)
		}
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&count); err != nil {
		return fmt.Errorf("failed to count subjects: %w", err)
	}
	if count == 0 {
		if err := seedSubjects(db); err != nil {
			return fmt.Errorf("failed to seed subjects: %w", err)
		}
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM materials").Scan(&count); err != nil {
		return fmt.Errorf("failed to count materials: %w", err)
	}
	if count == 0 {
		if err := seedMaterials(db); err != nil {
			return fmt.Errorf("failed to seed materials: %w", err)
		}
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM material_versions").Scan(&count); err != nil {
		return fmt.Errorf("failed to count material versions: %w", err)
	}
	if count == 0 {
		if err := seedMaterialVersions(db); err != nil {
			return fmt.Errorf("failed to seed material versions: %w", err)
		}
	}
	return nil
}
