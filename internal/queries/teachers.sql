-- name: ListTeachers :many
SELECT id,
	name,
	created_at
FROM teachers;
-- name: GetTeacherByID :one
SELECT id,
	name,
	created_at
FROM teachers
WHERE id = ?;