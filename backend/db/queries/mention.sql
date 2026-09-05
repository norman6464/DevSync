-- name: CreateMention :one
-- GORMの clause.OnConflict{DoNothing: true} に相当。衝突時は RETURNING 行が無くなるため
-- pgx.ErrNoRows を「作成されなかった」判定に使う。
INSERT INTO mentions (user_id, actor_id, post_id, comment_id, created_at)
VALUES ($1, $2, $3, $4, now())
ON CONFLICT DO NOTHING
RETURNING *;

-- name: ListMentionsByUser :many
-- GORMのPreload("Actor")に相当。actor_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(mentions), sqlc.embed(actor)
FROM mentions
JOIN users actor ON actor.id = mentions.actor_id
WHERE mentions.user_id = $1
ORDER BY mentions.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListMentionsByPostID :many
-- GORMのPreload("User").Preload("Actor")に相当。user_id/actor_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(mentions), sqlc.embed(mention_user), sqlc.embed(actor)
FROM mentions
JOIN users mention_user ON mention_user.id = mentions.user_id
JOIN users actor ON actor.id = mentions.actor_id
WHERE mentions.post_id = $1;

-- name: ListMentionsByCommentID :many
SELECT sqlc.embed(mentions), sqlc.embed(mention_user), sqlc.embed(actor)
FROM mentions
JOIN users mention_user ON mention_user.id = mentions.user_id
JOIN users actor ON actor.id = mentions.actor_id
WHERE mentions.comment_id = $1;

-- name: DeleteMentionsByPostID :exec
DELETE FROM mentions WHERE post_id = $1;

-- name: DeleteMentionsByCommentID :exec
DELETE FROM mentions WHERE comment_id = $1;
