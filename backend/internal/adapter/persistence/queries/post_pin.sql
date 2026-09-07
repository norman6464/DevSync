-- name: CreatePostPin :exec
INSERT INTO post_reactions (user_id, post_id, kind, pin_order, created_at)
VALUES ($1, $2, 'pin', $3, now());

-- name: DeletePostPin :exec
DELETE FROM post_reactions
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'pin';

-- name: ListPostPinsByUser :many
-- GORMのPreload("Post").Preload("Post.User")に相当。user_id/post_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_reactions), sqlc.embed(posts), sqlc.embed(users)
FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
JOIN users ON users.id = posts.user_id
WHERE post_reactions.user_id = $1 AND post_reactions.kind = 'pin'
ORDER BY post_reactions.pin_order ASC;

-- name: CountPostPinsByUser :one
SELECT COUNT(*) FROM post_reactions WHERE user_id = $1 AND kind = 'pin';

-- name: CountPostPinsByUserAndPost :one
SELECT COUNT(*) FROM post_reactions WHERE user_id = $1 AND post_id = $2 AND kind = 'pin';

-- name: UpdatePostPinOrder :exec
UPDATE post_reactions SET pin_order = $3
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2 AND post_reactions.kind = 'pin';
