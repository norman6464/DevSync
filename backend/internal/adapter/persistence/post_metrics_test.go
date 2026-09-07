package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostMetrics_LikeCommentViewFlow_UpdatesInSameStatement は、いいね・コメント・閲覧の
// 作成/削除が同一SQL文でpost_metricsを正しく増減させることを検証する（DEVSYNC-159）。
// post_metrics行が最初は存在しない投稿から始め、遅延生成される様子も含めて確認する。
func TestPostMetrics_LikeCommentViewFlow_UpdatesInSameStatement(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, otherID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "pm_author"},
		{&otherID, "pm_other"},
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

	var postID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, authorID).Scan(&postID)
	require.NoError(t, err)

	assertPostMetrics := func(wantLike, wantComment, wantView int64) {
		t.Helper()
		var like, comment, view int64
		row := pool.QueryRow(ctx, `
			SELECT COALESCE(like_count, 0), COALESCE(comment_count, 0), COALESCE(view_count, 0)
			FROM post_metrics WHERE post_id = $1
		`, postID)
		err := row.Scan(&like, &comment, &view)
		if err != nil {
			// post_metrics行がまだ無い（遅延生成前）場合は全て0として扱う。
			like, comment, view = 0, 0, 0
		}
		assert.Equal(t, wantLike, like, "like_count")
		assert.Equal(t, wantComment, comment, "comment_count")
		assert.Equal(t, wantView, view, "view_count")
	}

	// 最初はpost_metrics行が存在しない。
	assertPostMetrics(0, 0, 0)

	// いいね: 作成でpost_metrics行が遅延生成され1になる。
	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: otherID, PostID: postID}))
	assertPostMetrics(1, 0, 0)

	// いいね取り消し: 0未満にはならない。
	rowsAffected, err := q.DeletePostLike(ctx, sqlcgen.DeletePostLikeParams{UserID: otherID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)
	assertPostMetrics(0, 0, 0)

	// 既に取り消し済みのいいねをもう一度取り消しても0未満にならない（rowsAffected=0）。
	rowsAffected, err = q.DeletePostLike(ctx, sqlcgen.DeletePostLikeParams{UserID: otherID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 0, rowsAffected)
	assertPostMetrics(0, 0, 0)

	// コメント: 作成でcomment_countが増える。
	commentRow, err := q.CreatePostComment(ctx, sqlcgen.CreatePostCommentParams{
		UserID: otherID, PostID: postID, Content: "nice",
	})
	require.NoError(t, err)
	assertPostMetrics(0, 1, 0)

	// コメント削除でcomment_countが減る。
	require.NoError(t, q.DeletePostComment(ctx, commentRow.ID))
	assertPostMetrics(0, 0, 0)

	// 閲覧: 初回はview_countが増える。
	rowsAffected, err = q.CreatePostView(ctx, sqlcgen.CreatePostViewParams{UserID: otherID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)
	assertPostMetrics(0, 0, 1)

	// 同一ユーザーの2回目の閲覧はON CONFLICT DO NOTHINGで重複記録されず、view_countも増えない。
	rowsAffected, err = q.CreatePostView(ctx, sqlcgen.CreatePostViewParams{UserID: otherID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 0, rowsAffected)
	assertPostMetrics(0, 0, 1)
}

// TestReconcileAllPostMetrics_CorrectsDrift は、CASCADE削除等でpost_metricsが
// 実件数からずれたケースを意図的に作り、reconcileジョブが正しく補正することを検証する
// （DEVSYNC-159のチケット本文で明示されているテスト要件）。
func TestReconcileAllPostMetrics_CorrectsDrift(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)
	metricsRepo := NewPostMetricsRepository(q)

	var authorID, liker1ID, liker2ID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "pm_reconcile_author"},
		{&liker1ID, "pm_reconcile_liker1"},
		{&liker2ID, "pm_reconcile_liker2"},
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

	var postID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, authorID).Scan(&postID)
	require.NoError(t, err)

	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: liker1ID, PostID: postID}))
	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: liker2ID, PostID: postID}))

	var likeCount int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM post_metrics WHERE post_id = $1`, postID).Scan(&likeCount))
	require.EqualValues(t, 2, likeCount)

	// 意図的にpost_metricsをずらす（CASCADE削除等でIncrement/Decrementを経由しない
	// 変化が起きたケースを再現する。実際にはpost_reactions(kind='like')の行を直接消しても
	// FKのCASCADEはpost_metrics側では発生しないため、これはユーザー削除カスケードで
	// いいねが消えたのにpost_metricsだけ古い値のまま残る、という実際に起こりうる状況を
	// 模している）。
	_, err = pool.Exec(ctx, `DELETE FROM post_reactions WHERE user_id = $1 AND post_id = $2 AND kind = 'like'`, liker1ID, postID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `UPDATE post_metrics SET like_count = 999 WHERE post_id = $1`, postID)
	require.NoError(t, err)

	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM post_metrics WHERE post_id = $1`, postID).Scan(&likeCount))
	require.EqualValues(t, 999, likeCount, "drift is set up correctly before reconcile")

	require.NoError(t, metricsRepo.Reconcile(ctx))

	require.NoError(t, pool.QueryRow(ctx, `SELECT like_count FROM post_metrics WHERE post_id = $1`, postID).Scan(&likeCount))
	assert.EqualValues(t, 1, likeCount, "reconcile should correct like_count back to the real post_reactions row count")
}
