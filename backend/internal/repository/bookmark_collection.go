package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BookmarkCollectionRepository はブックマークコレクションデータへのアクセスを提供するリポジトリ実装。
type BookmarkCollectionRepository struct {
	db *gorm.DB
}

// NewBookmarkCollectionRepository は新しいBookmarkCollectionRepositoryインスタンスを生成する。
func NewBookmarkCollectionRepository(db *gorm.DB) *BookmarkCollectionRepository {
	return &BookmarkCollectionRepository{db: db}
}

// Create は新しいブックマークコレクションを作成する。
func (r *BookmarkCollectionRepository) Create(collection *model.BookmarkCollection) error {
	return r.db.Create(collection).Error
}

// FindByID は指定IDのコレクションを取得する。
func (r *BookmarkCollectionRepository) FindByID(id uint) (*model.BookmarkCollection, error) {
	var collection model.BookmarkCollection
	err := r.db.First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// FindByUserID は指定ユーザーのコレクション一覧を取得する。
func (r *BookmarkCollectionRepository) FindByUserID(userID uint) ([]model.BookmarkCollection, error) {
	var collections []model.BookmarkCollection
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&collections).Error
	return collections, err
}

// Update はコレクションを更新する。
func (r *BookmarkCollectionRepository) Update(collection *model.BookmarkCollection) error {
	return r.db.Save(collection).Error
}

// Delete はコレクションとそのアイテムを削除する。
func (r *BookmarkCollectionRepository) Delete(id uint) error {
	tx := r.db.Begin()
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
func (r *BookmarkCollectionRepository) AddPost(item *model.BookmarkCollectionItem) error {
	return r.db.Create(item).Error
}

// RemovePost はコレクションから投稿を削除する。
func (r *BookmarkCollectionRepository) RemovePost(collectionID, postID uint) error {
	return r.db.Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Delete(&model.BookmarkCollectionItem{}).Error
}

// GetPosts はコレクション内の投稿一覧を取得する。
func (r *BookmarkCollectionRepository) GetPosts(collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	var total int64
	r.db.Model(&model.BookmarkCollectionItem{}).Where("collection_id = ?", collectionID).Count(&total)

	var items []model.BookmarkCollectionItem
	err := r.db.Where("collection_id = ?", collectionID).
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
func (r *BookmarkCollectionRepository) HasPost(collectionID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BookmarkCollectionItem{}).
		Where("collection_id = ? AND post_id = ?", collectionID, postID).
		Count(&count).Error
	return count > 0, err
}
