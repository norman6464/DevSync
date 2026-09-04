-- name: CountCodeSnippetsByUser :one
SELECT count(*) FROM code_snippets
WHERE user_id = $1;

-- name: SumCodeSnippetCommentCountByUser :one
SELECT COALESCE(SUM(comment_count), 0)::bigint FROM code_snippets
WHERE user_id = $1;

-- name: CountCodeSnippetLanguagesByUser :one
SELECT count(DISTINCT language) FROM code_snippets
WHERE user_id = $1;
