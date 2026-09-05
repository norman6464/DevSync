-- name: CreatePost :one
INSERT INTO posts (
    user_id, title, content, image_urls, is_draft, like_count, comment_count,
    bookmark_count, view_count, estimated_read_time, scheduled_at, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now()
) RETURNING *;

-- name: UpdatePost :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE posts SET
    title = $2, content = $3, image_urls = $4, is_draft = $5, like_count = $6,
    comment_count = $7, bookmark_count = $8, view_count = $9, estimated_read_time = $10,
    scheduled_at = $11, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: LockPostForDelete :one
-- 投稿削除の直列化のための行ロック（GORMの clause.Locking{Strength: "UPDATE"} に相当）。
-- 存在しなければ pgx.ErrNoRows を返し、呼び出し側で「既に無ければ何もしない（冪等）」を判定する。
SELECT id FROM posts WHERE id = $1 FOR UPDATE;

-- name: DeleteCommentLikesByPostComments :exec
DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE post_id = $1);

-- name: DeleteMentionsByPostComments :exec
-- mentions自身もpost_id列を持つため、サブクエリのcommentsを明示的にエイリアス修飾しないと
-- post_idの参照先が曖昧になる（PostgreSQLの相関サブクエリの解決規則による）。
DELETE FROM mentions WHERE comment_id IN (SELECT c.id FROM comments c WHERE c.post_id = $1);

-- name: DeleteCommentsByPost :exec
DELETE FROM comments WHERE post_id = $1;

-- name: DeleteSnippetCommentsByPostSnippets :exec
DELETE FROM snippet_comments WHERE snippet_id IN (SELECT cs.id FROM code_snippets cs WHERE cs.post_id = $1);

-- name: DeleteCodeSnippetsByPost :exec
DELETE FROM code_snippets WHERE post_id = $1;

-- name: DeleteLikesByPost :exec
DELETE FROM likes WHERE post_id = $1;

-- name: DeleteReactionsByPost :exec
DELETE FROM reactions WHERE post_id = $1;

-- name: DeleteBookmarksByPost :exec
DELETE FROM bookmarks WHERE post_id = $1;

-- name: DeleteBookmarkCollectionItemsByPost :exec
DELETE FROM bookmark_collection_items WHERE post_id = $1;

-- name: DeletePostSeriesItemsByPost :exec
DELETE FROM post_series_items WHERE post_id = $1;

-- name: DeletePostCollectionItemsByPost :exec
DELETE FROM post_collection_items WHERE post_id = $1;

-- name: DeletePostTagsByPost :exec
DELETE FROM post_tags WHERE post_id = $1;

-- name: DeletePostPinsByPost :exec
DELETE FROM post_pins WHERE post_id = $1;

-- name: DeletePostViewsByPost :exec
DELETE FROM post_views WHERE post_id = $1;

-- name: DeleteNotificationsByPost :exec
DELETE FROM notifications WHERE post_id = $1;

-- name: DeleteMentionsByPost :exec
DELETE FROM mentions WHERE post_id = $1;

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
