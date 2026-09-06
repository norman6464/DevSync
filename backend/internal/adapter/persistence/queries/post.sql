-- name: CreatePost :one
INSERT INTO posts (
    user_id, title, content, image_urls, is_draft, like_count, comment_count,
    view_count, estimated_read_time, scheduled_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now()
) RETURNING *;

-- name: UpdatePost :one
-- GORMのSave（全カラム上書き）に相当。ただしlike_count/comment_count/bookmark_count/
-- view_countは対象外（Increment/Decrement系の専用クエリだけが更新する）。
-- ここに含めると、他リクエストによるカウンタ更新をこのUPDATEが読み取り時点の
-- 古い値で上書きする「ロストアップデート」を起こす。
UPDATE posts SET
    title = $2, content = $3, image_urls = $4, is_draft = $5,
    estimated_read_time = $6, scheduled_at = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: ListPublicPostsWithUser :many
-- GORMのPreload("User")に相当（CodeSnippetsは別途post_bookmark.sqlのListCodeSnippetsByPostIDsで取得する）。
-- user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.is_draft = false
ORDER BY posts.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicPosts :one
SELECT COUNT(*) FROM posts WHERE is_draft = false;

-- name: ListPublishedPostsByUserWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.user_id = $1 AND posts.is_draft = false
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDraftPostsByUserWithUser :many
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.user_id = $1 AND posts.is_draft = true
ORDER BY posts.updated_at DESC;

-- name: ListScheduledPostsByUserWithUser :many
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.user_id = $1 AND posts.is_draft = true AND posts.scheduled_at IS NOT NULL
ORDER BY posts.scheduled_at ASC;

-- name: CountScheduledPostsByUser :one
SELECT COUNT(*) FROM posts WHERE user_id = $1 AND is_draft = true AND scheduled_at IS NOT NULL;

-- name: ListTimelinePostsWithUser :many
-- フォロー中ユーザーと自分自身の公開済み投稿（移行前のGo実装のサブクエリをそのまま踏襲）。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE (posts.user_id IN (SELECT followee_id FROM follows WHERE follower_id = $1) OR posts.user_id = $1)
    AND posts.is_draft = false
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;
