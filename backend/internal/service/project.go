package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
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

// GetByID は指定IDのプロジェクトを取得する。所有権を検証する。
func (s *ProjectService) GetByID(id, userID uint) (*model.Project, error) {
	return s.findAndCheckOwnership(id, userID)
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

// Search はプロジェクトをキーワード検索する。
func (s *ProjectService) Search(query string, limit, offset int) ([]model.Project, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, "検索キーワードは必須です", nil)
	}
	return s.repo.Search(strings.TrimSpace(query), limit, offset)
}

// findAndCheckOwnership はプロジェクトを取得し、指定ユーザーが所有者かを検証する。
func (s *ProjectService) findAndCheckOwnership(id, userID uint) (*model.Project, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(p *model.Project) uint { return p.UserID })
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
		if len([]rune(strings.TrimSpace(updates.TechStack))) > 500 {
			return nil, domain.NewError(domain.ErrCodeValidation, "技術スタックは500文字以下である必要があります", nil)
		}
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
		if len([]rune(strings.TrimSpace(updates.Role))) > 100 {
			return nil, domain.NewError(domain.ErrCodeValidation, "役割は100文字以下である必要があります", nil)
		}
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

// Archive は所有権を検証した後、プロジェクトをアーカイブする。
func (s *ProjectService) Archive(id, userID uint) error {
	project, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	if project.IsArchived {
		return ErrBadRequest
	}
	return s.repo.Archive(id)
}

// Unarchive は所有権を検証した後、プロジェクトのアーカイブを解除する。
func (s *ProjectService) Unarchive(id, userID uint) error {
	project, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	if !project.IsArchived {
		return ErrBadRequest
	}
	return s.repo.Unarchive(id)
}

// GetArchivedByUserID は指定ユーザーのアーカイブ済みプロジェクトを取得する。
func (s *ProjectService) GetArchivedByUserID(userID uint, limit, offset int) ([]model.Project, int64, error) {
	return s.repo.FindArchivedByUserID(userID, limit, offset)
}

// Delete は所有権を検証した後、プロジェクトを削除する。
func (s *ProjectService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// CountByUserID は指定ユーザーのプロジェクト総数を返す。
func (s *ProjectService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}
