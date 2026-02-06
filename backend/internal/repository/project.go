package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) FindByID(id uint) (*model.Project, error) {
	var project model.Project
	err := r.db.Preload("User").Preload("GithubRepo").First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Preload("GithubRepo").
		Where("user_id = ?", userID).
		Order("featured DESC, created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) FindFeaturedByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Preload("GithubRepo").
		Where("user_id = ? AND featured = ?", userID, true).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&model.Project{}, id).Error
}

func (r *ProjectRepository) FindAll(limit, offset int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	r.db.Model(&model.Project{}).Count(&total)

	err := r.db.Preload("User").Preload("GithubRepo").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&projects).Error

	return projects, total, err
}
