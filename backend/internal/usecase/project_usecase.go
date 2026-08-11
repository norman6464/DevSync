package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// projectOwnerOf は所有権チェック用にプロジェクトの所有者 ID を取り出す。
func projectOwnerOf(p *model.Project) uint { return p.UserID }

// CreateProjectUseCase はプロジェクトを作成する。
type CreateProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewCreateProjectUseCase は CreateProjectUseCase を生成する。
func NewCreateProjectUseCase(projects repository.ProjectRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{projects: projects}
}

// Execute はタイトル・説明・各 URL を検証したうえでプロジェクトを作成する。
func (uc *CreateProjectUseCase) Execute(ctx context.Context, project *model.Project) error {
	v := validator.NewProjectValidator()
	if err := v.ValidateCreateProject(project.Title, project.Description, project.DemoURL, project.GithubURL); err != nil {
		return err
	}
	return uc.projects.Create(ctx, project)
}

// GetProjectUseCase はプロジェクトを 1 件取得する。
type GetProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewGetProjectUseCase は GetProjectUseCase を生成する。
func NewGetProjectUseCase(projects repository.ProjectRepository) *GetProjectUseCase {
	return &GetProjectUseCase{projects: projects}
}

// Execute は所有権を検証したうえでプロジェクトを返す。
func (uc *GetProjectUseCase) Execute(ctx context.Context, id, userID uint) (*model.Project, error) {
	return ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf)
}

// ListProjectsByUserUseCase は指定ユーザーのプロジェクト一覧を取得する。
type ListProjectsByUserUseCase struct {
	projects repository.ProjectRepository
}

// NewListProjectsByUserUseCase は ListProjectsByUserUseCase を生成する。
func NewListProjectsByUserUseCase(projects repository.ProjectRepository) *ListProjectsByUserUseCase {
	return &ListProjectsByUserUseCase{projects: projects}
}

// Execute は注目優先・作成日の新しい順でプロジェクトを返す。
func (uc *ListProjectsByUserUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	return uc.projects.FindByUserID(ctx, userID, limit, offset)
}

// ListFeaturedProjectsUseCase は注目プロジェクトの一覧を取得する。
type ListFeaturedProjectsUseCase struct {
	projects repository.ProjectRepository
}

// NewListFeaturedProjectsUseCase は ListFeaturedProjectsUseCase を生成する。
func NewListFeaturedProjectsUseCase(projects repository.ProjectRepository) *ListFeaturedProjectsUseCase {
	return &ListFeaturedProjectsUseCase{projects: projects}
}

// Execute は注目プロジェクトを作成日の新しい順で返す。
func (uc *ListFeaturedProjectsUseCase) Execute(ctx context.Context, userID uint) ([]model.Project, error) {
	return uc.projects.FindFeaturedByUserID(ctx, userID)
}

// ListArchivedProjectsUseCase はアーカイブ済みプロジェクトの一覧を取得する。
type ListArchivedProjectsUseCase struct {
	projects repository.ProjectRepository
}

// NewListArchivedProjectsUseCase は ListArchivedProjectsUseCase を生成する。
func NewListArchivedProjectsUseCase(projects repository.ProjectRepository) *ListArchivedProjectsUseCase {
	return &ListArchivedProjectsUseCase{projects: projects}
}

// Execute はアーカイブ済みプロジェクトを更新日の新しい順で返す。
func (uc *ListArchivedProjectsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	return uc.projects.FindArchivedByUserID(ctx, userID, limit, offset)
}

// ListAllProjectsUseCase は全プロジェクトの一覧を取得する。
type ListAllProjectsUseCase struct {
	projects repository.ProjectRepository
}

// NewListAllProjectsUseCase は ListAllProjectsUseCase を生成する。
func NewListAllProjectsUseCase(projects repository.ProjectRepository) *ListAllProjectsUseCase {
	return &ListAllProjectsUseCase{projects: projects}
}

// Execute は全プロジェクトを作成日の新しい順で返す。
func (uc *ListAllProjectsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.Project, int64, error) {
	return uc.projects.FindAll(ctx, limit, offset)
}

// SearchProjectsUseCase はプロジェクトをキーワード検索する。
type SearchProjectsUseCase struct {
	projects repository.ProjectRepository
}

// NewSearchProjectsUseCase は SearchProjectsUseCase を生成する。
func NewSearchProjectsUseCase(projects repository.ProjectRepository) *SearchProjectsUseCase {
	return &SearchProjectsUseCase{projects: projects}
}

// Execute は検索キーワードを検証したうえで検索する。
func (uc *SearchProjectsUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error) {
	q, err := validateSearchQuery(query)
	if err != nil {
		return nil, 0, err
	}
	return uc.projects.Search(ctx, q, limit, offset)
}

// UpdateProjectUseCase はプロジェクトを更新する。
type UpdateProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewUpdateProjectUseCase は UpdateProjectUseCase を生成する。
func NewUpdateProjectUseCase(projects repository.ProjectRepository) *UpdateProjectUseCase {
	return &UpdateProjectUseCase{projects: projects}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
// 開始日・終了日・GitHub リポジトリ ID はポインタが nil でないときだけ変更する。
func (uc *UpdateProjectUseCase) Execute(ctx context.Context, id, userID uint, updates *model.Project) (*model.Project, error) {
	project, err := ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf)
	if err != nil {
		return nil, err
	}

	v := validator.NewProjectValidator()
	if err := v.ValidateUpdateProject(updates.Title, updates.Description, updates.DemoURL, updates.GithubURL); err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		project.Title = title
	}
	if description := strings.TrimSpace(updates.Description); description != "" {
		project.Description = description
	}
	// 以降は移行前と同じく、上限だけをここで検証する。
	assignments := []struct {
		value  string
		max    int
		label  string
		target *string
	}{
		{strings.TrimSpace(updates.TechStack), 500, "技術スタック", &project.TechStack},
		{strings.TrimSpace(updates.DemoURL), 2000, "デモURL", &project.DemoURL},
		{strings.TrimSpace(updates.GithubURL), 2000, "GitHub URL", &project.GithubURL},
		{strings.TrimSpace(updates.ImageURL), 2000, "画像URL", &project.ImageURL},
		{strings.TrimSpace(updates.Role), 100, "役割", &project.Role},
	}
	for _, a := range assignments {
		if a.value == "" {
			continue
		}
		if err := domain.ValidateStringLength(a.value, 1, a.max, a.label); err != nil {
			return nil, err
		}
		*a.target = a.value
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

	if err := uc.projects.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// UpdateProjectFeaturedUseCase はプロジェクトの注目指定を切り替える。
type UpdateProjectFeaturedUseCase struct {
	projects repository.ProjectRepository
}

// NewUpdateProjectFeaturedUseCase は UpdateProjectFeaturedUseCase を生成する。
func NewUpdateProjectFeaturedUseCase(projects repository.ProjectRepository) *UpdateProjectFeaturedUseCase {
	return &UpdateProjectFeaturedUseCase{projects: projects}
}

// Execute は所有権を検証したうえで注目指定を書き換える。
func (uc *UpdateProjectFeaturedUseCase) Execute(ctx context.Context, id, userID uint, featured bool) (*model.Project, error) {
	project, err := ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf)
	if err != nil {
		return nil, err
	}

	project.Featured = featured
	if err := uc.projects.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// ArchiveProjectUseCase はプロジェクトをアーカイブする。
type ArchiveProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewArchiveProjectUseCase は ArchiveProjectUseCase を生成する。
func NewArchiveProjectUseCase(projects repository.ProjectRepository) *ArchiveProjectUseCase {
	return &ArchiveProjectUseCase{projects: projects}
}

// Execute は所有権を検証したうえでアーカイブする。既にアーカイブ済みなら 400 を返す。
func (uc *ArchiveProjectUseCase) Execute(ctx context.Context, id, userID uint) error {
	project, err := ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf)
	if err != nil {
		return err
	}
	if project.IsArchived {
		return domain.ErrBadRequest
	}
	return uc.projects.Archive(ctx, id)
}

// UnarchiveProjectUseCase はプロジェクトのアーカイブを解除する。
type UnarchiveProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewUnarchiveProjectUseCase は UnarchiveProjectUseCase を生成する。
func NewUnarchiveProjectUseCase(projects repository.ProjectRepository) *UnarchiveProjectUseCase {
	return &UnarchiveProjectUseCase{projects: projects}
}

// Execute は所有権を検証したうえでアーカイブを解除する。アーカイブされていなければ 400 を返す。
func (uc *UnarchiveProjectUseCase) Execute(ctx context.Context, id, userID uint) error {
	project, err := ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf)
	if err != nil {
		return err
	}
	if !project.IsArchived {
		return domain.ErrBadRequest
	}
	return uc.projects.Unarchive(ctx, id)
}

// DeleteProjectUseCase はプロジェクトを削除する。
type DeleteProjectUseCase struct {
	projects repository.ProjectRepository
}

// NewDeleteProjectUseCase は DeleteProjectUseCase を生成する。
func NewDeleteProjectUseCase(projects repository.ProjectRepository) *DeleteProjectUseCase {
	return &DeleteProjectUseCase{projects: projects}
}

// Execute は所有権を検証したうえでプロジェクトを削除する。
func (uc *DeleteProjectUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.projects.FindByID, id, userID, projectOwnerOf); err != nil {
		return err
	}
	return uc.projects.Delete(ctx, id)
}

// CountProjectsUseCase はプロジェクトの総数を取得する。
type CountProjectsUseCase struct {
	projects repository.ProjectRepository
}

// NewCountProjectsUseCase は CountProjectsUseCase を生成する。
func NewCountProjectsUseCase(projects repository.ProjectRepository) *CountProjectsUseCase {
	return &CountProjectsUseCase{projects: projects}
}

// Execute はプロジェクトの総数を返す。
func (uc *CountProjectsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.projects.CountByUserID(ctx, userID)
}
