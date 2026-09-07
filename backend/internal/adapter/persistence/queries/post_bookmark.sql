-- name: CreateBookmark :exec
INSERT INTO post_reactions (user_id, post_id, kind, created_at) VALUES ($1, $2, 'bookmark', now());

-- name: DeleteBookmark :execrows
DELETE FROM post_reactions
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'bookmark';

-- name: CountBookmarkByUserAndPost :one
SELECT COUNT(*) FROM post_reactions
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'bookmark';

-- name: ListBookmarkedPostsByUser :many
-- GORMのPreload("User")に相当（CodeSnippetsは別クエリで取得しGo側で結合する）。
-- user_id/post_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN post_reactions ON post_reactions.post_id = posts.id AND post_reactions.kind = 'bookmark'
JOIN users ON users.id = posts.user_id
WHERE post_reactions.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCodeSnippetsByPostIDs :many
-- GORMのPreload("CodeSnippets")に相当。1対多のため投稿IDのまとめ取りとGo側でのグルーピングで再現する。
SELECT * FROM code_snippets WHERE post_id = ANY($1::bigint[]) ORDER BY created_at ASC;

-- name: CountBookmarksByPostIDs :many
-- post.bookmark_countはORDER BYに使われない表示専用の値だったため列として持たず、
-- 都度post_reactionsからCOUNT(*)する。投稿IDのまとめ取りとGo側でのグルーピングで
-- ListCodeSnippetsByPostIDsと同じ形にする。ブックマーク0件の投稿はこの結果に現れない。
SELECT post_id, COUNT(*)::bigint AS bookmark_count
FROM post_reactions
WHERE post_id = ANY($1::bigint[]) AND kind = 'bookmark'
GROUP BY post_id;
