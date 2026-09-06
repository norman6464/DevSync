package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNotifications_TypeMustBeKnownVerb は、notifications.typeがnotification_verbsに
// 存在しないコードを拒否することを検証する（fk_notifications_type、DEVSYNC-159）。
func TestNotifications_TypeMustBeKnownVerb(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('nv_user', 'nv_user', 'nv_user@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, read, created_at)
		VALUES ($1, 'not-a-real-verb', $1, false, now())
	`, userID)
	require.Error(t, err, "notification_verbsに存在しないtypeはFK違反で拒否されるべき")

	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, read, created_at)
		VALUES ($1, 'follow', $1, false, now())
	`, userID)
	assert.NoError(t, err, "notification_verbsに存在するtypeは受け付けられるべき")
}

// TestNotifications_ExclusiveTargetCheck は、post_id/question_id/badge_idのうち
// 2つ以上が同時に非NULLの通知を拒否することを検証する
// （ck_notifications_exclusive_target、DEVSYNC-159）。
func TestNotifications_ExclusiveTargetCheck(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('nv_exclusive_user', 'nv_exclusive_user', 'nv_exclusive_user@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	var postID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, userID).Scan(&postID)
	require.NoError(t, err)

	badgeID := "first-commit"
	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, post_id, badge_id, read, created_at)
		VALUES ($1, 'post', $1, $2, $3, false, now())
	`, userID, postID, badgeID)
	require.Error(t, err, "post_idとbadge_idを同時に持つ通知は排他アークCHECKで拒否されるべき")

	_, err = pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, actor_id, post_id, read, created_at)
		VALUES ($1, 'post', $1, $2, false, now())
	`, userID, postID)
	assert.NoError(t, err, "target列が1つだけの通知は受け付けられるべき")
}

// TestSeedNotificationVerbs_Idempotent は、既知コードのシードを複数回実行しても
// エラーにならず重複行が作られないことを検証する。
func TestSeedNotificationVerbs_Idempotent(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)
	repo := NewNotificationVerbRepository(q)

	codes := []string{"post", "message", "like", "comment", "follow", "answer", "badge", "level_up", "mention"}
	require.NoError(t, repo.SeedKnownVerbs(ctx, codes))
	require.NoError(t, repo.SeedKnownVerbs(ctx, codes))

	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM notification_verbs WHERE code = ANY($1::text[])`, codes).Scan(&count))
	assert.EqualValues(t, len(codes), count)
}
