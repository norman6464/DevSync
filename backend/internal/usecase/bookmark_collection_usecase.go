package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// bookmarkCollectionOwnerOf は所有権チェック用にコレクションの所有者 ID を取り出す。
func bookmarkCollectionOwnerOf(c *model.BookmarkCollection) uint { return c.UserID }

// CreateBookmarkCollectionUseCase はブックマークコレクションを作成する。
type CreateBookmarkCollectionUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewCreateBookmarkCollectionUseCase は CreateBookmarkCollectionUseCase を生成する。
func NewCreateBookmarkCollectionUseCase(collections repository.BookmarkCollectionRepository) *CreateBookmarkCollectionUseCase {
	return &CreateBookmarkCollectionUseCase{collections: collections}
}

// Execute は各項目を検証し、前後の空白を落として作成する。
func (uc *CreateBookmarkCollectionUseCase) Execute(ctx context.Context, collection *model.BookmarkCollection) error {
	if err := domain.ValidateStringLength(collection.Name, 1, 100, "コレクション名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(collection.Description, 0, 500, "説明"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(collection.Color, 0, 20, "カラー"); err != nil {
		return err
	}
	collection.Name = strings.TrimSpace(collection.Name)
	collection.Description = strings.TrimSpace(collection.Description)
	collection.Color = strings.TrimSpace(collection.Color)
	return uc.collections.Create(ctx, collection)
}

// ListBookmarkCollectionsUseCase は指定ユーザーのコレクション一覧を取得する。
type ListBookmarkCollectionsUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewListBookmarkCollectionsUseCase は ListBookmarkCollectionsUseCase を生成する。
func NewListBookmarkCollectionsUseCase(collections repository.BookmarkCollectionRepository) *ListBookmarkCollectionsUseCase {
	return &ListBookmarkCollectionsUseCase{collections: collections}
}

// Execute はコレクション一覧を返す。
func (uc *ListBookmarkCollectionsUseCase) Execute(ctx context.Context, userID uint) ([]model.BookmarkCollection, error) {
	return uc.collections.FindByUserID(ctx, userID)
}

// UpdateBookmarkCollectionUseCase は所有者本人のコレクションを更新する。
type UpdateBookmarkCollectionUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewUpdateBookmarkCollectionUseCase は UpdateBookmarkCollectionUseCase を生成する。
func NewUpdateBookmarkCollectionUseCase(collections repository.BookmarkCollectionRepository) *UpdateBookmarkCollectionUseCase {
	return &UpdateBookmarkCollectionUseCase{collections: collections}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
func (uc *UpdateBookmarkCollectionUseCase) Execute(ctx context.Context, id, userID uint, updates *model.BookmarkCollection) (*model.BookmarkCollection, error) {
	collection, err := ensureOwner(ctx, uc.collections.FindByID, id, userID, bookmarkCollectionOwnerOf)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(updates.Name); name != "" {
		if err := domain.ValidateStringLength(name, 1, 100, "コレクション名"); err != nil {
			return nil, err
		}
		collection.Name = name
	}
	if desc := strings.TrimSpace(updates.Description); desc != "" {
		if err := domain.ValidateStringLength(desc, 1, 500, "説明"); err != nil {
			return nil, err
		}
		collection.Description = desc
	}
	if color := strings.TrimSpace(updates.Color); color != "" {
		if err := domain.ValidateStringLength(color, 1, 20, "カラー"); err != nil {
			return nil, err
		}
		collection.Color = color
	}

	if err := uc.collections.Update(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// DeleteBookmarkCollectionUseCase は所有者本人のコレクションを削除する。
type DeleteBookmarkCollectionUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewDeleteBookmarkCollectionUseCase は DeleteBookmarkCollectionUseCase を生成する。
func NewDeleteBookmarkCollectionUseCase(collections repository.BookmarkCollectionRepository) *DeleteBookmarkCollectionUseCase {
	return &DeleteBookmarkCollectionUseCase{collections: collections}
}

// Execute は所有権を検証してから削除する。
func (uc *DeleteBookmarkCollectionUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.collections.FindByID, id, userID, bookmarkCollectionOwnerOf); err != nil {
		return err
	}
	return uc.collections.Delete(ctx, id)
}

// AddPostToBookmarkCollectionUseCase はコレクションに投稿を追加する。
type AddPostToBookmarkCollectionUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewAddPostToBookmarkCollectionUseCase は AddPostToBookmarkCollectionUseCase を生成する。
func NewAddPostToBookmarkCollectionUseCase(collections repository.BookmarkCollectionRepository) *AddPostToBookmarkCollectionUseCase {
	return &AddPostToBookmarkCollectionUseCase{collections: collections}
}

// Execute は所有権を検証し、未追加の投稿だけを追加する。追加済みなら Conflict を返す。
func (uc *AddPostToBookmarkCollectionUseCase) Execute(ctx context.Context, collectionID, postID, userID uint) error {
	if _, err := ensureOwner(ctx, uc.collections.FindByID, collectionID, userID, bookmarkCollectionOwnerOf); err != nil {
		return err
	}

	added, err := uc.collections.AddPost(ctx, &model.BookmarkCollectionItem{
		CollectionID: collectionID,
		PostID:       postID,
	})
	if err != nil {
		return err
	}
	if !added {
		return domain.NewError(domain.ErrCodeConflict, "この投稿は既にコレクションに追加されています", nil)
	}
	return nil
}

// RemovePostFromBookmarkCollectionUseCase はコレクションから投稿を取り除く。
type RemovePostFromBookmarkCollectionUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewRemovePostFromBookmarkCollectionUseCase は RemovePostFromBookmarkCollectionUseCase を生成する。
func NewRemovePostFromBookmarkCollectionUseCase(collections repository.BookmarkCollectionRepository) *RemovePostFromBookmarkCollectionUseCase {
	return &RemovePostFromBookmarkCollectionUseCase{collections: collections}
}

// Execute は所有権を検証してから取り除く。
func (uc *RemovePostFromBookmarkCollectionUseCase) Execute(ctx context.Context, collectionID, postID, userID uint) error {
	if _, err := ensureOwner(ctx, uc.collections.FindByID, collectionID, userID, bookmarkCollectionOwnerOf); err != nil {
		return err
	}
	return uc.collections.RemovePost(ctx, collectionID, postID)
}

// ListBookmarkCollectionPostsUseCase はコレクション内の投稿一覧を取得する。
type ListBookmarkCollectionPostsUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewListBookmarkCollectionPostsUseCase は ListBookmarkCollectionPostsUseCase を生成する。
func NewListBookmarkCollectionPostsUseCase(collections repository.BookmarkCollectionRepository) *ListBookmarkCollectionPostsUseCase {
	return &ListBookmarkCollectionPostsUseCase{collections: collections}
}

// Execute はページネーション付きの投稿一覧と総数を返す。
func (uc *ListBookmarkCollectionPostsUseCase) Execute(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	return uc.collections.GetPosts(ctx, collectionID, limit, offset)
}

// CountBookmarkCollectionsUseCase は指定ユーザーのコレクション総数を返す。
type CountBookmarkCollectionsUseCase struct {
	collections repository.BookmarkCollectionRepository
}

// NewCountBookmarkCollectionsUseCase は CountBookmarkCollectionsUseCase を生成する。
func NewCountBookmarkCollectionsUseCase(collections repository.BookmarkCollectionRepository) *CountBookmarkCollectionsUseCase {
	return &CountBookmarkCollectionsUseCase{collections: collections}
}

// Execute はコレクション総数を返す。
func (uc *CountBookmarkCollectionsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.collections.CountByUserID(ctx, userID)
}
