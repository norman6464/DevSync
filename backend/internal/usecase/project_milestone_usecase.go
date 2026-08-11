package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// validMilestoneStatuses は有効なマイルストーンステータスの一覧。
var validMilestoneStatuses = map[string]bool{
	string(model.MilestoneNotStarted): true,
	string(model.MilestoneInProgress): true,
	string(model.MilestoneCompleted):  true,
}

// requireOwnedProject はプロジェクトを取得し、userID が所有者であることを検証する。
//
// 共通 helper の ensureOwner は使わない。このスライスは所有者でないときに専用のメッセージを
// 返しており、汎用メッセージに変わってしまうため（移行前の挙動を維持している）。
// 不在のときも DB 障害のときも同じ 404 を返す点も移行前どおり。
func requireOwnedProject(ctx context.Context, projects repository.ProjectReader, userID, projectID uint) error {
	project, err := projects.FindByID(ctx, projectID)
	if err != nil || project == nil {
		return domain.NewError(domain.ErrCodeNotFound, "プロジェクトが見つかりません", err)
	}
	if project.UserID != userID {
		return domain.NewError(domain.ErrCodeForbidden, "このプロジェクトを編集する権限がありません", nil)
	}
	return nil
}

// findMilestone はマイルストーンを取得する。不在なら 404 を返す。
func findMilestone(ctx context.Context, milestones repository.ProjectMilestoneRepository, milestoneID uint) (*model.ProjectMilestone, error) {
	milestone, err := milestones.FindByID(ctx, milestoneID)
	if err != nil || milestone == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "マイルストーンが見つかりません", err)
	}
	return milestone, nil
}

// CreateProjectMilestoneInput はマイルストーン作成の入力。
type CreateProjectMilestoneInput struct {
	UserID      uint
	ProjectID   uint
	Title       string
	Description string
	DueDate     *time.Time
}

// CreateProjectMilestoneUseCase はマイルストーンを作成する。
type CreateProjectMilestoneUseCase struct {
	milestones repository.ProjectMilestoneRepository
	projects   repository.ProjectReader
}

// NewCreateProjectMilestoneUseCase は CreateProjectMilestoneUseCase を生成する。
func NewCreateProjectMilestoneUseCase(
	milestones repository.ProjectMilestoneRepository,
	projects repository.ProjectReader,
) *CreateProjectMilestoneUseCase {
	return &CreateProjectMilestoneUseCase{milestones: milestones, projects: projects}
}

// Execute はプロジェクトの所有権と入力を検査したうえでマイルストーンを作成する。
func (uc *CreateProjectMilestoneUseCase) Execute(ctx context.Context, in CreateProjectMilestoneInput) error {
	if err := requireOwnedProject(ctx, uc.projects, in.UserID, in.ProjectID); err != nil {
		return err
	}

	if err := domain.ValidateStringLength(in.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(in.Description, 0, 1000, "説明"); err != nil {
		return err
	}

	return uc.milestones.Create(ctx, &model.ProjectMilestone{
		ProjectID:   in.ProjectID,
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		DueDate:     in.DueDate,
		Status:      model.MilestoneNotStarted,
	})
}

// ListProjectMilestonesUseCase はプロジェクトのマイルストーン一覧を取得する。
type ListProjectMilestonesUseCase struct {
	milestones repository.ProjectMilestoneRepository
}

// NewListProjectMilestonesUseCase は ListProjectMilestonesUseCase を生成する。
func NewListProjectMilestonesUseCase(milestones repository.ProjectMilestoneRepository) *ListProjectMilestonesUseCase {
	return &ListProjectMilestonesUseCase{milestones: milestones}
}

// Execute はマイルストーン一覧を返す。所有権は検証しない（移行前の挙動を維持している）。
func (uc *ListProjectMilestonesUseCase) Execute(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error) {
	return uc.milestones.FindByProjectID(ctx, projectID)
}

// UpdateProjectMilestoneInput はマイルストーン更新の入力。
// 前後の空白を除いて空になる項目と、nil の期日は据え置く部分更新。
type UpdateProjectMilestoneInput struct {
	UserID      uint
	MilestoneID uint
	Title       string
	Description string
	DueDate     *time.Time
	Status      string
}

// UpdateProjectMilestoneUseCase はマイルストーンを更新する。
type UpdateProjectMilestoneUseCase struct {
	milestones repository.ProjectMilestoneRepository
	projects   repository.ProjectReader
}

// NewUpdateProjectMilestoneUseCase は UpdateProjectMilestoneUseCase を生成する。
func NewUpdateProjectMilestoneUseCase(
	milestones repository.ProjectMilestoneRepository,
	projects repository.ProjectReader,
) *UpdateProjectMilestoneUseCase {
	return &UpdateProjectMilestoneUseCase{milestones: milestones, projects: projects}
}

// Execute はマイルストーンを部分更新する。完了へ遷移した場合のみ完了日時を記録する。
func (uc *UpdateProjectMilestoneUseCase) Execute(ctx context.Context, in UpdateProjectMilestoneInput) (*model.ProjectMilestone, error) {
	milestone, err := findMilestone(ctx, uc.milestones, in.MilestoneID)
	if err != nil {
		return nil, err
	}

	if err := requireOwnedProject(ctx, uc.projects, in.UserID, milestone.ProjectID); err != nil {
		return nil, err
	}

	if t := strings.TrimSpace(in.Title); t != "" {
		if err := domain.ValidateStringLength(t, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		milestone.Title = t
	}
	if d := strings.TrimSpace(in.Description); d != "" {
		if err := domain.ValidateStringLength(d, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		milestone.Description = d
	}
	if in.DueDate != nil {
		milestone.DueDate = in.DueDate
	}
	if in.Status != "" {
		if !validMilestoneStatuses[in.Status] {
			return nil, domain.NewError(domain.ErrCodeValidation,
				"無効なステータスです。not_started, in_progress, completed のいずれかを指定してください", nil)
		}
		milestone.Status = model.MilestoneStatus(in.Status)
		if in.Status == string(model.MilestoneCompleted) {
			now := time.Now()
			milestone.CompletedAt = &now
		} else {
			milestone.CompletedAt = nil
		}
	}

	if err := uc.milestones.Update(ctx, milestone); err != nil {
		return nil, err
	}
	return milestone, nil
}

// DeleteProjectMilestoneUseCase はマイルストーンを削除する。
type DeleteProjectMilestoneUseCase struct {
	milestones repository.ProjectMilestoneRepository
	projects   repository.ProjectReader
}

// NewDeleteProjectMilestoneUseCase は DeleteProjectMilestoneUseCase を生成する。
func NewDeleteProjectMilestoneUseCase(
	milestones repository.ProjectMilestoneRepository,
	projects repository.ProjectReader,
) *DeleteProjectMilestoneUseCase {
	return &DeleteProjectMilestoneUseCase{milestones: milestones, projects: projects}
}

// Execute はプロジェクトの所有権を検証したうえでマイルストーンを削除する。
func (uc *DeleteProjectMilestoneUseCase) Execute(ctx context.Context, userID, milestoneID uint) error {
	milestone, err := findMilestone(ctx, uc.milestones, milestoneID)
	if err != nil {
		return err
	}
	if err := requireOwnedProject(ctx, uc.projects, userID, milestone.ProjectID); err != nil {
		return err
	}
	return uc.milestones.Delete(ctx, milestoneID)
}
