-- name: CreatePostView :execrows
-- GORMの clause.OnConflict{DoNothing: true} に相当。post_reactions(kind='view')への記録と
-- post_metrics.view_countの加算を同一SQL文で行う（DEVSYNC-159）。既に閲覧済み
-- （ON CONFLICT DO NOTHINGでinsertedが0行）のときはpost_metricsも更新されず、
-- rowsAffectedは従来どおり「実際に新規閲覧として記録できたか」を表す。
-- DEVSYNC-160でpost_views単独テーブルからpost_reactionsへ統合。
WITH inserted AS (
    INSERT INTO post_reactions (user_id, post_id, kind, created_at)
    VALUES ($1, $2, 'view', now())
    ON CONFLICT (user_id, post_id, kind, value) DO NOTHING
    RETURNING post_reactions.post_id
)
INSERT INTO post_metrics (post_id, view_count)
SELECT inserted.post_id, 1 FROM inserted
ON CONFLICT (post_id) DO UPDATE SET view_count = post_metrics.view_count + 1;

-- name: CountPostViewsByPost :one
SELECT COUNT(*) FROM post_reactions WHERE post_reactions.post_id = $1 AND post_reactions.kind = 'view';

-- name: ListMostViewedPosts :many
SELECT post_reactions.post_id, COUNT(*)::bigint AS view_count FROM post_reactions
WHERE post_reactions.kind = 'view'
GROUP BY post_reactions.post_id
ORDER BY view_count DESC
LIMIT $1;
