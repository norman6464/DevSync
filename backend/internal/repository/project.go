package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ProjectRepository はプロジェクトショーケースデータへのアクセスを提供するリポジトリ実装。
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository は新しいProjectRepositoryインスタンスを生成する。
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create は新しいプロジェクトをデータベースに作成する。
func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

// FindByID は指定IDのプロジェクトをユーザー・GitHubリポジトリ情報付きで取得する。
func (r *ProjectRepository) FindByID(id uint) (*model.Project, error) {
	var project model.Project
	err := r.db.Preload("User").Preload("GithubRepo").First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindByUserID は指定ユーザーの全プロジェクトをGitHubリポジトリ情報付きで取得する。
// featured（注目）プロジェクトが先に、その後作成日降順でソートされる。
func (r *ProjectRepository) FindByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Preload("GithubRepo").
		Where("user_id = ?", userID).
		Order("featured DESC, created_at DESC").
		Find(&projects).Error
	return projects, err
}

// FindFeaturedByUserID は指定ユーザーの注目プロジェクトのみを取得する。
func (r *ProjectRepository) FindFeaturedByUserID(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Preload("GithubRepo").
		Where("user_id = ? AND featured = ?", userID, true).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

// Update は既存のプロジェクトを更新する。
func (r *ProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

// Delete は指定IDのプロジェクトを削除する。
func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&model.Project{}, id).Error
}

// FindAll は全プロジェクトをページネーション付きで取得する。
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
