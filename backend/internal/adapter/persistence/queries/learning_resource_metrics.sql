-- name: GetLearningResourceMetricsByResourceIDs :many
-- リソース一覧へlike_count/save_countを付与するためのまとめ取り。1件もいいね/保存が
-- 無いリソースはlearning_resource_metrics行が存在しない（遅延生成）ため、この結果に
-- 現れないリソースはGo側で0扱いにする（attachLearningResourceMetrics参照）。
SELECT resource_id, like_count, save_count
FROM learning_resource_metrics
WHERE resource_id = ANY($1::bigint[]);

-- name: ReconcileAllLearningResourceMetrics :exec
-- 夜次reconcileジョブ本体。resource_likes/resource_savesの実件数から
-- learning_resource_metrics全件をまとめて補正する。
INSERT INTO learning_resource_metrics (resource_id, like_count, save_count)
SELECT
    lr.id,
    COALESCE(l.cnt, 0),
    COALESCE(s.cnt, 0)
FROM learning_resources lr
LEFT JOIN (SELECT resource_id, COUNT(*) AS cnt FROM resource_likes GROUP BY resource_id) l ON l.resource_id = lr.id
LEFT JOIN (SELECT resource_id, COUNT(*) AS cnt FROM resource_saves GROUP BY resource_id) s ON s.resource_id = lr.id
ON CONFLICT (resource_id) DO UPDATE SET
    like_count = EXCLUDED.like_count,
    save_count = EXCLUDED.save_count;
