-- name: CountLearningResourcesByUser :one
-- learning_resources は GORM の論理削除（deleted_at）付きモデルのため、GORMの既定スコープに
-- 合わせて deleted_at IS NULL を明示する（Unscoped() されていない全クエリと同じ挙動）。
SELECT count(*) FROM learning_resources
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: SumLearningResourceLikeCountByUser :one
SELECT COALESCE(SUM(like_count), 0)::bigint FROM learning_resources
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: SumLearningResourceSaveCountByUser :one
SELECT COALESCE(SUM(save_count), 0)::bigint FROM learning_resources
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: CountLearningResourceCategoriesByUser :one
SELECT count(DISTINCT category) FROM learning_resources
WHERE user_id = $1 AND deleted_at IS NULL;
