-- name: CreateRoadmap :one
INSERT INTO roadmaps (
    user_id, title, description, category, is_public, is_template, step_count,
    completed_step_count, progress, status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now()
) RETURNING *;

-- name: UpdateRoadmap :one
-- GORMのSave（全カラム上書き）に相当。ただしstep_count/completed_step_count/progress/
-- status/completed_atは対象外。step_count等はIncrement/Decrement系の専用クエリだけが
-- 更新する。status/completed_atはUpdateRoadmapStatusに分離した
-- （ステップ完了によるrecalcRoadmapProgressの自動遷移を、このUPDATEが読み取り時点の
-- 古いstatus/completed_atで上書きする「ロストアップデート」を防ぐため）。
UPDATE roadmaps SET
    title = $2, description = $3, category = $4, is_public = $5, is_template = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateRoadmapStatus :one
-- ユーザーによる明示的なステータス変更専用（PUT /roadmaps/:id でstatus指定時のみ呼ぶ）。
-- 汎用UpdateRoadmapと経路を分けることで、status変更を伴わない更新がrecalcRoadmapProgress
-- による自動遷移を上書きしないようにする。
UPDATE roadmaps SET
    status = $2, completed_at = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRoadmap :exec
-- roadmap_stepsはFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM roadmaps WHERE id = $1;

-- name: GetRoadmapByID :one
SELECT * FROM roadmaps WHERE id = $1;

-- name: GetRoadmapWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(roadmaps), sqlc.embed(users)
FROM roadmaps
JOIN users ON users.id = roadmaps.user_id
WHERE roadmaps.id = $1;

-- name: ListRoadmapStepsByRoadmap :many
-- GORMのPreload("Steps", order_index ASC)に相当。単一ロードマップのステップ一覧。
SELECT * FROM roadmap_steps WHERE roadmap_id = $1 ORDER BY order_index ASC;

-- name: ListRoadmapStepsByRoadmapIDs :many
-- GetTemplatesなど複数ロードマップ分のステップをまとめ取りし、Go側でグルーピングする。
SELECT * FROM roadmap_steps WHERE roadmap_id = ANY($1::bigint[]) ORDER BY roadmap_id, order_index ASC;

-- name: ListRoadmapsByUser :many
SELECT * FROM roadmaps
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRoadmapsByStatus :many
SELECT * FROM roadmaps
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: ListPublicRoadmapsWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(roadmaps), sqlc.embed(users)
FROM roadmaps
JOIN users ON users.id = roadmaps.user_id
WHERE roadmaps.is_public = true
ORDER BY roadmaps.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicRoadmaps :one
SELECT COUNT(*) FROM roadmaps WHERE is_public = true;

-- name: ListTemplateRoadmaps :many
SELECT * FROM roadmaps WHERE is_template = true ORDER BY created_at ASC;

-- name: CreateRoadmapStep :one
INSERT INTO roadmap_steps (
    roadmap_id, title, description, order_index, is_completed, completed_at, resource_url,
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now(), now()
) RETURNING *;

-- name: IncrementRoadmapStepCount :exec
UPDATE roadmaps SET step_count = step_count + 1 WHERE id = $1;

-- name: GetRoadmapStepByID :one
SELECT * FROM roadmap_steps WHERE id = $1;

-- name: UpdateRoadmapStep :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE roadmap_steps SET
    title = $2, description = $3, order_index = $4, is_completed = $5, completed_at = $6,
    resource_url = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AdjustRoadmapCompletedStepCount :exec
UPDATE roadmaps SET completed_step_count = completed_step_count + sqlc.arg('delta')::bigint
WHERE id = sqlc.arg('id');

-- name: DecrementRoadmapStepCount :exec
UPDATE roadmaps SET step_count = step_count - 1 WHERE id = $1;

-- name: DeleteRoadmapStep :exec
DELETE FROM roadmap_steps WHERE id = $1;

-- name: UpdateRoadmapProgress :exec
UPDATE roadmaps SET progress = $2 WHERE id = $1;

-- name: UpdateRoadmapProgressCompleted :exec
-- 進捗100%到達での自動完了（GORMのUpdates({"progress":...,"status":"completed","completed_at":now()})に相当）。
UPDATE roadmaps SET progress = $2, status = 'completed', completed_at = now() WHERE id = $1;

-- name: UpdateRoadmapProgressReactivated :exec
-- 100%未満へ戻ったときのアクティブ復帰（GORMのUpdates({"progress":...,"status":"active","completed_at":nil})に相当）。
UPDATE roadmaps SET progress = $2, status = 'active', completed_at = NULL WHERE id = $1;

-- name: ReorderRoadmapStep :exec
UPDATE roadmap_steps SET order_index = $3 WHERE id = $1 AND roadmap_id = $2;
