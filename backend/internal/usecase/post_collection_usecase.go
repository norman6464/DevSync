package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// collectionOwnerOf は所有権チェック用にコレクションの所有者 ID を取り出す。
func collectionOwnerOf(c *model.PostCollection) uint { return c.UserID }

// CreatePostCollectionUseCase は投稿コレクションを作成する。
type CreatePostCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewCreatePostCollectionUseCase は CreatePostCollectionUseCase を生成する。
func NewCreatePostCollectionUseCase(collections repository.PostCollectionRepository) *CreatePostCollectionUseCase {
	return &CreatePostCollectionUseCase{collections: collections}
}

// Execute はタイトルと説明を検証し、前後の空白を落として作成する。
func (uc *CreatePostCollectionUseCase) Execute(ctx context.Context, collection *model.PostCollection) (*model.PostCollection, error) {
	if err := domain.ValidateStringLength(collection.Title, 1, 200, "タイトル"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(collection.Description, 0, 1000, "説明"); err != nil {
		return nil, err
	}
	collection.Title = strings.TrimSpace(collection.Title)
	collection.Description = strings.TrimSpace(collection.Description)
	if err := uc.collections.Create(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// GetPostCollectionUseCase は指定 ID のコレクションを取得する。
type GetPostCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewGetPostCollectionUseCase は GetPostCollectionUseCase を生成する。
func NewGetPostCollectionUseCase(collections repository.PostCollectionRepository) *GetPostCollectionUseCase {
	return &GetPostCollectionUseCase{collections: collections}
}

// Execute はコレクションを返す。
func (uc *GetPostCollectionUseCase) Execute(ctx context.Context, id uint) (*model.PostCollection, error) {
	return uc.collections.FindByID(ctx, id)
}

// ListPostCollectionsForViewerUseCase は閲覧者の立場に応じたコレクション一覧を返す。
type ListPostCollectionsForViewerUseCase struct {
	collections repository.PostCollectionRepository
}

// NewListPostCollectionsForViewerUseCase は ListPostCollectionsForViewerUseCase を生成する。
func NewListPostCollectionsForViewerUseCase(collections repository.PostCollectionRepository) *ListPostCollectionsForViewerUseCase {
	return &ListPostCollectionsForViewerUseCase{collections: collections}
}

// Execute は自分のコレクションなら全件（ページネーション付き）、他人のものは公開のみ返す。
func (uc *ListPostCollectionsForViewerUseCase) Execute(ctx context.Context, viewerID, targetUserID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	if viewerID == targetUserID {
		return uc.collections.FindByUserID(ctx, targetUserID, limit, offset)
	}
	collections, err := uc.collections.FindPublicByUserID(ctx, targetUserID)
	if err != nil {
		return nil, 0, err
	}
	return collections, int64(len(collections)), nil
}

// CountPostCollectionsUseCase は指定ユーザーのコレクション総数を返す。
type CountPostCollectionsUseCase struct {
	collections repository.PostCollectionRepository
}

// NewCountPostCollectionsUseCase は CountPostCollectionsUseCase を生成する。
func NewCountPostCollectionsUseCase(collections repository.PostCollectionRepository) *CountPostCollectionsUseCase {
	return &CountPostCollectionsUseCase{collections: collections}
}

// Execute はコレクション総数を返す。
func (uc *CountPostCollectionsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.collections.CountByUserID(ctx, userID)
}

// UpdatePostCollectionUseCase は所有者本人のコレクションを更新する。
type UpdatePostCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewUpdatePostCollectionUseCase は UpdatePostCollectionUseCase を生成する。
func NewUpdatePostCollectionUseCase(collections repository.PostCollectionRepository) *UpdatePostCollectionUseCase {
	return &UpdatePostCollectionUseCase{collections: collections}
}

// Execute は所有権を検証し、渡された値で上書きする。
func (uc *UpdatePostCollectionUseCase) Execute(ctx context.Context, id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error) {
	collection, err := ensureOwner(ctx, uc.collections.FindByID, id, userID, collectionOwnerOf)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(description, 0, 1000, "説明"); err != nil {
		return nil, err
	}

	collection.Title = strings.TrimSpace(title)
	collection.Description = strings.TrimSpace(description)
	collection.IsPublic = isPublic

	if err := uc.collections.Update(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// DeletePostCollectionUseCase は所有者本人のコレクションを削除する。
type DeletePostCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewDeletePostCollectionUseCase は DeletePostCollectionUseCase を生成する。
func NewDeletePostCollectionUseCase(collections repository.PostCollectionRepository) *DeletePostCollectionUseCase {
	return &DeletePostCollectionUseCase{collections: collections}
}

// Execute は所有権を検証してから削除する。
func (uc *DeletePostCollectionUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.collections.FindByID, id, userID, collectionOwnerOf); err != nil {
		return err
	}
	return uc.collections.Delete(ctx, id)
}

// AddPostToCollectionUseCase はコレクションに投稿を追加する。
type AddPostToCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewAddPostToCollectionUseCase は AddPostToCollectionUseCase を生成する。
func NewAddPostToCollectionUseCase(collections repository.PostCollectionRepository) *AddPostToCollectionUseCase {
	return &AddPostToCollectionUseCase{collections: collections}
}

// Execute はメモを検証し、所有権を確認したうえで未追加の投稿だけを追加する。
// メモの検証を所有権チェックより先に行う順序は移行前から変えていない。
func (uc *AddPostToCollectionUseCase) Execute(ctx context.Context, collectionID, userID, postID uint, note string) error {
	if err := domain.ValidateStringLength(note, 0, 500, "メモ"); err != nil {
		return err
	}

	if _, err := ensureOwner(ctx, uc.collections.FindByID, collectionID, userID, collectionOwnerOf); err != nil {
		return err
	}

	exists, err := uc.collections.HasPost(ctx, collectionID, postID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeBadRequest, "すでに追加済みの投稿です", nil)
	}

	return uc.collections.AddPost(ctx, &model.PostCollectionItem{
		CollectionID: collectionID,
		PostID:       postID,
		Note:         strings.TrimSpace(note),
	})
}

// RemovePostFromCollectionUseCase はコレクションから投稿を取り除く。
type RemovePostFromCollectionUseCase struct {
	collections repository.PostCollectionRepository
}

// NewRemovePostFromCollectionUseCase は RemovePostFromCollectionUseCase を生成する。
func NewRemovePostFromCollectionUseCase(collections repository.PostCollectionRepository) *RemovePostFromCollectionUseCase {
	return &RemovePostFromCollectionUseCase{collections: collections}
}

// Execute は所有権を検証してから取り除く。
func (uc *RemovePostFromCollectionUseCase) Execute(ctx context.Context, collectionID, userID, postID uint) error {
	if _, err := ensureOwner(ctx, uc.collections.FindByID, collectionID, userID, collectionOwnerOf); err != nil {
		return err
	}
	return uc.collections.RemovePost(ctx, collectionID, postID)
}

// ListPostCollectionPostsUseCase はコレクションに含まれる投稿一覧を取得する。
type ListPostCollectionPostsUseCase struct {
	collections repository.PostCollectionRepository
}

// NewListPostCollectionPostsUseCase は ListPostCollectionPostsUseCase を生成する。
func NewListPostCollectionPostsUseCase(collections repository.PostCollectionRepository) *ListPostCollectionPostsUseCase {
	return &ListPostCollectionPostsUseCase{collections: collections}
}

// Execute は順序付きの投稿一覧を返す。
func (uc *ListPostCollectionPostsUseCase) Execute(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error) {
	return uc.collections.GetPostsByCollectionID(ctx, collectionID)
}
