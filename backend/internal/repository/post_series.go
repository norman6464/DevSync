package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostSeriesRepository は投稿シリーズデータへのアクセスを提供するリポジトリ実装。
type PostSeriesRepository struct {
	db *gorm.DB
}

// NewPostSeriesRepository は新しいPostSeriesRepositoryインスタンスを生成する。
func NewPostSeriesRepository(db *gorm.DB) *PostSeriesRepository {
	return &PostSeriesRepository{db: db}
}

// Create は新しい投稿シリーズをデータベースに作成する。
func (r *PostSeriesRepository) Create(series *model.PostSeries) error {
	return r.db.Create(series).Error
}

// FindByID は指定IDのシリーズをユーザー情報付きで取得する。
func (r *PostSeriesRepository) FindByID(id uint) (*model.PostSeries, error) {
	var series model.PostSeries
	err := r.db.Preload("User").First(&series, id).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

// FindByUserID は指定ユーザーの全シリーズを取得する（新しい順）。
func (r *PostSeriesRepository) FindByUserID(userID uint) ([]model.PostSeries, error) {
	var series []model.PostSeries
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&series).Error
	return series, err
}

// Update は既存のシリーズを更新する。
func (r *PostSeriesRepository) Update(series *model.PostSeries) error {
	return r.db.Save(series).Error
}

// Delete は指定IDのシリーズとその関連アイテムを削除する。
func (r *PostSeriesRepository) Delete(id uint) error {
	if err := r.db.Where("series_id = ?", id).Delete(&model.PostSeriesItem{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&model.PostSeries{}, id).Error
}

// AddPost はシリーズに投稿を追加する。
func (r *PostSeriesRepository) AddPost(item *model.PostSeriesItem) error {
	return r.db.Create(item).Error
}

// RemovePost はシリーズから投稿を削除する。
func (r *PostSeriesRepository) RemovePost(seriesID, postID uint) error {
	return r.db.Where("series_id = ? AND post_id = ?", seriesID, postID).
		Delete(&model.PostSeriesItem{}).Error
}

// GetPostsBySeriesID は指定シリーズの投稿一覧を順序付きで取得する。
func (r *PostSeriesRepository) GetPostsBySeriesID(seriesID uint) ([]model.PostSeriesItem, error) {
	var items []model.PostSeriesItem
	err := r.db.Preload("Post").Preload("Post.User").
		Where("series_id = ?", seriesID).
		Order("order_index ASC").
		Find(&items).Error
	return items, err
}
