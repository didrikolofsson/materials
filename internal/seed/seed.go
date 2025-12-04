// Package seed provides seeding functionality for the database.
package seed

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

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

func Run(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&count); err != nil {
		return fmt.Errorf("failed to count subjects: %w", err)
	}
	if count == 0 {
		if err := seedSubjects(db); err != nil {
			return fmt.Errorf("failed to seed subjects: %w", err)
		}
	}
	return nil
}
