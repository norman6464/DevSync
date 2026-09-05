-- name: GetWidgetSettingsByUserID :one
SELECT * FROM widget_settings
WHERE user_id = $1;

-- name: UpsertWidgetSettings :exec
INSERT INTO widget_settings (user_id, settings, created_at, updated_at)
VALUES ($1, $2, now(), now())
ON CONFLICT (user_id) DO UPDATE
SET settings = EXCLUDED.settings,
    updated_at = now();
