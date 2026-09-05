-- name: GetNotificationSettingsByUserID :one
SELECT * FROM notification_settings
WHERE user_id = $1;

-- name: CreateNotificationSettings :one
INSERT INTO notification_settings (
  user_id, enable_likes, enable_comments, enable_follows, enable_messages,
  enable_mentions, enable_web_push, enable_email, enable_sound,
  created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
RETURNING *;

-- name: UpdateNotificationSettings :one
UPDATE notification_settings
SET enable_likes = $2,
    enable_comments = $3,
    enable_follows = $4,
    enable_messages = $5,
    enable_mentions = $6,
    enable_web_push = $7,
    enable_email = $8,
    enable_sound = $9,
    updated_at = now()
WHERE id = $1
RETURNING *;
