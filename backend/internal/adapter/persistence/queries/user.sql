-- name: CreateUser :one
-- パスワードハッシュが空文字でなければ同一SQL文でuser_credentialsへ登録する
-- （GitHub認証のみのユーザーはuser_credentials行を持たない。DEVSYNC-159でusersから分離）。
WITH inserted_user AS (
    INSERT INTO users (
        username, name, email, avatar_url, bio, git_hub_id, git_hub_username,
        git_hub_token, git_hub_connected, spotify_connected, spotify_token, spotify_refresh_token,
        spotify_token_expiry, zenn_username, qiita_username, at_coder_username, paiza_rank,
        skills_languages, skills_frameworks, onboarding_completed, email_weekly_report,
        email_language, created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
        $18, $19, $20, $21, $22, now(), now()
    ) RETURNING users.*
), credential_insert AS (
    INSERT INTO user_credentials (user_id, password_hash)
    SELECT inserted_user.id, sqlc.arg('password_hash')::text
    FROM inserted_user
    WHERE sqlc.arg('password_hash')::text != ''
)
SELECT inserted_user.* FROM inserted_user;

-- name: FindAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE users.id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE users.username = $1;

-- name: GetUserByEmail :one
-- パスワードハッシュも返す（ログイン処理での照合に使うため）。user_credentials行を
-- 持たないユーザー（GitHub認証のみ）はpassword_hashがNULLになる。
SELECT sqlc.embed(users), user_credentials.password_hash
FROM users
LEFT JOIN user_credentials ON user_credentials.user_id = users.id
WHERE users.email = $1;

-- name: GetUserByGitHubID :one
-- パスワードハッシュも返す（GitHub連携済みでもパスワードを別途設定している場合があるため）。
SELECT sqlc.embed(users), user_credentials.password_hash
FROM users
LEFT JOIN user_credentials ON user_credentials.user_id = users.id
WHERE users.git_hub_id = $1;

-- name: GetUserByIDWithPassword :one
-- パスワードハッシュも返す（退会時の本人確認に使うため）。
SELECT sqlc.embed(users), user_credentials.password_hash
FROM users
LEFT JOIN user_credentials ON user_credentials.user_id = users.id
WHERE users.id = $1;

-- name: SearchUsers :many
SELECT * FROM users WHERE users.name ILIKE $1 OR users.email ILIKE $1 LIMIT $2;

-- name: UpdateUser :one
-- GORMのSave（全カラム上書き）に相当。passwordは対象外（UpdateUserPasswordを使う。
-- user_credentials側にありusersの列でもないため、そもそも対象にできない。DEVSYNC-159）。
UPDATE users SET
    username = $2, name = $3, email = $4, avatar_url = $5, bio = $6,
    git_hub_id = $7, git_hub_username = $8, git_hub_token = $9, git_hub_connected = $10,
    spotify_connected = $11, spotify_token = $12, spotify_refresh_token = $13,
    spotify_token_expiry = $14, zenn_username = $15, qiita_username = $16,
    at_coder_username = $17, paiza_rank = $18, skills_languages = $19, skills_frameworks = $20,
    onboarding_completed = $21, email_weekly_report = $22, email_language = $23, updated_at = now()
WHERE users.id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
-- user_credentials行が無ければ新規作成し、あれば上書きする
-- （GitHubのみで登録したユーザーが後からパスワードを設定する場合に備える）。
INSERT INTO user_credentials (user_id, password_hash)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- name: DeleteUser :exec
-- user_credentialsはFKのON DELETE CASCADEでDBが自動的に削除する。
DELETE FROM users WHERE users.id = $1;
