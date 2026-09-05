-- name: UpsertWeeklyGoal :one
INSERT INTO weekly_goals (user_id, category, target_minutes, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
ON CONFLICT (user_id, category) DO UPDATE
SET target_minutes = EXCLUDED.target_minutes,
    updated_at = now()
RETURNING *;

-- name: ListWeeklyGoalsByUser :many
SELECT * FROM weekly_goals
WHERE user_id = $1
ORDER BY category ASC;

-- name: SumLearningLogDurationByUserCategorySince :one
SELECT COALESCE(SUM(duration), 0)::bigint FROM learning_logs
WHERE user_id = $1 AND category = $2 AND created_at >= $3;
