-- name: CreatePostCollection :one
INSERT INTO post_collections (user_id, title, description, is_public, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetPostCollectionWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_collections), sqlc.embed(users)
FROM post_collections
JOIN users ON users.id = post_collections.user_id
WHERE post_collections.id = $1;

-- name: ListPostCollectionsByUser :many
SELECT * FROM post_collections
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountPostCollectionsByUser :one
SELECT COUNT(*) FROM post_collections WHERE user_id = $1;

-- name: ListPublicPostCollectionsByUser :many
SELECT * FROM post_collections
WHERE user_id = $1 AND is_public = true
ORDER BY created_at DESC;

-- name: UpdatePostCollection :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE post_collections SET title = $2, description = $3, is_public = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePostCollectionItemsByCollectionID :exec
DELETE FROM post_collection_items WHERE collection_id = $1;

-- name: DeletePostCollection :exec
DELETE FROM post_collections WHERE id = $1;

-- name: CreatePostCollectionItem :one
INSERT INTO post_collection_items (collection_id, post_id, note, order_index, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: CountPostCollectionItemsByCollectionAndPost :one
SELECT COUNT(*) FROM post_collection_items WHERE collection_id = $1 AND post_id = $2;

-- name: DeletePostCollectionItem :exec
DELETE FROM post_collection_items WHERE collection_id = $1 AND post_id = $2;

-- name: ListPostCollectionItemsWithPostByCollectionID :many
-- GORMのPreload("Post").Preload("Post.User")に相当。post_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(post_collection_items), sqlc.embed(posts), sqlc.embed(users)
FROM post_collection_items
JOIN posts ON posts.id = post_collection_items.post_id
JOIN users ON users.id = posts.user_id
WHERE post_collection_items.collection_id = $1
ORDER BY post_collection_items.order_index ASC;
