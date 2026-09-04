-- name: CountRoadmapsByUser :one
SELECT count(*) FROM roadmaps
WHERE user_id = $1;

-- name: CountRoadmapsByUserAndStatus :one
SELECT count(*) FROM roadmaps
WHERE user_id = $1 AND status = $2;

-- name: SumRoadmapStepCountByUser :one
SELECT COALESCE(SUM(step_count), 0)::bigint FROM roadmaps
WHERE user_id = $1;

-- name: SumRoadmapCompletedStepCountByUser :one
SELECT COALESCE(SUM(completed_step_count), 0)::bigint FROM roadmaps
WHERE user_id = $1;
