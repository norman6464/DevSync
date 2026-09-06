-- name: UpsertZennArticle :one
INSERT INTO zenn_articles (user_id, zenn_id, title, slug, emoji, article_type, liked_count, comments_count, published_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
ON CONFLICT (user_id, zenn_id) DO UPDATE SET
    title = EXCLUDED.title,
    slug = EXCLUDED.slug,
    emoji = EXCLUDED.emoji,
    article_type = EXCLUDED.article_type,
    liked_count = EXCLUDED.liked_count,
    comments_count = EXCLUDED.comments_count,
    published_at = EXCLUDED.published_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: ListZennArticlesByUser :many
SELECT * FROM zenn_articles WHERE user_id = $1 ORDER BY published_at DESC;

-- name: GetZennStatsByUser :one
SELECT
    COUNT(*) AS total_articles,
    COALESCE(SUM(liked_count), 0)::bigint AS total_likes,
    COALESCE(SUM(comments_count), 0)::bigint AS total_comments
FROM zenn_articles WHERE user_id = $1;

-- name: DeleteZennArticlesByUser :exec
DELETE FROM zenn_articles WHERE user_id = $1;
