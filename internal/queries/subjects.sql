-- name: ListSubjects :many
SELECT id,
	name,
	created_at
FROM subjects;
-- name: SeedSubjects :exec
INSERT INTO subjects (name)
VALUES (?);