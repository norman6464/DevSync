-- name: CreateNoteFolder :one
INSERT INTO note_folders (user_id, parent_id, name, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetNoteFolderByID :one
SELECT * FROM note_folders
WHERE id = $1;

-- name: ListNoteFoldersByUser :many
SELECT * FROM note_folders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListNoteFoldersByParent :many
SELECT * FROM note_folders
WHERE parent_id = $1
ORDER BY created_at DESC;

-- name: ListRootNoteFoldersByUser :many
SELECT * FROM note_folders
WHERE user_id = $1 AND parent_id IS NULL
ORDER BY created_at DESC;

-- name: UpdateNoteFolder :one
UPDATE note_folders
SET parent_id = $2,
    name = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNoteFolder :exec
DELETE FROM note_folders
WHERE id = $1;
