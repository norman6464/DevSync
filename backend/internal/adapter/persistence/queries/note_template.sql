-- 新しい行がis_default=trueの場合のみ、同一ユーザーの既存デフォルトを同一文で外す
-- （$7=trueでなければ WHERE が一致せずclearedは0行、is_default IS TRUEでNULLを安全に除外）。
-- uq_note_templates_default（(user_id) WHERE is_default の部分UNIQUE索引）が
-- 最終的な「ユーザーごとデフォルトは高々1件」を保証する安全網になる。
-- name: CreateNoteTemplate :one
WITH cleared AS (
    UPDATE note_templates SET is_default = false
    WHERE note_templates.user_id = $1 AND note_templates.is_default IS TRUE AND $7 = true
)
INSERT INTO note_templates (user_id, name, description, default_title, content_template, default_tags, is_default, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateNoteTemplate :one
WITH cleared AS (
    UPDATE note_templates SET is_default = false
    WHERE note_templates.user_id = $8 AND note_templates.is_default IS TRUE AND note_templates.id != $1 AND $7 = true
)
UPDATE note_templates
SET name = $2,
    description = $3,
    default_title = $4,
    content_template = $5,
    default_tags = $6,
    is_default = $7,
    updated_at = now()
WHERE note_templates.id = $1
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

-- name: CountNoteTemplatesByUser :one
SELECT count(*) FROM note_templates
WHERE user_id = $1;
