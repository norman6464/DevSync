-- name: CountUserActivitiesByUser :one
SELECT count(*) FROM user_activities
WHERE user_id = $1;

-- name: CountUserActivitiesByUserAndType :one
SELECT count(*) FROM user_activities
WHERE user_id = $1 AND activity_type = $2;

-- name: ListUserActivitiesByUser :many
SELECT * FROM user_activities
WHERE user_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3;

-- name: ListUserActivitiesByUserAndType :many
SELECT * FROM user_activities
WHERE user_id = $1 AND activity_type = $2
ORDER BY created_at DESC, id DESC
LIMIT $3 OFFSET $4;
