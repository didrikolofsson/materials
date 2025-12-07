// Package models provides domain models shared across layers.
package models

import "time"

type GenericID int64

type Teacher struct {
	ID        GenericID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Subject struct {
	ID        GenericID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMaterialParams struct {
	TeacherID GenericID `json:"teacher_id" validate:"required,min=1"`
	SubjectID GenericID `json:"subject_id" validate:"required,min=1"`
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
	MaterialID        GenericID `json:"material_id"`
	MaterialVersionID GenericID `json:"material_version_id"`
}

type UpdateMainVersionForMaterialParams struct {
	MaterialID GenericID `json:"material_id" validate:"required,min=1"`
}
type UpdateMainVersionForMaterialBody struct {
	MaterialVersionID GenericID `json:"material_version_id" validate:"required,min=1"`
}

type UpdateMainVersionForMaterialRequest struct {
	Params UpdateMainVersionForMaterialParams `json:"params" validate:"required"`
	Body   UpdateMainVersionForMaterialBody   `json:"body" validate:"required"`
}

type UpdateMainVersionForMaterialResponse struct {
	MaterialID        GenericID `json:"material_id"`
	MaterialVersionID GenericID `json:"material_version_id"`
}

type Material struct {
	ID                 GenericID  `json:"id"`
	TeacherID          GenericID  `json:"teacher_id"`
	SubjectID          *GenericID `json:"subject_id"`
	OriginalMaterialID GenericID  `json:"original_material_id"`
	CurrentVersionID   *GenericID `json:"current_version_id"`
	CreatedAt          *time.Time `json:"created_at"`
}

type MaterialVersion struct {
	ID            GenericID `json:"id"`
	MaterialID    GenericID `json:"material_id"`
	Title         string    `json:"title"`
	Summary       *string   `json:"summary"`
	Description   *string   `json:"description"`
	VersionNumber int       `json:"version_number"`
	Content       string    `json:"content"`
	IsMain        bool      `json:"is_main"`
	CreatedAt     time.Time `json:"created_at"`
}
