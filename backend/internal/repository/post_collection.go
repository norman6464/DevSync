package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostCollectionRepository は投稿コレクションデータへのアクセスを提供するリポジトリ実装。
type PostCollectionRepository struct {
	db *gorm.DB
}

// NewPostCollectionRepository は新しいPostCollectionRepositoryインスタンスを生成する。
func NewPostCollectionRepository(db *gorm.DB) *PostCollectionRepository {
	return &PostCollectionRepository{db: db}
}

// Create は新しい投稿コレクションをデータベースに作成する。
func (r *PostCollectionRepository) Create(collection *model.PostCollection) error {
	return r.db.Create(collection).Error
}

// FindByID は指定IDのコレクションをユーザー情報付きで取得する。
func (r *PostCollectionRepository) FindByID(id uint) (*model.PostCollection, error) {
	var collection model.PostCollection
	err := r.db.Preload("User").First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// FindByUserID は指定ユーザーの全コレクションをページネーション付きで取得する（新しい順）。
func (r *PostCollectionRepository) FindByUserID(userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	var collections []model.PostCollection
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.PostCollection{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&collections).Error
	return collections, total, err
}

// FindPublicByUserID は指定ユーザーの公開コレクションを取得する（新しい順）。
func (r *PostCollectionRepository) FindPublicByUserID(userID uint) ([]model.PostCollection, error) {
	var collections []model.PostCollection
	err := r.db.Where("user_id = ? AND is_public = ?", userID, true).
		Order("created_at DESC").
		Find(&collections).Error
	return collections, err
}

// Update は既存のコレクションを更新する。
func (r *PostCollectionRepository) Update(collection *model.PostCollection) error {
	return r.db.Save(collection).Error
}

// Delete は指定IDのコレクションとその関連アイテムを削除する。
func (r *PostCollectionRepository) Delete(id uint) error {
	if err := r.db.Where("collection_id = ?", id).Delete(&model.PostCollectionItem{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&model.PostCollection{}, id).Error
}

// AddPost はコレクションに投稿を追加する。
func (r *PostCollectionRepository) AddPost(item *model.PostCollectionItem) error {
	return r.db.Create(item).Error
}

// HasPost は指定コレクションに指定投稿が存在するかを確認する。
func (r *PostCollectionRepository) HasPost(collectionID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostCollectionItem{}).
		Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Count(&count).Error
	return count > 0, err
}

// RemovePost はコレクションから投稿を削除する。
func (r *PostCollectionRepository) RemovePost(collectionID, postID uint) error {
	return r.db.Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Delete(&model.PostCollectionItem{}).Error
}

// GetPostsByCollectionID は指定コレクションの投稿一覧を順序付きで取得する。
func (r *PostCollectionRepository) GetPostsByCollectionID(collectionID uint) ([]model.PostCollectionItem, error) {
	var items []model.PostCollectionItem
	err := r.db.Preload("Post").Preload("Post.User").
		Where("collection_id = ?", collectionID).
		Order("order_index ASC").
		Find(&items).Error
	return items, err
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (r *PostCollectionRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.PostCollection{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
