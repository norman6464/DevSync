-- name: GetLearningResourceByID :one
SELECT * FROM learning_resources WHERE id = $1;
