package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postCollectionRepository は [repository.PostCollectionRepository] の GORM 実装。
type postCollectionRepository struct {
	db *gorm.DB
}

// NewPostCollectionRepository は PostCollectionRepository の GORM 実装を返す。
func NewPostCollectionRepository(db *gorm.DB) repository.PostCollectionRepository {
	return &postCollectionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostCollectionRepository = (*postCollectionRepository)(nil)

// Create は新しい投稿コレクションをデータベースに作成する。
func (r *postCollectionRepository) Create(ctx context.Context, collection *model.PostCollection) error {
	return r.db.WithContext(ctx).Create(collection).Error
}

// FindByID は指定IDのコレクションをユーザー情報付きで取得する。
func (r *postCollectionRepository) FindByID(ctx context.Context, id uint) (*model.PostCollection, error) {
	var collection model.PostCollection
	err := r.db.WithContext(ctx).Preload("User").First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// FindByUserID は指定ユーザーの全コレクションをページネーション付きで取得する（新しい順）。
func (r *postCollectionRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	var collections []model.PostCollection
	var total int64
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query.Model(&model.PostCollection{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&collections).Error
	return collections, total, err
}

// FindPublicByUserID は指定ユーザーの公開コレクションを取得する（新しい順）。
func (r *postCollectionRepository) FindPublicByUserID(ctx context.Context, userID uint) ([]model.PostCollection, error) {
	var collections []model.PostCollection
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_public = ?", userID, true).
		Order("created_at DESC").
		Find(&collections).Error
	return collections, err
}

// Update は既存のコレクションを更新する。
func (r *postCollectionRepository) Update(ctx context.Context, collection *model.PostCollection) error {
	return r.db.WithContext(ctx).Save(collection).Error
}

// Delete は指定IDのコレクションとその関連アイテムを削除する。
func (r *postCollectionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Where("collection_id = ?", id).Delete(&model.PostCollectionItem{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model.PostCollection{}, id).Error
}

// AddPost はコレクションに投稿を追加する。
func (r *postCollectionRepository) AddPost(ctx context.Context, item *model.PostCollectionItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// HasPost は指定コレクションに指定投稿が存在するかを確認する。
func (r *postCollectionRepository) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostCollectionItem{}).
		Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Count(&count).Error
	return count > 0, err
}

// RemovePost はコレクションから投稿を削除する。
func (r *postCollectionRepository) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return r.db.WithContext(ctx).Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Delete(&model.PostCollectionItem{}).Error
}

// GetPostsByCollectionID は指定コレクションの投稿一覧を順序付きで取得する。
func (r *postCollectionRepository) GetPostsByCollectionID(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error) {
	var items []model.PostCollectionItem
	err := r.db.WithContext(ctx).Preload("Post").Preload("Post.User").
		Where("collection_id = ?", collectionID).
		Order("order_index ASC").
		Find(&items).Error
	return items, err
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (r *postCollectionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostCollection{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
