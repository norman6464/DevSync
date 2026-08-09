package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// RecordViewIfAbsent は (user_id, post_id) のユニーク制約と ON CONFLICT DO NOTHING を用いて
// 未閲覧のときだけ閲覧を記録し、実際に挿入できた場合のみ view_count を加算する。
// 記録と加算を単一トランザクションで行い、並行リクエストによる二重記録・二重加算・重複エラーを防ぐ。
func (r *postViewRepository) RecordViewIfAbsent(ctx context.Context, view *model.PostView) (bool, error) {
	var inserted bool
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(view)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 既に閲覧済み（競合）→ 何もしない。
			return nil
		}
		inserted = true
		return tx.Model(&model.Post{}).Where("id = ?", view.PostID).
			UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
	})
	return inserted, err
}

func (r *postViewRepository) GetViewCount(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostView{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
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
