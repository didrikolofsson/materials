// Package school provides models for the school domain.
package school

import "time"

type TeacherID string
type SubjectID string
type MaterialID string
type MaterialVersionID string

type Teacher struct {
	ID        TeacherID
	Name      string
	CreatedAt time.Time
}

type Subject struct {
	ID        SubjectID
	Name      string
	CreatedAt time.Time
}

type Material struct {
	ID                 MaterialID
	TeacherID          TeacherID
	SubjectID          *SubjectID
	OriginalMaterialID MaterialID
	CurrentVersionID   *MaterialVersionID
	CreatedAt          *time.Time
}

type MaterialVersion struct {
	ID            MaterialVersionID
	MaterialID    MaterialID
	Title         string
	Summary       *string
	Description   *string
	VersionNumber int
	Content       string
	IsMain        bool
	CreatedAt     time.Time
}
