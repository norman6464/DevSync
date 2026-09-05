-- name: UpsertResourceProgress :one
INSERT INTO resource_progresses (user_id, resource_id, status, completion_percent, note, started_at, completed_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
ON CONFLICT (user_id, resource_id) DO UPDATE SET
    status = EXCLUDED.status,
    completion_percent = EXCLUDED.completion_percent,
    note = EXCLUDED.note,
    started_at = EXCLUDED.started_at,
    completed_at = EXCLUDED.completed_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetResourceProgressByUserAndResource :one
SELECT * FROM resource_progresses WHERE user_id = $1 AND resource_id = $2;

-- name: CountResourceProgressByUser :one
SELECT COUNT(*) FROM resource_progresses
WHERE user_id = $1
    AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));

-- name: ListResourceProgressByUser :many
-- GORMのPreload("Resource")に相当。resourceは論理削除（deleted_at）付きモデルのため、
-- ソフトデリート済みリソースはJOIN条件で除外しResourceがnilになるようLEFT JOINにする
-- （resource_idの列自体はNOT NULLだが、論理削除の有無でJOINの成否が変わるため）。
-- id を第2ソートキーにして、updated_at 同値の行でもページングが安定するようにする
-- （移行前のGORM実装と同じ）。
SELECT sqlc.embed(resource_progresses),
    lr.id AS resource_id_2,
    lr.user_id AS resource_user_id,
    lr.title AS resource_title,
    lr.description AS resource_description,
    lr.url AS resource_url,
    lr.category AS resource_category,
    lr.difficulty AS resource_difficulty,
    lr.tags AS resource_tags,
    lr.image_url AS resource_image_url,
    lr.is_public AS resource_is_public,
    lr.like_count AS resource_like_count,
    lr.save_count AS resource_save_count,
    lr.created_at AS resource_created_at,
    lr.updated_at AS resource_updated_at
FROM resource_progresses
LEFT JOIN learning_resources lr ON lr.id = resource_progresses.resource_id AND lr.deleted_at IS NULL
WHERE resource_progresses.user_id = $1
    AND (sqlc.narg('status')::text IS NULL OR resource_progresses.status = sqlc.narg('status'))
ORDER BY resource_progresses.updated_at DESC, resource_progresses.id DESC
LIMIT $2 OFFSET $3;
