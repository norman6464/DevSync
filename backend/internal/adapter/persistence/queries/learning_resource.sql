-- name: CreateLearningResource :one
INSERT INTO learning_resources (
    user_id, title, description, url, category, difficulty, tags, image_url, is_public,
    like_count, save_count, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now(), now()
) RETURNING *;

-- name: UpdateLearningResource :one
-- GORMのSave（全カラム上書き）に相当。learning_resourcesは論理削除があるため、
-- GORMが自動付与するdeleted_at IS NULLスコープをUPDATEにも明示する。
UPDATE learning_resources SET
    title = $2, description = $3, url = $4, category = $5, difficulty = $6, tags = $7,
    image_url = $8, is_public = $9, like_count = $10, save_count = $11, updated_at = now()
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL
RETURNING *;

-- name: DeleteLearningResource :exec
-- GORMのDelete（論理削除）に相当。
UPDATE learning_resources SET deleted_at = now() WHERE learning_resources.id = $1;

-- name: GetLearningResourceWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL;

-- name: ListLearningResourcesByUser :many
-- Userは含めない（移行前からの挙動。一覧系の中でこれだけPreloadしない）。
-- include_privateがfalseのときだけis_public=trueで絞り込む。
SELECT * FROM learning_resources
WHERE learning_resources.user_id = $1
    AND (sqlc.arg('include_private')::bool OR learning_resources.is_public = true)
    AND learning_resources.deleted_at IS NULL
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserVisibleLearningResources :one
-- FindByUserID の総数取得専用。include_private=falseなら公開分のみで数える
-- （全件カウントはlearning_resource_stats.sqlの既存クエリCountLearningResourcesByUserを使う）。
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.user_id = $1
    AND (sqlc.arg('include_private')::bool OR learning_resources.is_public = true)
    AND learning_resources.deleted_at IS NULL;

-- name: ListPublicLearningResourcesWithUser :many
-- GORMのPreload("User")に相当。categoryとdifficultyは空文字なら絞り込まない。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.is_public = true
    AND (sqlc.narg('category')::text IS NULL OR learning_resources.category = sqlc.narg('category'))
    AND (sqlc.narg('difficulty')::text IS NULL OR learning_resources.difficulty = sqlc.narg('difficulty'))
    AND learning_resources.deleted_at IS NULL
ORDER BY learning_resources.like_count DESC, learning_resources.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true
    AND (sqlc.narg('category')::text IS NULL OR learning_resources.category = sqlc.narg('category'))
    AND (sqlc.narg('difficulty')::text IS NULL OR learning_resources.difficulty = sqlc.narg('difficulty'))
    AND learning_resources.deleted_at IS NULL;

-- name: ListLearningResourcesByDifficultyWithUser :many
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.is_public = true AND learning_resources.difficulty = $1
    AND learning_resources.deleted_at IS NULL
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLearningResourcesByDifficulty :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true AND learning_resources.difficulty = $1
    AND learning_resources.deleted_at IS NULL;

-- name: SearchLearningResourcesWithUser :many
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.is_public = true
    AND (learning_resources.title ILIKE $1 OR learning_resources.description ILIKE $1 OR learning_resources.tags ILIKE $1)
    AND learning_resources.deleted_at IS NULL
ORDER BY learning_resources.like_count DESC, learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true
    AND (learning_resources.title ILIKE $1 OR learning_resources.description ILIKE $1 OR learning_resources.tags ILIKE $1)
    AND learning_resources.deleted_at IS NULL;

-- name: ListSavedLearningResourcesWithUser :many
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.id IN (
    SELECT rs.resource_id FROM resource_saves rs WHERE rs.user_id = $1
) AND learning_resources.deleted_at IS NULL
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSavedLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.id IN (
    SELECT rs.resource_id FROM resource_saves rs WHERE rs.user_id = $1
) AND learning_resources.deleted_at IS NULL;

-- name: CreateResourceLike :exec
INSERT INTO resource_likes (user_id, resource_id, created_at) VALUES ($1, $2, now());

-- name: DeleteResourceLike :exec
DELETE FROM resource_likes
WHERE resource_likes.user_id = $1 AND resource_likes.resource_id = $2;

-- name: CountResourceLike :one
SELECT COUNT(*) FROM resource_likes
WHERE resource_likes.user_id = $1 AND resource_likes.resource_id = $2;

-- name: IncrementResourceLikeCount :exec
-- deleted_at IS NULLを明示（GORMは論理削除モデルへのUPDATEにも自動でこのスコープを付与するため、
-- Like/Unlikeがトランザクションで括られていない今の実装では、この条件がないと
-- 削除確定後に届いた更新が論理削除済みの行を書き換えてしまう）。
UPDATE learning_resources SET like_count = like_count + 1
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL;

-- name: DecrementResourceLikeCountFloored :exec
-- 0未満にはしない（GORMのGREATEST(like_count - 1, 0)に相当）。deleted_at IS NULLの理由は
-- IncrementResourceLikeCountと同じ。
UPDATE learning_resources SET like_count = GREATEST(like_count - 1, 0)
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL;

-- name: CreateResourceSave :exec
INSERT INTO resource_saves (user_id, resource_id, created_at) VALUES ($1, $2, now());

-- name: DeleteResourceSave :exec
DELETE FROM resource_saves
WHERE resource_saves.user_id = $1 AND resource_saves.resource_id = $2;

-- name: CountResourceSave :one
SELECT COUNT(*) FROM resource_saves
WHERE resource_saves.user_id = $1 AND resource_saves.resource_id = $2;

-- name: IncrementResourceSaveCount :exec
-- deleted_at IS NULLを明示する理由はIncrementResourceLikeCountと同じ。
UPDATE learning_resources SET save_count = save_count + 1
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL;

-- name: DecrementResourceSaveCountFloored :exec
-- 0未満にはしない（GORMのGREATEST(save_count - 1, 0)に相当）。deleted_at IS NULLの理由は
-- IncrementResourceLikeCountと同じ。
UPDATE learning_resources SET save_count = GREATEST(save_count - 1, 0)
WHERE learning_resources.id = $1 AND learning_resources.deleted_at IS NULL;
