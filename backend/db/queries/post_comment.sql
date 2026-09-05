-- name: CreatePostComment :one
INSERT INTO comments (user_id, post_id, parent_id, content, like_count, is_hidden, created_at, updated_at)
VALUES ($1, $2, $3, $4, 0, false, now(), now())
RETURNING *;

-- name: IncrementPostCommentCount :exec
UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1;

-- name: DecrementPostCommentCount :exec
UPDATE posts SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = $1;

-- name: UpdatePostComment :one
-- GORMのSave（全カラム上書き）に相当。呼び出し側は必ずDBから読み込んだcommentの
-- content/is_hiddenだけを変更してから呼ぶため、この2カラムの書き戻しで等価になる。
UPDATE comments SET content = $2, is_hidden = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePostComment :exec
DELETE FROM comments WHERE id = $1;

-- name: ListTopLevelCommentsByPost :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(comments), sqlc.embed(users)
FROM comments
JOIN users ON users.id = comments.user_id
WHERE comments.post_id = $1 AND comments.parent_id IS NULL
ORDER BY comments.created_at ASC;

-- name: ListCommentRepliesByParentIDs :many
-- GORMのPreload("User")に相当（返信の取得）。ListReplies単体呼び出しと
-- ListByPostIDのバッチ取得の両方で、親IDの配列（1件でも複数件でも）として再利用する。
SELECT sqlc.embed(comments), sqlc.embed(users)
FROM comments
JOIN users ON users.id = comments.user_id
WHERE comments.parent_id = ANY(sqlc.arg('parent_ids')::bigint[])
ORDER BY comments.created_at ASC;
