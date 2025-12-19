-- name: ListAllMaterials :many
SELECT m.id,
	t.name as teacher_name,
	s.name as subject_name,
	m.created_at,
	mv.title,
	mv.description,
	mv.summary
FROM materials m
	INNER JOIN teachers t ON m.teacher_id = t.id
	INNER JOIN subjects s ON m.subject_id = s.id
	INNER JOIN material_versions mv ON m.current_version_id = mv.id;
-- name: GetTeacherMaterials :many
SELECT m.id,
	t.name as teacher_name,
	s.name as subject_name,
	m.created_at,
	mv.title,
	mv.description,
	mv.summary
FROM materials m
	INNER JOIN teachers t ON m.teacher_id = t.id
	INNER JOIN subjects s ON m.subject_id = s.id
	INNER JOIN material_versions mv ON m.current_version_id = mv.id
WHERE m.teacher_id = ?;
-- name: GetTeacherMaterialByID :one
SELECT m.id,
	t.name as teacher_name,
	s.name as subject_name,
	m.created_at,
	mv.title,
	mv.description,
	mv.summary
FROM materials m
	INNER JOIN teachers t ON m.teacher_id = t.id
	INNER JOIN subjects s ON m.subject_id = s.id
	INNER JOIN material_versions mv ON m.current_version_id = mv.id
WHERE m.teacher_id = ?
	AND m.id = ?;
-- name: CreateMaterial :execresult
INSERT INTO materials (
		teacher_id,
		subject_id,
		original_material_id,
		current_version_id
	)
VALUES (?, ?, ?, ?);
-- name: UpdateMaterialCurrentVersion :exec
UPDATE materials
SET current_version_id = ?
WHERE id = ?;
-- name: DeleteMaterial :exec
DELETE FROM materials
WHERE id = ?;
-- name: UpdateMaterialOriginalID :exec
UPDATE materials
SET original_material_id = ?
WHERE id = ?;