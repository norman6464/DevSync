-- name: DeletePostTagsByPostID :exec
DELETE FROM post_tags WHERE post_id = $1;

-- name: CreatePostTag :exec
INSERT INTO post_tags (post_id, tag) VALUES ($1, $2);

-- name: ListPostTagsByPostID :many
SELECT tag FROM post_tags WHERE post_id = $1 ORDER BY id ASC;

-- name: CountPostsByTag :one
-- 下書きは除外する。
SELECT COUNT(*) FROM post_tags
JOIN posts ON posts.id = post_tags.post_id AND posts.is_draft = false
WHERE post_tags.tag = $1;

-- name: ListPostIDsByTag :many
-- ページングはpost_tags.id順で決め、実データ（Post本体）は別クエリでcreated_at順に取得する
-- （移行前のGORM実装と同じ2段構え）。
SELECT post_tags.post_id FROM post_tags
JOIN posts ON posts.id = post_tags.post_id AND posts.is_draft = false
WHERE post_tags.tag = $1
ORDER BY post_tags.id DESC
LIMIT $2 OFFSET $3;

-- name: ListPostsByIDs :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.id = ANY($1::bigint[])
ORDER BY posts.created_at DESC;

-- name: ListPopularTags :many
SELECT tag, COUNT(*)::bigint AS count FROM post_tags
GROUP BY tag
ORDER BY count DESC
LIMIT $1;
