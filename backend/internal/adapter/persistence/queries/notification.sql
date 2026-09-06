-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, actor_id, post_id, question_id, badge_id, read, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now())
RETURNING *;

-- name: FindFollowerIDsByFollowee :many
-- 指定ユーザーをフォローしているユーザーのIDを返す。
SELECT follower_id FROM follows WHERE followee_id = $1;

-- name: FindNotificationsByUserID :many
-- GORMのPreload("Actor").Preload("Post").Preload("Question")に相当。
-- actor_idはNOT NULLのためINNER JOINでよいが、post_id/question_idはNULL許容のためLEFT JOIN。
-- LEFT JOIN側のpost/questionはsqlc.embedを使わず個別カラム選択にすることで、
-- テーブル自体のスキーマではなくJOINコンテキストからNULL許容性を正しく推論させる。
SELECT
    sqlc.embed(notifications),
    sqlc.embed(actor),
    posts.id AS post_id_2,
    posts.user_id AS post_user_id,
    posts.title AS post_title,
    posts.content AS post_content,
    posts.image_urls AS post_image_urls,
    posts.is_draft AS post_is_draft,
    posts.like_count AS post_like_count,
    posts.comment_count AS post_comment_count,
    (SELECT COUNT(*) FROM bookmarks WHERE bookmarks.post_id = posts.id) AS post_bookmark_count,
    posts.view_count AS post_view_count,
    posts.estimated_read_time AS post_estimated_read_time,
    posts.scheduled_at AS post_scheduled_at,
    posts.created_at AS post_created_at,
    posts.updated_at AS post_updated_at,
    questions.id AS question_id_2,
    questions.user_id AS question_user_id,
    questions.title AS question_title,
    questions.body AS question_body,
    questions.tags AS question_tags,
    questions.vote_count AS question_vote_count,
    questions.answer_count AS question_answer_count,
    questions.is_solved AS question_is_solved,
    questions.created_at AS question_created_at,
    questions.updated_at AS question_updated_at
FROM notifications
JOIN users actor ON actor.id = notifications.actor_id
LEFT JOIN posts ON posts.id = notifications.post_id
LEFT JOIN questions ON questions.id = notifications.question_id
WHERE notifications.user_id = $1
    AND (sqlc.narg('notification_type')::text IS NULL OR notifications.type = sqlc.narg('notification_type'))
ORDER BY notifications.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountNotificationsByUserID :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1
    AND (sqlc.narg('notification_type')::text IS NULL OR type = sqlc.narg('notification_type'));

-- name: CountUnreadNotifications :one
SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false;

-- name: MarkNotificationAsRead :exec
UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications SET read = true WHERE user_id = $1 AND read = false;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1 AND user_id = $2;
