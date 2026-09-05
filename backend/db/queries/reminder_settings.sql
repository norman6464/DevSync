-- name: GetReminderSettingsByUserID :one
SELECT * FROM reminder_settings
WHERE user_id = $1;

-- name: CreateDefaultReminderSettings :many
-- 同時に複数リクエストが「不在」と判定してもuser_idの一意制約とDO NOTHINGで
-- 競合を無害化する。競合で挿入されなかった場合は0行が返る（エラーにはしない）。
INSERT INTO reminder_settings (
  user_id, enabled, frequency, notification_time, inactive_days,
  enable_web, enable_email, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
ON CONFLICT (user_id) DO NOTHING
RETURNING *;

-- name: UpdateReminderSettings :one
UPDATE reminder_settings
SET enabled = $2,
    frequency = $3,
    notification_time = $4,
    inactive_days = $5,
    enable_web = $6,
    enable_email = $7,
    last_reminded_at = $8,
    updated_at = now()
WHERE id = $1
RETURNING *;
