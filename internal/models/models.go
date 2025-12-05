// Package models provides domain models shared across layers.
package models

import "time"

type TeacherID int64
type SubjectID int64
type MaterialID int64
type MaterialVersionID int64

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

type CreateMaterialParams struct {
	TeacherID TeacherID `json:"teacher_id" validate:"required,min=1"`
	SubjectID SubjectID `json:"subject_id" validate:"required,min=1"`
}

type CreateMaterialBody struct {
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Summary     *string `json:"summary" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Content     string  `json:"content" validate:"required,min=1,max=10000"`
}

type CreateMaterialRequest struct {
	Params CreateMaterialParams `json:"params" validate:"required"`
	Body   CreateMaterialBody   `json:"body" validate:"required"`
}

type CreateMaterialResponse struct {
	MaterialID        MaterialID        `json:"material_id"`
	MaterialVersionID MaterialVersionID `json:"material_version_id"`
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
