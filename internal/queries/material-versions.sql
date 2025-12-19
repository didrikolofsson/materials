-- name: ListAllMaterialVersions :many
SELECT id,
	title,
	summary,
	description,
	content,
	version_number,
	is_main,
	created_at
FROM material_versions;
-- name: ListMaterialVersionsByMaterialID :many
SELECT id,
	title,
	summary,
	description,
	content,
	version_number,
	is_main,
	created_at
FROM material_versions
WHERE material_id = ?
ORDER BY version_number DESC;
-- name: CreateMaterialVersion :execresult
INSERT INTO material_versions (
		material_id,
		title,
		summary,
		description,
		content,
		version_number,
		is_main
	)
VALUES (?, ?, ?, ?, ?, ?, ?);
-- name: UpdateMaterialVersionMain :exec
UPDATE material_versions
SET is_main = CASE
		WHEN id = ? THEN TRUE
		ELSE FALSE
	END
WHERE material_id = ?;
-- name: GetMaxVersionNumberByMaterialID :one
SELECT version_number
FROM material_versions
WHERE material_id = ?
ORDER BY version_number DESC
LIMIT 1;