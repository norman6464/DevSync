-- name: CountPostsByUser :one
SELECT COUNT(*) FROM posts WHERE user_id = $1;

-- name: SumPostLikesReceivedByUser :one
SELECT COALESCE(SUM(like_count), 0)::bigint FROM posts WHERE user_id = $1;

-- name: CountGitHubContributionDaysByUser :one
SELECT COUNT(DISTINCT date) FROM git_hub_contributions WHERE user_id = $1 AND count > 0;

-- name: CountCompletedLearningGoalsByUser :one
SELECT COUNT(*) FROM learning_goals WHERE user_id = $1 AND status = $2;

-- name: CountCommentsByUser :one
SELECT COUNT(*) FROM comments WHERE user_id = $1;

-- name: ListLearningLogDatesByUser :many
-- 学習ログの連続記録日数（ストリーク）算出に使う、記録のある日付一覧（新しい順）。
SELECT DISTINCT DATE(created_at) AS log_date FROM learning_logs
WHERE user_id = $1
ORDER BY log_date DESC;
