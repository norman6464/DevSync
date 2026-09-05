-- name: CreateBookmarkCollection :one
INSERT INTO bookmark_collections (user_id, name, description, color, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetBookmarkCollectionByID :one
SELECT * FROM bookmark_collections WHERE id = $1;

-- name: ListBookmarkCollectionsByUser :many
SELECT * FROM bookmark_collections WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateBookmarkCollection :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE bookmark_collections SET name = $2, description = $3, color = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBookmarkCollectionItemsByCollectionID :exec
DELETE FROM bookmark_collection_items WHERE collection_id = $1;

-- name: DeleteBookmarkCollection :exec
DELETE FROM bookmark_collections WHERE id = $1;

-- name: CreateBookmarkCollectionItem :execrows
-- (collection_id, post_id) の一意制約に任せてON CONFLICT DO NOTHINGで挿入し、
-- 実際に挿入できた行数を返す（呼び出し側で「既に入っていたか」を判定する）。
INSERT INTO bookmark_collection_items (collection_id, post_id, created_at)
VALUES ($1, $2, now())
ON CONFLICT (collection_id, post_id) DO NOTHING;

-- name: DeleteBookmarkCollectionItem :exec
DELETE FROM bookmark_collection_items WHERE collection_id = $1 AND post_id = $2;

-- name: CountBookmarkCollectionItemsByCollection :one
SELECT COUNT(*) FROM bookmark_collection_items WHERE collection_id = $1;

-- name: ListBookmarkCollectionItemsWithPostByCollection :many
-- GORMのPreload("Post")に相当（Post.Userは移行前もPreloadされていないため含めない）。
-- post_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(bookmark_collection_items), sqlc.embed(posts)
FROM bookmark_collection_items
JOIN posts ON posts.id = bookmark_collection_items.post_id
WHERE bookmark_collection_items.collection_id = $1
ORDER BY bookmark_collection_items.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountBookmarkCollectionsByUser :one
SELECT COUNT(*) FROM bookmark_collections WHERE user_id = $1;
