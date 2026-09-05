-- name: CreateNoteVersion :one
INSERT INTO note_versions (note_id, version_number, title, content, tags, created_at)
VALUES ($1, $2, $3, $4, $5, now())
RETURNING *;

-- name: CountNoteVersionsByNote :one
SELECT count(*) FROM note_versions
WHERE note_id = $1;

-- name: ListNoteVersionsByNote :many
SELECT * FROM note_versions
WHERE note_id = $1
ORDER BY version_number DESC
LIMIT $2 OFFSET $3;

-- name: GetNoteVersionByID :one
SELECT * FROM note_versions
WHERE id = $1;

-- name: GetLatestNoteVersionNumber :one
SELECT COALESCE(MAX(version_number), 0)::bigint FROM note_versions
WHERE note_id = $1;
