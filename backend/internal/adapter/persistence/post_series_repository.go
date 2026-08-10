package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postSeriesRepository は [repository.PostSeriesRepository] の GORM 実装。
type postSeriesRepository struct {
	db *gorm.DB
}

// NewPostSeriesRepository は PostSeriesRepository の GORM 実装を返す。
func NewPostSeriesRepository(db *gorm.DB) repository.PostSeriesRepository {
	return &postSeriesRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostSeriesRepository = (*postSeriesRepository)(nil)

// Create は新しい投稿シリーズをデータベースに作成する。
func (r *postSeriesRepository) Create(ctx context.Context, series *model.PostSeries) error {
	return r.db.WithContext(ctx).Create(series).Error
}

// FindByID は指定IDのシリーズをユーザー情報付きで取得する。
func (r *postSeriesRepository) FindByID(ctx context.Context, id uint) (*model.PostSeries, error) {
	var series model.PostSeries
	if err := r.db.WithContext(ctx).Preload("User").First(&series, id).Error; err != nil {
		return nil, err
	}
	return &series, nil
}

// FindByUserID は指定ユーザーのシリーズをページネーション付きで取得する（新しい順）。
func (r *postSeriesRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.PostSeries, error) {
	var series []model.PostSeries
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&series).Error
	return series, err
}

// CountByUser は指定ユーザーのシリーズ数をカウントする。
func (r *postSeriesRepository) CountByUser(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostSeries{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Update は既存のシリーズを更新する。
func (r *postSeriesRepository) Update(ctx context.Context, series *model.PostSeries) error {
	return r.db.WithContext(ctx).Save(series).Error
}

// Delete は指定IDのシリーズとその関連アイテムを削除する。
func (r *postSeriesRepository) Delete(ctx context.Context, id uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Where("series_id = ?", id).Delete(&model.PostSeriesItem{}).Error; err != nil {
		return err
	}
	return db.Delete(&model.PostSeries{}, id).Error
}

// AddPost はシリーズに投稿を追加する。
func (r *postSeriesRepository) AddPost(ctx context.Context, item *model.PostSeriesItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// HasPost は指定シリーズに指定投稿が存在するかを確認する。
func (r *postSeriesRepository) HasPost(ctx context.Context, seriesID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostSeriesItem{}).
		Where("series_id = ? AND post_id = ?", seriesID, postID).
		Count(&count).Error
	return count > 0, err
}

// RemovePost はシリーズから投稿を削除する。
func (r *postSeriesRepository) RemovePost(ctx context.Context, seriesID, postID uint) error {
	return r.db.WithContext(ctx).Where("series_id = ? AND post_id = ?", seriesID, postID).
		Delete(&model.PostSeriesItem{}).Error
}

// GetPostsBySeriesID は指定シリーズの投稿一覧を順序付きで取得する。
func (r *postSeriesRepository) GetPostsBySeriesID(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error) {
	var items []model.PostSeriesItem
	err := r.db.WithContext(ctx).Preload("Post").Preload("Post.User").
		Where("series_id = ?", seriesID).
		Order("order_index ASC").
		Find(&items).Error
	return items, err
}
