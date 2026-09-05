-- name: GetPostWithUserByID :one
-- GORMのPreload("User")に相当（CodeSnippetsはpost_bookmark.sqlのListCodeSnippetsByPostIDsで別途取得する）。
-- user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(posts), sqlc.embed(users)
FROM posts
JOIN users ON users.id = posts.user_id
WHERE posts.id = $1;
