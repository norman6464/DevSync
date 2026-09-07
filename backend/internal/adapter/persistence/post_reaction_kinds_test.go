package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostReactions_UserCanHoldMultipleKindsOnSamePost は、post_reactionsが
// likes/bookmarks/post_pins/post_views/reactionsの5テーブルを対象単位で統合した後も
// （DEVSYNC-160）、同一ユーザーが同一投稿に対してlike・bookmark・pin・viewを
// 同時に持てることを検証する（UNIQUE(user_id, post_id, kind, value)によりkindごとに
// 独立して重複排除されるべきで、統合前と同じ挙動を保つ）。
func TestPostReactions_UserCanHoldMultipleKindsOnSamePost(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, actorID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "prk_author"},
		{&actorID, "prk_actor"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, authorID, actorID)
	})

	var postID int64
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, authorID).Scan(&postID))

	require.NoError(t, q.CreatePostLike(ctx, sqlcgen.CreatePostLikeParams{UserID: actorID, PostID: postID}))
	require.NoError(t, q.CreateBookmark(ctx, sqlcgen.CreateBookmarkParams{UserID: actorID, PostID: postID}))
	pinOrder := int64(0)
	require.NoError(t, q.CreatePostPin(ctx, sqlcgen.CreatePostPinParams{UserID: actorID, PostID: postID, PinOrder: &pinOrder}))
	rowsAffected, err := q.CreatePostView(ctx, sqlcgen.CreatePostViewParams{UserID: actorID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected)
	require.NoError(t, q.CreateReaction(ctx, sqlcgen.CreateReactionParams{UserID: actorID, PostID: postID, Value: "👍"}))

	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT COUNT(*) FROM post_reactions WHERE user_id = $1 AND post_id = $2`, actorID, postID).Scan(&count))
	assert.EqualValues(t, 5, count, "like/bookmark/pin/view/emojiの5行が共存できるべき")

	hasLiked, err := q.CountPostLikeByUserAndPost(ctx, sqlcgen.CountPostLikeByUserAndPostParams{UserID: actorID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, hasLiked)

	hasBookmarked, err := q.CountBookmarkByUserAndPost(ctx, sqlcgen.CountBookmarkByUserAndPostParams{UserID: actorID, PostID: postID})
	require.NoError(t, err)
	assert.EqualValues(t, 1, hasBookmarked)
}

// TestPostReactions_EmojiAllowsMultipleValuesButNotDuplicates は、絵文字リアクション
// （kind='emoji'）だけが1ユーザー・1投稿に対して複数行（絵文字違い）を持てる一方、
// 同じ絵文字の重複は拒否されることを検証する（元のreactionsテーブルの
// UNIQUE(user_id, post_id, emoji)と同じ制約をUNIQUE(user_id, post_id, kind, value)で
// 再現できているかの確認）。
func TestPostReactions_EmojiAllowsMultipleValuesButNotDuplicates(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)

	var authorID, actorID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "prk_emoji_author"},
		{&actorID, "prk_emoji_actor"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2)`, authorID, actorID)
	})

	var postID int64
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, title, content, is_draft, created_at, updated_at)
		VALUES ($1, 'title', 'content', false, now(), now())
		RETURNING id
	`, authorID).Scan(&postID))

	require.NoError(t, q.CreateReaction(ctx, sqlcgen.CreateReactionParams{UserID: actorID, PostID: postID, Value: "👍"}))
	require.NoError(t, q.CreateReaction(ctx, sqlcgen.CreateReactionParams{UserID: actorID, PostID: postID, Value: "🎉"}))

	emojis, err := q.ListUserReactionEmojisByPost(ctx, sqlcgen.ListUserReactionEmojisByPostParams{UserID: actorID, PostID: postID})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"👍", "🎉"}, emojis)

	err = q.CreateReaction(ctx, sqlcgen.CreateReactionParams{UserID: actorID, PostID: postID, Value: "👍"})
	assert.Error(t, err, "同じ絵文字の重複はUNIQUE制約で拒否されるべき")
}
