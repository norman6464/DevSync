-- name: GetPostMetricsByPostIDs :many
-- 投稿一覧へlike_count/comment_count/view_countを付与するためのまとめ取り。
-- 1件もいいね/コメント/閲覧が無い投稿はpost_metrics行が存在しない（遅延生成）ため、
-- この結果に現れない投稿はGo側で0扱いにする（attachMetricsToPosts参照）。
SELECT post_id, like_count, comment_count, view_count
FROM post_metrics
WHERE post_id = ANY($1::bigint[]);

-- name: IncrementPostLikeMetric :exec
-- likesへのINSERTと同じ呼び出し元から同一トランザクション相当で呼ぶ。post_metrics行が
-- 無ければ1から作る（遅延生成）。
INSERT INTO post_metrics (post_id, like_count) VALUES ($1, 1)
ON CONFLICT (post_id) DO UPDATE SET like_count = post_metrics.like_count + 1;

-- name: DecrementPostLikeMetric :exec
-- 0未満にはしない。post_metrics行がまだ無い場合（理論上は起こらないが防御的に）は
-- 0のまま新規作成する。
INSERT INTO post_metrics (post_id, like_count) VALUES ($1, 0)
ON CONFLICT (post_id) DO UPDATE SET like_count = GREATEST(post_metrics.like_count - 1, 0);

-- name: IncrementPostCommentMetric :exec
INSERT INTO post_metrics (post_id, comment_count) VALUES ($1, 1)
ON CONFLICT (post_id) DO UPDATE SET comment_count = post_metrics.comment_count + 1;

-- name: DecrementPostCommentMetric :exec
INSERT INTO post_metrics (post_id, comment_count) VALUES ($1, 0)
ON CONFLICT (post_id) DO UPDATE SET comment_count = GREATEST(post_metrics.comment_count - 1, 0);

-- name: IncrementPostViewMetric :exec
INSERT INTO post_metrics (post_id, view_count) VALUES ($1, 1)
ON CONFLICT (post_id) DO UPDATE SET view_count = post_metrics.view_count + 1;

-- name: ReconcileAllPostMetrics :exec
-- 夜次reconcileジョブ本体。likes/comments/post_viewsの実件数からpost_metricsを
-- 全件まとめて補正する。CASCADE削除等でIncrement/Decrementを経由しない変化を吸収する。
-- 1件も無いカウンタは0で確定させる（COALESCE）。
INSERT INTO post_metrics (post_id, like_count, comment_count, view_count)
SELECT
    p.id,
    COALESCE(l.cnt, 0),
    COALESCE(c.cnt, 0),
    COALESCE(v.cnt, 0)
FROM posts p
LEFT JOIN (SELECT post_id, COUNT(*) AS cnt FROM likes GROUP BY post_id) l ON l.post_id = p.id
LEFT JOIN (SELECT post_id, COUNT(*) AS cnt FROM comments GROUP BY post_id) c ON c.post_id = p.id
LEFT JOIN (SELECT post_id, COUNT(*) AS cnt FROM post_views GROUP BY post_id) v ON v.post_id = p.id
ON CONFLICT (post_id) DO UPDATE SET
    like_count = EXCLUDED.like_count,
    comment_count = EXCLUDED.comment_count,
    view_count = EXCLUDED.view_count;
