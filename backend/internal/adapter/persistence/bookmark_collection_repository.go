package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// bookmarkCollectionRepository は [repository.BookmarkCollectionRepository] の GORM 実装。
type bookmarkCollectionRepository struct {
	db *gorm.DB
}

// NewBookmarkCollectionRepository は BookmarkCollectionRepository の GORM 実装を返す。
func NewBookmarkCollectionRepository(db *gorm.DB) repository.BookmarkCollectionRepository {
	return &bookmarkCollectionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookmarkCollectionRepository = (*bookmarkCollectionRepository)(nil)

// Create は新しいブックマークコレクションを作成する。
func (r *bookmarkCollectionRepository) Create(ctx context.Context, collection *model.BookmarkCollection) error {
	return r.db.WithContext(ctx).Create(collection).Error
}

// FindByID は指定IDのコレクションを取得する。
func (r *bookmarkCollectionRepository) FindByID(ctx context.Context, id uint) (*model.BookmarkCollection, error) {
	var collection model.BookmarkCollection
	err := r.db.WithContext(ctx).First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// FindByUserID は指定ユーザーのコレクション一覧を取得する。
func (r *bookmarkCollectionRepository) FindByUserID(ctx context.Context, userID uint) ([]model.BookmarkCollection, error) {
	var collections []model.BookmarkCollection
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&collections).Error
	return collections, err
}

// Update はコレクションを更新する。
func (r *bookmarkCollectionRepository) Update(ctx context.Context, collection *model.BookmarkCollection) error {
	return r.db.WithContext(ctx).Save(collection).Error
}

// Delete はコレクションとそのアイテムを削除する。
func (r *bookmarkCollectionRepository) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("collection_id = ?", id).Delete(&model.BookmarkCollectionItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&model.BookmarkCollection{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// AddPost はコレクションに投稿を追加する。
func (r *bookmarkCollectionRepository) AddPost(ctx context.Context, item *model.BookmarkCollectionItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// RemovePost はコレクションから投稿を削除する。
func (r *bookmarkCollectionRepository) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return r.db.WithContext(ctx).Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Delete(&model.BookmarkCollectionItem{}).Error
}

// GetPosts はコレクション内の投稿一覧を取得する。
func (r *bookmarkCollectionRepository) GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.BookmarkCollectionItem{}).Where("collection_id = ?", collectionID).Count(&total)

	var items []model.BookmarkCollectionItem
	err := r.db.WithContext(ctx).Where("collection_id = ?", collectionID).
		Preload("Post").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	posts := make([]model.Post, 0, len(items))
	for _, item := range items {
		posts = append(posts, item.Post)
	}
	return posts, total, nil
}

// HasPost はコレクションに指定投稿が含まれているかを返す。
func (r *bookmarkCollectionRepository) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BookmarkCollectionItem{}).
		Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Count(&count).Error
	return count > 0, err
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (r *bookmarkCollectionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BookmarkCollection{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
