-- name: CreateBookmark :exec
INSERT INTO bookmarks (user_id, post_id, created_at) VALUES ($1, $2, now());

-- name: DeleteBookmark :execrows
DELETE FROM bookmarks WHERE user_id = $1 AND post_id = $2;

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

-- name: CountBookmarksByPostIDs :many
-- post.bookmark_countはORDER BYに使われない表示専用の値だったため列として持たず、
-- 都度bookmarksからCOUNT(*)する。投稿IDのまとめ取りとGo側でのグルーピングで
-- ListCodeSnippetsByPostIDsと同じ形にする。ブックマーク0件の投稿はこの結果に現れない。
SELECT post_id, COUNT(*)::bigint AS bookmark_count
FROM bookmarks
WHERE post_id = ANY($1::bigint[])
GROUP BY post_id;
