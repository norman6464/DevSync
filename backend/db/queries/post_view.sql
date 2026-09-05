-- name: CreatePostView :execrows
-- GORMの clause.OnConflict{DoNothing: true} に相当。実際に挿入できた行数を返し、
-- 呼び出し側で「既に閲覧済みだったか」を判定する。
INSERT INTO post_views (user_id, post_id, created_at)
VALUES ($1, $2, now())
ON CONFLICT DO NOTHING;

-- name: IncrementPostViewCount :exec
UPDATE posts SET view_count = view_count + 1 WHERE id = $1;

-- name: CountPostViewsByPost :one
SELECT COUNT(*) FROM post_views WHERE post_id = $1;

-- name: ListMostViewedPosts :many
SELECT post_id, COUNT(*)::bigint AS view_count FROM post_views
GROUP BY post_id
ORDER BY view_count DESC
LIMIT $1;
