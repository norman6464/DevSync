package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ProjectMilestoneRepository はプロジェクトマイルストーンのデータアクセスを提供する。
type ProjectMilestoneRepository struct {
	db *gorm.DB
}

// NewProjectMilestoneRepository は新しいProjectMilestoneRepositoryインスタンスを生成する。
func NewProjectMilestoneRepository(db *gorm.DB) *ProjectMilestoneRepository {
	return &ProjectMilestoneRepository{db: db}
}

// Create はマイルストーンを作成する。
func (r *ProjectMilestoneRepository) Create(milestone *model.ProjectMilestone) error {
	return r.db.Create(milestone).Error
}

// FindByID は指定IDのマイルストーンを取得する。
func (r *ProjectMilestoneRepository) FindByID(id uint) (*model.ProjectMilestone, error) {
	var milestone model.ProjectMilestone
	if err := r.db.First(&milestone, id).Error; err != nil {
		return nil, err
	}
	return &milestone, nil
}

// FindByProjectID は指定プロジェクトのマイルストーン一覧を期日順で取得する。
func (r *ProjectMilestoneRepository) FindByProjectID(projectID uint) ([]model.ProjectMilestone, error) {
	var milestones []model.ProjectMilestone
	if err := r.db.Where("project_id = ?", projectID).Order("due_date ASC NULLS LAST, created_at ASC").Find(&milestones).Error; err != nil {
		return nil, err
	}
	return milestones, nil
}

// Update はマイルストーンを更新する。
func (r *ProjectMilestoneRepository) Update(milestone *model.ProjectMilestone) error {
	return r.db.Save(milestone).Error
}

// Delete はマイルストーンを削除する。
func (r *ProjectMilestoneRepository) Delete(id uint) error {
	return r.db.Delete(&model.ProjectMilestone{}, id).Error
}
