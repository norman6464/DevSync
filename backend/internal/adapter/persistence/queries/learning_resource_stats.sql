-- name: CountLearningResourcesByUser :one
SELECT count(*) FROM learning_resources
WHERE user_id = $1;

-- name: SumLearningResourceLikeCountByUser :one
-- like_countはlearning_resource_metrics側（DEVSYNC-159）。LEFT JOIN + COALESCEで0扱いにする。
SELECT COALESCE(SUM(lrm.like_count), 0)::bigint
FROM learning_resources lr
LEFT JOIN learning_resource_metrics lrm ON lrm.resource_id = lr.id
WHERE lr.user_id = $1;

-- name: SumLearningResourceSaveCountByUser :one
SELECT COALESCE(SUM(lrm.save_count), 0)::bigint
FROM learning_resources lr
LEFT JOIN learning_resource_metrics lrm ON lrm.resource_id = lr.id
WHERE lr.user_id = $1;

-- name: CountLearningResourceCategoriesByUser :one
SELECT count(DISTINCT category) FROM learning_resources
WHERE user_id = $1;
