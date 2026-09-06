-- name: CountRoadmapsByUser :one
SELECT count(*) FROM roadmaps
WHERE user_id = $1;

-- name: CountRoadmapsByUserAndStatus :one
SELECT count(*) FROM roadmaps
WHERE user_id = $1 AND status = $2;

-- name: SumRoadmapStepCountByUser :one
-- step_countはroadmap_metrics側（DEVSYNC-159）。LEFT JOIN + COALESCEで0扱いにする。
SELECT COALESCE(SUM(rm.step_count), 0)::bigint
FROM roadmaps r
LEFT JOIN roadmap_metrics rm ON rm.roadmap_id = r.id
WHERE r.user_id = $1;

-- name: SumRoadmapCompletedStepCountByUser :one
SELECT COALESCE(SUM(rm.completed_step_count), 0)::bigint
FROM roadmaps r
LEFT JOIN roadmap_metrics rm ON rm.roadmap_id = r.id
WHERE r.user_id = $1;
