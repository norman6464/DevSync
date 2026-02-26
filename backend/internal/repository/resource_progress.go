package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ResourceProgressRepository は学習リソース進捗のデータアクセスを提供する。
type ResourceProgressRepository struct {
	db *gorm.DB
}

// NewResourceProgressRepository は新しいResourceProgressRepositoryインスタンスを生成する。
func NewResourceProgressRepository(db *gorm.DB) *ResourceProgressRepository {
	return &ResourceProgressRepository{db: db}
}

// Upsert は進捗をUPSERTする（存在しなければ作成、存在すれば更新）。
func (r *ResourceProgressRepository) Upsert(progress *model.ResourceProgress) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "resource_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "completion_percent", "note", "started_at", "completed_at", "updated_at"}),
	}).Create(progress).Error
}

// FindByUserAndResource は指定ユーザー・リソースの進捗を取得する。
func (r *ResourceProgressRepository) FindByUserAndResource(userID, resourceID uint) (*model.ResourceProgress, error) {
	var progress model.ResourceProgress
	if err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// FindByUserID は指定ユーザーの進捗一覧を取得する。statusが空でなければフィルタする。
func (r *ResourceProgressRepository) FindByUserID(userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	var progresses []model.ResourceProgress
	var total int64

	query := r.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Model(&model.ResourceProgress{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Resource").Order("updated_at DESC").Limit(limit).Offset(offset).Find(&progresses).Error; err != nil {
		return nil, 0, err
	}

	return progresses, total, nil
}
