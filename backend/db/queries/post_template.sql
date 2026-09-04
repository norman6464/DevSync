-- name: CreatePostTemplate :one
INSERT INTO post_templates (user_id, name, title_template, content_template, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetPostTemplateByID :one
SELECT * FROM post_templates
WHERE id = $1;

-- name: ListPostTemplatesByUserID :many
SELECT * FROM post_templates
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPostTemplatesByUserID :one
SELECT count(*) FROM post_templates
WHERE user_id = $1;

-- name: UpdatePostTemplate :one
UPDATE post_templates
SET name = $2,
    title_template = $3,
    content_template = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePostTemplate :exec
DELETE FROM post_templates
WHERE id = $1;
