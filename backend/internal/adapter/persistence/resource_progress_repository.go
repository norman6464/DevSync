package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// resourceProgressRepository は [repository.ResourceProgressRepository] の GORM 実装。
type resourceProgressRepository struct {
	db *gorm.DB
}

// NewResourceProgressRepository は ResourceProgressRepository の GORM 実装を返す。
func NewResourceProgressRepository(db *gorm.DB) repository.ResourceProgressRepository {
	return &resourceProgressRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ResourceProgressRepository = (*resourceProgressRepository)(nil)

// Upsert は (user_id, resource_id) をキーに進捗を作成または更新する。
func (r *resourceProgressRepository) Upsert(ctx context.Context, progress *model.ResourceProgress) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "resource_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "completion_percent", "note", "started_at", "completed_at", "updated_at"}),
	}).Create(progress).Error
}

// FindByUserAndResource は指定ユーザー・リソースの進捗を取得する。
func (r *resourceProgressRepository) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error) {
	var progress model.ResourceProgress
	if err := r.db.WithContext(ctx).Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// FindByUserID は指定ユーザーの進捗一覧を取得する。status が空でなければ絞り込む。
func (r *resourceProgressRepository) FindByUserID(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	var progresses []model.ResourceProgress
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&model.ResourceProgress{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// updated_at 同値の行でもページングが安定するよう id を第 2 ソートキーにして順序を決定的にする。
	if err := query.Preload("Resource").Order("updated_at DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&progresses).Error; err != nil {
		return nil, 0, err
	}

	return progresses, total, nil
}
