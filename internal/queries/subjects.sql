-- name: ListSubjects :many
SELECT id,
	name,
	created_at
FROM subjects;