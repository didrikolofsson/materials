// Package seed provides seeding functionality for the database.
package seed

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/didrikolofsson/materials/internal/models"
)

func seedTeachers(db *sql.DB) error {
	log.Println("Seeding teachers...")
	teacherNames := []string{
		"John Doe",
		"Jane Doe",
		"John Smith",
		"Jane Smith",
	}
	for _, name := range teacherNames {
		if _, err := db.Exec("INSERT INTO teachers (name) VALUES (?)", name); err != nil {
			return fmt.Errorf("failed to insert teacher: %w", err)
		}
	}
	return nil
}

func seedSubjects(db *sql.DB) error {
	log.Println("Seeding subjects...")
	subjectNames := []string{
		"Mathematics",
		"Science",
		"History",
		"English",
		"Art",
		"Music",
	}
	for _, name := range subjectNames {
		if _, err := db.Exec("INSERT INTO subjects (name) VALUES (?)", name); err != nil {
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
			// Insert material first, let AUTO_INCREMENT generate the ID
			result, err := db.Exec(
				"INSERT INTO materials (teacher_id, subject_id) VALUES (?, ?)",
				teacher.ID, subject.ID,
			)
			if err != nil {
				return fmt.Errorf("failed to insert material: %w", err)
			}

			// Get the generated ID
			materialID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get material ID: %w", err)
			}

			// Update original_material_id to reference itself
			if _, err := db.Exec(
				"UPDATE materials SET original_material_id = ? WHERE id = ?",
				materialID, materialID,
			); err != nil {
				return fmt.Errorf("failed to update original_material_id: %w", err)
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
		if err := materials.Scan(&material.ID, &material.TeacherName, &material.SubjectName, &material.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan material: %w", err)
		}
		nVersions := 3
		for versionNumber := range nVersions {
			versionNumber++ // range starts at 0, so increment to start at 1
			result, err := db.Exec(
				"INSERT INTO material_versions (material_id, title, summary, description, content, version_number, is_main) VALUES (?, ?, ?, ?, ?, ?, ?)",
				material.ID, "Title", "Summary", "Description", "Content", versionNumber, versionNumber == 1,
			)
			if err != nil {
				return fmt.Errorf("failed to insert material version: %w", err)
			}
			if versionNumber == 1 {
				versionID, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("failed to get material version ID: %w", err)
				}
				if _, err := db.Exec(
					"UPDATE materials SET current_version_id = ? WHERE id = ?",
					versionID, material.ID,
				); err != nil {
					return fmt.Errorf("failed to update material current version: %w", err)
				}
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
