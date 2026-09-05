-- name: CreateResourceReview :one
INSERT INTO resource_reviews (user_id, resource_id, rating, comment, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetResourceReviewWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(resource_reviews), sqlc.embed(users)
FROM resource_reviews
JOIN users ON users.id = resource_reviews.user_id
WHERE resource_reviews.id = $1;

-- name: CountResourceReviewsByResource :one
SELECT COUNT(*) FROM resource_reviews WHERE resource_id = $1;

-- name: ListResourceReviewsByResource :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(resource_reviews), sqlc.embed(users)
FROM resource_reviews
JOIN users ON users.id = resource_reviews.user_id
WHERE resource_reviews.resource_id = $1
ORDER BY resource_reviews.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetResourceReviewByUserAndResource :one
SELECT * FROM resource_reviews WHERE user_id = $1 AND resource_id = $2;

-- name: UpdateResourceReview :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE resource_reviews SET rating = $2, comment = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteResourceReview :exec
DELETE FROM resource_reviews WHERE id = $1;
