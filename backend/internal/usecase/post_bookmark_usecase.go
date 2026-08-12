package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// BookmarkPostUseCase は投稿をブックマークする。
type BookmarkPostUseCase struct {
	bookmarks repository.PostBookmarkRepository
	posts     repository.PostAuthorReader
}

// NewBookmarkPostUseCase は BookmarkPostUseCase を生成する。
func NewBookmarkPostUseCase(
	bookmarks repository.PostBookmarkRepository,
	posts repository.PostAuthorReader,
) *BookmarkPostUseCase {
	return &BookmarkPostUseCase{bookmarks: bookmarks, posts: posts}
}

// Execute は投稿者を検証したうえでブックマークする。自分の投稿はブックマークできない。
func (uc *BookmarkPostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.bookmarks.Bookmark(ctx, userID, postID)
}

// UnbookmarkPostUseCase は投稿のブックマークを解除する。
type UnbookmarkPostUseCase struct {
	bookmarks repository.PostBookmarkRepository
	posts     repository.PostAuthorReader
}

// NewUnbookmarkPostUseCase は UnbookmarkPostUseCase を生成する。
func NewUnbookmarkPostUseCase(
	bookmarks repository.PostBookmarkRepository,
	posts repository.PostAuthorReader,
) *UnbookmarkPostUseCase {
	return &UnbookmarkPostUseCase{bookmarks: bookmarks, posts: posts}
}

// Execute は投稿者を検証したうえでブックマークを解除する。
// 自分の投稿はそもそもブックマークできないため、解除も同じ条件で弾く。
func (uc *UnbookmarkPostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.bookmarks.Unbookmark(ctx, userID, postID)
}

// HasBookmarkedPostUseCase は投稿をブックマーク済みかを判定する。
type HasBookmarkedPostUseCase struct {
	bookmarks repository.PostBookmarkRepository
}

// NewHasBookmarkedPostUseCase は HasBookmarkedPostUseCase を生成する。
func NewHasBookmarkedPostUseCase(bookmarks repository.PostBookmarkRepository) *HasBookmarkedPostUseCase {
	return &HasBookmarkedPostUseCase{bookmarks: bookmarks}
}

// Execute は指定ユーザーが投稿をブックマーク済みかを返す。
func (uc *HasBookmarkedPostUseCase) Execute(ctx context.Context, userID, postID uint) (bool, error) {
	return uc.bookmarks.HasBookmarked(ctx, userID, postID)
}

// ListBookmarkedPostsUseCase はブックマーク済み投稿の一覧を取得する。
type ListBookmarkedPostsUseCase struct {
	bookmarks repository.PostBookmarkRepository
}

// NewListBookmarkedPostsUseCase は ListBookmarkedPostsUseCase を生成する。
func NewListBookmarkedPostsUseCase(bookmarks repository.PostBookmarkRepository) *ListBookmarkedPostsUseCase {
	return &ListBookmarkedPostsUseCase{bookmarks: bookmarks}
}

// Execute はブックマーク済み投稿を新しい順に返す（総件数も返す）。
func (uc *ListBookmarkedPostsUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error) {
	return uc.bookmarks.FindBookmarkedByUserID(ctx, userID, page, limit)
}

// CountBookmarkedPostsUseCase はブックマーク済み投稿の件数を取得する。
type CountBookmarkedPostsUseCase struct {
	bookmarks repository.PostBookmarkRepository
}

// NewCountBookmarkedPostsUseCase は CountBookmarkedPostsUseCase を生成する。
func NewCountBookmarkedPostsUseCase(bookmarks repository.PostBookmarkRepository) *CountBookmarkedPostsUseCase {
	return &CountBookmarkedPostsUseCase{bookmarks: bookmarks}
}

// Execute はブックマーク件数を返す。
func (uc *CountBookmarkedPostsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.bookmarks.CountBookmarkedByUserID(ctx, userID)
}
