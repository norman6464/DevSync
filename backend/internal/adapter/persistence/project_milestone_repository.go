package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// projectMilestoneRepository は [repository.ProjectMilestoneRepository] の GORM 実装。
type projectMilestoneRepository struct {
	db *gorm.DB
}

// NewProjectMilestoneRepository は ProjectMilestoneRepository の GORM 実装を返す。
func NewProjectMilestoneRepository(db *gorm.DB) repository.ProjectMilestoneRepository {
	return &projectMilestoneRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectMilestoneRepository = (*projectMilestoneRepository)(nil)

// Create はマイルストーンを作成する。
func (r *projectMilestoneRepository) Create(ctx context.Context, milestone *model.ProjectMilestone) error {
	return r.db.WithContext(ctx).Create(milestone).Error
}

// FindByID は指定 ID のマイルストーンを取得する。不在の場合は (nil, nil) を返す。
func (r *projectMilestoneRepository) FindByID(ctx context.Context, id uint) (*model.ProjectMilestone, error) {
	var milestone model.ProjectMilestone
	err := r.db.WithContext(ctx).First(&milestone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// FindByProjectID は指定プロジェクトのマイルストーン一覧を期日順で取得する。
func (r *projectMilestoneRepository) FindByProjectID(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error) {
	var milestones []model.ProjectMilestone
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("due_date ASC NULLS LAST, created_at ASC").
		Find(&milestones).Error; err != nil {
		return nil, err
	}
	return milestones, nil
}

// Update はマイルストーンを更新する。
func (r *projectMilestoneRepository) Update(ctx context.Context, milestone *model.ProjectMilestone) error {
	return r.db.WithContext(ctx).Save(milestone).Error
}

// Delete はマイルストーンを削除する。
func (r *projectMilestoneRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ProjectMilestone{}, id).Error
}
