-- name: GetLearningResourceByID :one
-- learning_resources は GORM の論理削除（deleted_at）付きモデルのため deleted_at IS NULL を明示する。
SELECT * FROM learning_resources WHERE id = $1 AND deleted_at IS NULL;
