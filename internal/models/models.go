// Package models provides the models for the application.
package models

import (
	"time"
)

type Teacher struct {
	ID        int64     `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1,max=255"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type Subject struct {
	ID        int64     `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1,max=255"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type Material struct {
	ID          int64     `json:"id" validate:"required"`
	TeacherName string    `json:"teacher_name" validate:"required,min=1,max=255"`
	SubjectName string    `json:"subject_name" validate:"required,min=1,max=255"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description *string   `json:"description" validate:"omitempty,min=1,max=1000"`
	Summary     *string   `json:"summary" validate:"omitempty,min=1,max=255"`
}

type MaterialVersion struct {
	ID            int64     `json:"id" validate:"required"`
	Title         string    `json:"title" validate:"required,min=1,max=255"`
	Description   *string   `json:"description" validate:"omitempty,min=1,max=1000"`
	Summary       *string   `json:"summary" validate:"omitempty,min=1,max=255"`
	Content       string    `json:"content" validate:"required,min=1"`
	VersionNumber int       `json:"version_number" validate:"required,min=1"`
	IsMain        bool      `json:"is_main" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
}

type CreateMaterialRequest struct {
	SubjectID   *int64  `json:"subject_id" validate:"omitempty,min=1"`
	Title       string  `json:"title" validate:"required,min=1,max=255"`
	Summary     *string `json:"summary" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Content     string  `json:"content" validate:"required,min=1"`
}

type UpdateMaterialRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=255"`
	Summary     *string `json:"summary" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	Content     *string `json:"content" validate:"omitempty,min=1"`
}
