-- name: GetLearningHeatmapData :many
SELECT
    EXTRACT(DOW FROM created_at)::int AS day_of_week,
    EXTRACT(HOUR FROM created_at)::int AS hour,
    COALESCE(SUM(duration), 0)::bigint AS total_minutes
FROM learning_logs
WHERE user_id = $1
GROUP BY day_of_week, hour
ORDER BY day_of_week, hour;

-- name: GetLearningCategoryBreakdown :many
-- 割合はusecase側で計算するため、ここでは合計時間とログ件数のみを返す。
SELECT
    category,
    COALESCE(SUM(duration), 0)::bigint AS total_minutes,
    COUNT(*) AS log_count
FROM learning_logs
WHERE user_id = $1
GROUP BY category
ORDER BY total_minutes DESC;

-- name: GetLearningWeeklyTrends :many
SELECT
    TO_CHAR(DATE_TRUNC('week', created_at), 'YYYY-MM-DD') AS week_start,
    COALESCE(SUM(duration), 0)::bigint AS total_minutes,
    COUNT(*) AS log_count
FROM learning_logs
WHERE user_id = $1 AND created_at >= $2
GROUP BY week_start
ORDER BY week_start;

-- name: CountLearningLogsBySourceSince :one
SELECT COUNT(*) FROM learning_logs WHERE user_id = $1 AND source = $2 AND created_at >= $3;

-- name: CountLearningLogDaysSince :one
SELECT COUNT(DISTINCT DATE(created_at)) FROM learning_logs WHERE user_id = $1 AND created_at >= $2;

-- name: ListDistinctLearningLogDates :many
-- ストリーク（現在・最長連続日数）の算出に使う、学習記録がある日付の一覧（新しい順）。
SELECT DISTINCT DATE(created_at) AS log_date
FROM learning_logs
WHERE user_id = $1
ORDER BY log_date DESC;
