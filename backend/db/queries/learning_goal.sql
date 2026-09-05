-- name: CreateLearningGoal :one
INSERT INTO learning_goals (
    user_id, title, description, category, target_date, progress, target_hours, status, is_public,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()
) RETURNING *;

-- name: UpdateLearningGoal :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE learning_goals SET
    title = $2, description = $3, category = $4, target_date = $5, progress = $6,
    target_hours = $7, status = $8, is_public = $9, completed_at = $10, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLearningGoal :exec
DELETE FROM learning_goals WHERE id = $1;

-- name: GetLearningGoalByID :one
SELECT * FROM learning_goals WHERE id = $1;

-- name: ListLearningGoalsByUser :many
SELECT * FROM learning_goals
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLearningGoalsByUser :one
-- 件数は GetByUserID の総数取得と CountByUserID(単体メソッド) の両方から利用する。
SELECT COUNT(*) FROM learning_goals WHERE user_id = $1;

-- name: ListActiveLearningGoalsByUser :many
SELECT * FROM learning_goals
WHERE user_id = $1 AND status = 'active'
ORDER BY created_at DESC;

-- name: ListLearningGoalsByCategory :many
SELECT * FROM learning_goals
WHERE user_id = $1 AND category = $2
ORDER BY created_at DESC;

-- name: ListLearningGoalsByStatus :many
SELECT * FROM learning_goals
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: ListPublicLearningGoalsByUser :many
SELECT * FROM learning_goals
WHERE user_id = $1 AND is_public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPublicLearningGoalsByUser :one
SELECT COUNT(*) FROM learning_goals WHERE user_id = $1 AND is_public = true;

-- name: ListPublicLearningGoals :many
SELECT * FROM learning_goals
WHERE is_public = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicLearningGoals :one
SELECT COUNT(*) FROM learning_goals WHERE is_public = true;

-- name: CountActiveLearningGoalsByUser :one
SELECT COUNT(*) FROM learning_goals WHERE user_id = $1 AND status = 'active';

-- name: GetAverageActiveProgressByUser :one
SELECT COALESCE(AVG(progress), 0)::float8 FROM learning_goals
WHERE user_id = $1 AND status = 'active';
