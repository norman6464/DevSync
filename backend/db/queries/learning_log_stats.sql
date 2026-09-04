-- name: CountLearningLogsByUser :one
SELECT count(*) FROM learning_logs
WHERE user_id = $1;

-- name: SumLearningLogDurationByUser :one
SELECT COALESCE(SUM(duration), 0)::bigint FROM learning_logs
WHERE user_id = $1;

-- name: CountLearningLogCategoriesByUser :one
SELECT count(DISTINCT category) FROM learning_logs
WHERE user_id = $1;

-- name: CountLearningLogsByUserSince :one
SELECT count(*) FROM learning_logs
WHERE user_id = $1 AND created_at >= $2;
