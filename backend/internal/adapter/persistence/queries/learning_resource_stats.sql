-- name: CountLearningResourcesByUser :one
SELECT count(*) FROM learning_resources
WHERE user_id = $1;

-- name: SumLearningResourceLikeCountByUser :one
SELECT COALESCE(SUM(like_count), 0)::bigint FROM learning_resources
WHERE user_id = $1;

-- name: SumLearningResourceSaveCountByUser :one
SELECT COALESCE(SUM(save_count), 0)::bigint FROM learning_resources
WHERE user_id = $1;

-- name: CountLearningResourceCategoriesByUser :one
SELECT count(DISTINCT category) FROM learning_resources
WHERE user_id = $1;
