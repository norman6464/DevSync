-- name: UpsertQiitaArticle :one
INSERT INTO qiita_articles (user_id, qiita_id, title, url, likes_count, comments_count, tags, published_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
ON CONFLICT (qiita_id) DO UPDATE SET
    title = EXCLUDED.title,
    url = EXCLUDED.url,
    likes_count = EXCLUDED.likes_count,
    comments_count = EXCLUDED.comments_count,
    tags = EXCLUDED.tags,
    published_at = EXCLUDED.published_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListQiitaArticlesByUser :many
SELECT * FROM qiita_articles WHERE user_id = $1 ORDER BY published_at DESC;

-- name: GetQiitaStatsByUser :one
SELECT
    COUNT(*) AS total_articles,
    COALESCE(SUM(likes_count), 0)::bigint AS total_likes,
    COALESCE(SUM(comments_count), 0)::bigint AS total_comments
FROM qiita_articles WHERE user_id = $1;

-- name: DeleteQiitaArticlesByUser :exec
DELETE FROM qiita_articles WHERE user_id = $1;
