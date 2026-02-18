package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostViewRepository は投稿閲覧数のGORM実装。
type PostViewRepository struct {
	db *gorm.DB
}

// NewPostViewRepository は新しいPostViewRepositoryインスタンスを生成する。
func NewPostViewRepository(db *gorm.DB) *PostViewRepository {
	return &PostViewRepository{db: db}
}

func (r *PostViewRepository) RecordView(view *model.PostView) error {
	return r.db.Create(view).Error
}

func (r *PostViewRepository) GetViewCount(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.PostView{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (r *PostViewRepository) HasViewed(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostView{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *PostViewRepository) GetMostViewed(limit int) ([]model.ViewCount, error) {
	var results []model.ViewCount
	err := r.db.Model(&model.PostView{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}
