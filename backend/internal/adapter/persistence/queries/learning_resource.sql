-- name: CreateLearningResource :one
INSERT INTO learning_resources (
    user_id, title, description, url, category, difficulty, tags, image_url, is_public,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()
) RETURNING *;

-- name: UpdateLearningResource :one
-- GORMのSave（全カラム上書き）に相当。like_count/save_countはlearning_resources自体の
-- 列ではなくlearning_resource_metrics側（DEVSYNC-159）にあるため、ここには含まれない。
UPDATE learning_resources SET
    title = $2, description = $3, url = $4, category = $5, difficulty = $6, tags = $7,
    image_url = $8, is_public = $9, updated_at = now()
WHERE learning_resources.id = $1
RETURNING *;

-- name: DeleteLearningResource :exec
-- 依存するいいね・保存等はFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM learning_resources WHERE learning_resources.id = $1;

-- name: GetLearningResourceWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.id = $1;

-- name: ListLearningResourcesByUser :many
-- Userは含めない（移行前からの挙動。一覧系の中でこれだけPreloadしない）。
-- include_privateがfalseのときだけis_public=trueで絞り込む。
SELECT * FROM learning_resources
WHERE learning_resources.user_id = $1
    AND (sqlc.arg('include_private')::bool OR learning_resources.is_public = true)
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserVisibleLearningResources :one
-- FindByUserID の総数取得専用。include_private=falseなら公開分のみで数える
-- （全件カウントはlearning_resource_stats.sqlの既存クエリCountLearningResourcesByUserを使う）。
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.user_id = $1
    AND (sqlc.arg('include_private')::bool OR learning_resources.is_public = true);

-- name: ListPublicLearningResourcesWithUser :many
-- GORMのPreload("User")に相当。categoryとdifficultyは空文字なら絞り込まない。
-- like_countはlearning_resource_metrics側（DEVSYNC-159）。LEFT JOIN + COALESCEで0扱いにする。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
LEFT JOIN learning_resource_metrics lrm ON lrm.resource_id = learning_resources.id
WHERE learning_resources.is_public = true
    AND (sqlc.narg('category')::text IS NULL OR learning_resources.category = sqlc.narg('category'))
    AND (sqlc.narg('difficulty')::text IS NULL OR learning_resources.difficulty = sqlc.narg('difficulty'))
ORDER BY COALESCE(lrm.like_count, 0) DESC, learning_resources.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true
    AND (sqlc.narg('category')::text IS NULL OR learning_resources.category = sqlc.narg('category'))
    AND (sqlc.narg('difficulty')::text IS NULL OR learning_resources.difficulty = sqlc.narg('difficulty'));

-- name: ListLearningResourcesByDifficultyWithUser :many
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.is_public = true AND learning_resources.difficulty = $1
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLearningResourcesByDifficulty :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true AND learning_resources.difficulty = $1;

-- name: SearchLearningResourcesWithUser :many
-- like_countはlearning_resource_metrics側（DEVSYNC-159）。LEFT JOIN + COALESCEで0扱いにする。
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
LEFT JOIN learning_resource_metrics lrm ON lrm.resource_id = learning_resources.id
WHERE learning_resources.is_public = true
    AND (learning_resources.title ILIKE $1 OR learning_resources.description ILIKE $1 OR learning_resources.tags ILIKE $1)
ORDER BY COALESCE(lrm.like_count, 0) DESC, learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.is_public = true
    AND (learning_resources.title ILIKE $1 OR learning_resources.description ILIKE $1 OR learning_resources.tags ILIKE $1);

-- name: ListSavedLearningResourcesWithUser :many
SELECT sqlc.embed(learning_resources), sqlc.embed(users)
FROM learning_resources
JOIN users ON users.id = learning_resources.user_id
WHERE learning_resources.id IN (
    SELECT rs.resource_id FROM resource_saves rs WHERE rs.user_id = $1
)
ORDER BY learning_resources.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSavedLearningResources :one
SELECT COUNT(*) FROM learning_resources
WHERE learning_resources.id IN (
    SELECT rs.resource_id FROM resource_saves rs WHERE rs.user_id = $1
);

-- name: CreateResourceLike :exec
-- resource_likesへのINSERTとlearning_resource_metrics.like_countの加算を
-- 同一SQL文で行う（DEVSYNC-159）。CTEの出力列を別名にし、sqlcの列名衝突による
-- 誤検出を避ける。
WITH inserted AS (
    INSERT INTO resource_likes (user_id, resource_id, created_at) VALUES ($1, $2, now())
    RETURNING resource_likes.resource_id AS liked_resource_id
)
INSERT INTO learning_resource_metrics (resource_id, like_count)
SELECT liked_resource_id, 1 FROM inserted
ON CONFLICT (resource_id) DO UPDATE SET like_count = learning_resource_metrics.like_count + 1;

-- name: DeleteResourceLike :execrows
-- resource_likesの削除とlearning_resource_metrics.like_countの減算を同一SQL文で行う。
-- 実際に削除できた（＝deletedに1行ある）ときだけmetricsを更新するため、
-- rowsAffectedは従来どおり「実際にいいねを取り消せたか」を表す。
WITH deleted AS (
    DELETE FROM resource_likes WHERE resource_likes.user_id = $1 AND resource_likes.resource_id = $2
    RETURNING resource_likes.resource_id AS liked_resource_id
)
INSERT INTO learning_resource_metrics (resource_id, like_count)
SELECT liked_resource_id, 0 FROM deleted
ON CONFLICT (resource_id) DO UPDATE SET like_count = GREATEST(learning_resource_metrics.like_count - 1, 0);

-- name: CountResourceLike :one
SELECT COUNT(*) FROM resource_likes
WHERE resource_likes.user_id = $1 AND resource_likes.resource_id = $2;

-- name: CreateResourceSave :exec
-- resource_savesへのINSERTとlearning_resource_metrics.save_countの加算を同一SQL文で行う。
WITH inserted AS (
    INSERT INTO resource_saves (user_id, resource_id, created_at) VALUES ($1, $2, now())
    RETURNING resource_saves.resource_id AS saved_resource_id
)
INSERT INTO learning_resource_metrics (resource_id, save_count)
SELECT saved_resource_id, 1 FROM inserted
ON CONFLICT (resource_id) DO UPDATE SET save_count = learning_resource_metrics.save_count + 1;

-- name: DeleteResourceSave :execrows
-- resource_savesの削除とlearning_resource_metrics.save_countの減算を同一SQL文で行う。
WITH deleted AS (
    DELETE FROM resource_saves WHERE resource_saves.user_id = $1 AND resource_saves.resource_id = $2
    RETURNING resource_saves.resource_id AS saved_resource_id
)
INSERT INTO learning_resource_metrics (resource_id, save_count)
SELECT saved_resource_id, 0 FROM deleted
ON CONFLICT (resource_id) DO UPDATE SET save_count = GREATEST(learning_resource_metrics.save_count - 1, 0);

-- name: CountResourceSave :one
SELECT COUNT(*) FROM resource_saves
WHERE resource_saves.user_id = $1 AND resource_saves.resource_id = $2;
