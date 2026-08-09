package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postViewRepository は [repository.PostViewRepository] の GORM 実装。
type postViewRepository struct {
	db *gorm.DB
}

// NewPostViewRepository は PostViewRepository の GORM 実装を返す。
func NewPostViewRepository(db *gorm.DB) repository.PostViewRepository {
	return &postViewRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostViewRepository = (*postViewRepository)(nil)

// RecordView は閲覧を記録し、対象投稿の view_count をインクリメントする。
func (r *postViewRepository) RecordView(ctx context.Context, view *model.PostView) error {
	if err := r.db.WithContext(ctx).Create(view).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", view.PostID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *postViewRepository) GetViewCount(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostView{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (r *postViewRepository) HasViewed(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *postViewRepository) GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error) {
	var results []model.ViewCount
	err := r.db.WithContext(ctx).Model(&model.PostView{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}
