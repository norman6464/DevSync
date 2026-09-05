-- name: CreateCodeSnippet :one
INSERT INTO code_snippets (
    post_id, user_id, language, file_name, code, comment_count, forked_from_id, fork_count,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, now(), now()
) RETURNING *;

-- name: GetCodeSnippetByID :one
SELECT * FROM code_snippets WHERE code_snippets.id = $1;

-- name: FindCodeSnippetsByUserAndLanguage :many
SELECT * FROM code_snippets
WHERE code_snippets.user_id = $1 AND code_snippets.language = $2
ORDER BY code_snippets.created_at DESC;

-- name: SearchCodeSnippets :many
-- CodeSnippet は投稿者の ID しか持たず User の関連を張っていないため、Preload しない
-- （移行前のGORM実装のコメントの通り、Preloadするとunsupported relationsで失敗する）。
SELECT * FROM code_snippets
WHERE code_snippets.language ILIKE $1 OR code_snippets.file_name ILIKE $1 OR code_snippets.code ILIKE $1
ORDER BY code_snippets.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchCodeSnippets :one
SELECT COUNT(*) FROM code_snippets
WHERE code_snippets.language ILIKE $1 OR code_snippets.file_name ILIKE $1 OR code_snippets.code ILIKE $1;

-- name: UpdateCodeSnippet :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE code_snippets SET
    language = $2, file_name = $3, code = $4, comment_count = $5, forked_from_id = $6,
    fork_count = $7, updated_at = now()
WHERE code_snippets.id = $1
RETURNING *;

-- name: DeleteCodeSnippetByID :exec
DELETE FROM code_snippets WHERE code_snippets.id = $1;

-- name: IncrementSnippetForkCount :exec
UPDATE code_snippets SET fork_count = fork_count + 1 WHERE code_snippets.id = $1;

-- name: CreateSnippetComment :one
INSERT INTO snippet_comments (snippet_id, user_id, line_number, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: IncrementSnippetCommentCount :exec
UPDATE code_snippets SET comment_count = comment_count + 1 WHERE code_snippets.id = $1;

-- name: DecrementSnippetCommentCountFloored :exec
-- 0未満にはしない（GORMのGREATEST(comment_count - 1, 0)に相当）。
UPDATE code_snippets SET comment_count = GREATEST(comment_count - 1, 0) WHERE code_snippets.id = $1;

-- name: GetSnippetCommentsWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(snippet_comments), sqlc.embed(users)
FROM snippet_comments
JOIN users ON users.id = snippet_comments.user_id
WHERE snippet_comments.snippet_id = $1
ORDER BY snippet_comments.line_number ASC, snippet_comments.created_at ASC;

-- name: GetSnippetCommentByID :one
SELECT * FROM snippet_comments WHERE snippet_comments.id = $1;

-- name: DeleteSnippetCommentByID :exec
DELETE FROM snippet_comments WHERE snippet_comments.id = $1;

-- name: CreateSnippetFavorite :exec
INSERT INTO code_snippet_favorites (user_id, snippet_id, created_at) VALUES ($1, $2, now());

-- name: DeleteSnippetFavorite :exec
DELETE FROM code_snippet_favorites
WHERE code_snippet_favorites.user_id = $1 AND code_snippet_favorites.snippet_id = $2;

-- name: CountSnippetFavorite :one
SELECT COUNT(*) FROM code_snippet_favorites
WHERE code_snippet_favorites.user_id = $1 AND code_snippet_favorites.snippet_id = $2;

-- name: ListFavoritedCodeSnippets :many
SELECT * FROM code_snippets
WHERE code_snippets.id IN (
    SELECT f.snippet_id FROM code_snippet_favorites f WHERE f.user_id = $1
)
ORDER BY code_snippets.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountFavoritedCodeSnippets :one
SELECT COUNT(*) FROM code_snippets
WHERE code_snippets.id IN (
    SELECT f.snippet_id FROM code_snippet_favorites f WHERE f.user_id = $1
);
