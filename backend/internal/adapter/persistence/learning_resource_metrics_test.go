package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLearningResourceMetrics_LikeSaveFlow_UpdatesInSameStatement は、いいね・保存の
// 作成/削除が同一SQL文でlearning_resource_metricsを正しく増減させることを検証する
// （DEVSYNC-159）。learning_resource_metrics行が最初は存在しないリソースから始め、
// 遅延生成される様子も含めて確認する。
func TestLearningResourceMetrics_LikeSaveFlow_UpdatesInSameStatement(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, otherID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "lrm_author"},
		{&otherID, "lrm_other"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, authorID, otherID)
	})

	var resourceID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO learning_resources (user_id, title, category, is_public, created_at, updated_at)
		VALUES ($1, 'title', 'article', true, now(), now())
		RETURNING id
	`, authorID).Scan(&resourceID)
	require.NoError(t, err)

	assertMetrics := func(wantLike, wantSave int64) {
		t.Helper()
		var like, save int64
		row := pool.QueryRow(ctx, `
			SELECT COALESCE(like_count, 0), COALESCE(save_count, 0)
			FROM learning_resource_metrics WHERE resource_id = $1
		`, resourceID)
		err := row.Scan(&like, &save)
		if err != nil {
			// learning_resource_metrics行がまだ無い（遅延生成前）場合は全て0として扱う。
			like, save = 0, 0
		}
		assert.Equal(t, wantLike, like, "like_count")
		assert.Equal(t, wantSave, save, "save_count")
	}

	// 最初はlearning_resource_metrics行が存在しない。
	assertMetrics(0, 0)

	// いいね: 作成でlearning_resource_metrics行が遅延生成され1になる。
	require.NoError(t, q.CreateResourceLike(ctx, sqlcgen.CreateResourceLikeParams{UserID: otherID, ResourceID: resourceID}))
	assertMetrics(1, 0)

	// いいね取り消し: 0未満にはならない。
	rowsAffected, err := q.DeleteResourceLike(ctx, sqlcgen.DeleteResourceLikeParams{UserID: otherID, ResourceID: resourceID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)
	assertMetrics(0, 0)

	// 既に取り消し済みのいいねをもう一度取り消しても0未満にならない（rowsAffected=0）。
	rowsAffected, err = q.DeleteResourceLike(ctx, sqlcgen.DeleteResourceLikeParams{UserID: otherID, ResourceID: resourceID})
	require.NoError(t, err)
	assert.EqualValues(t, 0, rowsAffected)
	assertMetrics(0, 0)

	// 保存: 作成でsave_countが増える。
	require.NoError(t, q.CreateResourceSave(ctx, sqlcgen.CreateResourceSaveParams{UserID: otherID, ResourceID: resourceID}))
	assertMetrics(0, 1)

	// 保存取り消しでsave_countが減る。
	rowsAffected, err = q.DeleteResourceSave(ctx, sqlcgen.DeleteResourceSaveParams{UserID: otherID, ResourceID: resourceID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)
	assertMetrics(0, 0)
}

// TestReconcileAllLearningResourceMetrics_CorrectsDrift は、CASCADE削除等で
// learning_resource_metricsが実件数からずれたケースを意図的に作り、reconcileジョブが
// 正しく補正することを検証する（DEVSYNC-159のチケット本文で明示されているテスト要件）。
func TestReconcileAllLearningResourceMetrics_CorrectsDrift(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)
	metricsRepo := NewLearningResourceMetricsRepository(q)

	var authorID, liker1ID, liker2ID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "lrm_reconcile_author"},
		{&liker1ID, "lrm_reconcile_liker1"},
		{&liker2ID, "lrm_reconcile_liker2"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2, $3)`, authorID, liker1ID, liker2ID)
	})

	var resourceID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO learning_resources (user_id, title, category, is_public, created_at, updated_at)
		VALUES ($1, 'title', 'article', true, now(), now())
		RETURNING id
	`, authorID).Scan(&resourceID)
	require.NoError(t, err)

	require.NoError(t, q.CreateResourceLike(ctx, sqlcgen.CreateResourceLikeParams{UserID: liker1ID, ResourceID: resourceID}))
	require.NoError(t, q.CreateResourceLike(ctx, sqlcgen.CreateResourceLikeParams{UserID: liker2ID, ResourceID: resourceID}))

	var likeCount int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM learning_resource_metrics WHERE resource_id = $1`, resourceID).Scan(&likeCount))
	require.EqualValues(t, 2, likeCount)

	// 意図的にlearning_resource_metricsをずらす（ユーザー削除カスケードでresource_likes
	// が消えたのにlearning_resource_metricsだけ古い値のまま残る、という実際に起こり
	// うる状況を模している）。
	_, err = pool.Exec(ctx, `DELETE FROM resource_likes WHERE user_id = $1 AND resource_id = $2`, liker1ID, resourceID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `UPDATE learning_resource_metrics SET like_count = 999 WHERE resource_id = $1`, resourceID)
	require.NoError(t, err)

	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM learning_resource_metrics WHERE resource_id = $1`, resourceID).Scan(&likeCount))
	require.EqualValues(t, 999, likeCount, "drift is set up correctly before reconcile")

	require.NoError(t, metricsRepo.Reconcile(ctx))

	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM learning_resource_metrics WHERE resource_id = $1`, resourceID).Scan(&likeCount))
	assert.EqualValues(t, 1, likeCount, "reconcile should correct like_count back to the real resource_likes row count")
}
