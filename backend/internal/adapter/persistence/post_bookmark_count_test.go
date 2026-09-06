package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFindAll_AttachesAccurateBookmarkCounts は、posts.bookmark_count列を撤去した後、
// attachBookmarkCountsToPostsがbookmarksテーブルから実際の件数を正しく集計することを検証する
// （ブックマーク0件の投稿がGROUP BYの結果に現れなくても0になることを含む）。
func TestFindAll_AttachesAccurateBookmarkCounts(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()

	var authorID, bookmarker1ID, bookmarker2ID int64
	for _, dst := range []struct {
		id       *int64
		username string
	}{
		{&authorID, "bmc_author"},
		{&bookmarker1ID, "bmc_bookmarker1"},
		{&bookmarker2ID, "bmc_bookmarker2"},
	} {
		err := pool.QueryRow(ctx, `
			INSERT INTO users (username, name, email, created_at, updated_at)
			VALUES ($1, $1, $1 || '@example.com', now(), now())
			RETURNING id
		`, dst.username).Scan(dst.id)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1, $2, $3)`, authorID, bookmarker1ID, bookmarker2ID)
	})

	postRepo := NewPostRepository(pool)

	postTwoBookmarks := &model.Post{UserID: uint(authorID), Title: "2件ブックマーク", Content: "本文"}
	require.NoError(t, postRepo.Create(ctx, postTwoBookmarks))
	postOneBookmark := &model.Post{UserID: uint(authorID), Title: "1件ブックマーク", Content: "本文"}
	require.NoError(t, postRepo.Create(ctx, postOneBookmark))
	postNoBookmarks := &model.Post{UserID: uint(authorID), Title: "0件ブックマーク", Content: "本文"}
	require.NoError(t, postRepo.Create(ctx, postNoBookmarks))

	for _, b := range []struct{ userID, postID int64 }{
		{bookmarker1ID, int64(postTwoBookmarks.ID)},
		{bookmarker2ID, int64(postTwoBookmarks.ID)},
		{bookmarker1ID, int64(postOneBookmark.ID)},
	} {
		_, err := pool.Exec(ctx, `INSERT INTO bookmarks (user_id, post_id, created_at) VALUES ($1, $2, now())`, b.userID, b.postID)
		require.NoError(t, err)
	}

	posts, _, err := postRepo.FindByUserID(ctx, uint(authorID), 10, 0)
	require.NoError(t, err)
	countByTitle := make(map[string]int, len(posts))
	for _, p := range posts {
		countByTitle[p.Title] = p.BookmarkCount
	}

	assert.Equal(t, 2, countByTitle["2件ブックマーク"])
	assert.Equal(t, 1, countByTitle["1件ブックマーク"])
	assert.Equal(t, 0, countByTitle["0件ブックマーク"])
}
