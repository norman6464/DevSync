package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ProjectService はプロジェクトショーケースのビジネスロジックを提供する。
// プロジェクトのCRUD操作と注目（featured）ステータス管理を担当する。
type ProjectService struct {
	repo repository.ProjectRepositoryInterface
}

// NewProjectService は新しいProjectServiceインスタンスを生成する。
func NewProjectService(repo repository.ProjectRepositoryInterface) *ProjectService {
	return &ProjectService{repo: repo}
}

// Create は新しいプロジェクトを作成する。
func (s *ProjectService) Create(project *model.Project) error {
	v := validator.NewProjectValidator()
	if err := v.ValidateCreateProject(project.Title, project.Description, project.DemoURL, project.GithubURL); err != nil {
		return err
	}
	return s.repo.Create(project)
}

// GetByID は指定IDのプロジェクトを取得する。
func (s *ProjectService) GetByID(id uint) (*model.Project, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーのプロジェクトをページネーション付きで取得する。
func (s *ProjectService) GetByUserID(userID uint, limit, offset int) ([]model.Project, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetFeaturedByUserID は指定ユーザーの注目プロジェクトのみを取得する。
func (s *ProjectService) GetFeaturedByUserID(userID uint) ([]model.Project, error) {
	return s.repo.FindFeaturedByUserID(userID)
}

// GetAll はプロジェクト一覧をページネーション付きで取得する。
func (s *ProjectService) GetAll(limit, offset int) ([]model.Project, int64, error) {
	return s.repo.FindAll(limit, offset)
}

// findAndCheckOwnership はプロジェクトを取得し、指定ユーザーが所有者かを検証する。
func (s *ProjectService) findAndCheckOwnership(id, userID uint) (*model.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrForbidden
	}
	return project, nil
}

// Update は所有権を検証した後、プロジェクトを更新する。
func (s *ProjectService) Update(id, userID uint, updates *model.Project) (*model.Project, error) {
	project, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	v := validator.NewProjectValidator()
	if err := v.ValidateUpdateProject(updates.Title, updates.Description, updates.DemoURL, updates.GithubURL); err != nil {
		return nil, err
	}

	if strings.TrimSpace(updates.Title) != "" {
		project.Title = strings.TrimSpace(updates.Title)
	}
	if strings.TrimSpace(updates.Description) != "" {
		project.Description = strings.TrimSpace(updates.Description)
	}
	if strings.TrimSpace(updates.TechStack) != "" {
		project.TechStack = strings.TrimSpace(updates.TechStack)
	}
	if strings.TrimSpace(updates.DemoURL) != "" {
		project.DemoURL = strings.TrimSpace(updates.DemoURL)
	}
	if strings.TrimSpace(updates.GithubURL) != "" {
		project.GithubURL = strings.TrimSpace(updates.GithubURL)
	}
	if strings.TrimSpace(updates.ImageURL) != "" {
		project.ImageURL = strings.TrimSpace(updates.ImageURL)
	}
	if strings.TrimSpace(updates.Role) != "" {
		project.Role = strings.TrimSpace(updates.Role)
	}
	if updates.StartDate != nil {
		project.StartDate = updates.StartDate
	}
	if updates.EndDate != nil {
		project.EndDate = updates.EndDate
	}
	if updates.GithubRepoID != nil {
		project.GithubRepoID = updates.GithubRepoID
	}

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

// UpdateFeatured は所有権を検証した後、プロジェクトの注目ステータスを更新する。
func (s *ProjectService) UpdateFeatured(id, userID uint, featured bool) (*model.Project, error) {
	project, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	project.Featured = featured

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

// Delete は所有権を検証した後、プロジェクトを削除する。
func (s *ProjectService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
