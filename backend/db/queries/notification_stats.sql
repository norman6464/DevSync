-- name: CountNotificationsByUser :one
SELECT count(*) FROM notifications
WHERE user_id = $1;

-- name: CountUnreadNotificationsByUser :one
SELECT count(*) FROM notifications
WHERE user_id = $1 AND read = false;

-- name: CountNotificationsByUserSince :one
SELECT count(*) FROM notifications
WHERE user_id = $1 AND created_at >= $2;
