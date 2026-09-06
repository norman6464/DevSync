-- name: LockUserForStreakFreeze :exec
-- 同一ユーザーの CreateWithinLimits 同時実行を直列化するための行ロック（GORMの clause.Locking{Strength: "UPDATE"} に相当）。
SELECT id FROM users WHERE id = $1 FOR UPDATE;

-- name: CountStreakFreezesOnDate :one
SELECT COUNT(*) FROM streak_freezes WHERE user_id = $1 AND used_on = $2;

-- name: CountStreakFreezesInMonth :one
SELECT COUNT(*) FROM streak_freezes WHERE user_id = $1 AND year = $2 AND month = $3;

-- name: CreateStreakFreeze :one
INSERT INTO streak_freezes (user_id, used_on, month, year, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: ListStreakFreezesByUserAndMonth :many
SELECT * FROM streak_freezes
WHERE user_id = $1 AND year = $2 AND month = $3
ORDER BY used_on ASC;

-- name: HasStreakFreezeOnDate :one
SELECT COUNT(*) > 0 FROM streak_freezes WHERE user_id = $1 AND used_on = $2;
