-- name: CreateNoteTemplate :one
INSERT INTO note_templates (user_id, name, description, default_title, content_template, default_tags, is_default, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateNoteTemplate :one
UPDATE note_templates
SET name = $2,
    description = $3,
    default_title = $4,
    content_template = $5,
    default_tags = $6,
    is_default = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteNoteTemplate :exec
DELETE FROM note_templates
WHERE id = $1;

-- name: GetNoteTemplateByID :one
SELECT * FROM note_templates
WHERE id = $1;

-- name: ListNoteTemplatesByUser :many
SELECT * FROM note_templates
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetDefaultNoteTemplateByUser :one
SELECT * FROM note_templates
WHERE user_id = $1 AND is_default = true
ORDER BY id ASC
LIMIT 1;

-- name: ClearNoteTemplateDefaultFlag :exec
UPDATE note_templates
SET is_default = false
WHERE user_id = $1;

-- name: CountNoteTemplatesByUser :one
SELECT count(*) FROM note_templates
WHERE user_id = $1;
