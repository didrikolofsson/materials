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
	Title              string
	Description        string
	OwnerTeacherID     TeacherID
	SubjectID          SubjectID
	CurrentVersionID   MaterialVersionID
	OriginalMaterialID MaterialID
	CreatedAt          time.Time
}

type MaterialVersion struct {
	ID              MaterialVersionID
	MaterialID      MaterialID
	TeacherID       TeacherID
	ParentVersionID MaterialVersionID
	VersionNumber   int
	Summary         string
	Content         string
	IsMain          bool
	CreatedAt       time.Time
}
