-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at, used, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: GetPasswordResetTokenByToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1;

-- name: MarkPasswordResetTokenAsUsed :exec
UPDATE password_reset_tokens
SET used = true
WHERE id = $1;

-- name: InvalidateUserPasswordResetTokens :exec
UPDATE password_reset_tokens
SET used = true
WHERE user_id = $1 AND used = false;
