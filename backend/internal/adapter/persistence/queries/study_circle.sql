-- name: CreateStudyCircle :one
INSERT INTO study_circles (name, topic, description, owner_id, max_members, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *;

-- name: CreateStudyCircleMember :exec
INSERT INTO study_circle_members (circle_id, user_id, role, joined_at)
VALUES ($1, $2, $3, $4);

-- name: GetStudyCircleWithOwnerByID :one
-- GORMのPreload("Owner")に相当。owner_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(study_circles), sqlc.embed(users)
FROM study_circles
JOIN users ON users.id = study_circles.owner_id
WHERE study_circles.id = $1;

-- name: ListStudyCircleStepsByCircle :many
-- GORMのPreload("Steps", order_index ASC)に相当。
SELECT * FROM study_circle_steps WHERE circle_id = $1 ORDER BY order_index ASC;

-- name: ListStudyCircleMembersWithUserByCircle :many
-- GORMのPreload("Members").Preload("Members.User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(study_circle_members), sqlc.embed(users)
FROM study_circle_members
JOIN users ON users.id = study_circle_members.user_id
WHERE study_circle_members.circle_id = $1;

-- name: UpdateStudyCircle :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE study_circles SET
    name = $2, topic = $3, description = $4, max_members = $5, status = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteStudyCircleCheckinsByCircle :exec
DELETE FROM study_circle_checkins WHERE circle_id = $1;

-- name: DeleteStudyCircleMemberProgressByCircle :exec
DELETE FROM study_circle_member_progresses WHERE circle_id = $1;

-- name: DeleteStudyCircleStepsByCircle :exec
DELETE FROM study_circle_steps WHERE circle_id = $1;

-- name: DeleteStudyCircleMembersByCircle :exec
DELETE FROM study_circle_members WHERE circle_id = $1;

-- name: DeleteStudyCircle :exec
DELETE FROM study_circles WHERE id = $1;

-- name: ListStudyCirclesByUserAndStatus :many
-- GORMのPreload("Owner").Preload("Members").Preload("Members.User")に相当。
SELECT sqlc.embed(study_circles), sqlc.embed(users)
FROM study_circles
JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id
JOIN users ON users.id = study_circles.owner_id
WHERE study_circle_members.user_id = $1 AND study_circles.status = $2
ORDER BY study_circles.updated_at DESC;

-- name: ListStudyCirclesByUser :many
-- GORMのPreload("Owner").Preload("Members").Preload("Members.User")に相当。
SELECT sqlc.embed(study_circles), sqlc.embed(users)
FROM study_circles
JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id
JOIN users ON users.id = study_circles.owner_id
WHERE study_circle_members.user_id = $1
ORDER BY study_circles.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CountStudyCirclesByUserMembership :one
SELECT COUNT(*) FROM study_circles
JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id
WHERE study_circle_members.user_id = $1;

-- name: SearchStudyCircles :many
-- GORMのPreload("Owner").Preload("Members").Preload("Members.User")に相当。
SELECT sqlc.embed(study_circles), sqlc.embed(users)
FROM study_circles
JOIN users ON users.id = study_circles.owner_id
WHERE study_circles.name ILIKE $1 OR study_circles.topic ILIKE $1 OR study_circles.description ILIKE $1
ORDER BY study_circles.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchStudyCircles :one
SELECT COUNT(*) FROM study_circles
WHERE name ILIKE $1 OR topic ILIKE $1 OR description ILIKE $1;

-- name: LockStudyCircleForMemberLimit :one
-- 「数える→追加する」を同時実行しても上限を超えないようにするための行ロック
-- （GORMの clause.Locking{Strength: "UPDATE"} に相当）。
SELECT max_members FROM study_circles WHERE id = $1 FOR UPDATE;

-- name: DeleteStudyCircleMember :exec
DELETE FROM study_circle_members WHERE circle_id = $1 AND user_id = $2;

-- name: UpdateStudyCircleMemberRole :exec
UPDATE study_circle_members SET role = $3 WHERE circle_id = $1 AND user_id = $2;

-- name: CountStudyCircleMembership :one
SELECT COUNT(*) FROM study_circle_members WHERE circle_id = $1 AND user_id = $2;

-- name: CountStudyCircleMembershipsByUser :one
SELECT COUNT(*) FROM study_circle_members WHERE user_id = $1;

-- name: CreateStudyCircleStep :one
INSERT INTO study_circle_steps (circle_id, title, description, order_index, resource_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *;

-- name: UpdateStudyCircleStep :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE study_circle_steps SET
    title = $2, description = $3, order_index = $4, resource_url = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteStudyCircleStep :exec
DELETE FROM study_circle_steps WHERE id = $1;

-- name: GetStudyCircleStepByID :one
SELECT * FROM study_circle_steps WHERE id = $1;

-- name: ReorderStudyCircleStep :exec
UPDATE study_circle_steps SET order_index = $3 WHERE id = $1 AND circle_id = $2;

-- name: UpsertStudyCircleMemberProgress :one
INSERT INTO study_circle_member_progresses (circle_id, step_id, user_id, is_completed, completed_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (circle_id, step_id, user_id) DO UPDATE SET
    is_completed = EXCLUDED.is_completed,
    completed_at = EXCLUDED.completed_at
RETURNING *;

-- name: ListStudyCircleMemberProgressWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(study_circle_member_progresses), sqlc.embed(users)
FROM study_circle_member_progresses
JOIN users ON users.id = study_circle_member_progresses.user_id
WHERE study_circle_member_progresses.circle_id = $1;

-- name: CreateStudyCircleCheckin :one
INSERT INTO study_circle_checkins (circle_id, user_id, checked_on, content, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: ListStudyCircleCheckinsWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(study_circle_checkins), sqlc.embed(users)
FROM study_circle_checkins
JOIN users ON users.id = study_circle_checkins.user_id
WHERE study_circle_checkins.circle_id = $1
ORDER BY study_circle_checkins.created_at DESC;

-- name: CountStudyCircleCheckinsToday :one
SELECT COUNT(*) FROM study_circle_checkins
WHERE circle_id = $1 AND user_id = $2 AND checked_on = $3;

-- name: ListStudyCircleCheckinDatesByCircle :many
-- GetStreakRankingのN+1回避のため、サークル分のチェックインをまとめて取得し、
-- Go側でuser_idごとにグルーピングする。checked_onの降順はcalculateCheckinStreakの前提。
SELECT user_id, checked_on FROM study_circle_checkins
WHERE circle_id = $1
ORDER BY checked_on DESC;
