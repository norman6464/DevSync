-- name: CreatePostPin :exec
INSERT INTO post_pins (user_id, post_id, pin_order, created_at)
VALUES ($1, $2, $3, now());

-- name: DeletePostPin :exec
DELETE FROM post_pins WHERE user_id = $1 AND post_id = $2;

-- name: ListPostPinsByUser :many
-- GORMのPreload("Post").Preload("Post.User")に相当。user_id/post_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_pins), sqlc.embed(posts), sqlc.embed(users)
FROM post_pins
JOIN posts ON posts.id = post_pins.post_id
JOIN users ON users.id = posts.user_id
WHERE post_pins.user_id = $1
ORDER BY post_pins.pin_order ASC;

-- name: CountPostPinsByUser :one
SELECT COUNT(*) FROM post_pins WHERE user_id = $1;

-- name: CountPostPinsByUserAndPost :one
SELECT COUNT(*) FROM post_pins WHERE user_id = $1 AND post_id = $2;

-- name: UpdatePostPinOrder :exec
UPDATE post_pins SET pin_order = $3 WHERE user_id = $1 AND post_id = $2;
