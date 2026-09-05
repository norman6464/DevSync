-- name: CreateUser :one
INSERT INTO users (
    username, name, email, password, avatar_url, bio, git_hub_id, git_hub_username,
    git_hub_token, git_hub_connected, spotify_connected, spotify_token, spotify_refresh_token,
    spotify_token_expiry, zenn_username, qiita_username, at_coder_username, paiza_rank,
    skills_languages, skills_frameworks, onboarding_completed, email_weekly_report,
    email_language, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
    $19, $20, $21, $22, $23, now(), now()
) RETURNING *;

-- name: FindAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE users.id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE users.username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE users.email = $1;

-- name: GetUserByGitHubID :one
SELECT * FROM users WHERE users.git_hub_id = $1;

-- name: SearchUsers :many
SELECT * FROM users WHERE users.name ILIKE $1 OR users.email ILIKE $1 LIMIT $2;

-- name: UpdateUser :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE users SET
    username = $2, name = $3, email = $4, password = $5, avatar_url = $6, bio = $7,
    git_hub_id = $8, git_hub_username = $9, git_hub_token = $10, git_hub_connected = $11,
    spotify_connected = $12, spotify_token = $13, spotify_refresh_token = $14,
    spotify_token_expiry = $15, zenn_username = $16, qiita_username = $17,
    at_coder_username = $18, paiza_rank = $19, skills_languages = $20, skills_frameworks = $21,
    onboarding_completed = $22, email_weekly_report = $23, email_language = $24, updated_at = now()
WHERE users.id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $2 WHERE users.id = $1;

-- name: DeleteCommentLikesByUserCommentsSet :exec
-- 退会者の投稿配下のコメント、および退会者自身が他の投稿へ書いたコメントの
-- 従属行（コメントいいね）をまとめて消す。
DELETE FROM comment_likes WHERE comment_likes.comment_id IN (
    SELECT c.id FROM comments c
    WHERE c.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1) OR c.user_id = $1
);

-- name: DeleteMentionsByUserCommentsSet :exec
-- mentions自身もpost_id列を持つため、サブクエリのcommentsを明示的にエイリアス修飾しないと
-- post_idの参照先が曖昧になる（PostgreSQLの相関サブクエリの解決規則による）。
DELETE FROM mentions WHERE mentions.comment_id IN (
    SELECT c.id FROM comments c
    WHERE c.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1) OR c.user_id = $1
);

-- name: DeleteCommentsByUserPostsAndSelf :exec
DELETE FROM comments
WHERE comments.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1) OR comments.user_id = $1;

-- name: DeleteSnippetCommentsByUserPostSnippets :exec
DELETE FROM snippet_comments WHERE snippet_comments.snippet_id IN (
    SELECT cs.id FROM code_snippets cs
    WHERE cs.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1)
);

-- name: DeleteCodeSnippetsByUserPosts :exec
DELETE FROM code_snippets
WHERE code_snippets.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteLikesByUserPosts :exec
DELETE FROM likes WHERE likes.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteReactionsByUserPosts :exec
DELETE FROM reactions WHERE reactions.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteBookmarksByUserPosts :exec
DELETE FROM bookmarks WHERE bookmarks.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteBookmarkCollectionItemsByUserPosts :exec
DELETE FROM bookmark_collection_items
WHERE bookmark_collection_items.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeletePostSeriesItemsByUserPosts :exec
DELETE FROM post_series_items
WHERE post_series_items.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeletePostCollectionItemsByUserPosts :exec
DELETE FROM post_collection_items
WHERE post_collection_items.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeletePostTagsByUserPosts :exec
DELETE FROM post_tags WHERE post_tags.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeletePostPinsByUserPosts :exec
DELETE FROM post_pins WHERE post_pins.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeletePostViewsByUserPosts :exec
DELETE FROM post_views WHERE post_views.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteMentionsByUserPosts :exec
DELETE FROM mentions WHERE mentions.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteNotificationsForUserDeletion :exec
DELETE FROM notifications
WHERE notifications.user_id = $1 OR notifications.actor_id = $1
    OR notifications.post_id IN (SELECT p.id FROM posts p WHERE p.user_id = $1);

-- name: DeleteMessagesByUser :exec
DELETE FROM messages WHERE messages.sender_id = $1 OR messages.receiver_id = $1;

-- name: DeleteLikesByUser :exec
DELETE FROM likes WHERE likes.user_id = $1;

-- name: DeletePostsByUser :exec
DELETE FROM posts WHERE posts.user_id = $1;

-- name: DeleteFollowsByUser :exec
DELETE FROM follows WHERE follows.follower_id = $1 OR follows.followee_id = $1;

-- name: DeletePasswordResetTokensByUser :exec
DELETE FROM password_reset_tokens WHERE password_reset_tokens.user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE users.id = $1;
