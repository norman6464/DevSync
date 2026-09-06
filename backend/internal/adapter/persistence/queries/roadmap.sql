-- name: CreateRoadmap :one
INSERT INTO roadmaps (
    user_id, title, description, category, is_public, is_template, status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now(), now()
) RETURNING *;

-- name: UpdateRoadmap :one
-- GORMのSave（全カラム上書き）に相当。ただしstatus/completed_atは対象外
-- （UpdateRoadmapStatusに分離済み。ステップ完了によるステータス自動遷移を、この
-- UPDATEが読み取り時点の古いstatus/completed_atで上書きする「ロストアップデート」を
-- 防ぐため）。step_count/completed_step_count/progressはroadmaps自体の列ではなく
-- roadmap_metrics側（DEVSYNC-159。progressはSELECT側の算出式）にあるため対象外。
UPDATE roadmaps SET
    title = $2, description = $3, category = $4, is_public = $5, is_template = $6,
    updated_at = now()
WHERE roadmaps.id = $1
RETURNING *;

-- name: UpdateRoadmapStatus :one
-- ユーザーによる明示的なステータス変更専用（PUT /roadmaps/:id でstatus指定時のみ呼ぶ）。
-- 汎用UpdateRoadmapと経路を分けることで、status変更を伴わない更新がステップ完了による
-- ステータス自動遷移を上書きしないようにする。
UPDATE roadmaps SET
    status = $2, completed_at = $3, updated_at = now()
WHERE roadmaps.id = $1
RETURNING *;

-- name: DeleteRoadmap :exec
-- roadmap_stepsはFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM roadmaps WHERE roadmaps.id = $1;

-- name: GetRoadmapByID :one
SELECT * FROM roadmaps WHERE roadmaps.id = $1;

-- name: GetRoadmapWithUserByID :one
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(roadmaps), sqlc.embed(users)
FROM roadmaps
JOIN users ON users.id = roadmaps.user_id
WHERE roadmaps.id = $1;

-- name: ListRoadmapStepsByRoadmap :many
-- GORMのPreload("Steps", order_index ASC)に相当。単一ロードマップのステップ一覧。
SELECT * FROM roadmap_steps WHERE roadmap_steps.roadmap_id = $1 ORDER BY roadmap_steps.order_index ASC;

-- name: ListRoadmapStepsByRoadmapIDs :many
-- GetTemplatesなど複数ロードマップ分のステップをまとめ取りし、Go側でグルーピングする。
SELECT * FROM roadmap_steps WHERE roadmap_steps.roadmap_id = ANY($1::bigint[]) ORDER BY roadmap_steps.roadmap_id, roadmap_steps.order_index ASC;

-- name: ListRoadmapsByUser :many
SELECT * FROM roadmaps
WHERE roadmaps.user_id = $1
ORDER BY roadmaps.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRoadmapsByStatus :many
SELECT * FROM roadmaps
WHERE roadmaps.user_id = $1 AND roadmaps.status = $2
ORDER BY roadmaps.created_at DESC;

-- name: ListPublicRoadmapsWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(roadmaps), sqlc.embed(users)
FROM roadmaps
JOIN users ON users.id = roadmaps.user_id
WHERE roadmaps.is_public = true
ORDER BY roadmaps.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublicRoadmaps :one
SELECT COUNT(*) FROM roadmaps WHERE roadmaps.is_public = true;

-- name: ListTemplateRoadmaps :many
SELECT * FROM roadmaps WHERE roadmaps.is_template = true ORDER BY roadmaps.created_at ASC;

-- name: CreateRoadmapStep :one
-- roadmap_stepsへのINSERTとroadmap_metrics.step_countの加算を同一SQL文で行う
-- （DEVSYNC-159）。metrics_upsertの結果自体は読まないが、データ変更を伴うCTEは
-- 参照の有無に関わらず必ず実行されるため、確実にstep_countへ反映される。
WITH inserted_step AS (
    INSERT INTO roadmap_steps (
        roadmap_id, title, description, order_index, is_completed, completed_at, resource_url,
        created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, now(), now()
    ) RETURNING roadmap_steps.*
), metrics_upsert AS (
    INSERT INTO roadmap_metrics (roadmap_id, step_count)
    SELECT inserted_step.roadmap_id, 1 FROM inserted_step
    ON CONFLICT (roadmap_id) DO UPDATE SET step_count = roadmap_metrics.step_count + 1
)
SELECT inserted_step.* FROM inserted_step;

-- name: GetRoadmapStepByID :one
SELECT * FROM roadmap_steps WHERE roadmap_steps.id = $1;

-- name: UpdateRoadmapStep :one
-- GORMのSave（全カラム上書き）に相当。completed_deltaが0でなければ
-- roadmap_metrics.completed_step_countも同一SQL文で加減算する（0未満にはしない）。
-- completed_deltaは呼び出し側（Go）が更新前後のis_completedを比較して算出する。
WITH updated_step AS (
    UPDATE roadmap_steps SET
        title = $2, description = $3, order_index = $4, is_completed = $5, completed_at = $6,
        resource_url = $7, updated_at = now()
    WHERE roadmap_steps.id = $1
    RETURNING roadmap_steps.*
), metrics_upsert AS (
    INSERT INTO roadmap_metrics (roadmap_id, completed_step_count)
    SELECT updated_step.roadmap_id, sqlc.arg('completed_delta')::bigint
    FROM updated_step
    WHERE sqlc.arg('completed_delta')::bigint != 0
    ON CONFLICT (roadmap_id) DO UPDATE SET
        completed_step_count = GREATEST(roadmap_metrics.completed_step_count + EXCLUDED.completed_step_count, 0)
)
SELECT updated_step.* FROM updated_step;

-- name: DeleteRoadmapStep :exec
-- roadmap_stepsの削除とroadmap_metrics.step_count/completed_step_countの減算を
-- 同一SQL文で行う（DEVSYNC-159）。completed_deltaは削除されたステップが完了済み
-- だったかどうかを呼び出し側（Go）が判定して渡す（0 または -1）。
WITH deleted_step AS (
    DELETE FROM roadmap_steps WHERE roadmap_steps.id = $1
    RETURNING roadmap_steps.roadmap_id AS deleted_roadmap_id
)
INSERT INTO roadmap_metrics (roadmap_id, step_count, completed_step_count)
SELECT deleted_roadmap_id, -1, sqlc.arg('completed_delta')::bigint FROM deleted_step
ON CONFLICT (roadmap_id) DO UPDATE SET
    step_count = GREATEST(roadmap_metrics.step_count - 1, 0),
    completed_step_count = GREATEST(roadmap_metrics.completed_step_count + EXCLUDED.completed_step_count, 0);

-- name: ReorderRoadmapStep :exec
UPDATE roadmap_steps SET order_index = $3 WHERE roadmap_steps.id = $1 AND roadmap_steps.roadmap_id = $2;
