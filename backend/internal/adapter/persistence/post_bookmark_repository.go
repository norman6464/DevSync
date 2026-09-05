package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postBookmarkRepository は [repository.PostBookmarkRepository] の sqlc(pgx) 実装。
type postBookmarkRepository struct {
	q *sqlcgen.Queries
}

// NewPostBookmarkRepository は PostBookmarkRepository の sqlc(pgx) 実装を返す。
func NewPostBookmarkRepository(q *sqlcgen.Queries) repository.PostBookmarkRepository {
	return &postBookmarkRepository{q: q}
}

var _ repository.PostBookmarkRepository = (*postBookmarkRepository)(nil)

// toModelCodeSnippet は sqlc の生成行を model.CodeSnippet へ変換する。
func toModelCodeSnippet(row sqlcgen.CodeSnippet) model.CodeSnippet {
	return model.CodeSnippet{
		ID:           uint(row.ID),
		PostID:       uint(row.PostID),
		UserID:       uint(row.UserID),
		Language:     row.Language,
		FileName:     fromStringPtr(row.FileName),
		Code:         row.Code,
		CommentCount: int(fromInt64PtrValue(row.CommentCount)),
		ForkedFromID: fromInt64PtrToUintPtr(row.ForkedFromID),
		ForkCount:    int(fromInt64PtrValue(row.ForkCount)),
		CreatedAt:    timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:    timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Bookmark は投稿をブックマークし、投稿の bookmark_count を加算する。
// 移行前の GORM 実装と同じくトランザクションでは括らない（元実装も2つの独立した操作だったため）。
func (r *postBookmarkRepository) Bookmark(ctx context.Context, userID, postID uint) error {
	if err := r.q.CreateBookmark(ctx, sqlcgen.CreateBookmarkParams{
		UserID: int64(userID),
		PostID: int64(postID),
	}); err != nil {
		return err
	}
	return r.q.IncrementPostBookmarkCount(ctx, int64(postID))
}

// Unbookmark はブックマークを解除し、実際に削除できたときだけ bookmark_count をデクリメントする。
// 移行前の GORM 実装と同じく、デクリメント自体のエラーは呼び出し元へ返さない。
func (r *postBookmarkRepository) Unbookmark(ctx context.Context, userID, postID uint) error {
	rowsAffected, err := r.q.DeleteBookmark(ctx, sqlcgen.DeleteBookmarkParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
	if rowsAffected > 0 {
		_ = r.q.DecrementPostBookmarkCount(ctx, int64(postID))
	}
	return err
}

// HasBookmarked は指定ユーザーが投稿をブックマーク済みかを返す。
func (r *postBookmarkRepository) HasBookmarked(ctx context.Context, userID, postID uint) (bool, error) {
	count, err := r.q.CountBookmarkByUserAndPost(ctx, sqlcgen.CountBookmarkByUserAndPostParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindBookmarkedByUserID は指定ユーザーのブックマーク済み投稿をページネーション付きで取得する。
func (r *postBookmarkRepository) FindBookmarkedByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error) {
	offset := (page - 1) * limit

	total, err := r.q.CountBookmarksMadeByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListBookmarkedPostsByUser(ctx, sqlcgen.ListBookmarkedPostsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	posts := make([]model.Post, len(rows))
	postIDs := make([]int64, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
		postIDs[i] = row.Post.ID
	}

	if len(postIDs) > 0 {
		snippetRows, err := r.q.ListCodeSnippetsByPostIDs(ctx, postIDs)
		if err != nil {
			return nil, 0, err
		}
		snippetsByPostID := make(map[uint][]model.CodeSnippet)
		for _, row := range snippetRows {
			postID := uint(row.PostID)
			snippetsByPostID[postID] = append(snippetsByPostID[postID], toModelCodeSnippet(row))
		}
		for i := range posts {
			posts[i].CodeSnippets = snippetsByPostID[posts[i].ID]
		}
	}

	return posts, total, nil
}

// CountBookmarkedByUserID は指定ユーザーのブックマーク件数を返す。
func (r *postBookmarkRepository) CountBookmarkedByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountBookmarksMadeByUser(ctx, int64(userID))
}
