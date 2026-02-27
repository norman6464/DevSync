package service

import (
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// validMilestoneStatuses は有効なマイルストーンステータスの一覧。
var validMilestoneStatuses = map[string]bool{
	string(model.MilestoneNotStarted): true,
	string(model.MilestoneInProgress): true,
	string(model.MilestoneCompleted):  true,
}

// ProjectMilestoneService はプロジェクトマイルストーンのビジネスロジックを提供する。
type ProjectMilestoneService struct {
	milestoneRepo repository.ProjectMilestoneRepositoryInterface
	projectRepo   repository.ProjectRepositoryInterface
}

// NewProjectMilestoneService は新しいProjectMilestoneServiceインスタンスを生成する。
func NewProjectMilestoneService(milestoneRepo repository.ProjectMilestoneRepositoryInterface, projectRepo repository.ProjectRepositoryInterface) *ProjectMilestoneService {
	return &ProjectMilestoneService{milestoneRepo: milestoneRepo, projectRepo: projectRepo}
}

// checkProjectOwnership はプロジェクトの所有権を検証する。
func (s *ProjectMilestoneService) checkProjectOwnership(userID, projectID uint) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "プロジェクトが見つかりません", err)
	}
	if project.UserID != userID {
		return domain.NewError(domain.ErrCodeForbidden, "このプロジェクトを編集する権限がありません", nil)
	}
	return nil
}

// Create はマイルストーンを作成する。
func (s *ProjectMilestoneService) Create(userID, projectID uint, title, description string, dueDate *time.Time) error {
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
		return err
	}

	if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(description, 0, 1000, "説明"); err != nil {
		return err
	}

	milestone := &model.ProjectMilestone{
		ProjectID:   projectID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		DueDate:     dueDate,
		Status:      model.MilestoneNotStarted,
	}

	return s.milestoneRepo.Create(milestone)
}

// GetByProjectID はプロジェクトのマイルストーン一覧を取得する。
func (s *ProjectMilestoneService) GetByProjectID(projectID uint) ([]model.ProjectMilestone, error) {
	return s.milestoneRepo.FindByProjectID(projectID)
}

// Update はマイルストーンを更新する。
func (s *ProjectMilestoneService) Update(userID, milestoneID uint, title, description string, dueDate *time.Time, status string) (*model.ProjectMilestone, error) {
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "マイルストーンが見つかりません", err)
	}

	if err := s.checkProjectOwnership(userID, milestone.ProjectID); err != nil {
		return nil, err
	}

	if t := strings.TrimSpace(title); t != "" {
		if err := domain.ValidateStringLength(t, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		milestone.Title = t
	}
	if d := strings.TrimSpace(description); d != "" {
		if err := domain.ValidateStringLength(d, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		milestone.Description = d
	}
	if dueDate != nil {
		milestone.DueDate = dueDate
	}
	if status != "" {
		if !validMilestoneStatuses[status] {
			return nil, domain.NewError(domain.ErrCodeValidation, "無効なステータスです。not_started, in_progress, completed のいずれかを指定してください", nil)
		}
		milestone.Status = model.MilestoneStatus(status)
		if status == string(model.MilestoneCompleted) {
			now := time.Now()
			milestone.CompletedAt = &now
		} else {
			milestone.CompletedAt = nil
		}
	}

	if err := s.milestoneRepo.Update(milestone); err != nil {
		return nil, err
	}
	return milestone, nil
}

// Delete はマイルストーンを削除する。
func (s *ProjectMilestoneService) Delete(userID, milestoneID uint) error {
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "マイルストーンが見つかりません", err)
	}

	if err := s.checkProjectOwnership(userID, milestone.ProjectID); err != nil {
		return err
	}

	return s.milestoneRepo.Delete(milestoneID)
}
