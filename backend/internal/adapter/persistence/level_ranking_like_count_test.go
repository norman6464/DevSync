package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSumPostLikesReceivedByUser_ReadsFromPostMetrics は、SumPostLikesReceivedByUserが
// post_metrics側のlike_count（DEVSYNC-159で分離済み）を正しく合算することを検証する。
// posts自体にはlike_count列が無いため、この検証が無いとJOIN先を誤って実行時エラーに
// なる退行（DEVSYNC-162）に気づけない。
func TestSumPostLikesReceivedByUser_ReadsFromPostMetrics(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, likerID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "lr_author"},
		{&likerID, "lr_liker"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, authorID, likerID)
	})

	var post1ID, post2ID int64
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title1', 'content1', false, now(), now())
		RETURNING id
	`, authorID).Scan(&post1ID))
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title2', 'content2', false, now(), now())
		RETURNING id
	`, authorID).Scan(&post2ID))

	// 投稿がまだ1件もいいねされていない状態でも0を返すこと（post_metrics行が無い場合）。
	total, err := q.SumPostLikesReceivedByUser(ctx, authorID)
	require.NoError(t, err)
	assert.EqualValues(t, 0, total)

	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: likerID, PostID: post1ID}))
	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: likerID, PostID: post2ID}))

	total, err = q.SumPostLikesReceivedByUser(ctx, authorID)
	require.NoError(t, err)
	assert.EqualValues(t, 2, total, "2つの投稿への合計いいね数が合算されるべき")
}

// TestGetLevelRanking_IncludesPostLikesWithoutError は、GetLevelRankingがlike_countの
// JOINを含めて実行時エラーなく完走し、いいねを含めたスコアを返すことを検証する
// （DEVSYNC-162: post_metrics移行の取り残しでランキング機能が実行時エラーになっていた）。
func TestGetLevelRanking_IncludesPostLikesWithoutError(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, likerID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "lr_rank_author"},
		{&likerID, "lr_rank_liker"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, authorID, likerID)
	})

	var postID int64
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, authorID).Scan(&postID))
	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: likerID, PostID: postID}))

	rows, err := q.GetLevelRanking(ctx)
	require.NoError(t, err, "post_metricsへのJOINが壊れていると実行時エラーになる")

	var authorScore int32
	found := false
	for _, row := range rows {
		if row.UserID == authorID {
			authorScore = row.Score
			found = true
		}
	}
	require.True(t, found, "投稿+いいねがあるユーザーはスコア0で足切りされずランキングに現れるべき")
	assert.EqualValues(t, 30+3, authorScore, "投稿1件(xp30)+いいね1件(xp3)のスコアになるべき")
}
