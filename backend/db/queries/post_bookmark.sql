-- name: CreateBookmark :exec
INSERT INTO bookmarks (user_id, post_id, created_at) VALUES ($1, $2, now());

-- name: IncrementPostBookmarkCount :exec
UPDATE posts SET bookmark_count = bookmark_count + 1 WHERE id = $1;

-- name: DeleteBookmark :execrows
DELETE FROM bookmarks WHERE user_id = $1 AND post_id = $2;

-- name: DecrementPostBookmarkCount :exec
UPDATE posts SET bookmark_count = GREATEST(bookmark_count - 1, 0) WHERE id = $1;

-- name: CountBookmarkByUserAndPost :one
SELECT COUNT(*) FROM bookmarks WHERE user_id = $1 AND post_id = $2;

-- name: ListBookmarkedPostsByUser :many
-- GORMのPreload("User")に相当（CodeSnippetsは別クエリで取得しGo側で結合する）。
-- user_id/post_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN bookmarks ON bookmarks.post_id = posts.id
JOIN users ON users.id = posts.user_id
WHERE bookmarks.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCodeSnippetsByPostIDs :many
-- GORMのPreload("CodeSnippets")に相当。1対多のため投稿IDのまとめ取りとGo側でのグルーピングで再現する。
SELECT * FROM code_snippets WHERE post_id = ANY($1::bigint[]) ORDER BY created_at ASC;
