-- name: GetRoadmapMetricsByRoadmapIDs :many
-- ロードマップ一覧へstep_count/completed_step_count/progressを付与するためのまとめ取り。
-- progressは列として持たず、ここでSELECT側の算出式として都度計算する（生成列にすると
-- 既存のUPDATE系クエリを壊すため不採用。DEVSYNC-159）。1件もステップが無い
-- ロードマップはroadmap_metrics行が存在しない（遅延生成）ため、この結果に
-- 現れないロードマップはGo側で0扱いにする（attachRoadmapMetrics参照）。
SELECT
    roadmap_id,
    step_count,
    completed_step_count,
    CASE WHEN step_count > 0 THEN (completed_step_count * 100 / step_count) ELSE 0 END::bigint AS progress
FROM roadmap_metrics
WHERE roadmap_id = ANY($1::bigint[]);

-- name: ReconcileAllRoadmapMetrics :exec
-- 夜次reconcileジョブ本体。roadmap_stepsの実件数からroadmap_metrics全件を補正する。
INSERT INTO roadmap_metrics (roadmap_id, step_count, completed_step_count)
SELECT
    r.id,
    COALESCE(s.cnt, 0),
    COALESCE(sc.cnt, 0)
FROM roadmaps r
LEFT JOIN (SELECT roadmap_id, COUNT(*) AS cnt FROM roadmap_steps GROUP BY roadmap_id) s ON s.roadmap_id = r.id
LEFT JOIN (SELECT roadmap_id, COUNT(*) AS cnt FROM roadmap_steps WHERE is_completed = true GROUP BY roadmap_id) sc ON sc.roadmap_id = r.id
ON CONFLICT (roadmap_id) DO UPDATE SET
    step_count = EXCLUDED.step_count,
    completed_step_count = EXCLUDED.completed_step_count;
