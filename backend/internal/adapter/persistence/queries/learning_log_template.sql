-- name: CreateLearningLogTemplate :one
INSERT INTO learning_log_templates (user_id, name, default_title, default_content, default_category, default_duration, is_default, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
RETURNING *;

-- name: UpdateLearningLogTemplate :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE learning_log_templates SET
    name = $2,
    default_title = $3,
    default_content = $4,
    default_category = $5,
    default_duration = $6,
    is_default = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLearningLogTemplate :exec
DELETE FROM learning_log_templates WHERE id = $1;

-- name: GetLearningLogTemplateByID :one
SELECT * FROM learning_log_templates WHERE id = $1;

-- name: ListLearningLogTemplatesByUser :many
SELECT * FROM learning_log_templates WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetDefaultLearningLogTemplateByUser :one
SELECT * FROM learning_log_templates WHERE user_id = $1 AND is_default = true;

-- name: ClearLearningLogTemplateDefaultFlag :exec
UPDATE learning_log_templates SET is_default = false, updated_at = now() WHERE user_id = $1;

-- name: CountLearningLogTemplatesByUser :one
SELECT COUNT(*) FROM learning_log_templates WHERE user_id = $1;
