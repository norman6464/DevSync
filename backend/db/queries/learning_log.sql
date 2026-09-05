-- name: CreateLearningLog :one
INSERT INTO learning_logs (
    user_id, title, content, category, duration, goal_id, source, is_favorite,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, now(), now()
) RETURNING *;

-- name: UpdateLearningLog :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE learning_logs SET
    title = $2, content = $3, category = $4, duration = $5, goal_id = $6, source = $7,
    is_favorite = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLearningLog :exec
-- 所有者本人のログだけを削除する。
DELETE FROM learning_logs WHERE id = $1 AND user_id = $2;

-- name: GetLearningLogByID :one
SELECT * FROM learning_logs WHERE id = $1;

-- name: ListLearningLogsByUser :many
SELECT * FROM learning_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFavoriteLearningLogs :many
SELECT * FROM learning_logs
WHERE user_id = $1 AND is_favorite = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountFavoriteLearningLogs :one
SELECT COUNT(*) FROM learning_logs WHERE user_id = $1 AND is_favorite = true;

-- name: ListLearningLogsByGoal :many
SELECT * FROM learning_logs
WHERE goal_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLearningLogsByGoal :one
SELECT COUNT(*) FROM learning_logs WHERE goal_id = $1;

-- name: ListLearningLogsByCategory :many
SELECT * FROM learning_logs
WHERE user_id = $1 AND category = $2
ORDER BY created_at DESC;

-- name: ListLearningLogsBySource :many
SELECT * FROM learning_logs
WHERE user_id = $1 AND source = $2
ORDER BY created_at DESC;

-- name: ListLearningLogsByPeriod :many
-- daysが0以下のときは呼び出し側がsinceにnilを渡し、全期間を対象にする。
SELECT * FROM learning_logs
WHERE user_id = $1 AND (sqlc.narg('since')::timestamptz IS NULL OR created_at >= sqlc.narg('since'))
ORDER BY created_at DESC;

-- name: SumLearningLogDurationSince :one
SELECT COALESCE(SUM(duration), 0)::bigint FROM learning_logs
WHERE user_id = $1 AND created_at >= $2;

-- name: SumLearningLogDurationByGoal :one
SELECT COALESCE(SUM(duration), 0)::bigint FROM learning_logs WHERE goal_id = $1;

-- name: ListRecentLearningLogCategories :many
SELECT category FROM learning_logs
WHERE user_id = $1
GROUP BY category
ORDER BY COUNT(*) DESC
LIMIT $2;

-- name: ListLearningLogCalendarData :many
SELECT DATE(created_at) AS date, COUNT(*) AS count
FROM learning_logs
WHERE user_id = $1
GROUP BY DATE(created_at)
ORDER BY date ASC;

-- name: ListLearningLogMonthlySummary :many
SELECT
    TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM-DD') AS month,
    COALESCE(SUM(duration), 0)::bigint AS total_minutes,
    COUNT(*) AS log_count
FROM learning_logs
WHERE user_id = $1 AND created_at >= $2
GROUP BY month
ORDER BY month;
