-- name: CreatePostSeries :one
INSERT INTO post_series (user_id, title, description, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetPostSeriesWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_series), sqlc.embed(users)
FROM post_series
JOIN users ON users.id = post_series.user_id
WHERE post_series.id = $1;

-- name: ListPostSeriesByUser :many
SELECT * FROM post_series
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPostSeriesByUser :one
SELECT COUNT(*) FROM post_series WHERE user_id = $1;

-- name: UpdatePostSeries :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE post_series SET title = $2, description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePostSeriesItemsBySeriesID :exec
DELETE FROM post_series_items WHERE series_id = $1;

-- name: DeletePostSeries :exec
DELETE FROM post_series WHERE id = $1;

-- name: CreatePostSeriesItem :one
INSERT INTO post_series_items (series_id, post_id, order_index)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CountPostSeriesItemsBySeriesAndPost :one
SELECT COUNT(*) FROM post_series_items WHERE series_id = $1 AND post_id = $2;

-- name: DeletePostSeriesItem :exec
DELETE FROM post_series_items WHERE series_id = $1 AND post_id = $2;

-- name: ListPostSeriesItemsWithPostBySeriesID :many
-- GORMのPreload("Post").Preload("Post.User")に相当。post_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_series_items), sqlc.embed(posts), sqlc.embed(users)
FROM post_series_items
JOIN posts ON posts.id = post_series_items.post_id
JOIN users ON users.id = posts.user_id
WHERE post_series_items.series_id = $1
ORDER BY post_series_items.order_index ASC;
