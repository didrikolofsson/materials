// Package models provides domain models shared across layers.
package models

import "time"

type TeacherID string
type SubjectID string
type MaterialID string
type MaterialVersionID string

type Teacher struct {
	ID        TeacherID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Subject struct {
	ID        SubjectID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Material struct {
	ID                 MaterialID         `json:"id"`
	TeacherID          TeacherID          `json:"teacher_id"`
	SubjectID          *SubjectID         `json:"subject_id"`
	OriginalMaterialID MaterialID         `json:"original_material_id"`
	CurrentVersionID   *MaterialVersionID `json:"current_version_id"`
	CreatedAt          *time.Time         `json:"created_at"`
}

type MaterialVersion struct {
	ID            MaterialVersionID `json:"id"`
	MaterialID    MaterialID        `json:"material_id"`
	Title         string            `json:"title"`
	Summary       *string           `json:"summary"`
	Description   *string           `json:"description"`
	VersionNumber int               `json:"version_number"`
	Content       string            `json:"content"`
	IsMain        bool              `json:"is_main"`
	CreatedAt     time.Time         `json:"created_at"`
}
