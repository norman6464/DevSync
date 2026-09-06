-- name: SeedNotificationVerbs :exec
-- notifications.typeのFK制約（fk_notifications_type、DEVSYNC-159）を満たすため、
-- アプリ起動時に既知の通知種別コードをまとめて登録する（冪等。既存分はスキップ）。
INSERT INTO notification_verbs (code)
SELECT unnest($1::text[])
ON CONFLICT (code) DO NOTHING;
