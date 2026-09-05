package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// projectRepository は [repository.ProjectRepository] の sqlc(pgx) 実装。
// projects は論理削除を使うため、全クエリで deleted_at IS NULL を明示する
// （GORM が自動的に付与していたスコープ相当）。
type projectRepository struct {
	q *sqlcgen.Queries
}

// NewProjectRepository は ProjectRepository の sqlc(pgx) 実装を返す。
func NewProjectRepository(q *sqlcgen.Queries) repository.ProjectRepository {
	return &projectRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectRepository = (*projectRepository)(nil)

// attachProjectGithubRepo は LEFT JOIN で取得した git_hub_repositories の個別カラムを
// model.Project へ GithubRepo として付与する。
func attachProjectGithubRepo(
	project *model.Project,
	repoID, repoUserID, repoGitHubRepoID *int64,
	repoName, repoFullName, repoDescription, repoLanguage *string,
	repoStars, repoForks *int64,
	repoIsPrivate *bool,
	repoUpdatedAt pgtype.Timestamptz,
) {
	if repoID == nil {
		return
	}
	project.GithubRepo = &model.GitHubRepository{
		ID:           uint(*repoID),
		UserID:       uint(fromInt64PtrValue(repoUserID)),
		GitHubRepoID: fromInt64PtrValue(repoGitHubRepoID),
		Name:         fromStringPtr(repoName),
		FullName:     fromStringPtr(repoFullName),
		Description:  fromStringPtr(repoDescription),
		Language:     fromStringPtr(repoLanguage),
		Stars:        int(fromInt64PtrValue(repoStars)),
		Forks:        int(fromInt64PtrValue(repoForks)),
		IsPrivate:    fromBoolPtr(repoIsPrivate),
		UpdatedAt:    timeValue(fromTimestamptz(repoUpdatedAt)),
	}
}

// Create は新しいプロジェクトを作成する。
func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	row, err := r.q.CreateProject(ctx, sqlcgen.CreateProjectParams{
		UserID:       int64(project.UserID),
		Title:        project.Title,
		Description:  &project.Description,
		TechStack:    &project.TechStack,
		DemoUrl:      &project.DemoURL,
		GithubUrl:    &project.GithubURL,
		ImageUrl:     &project.ImageURL,
		Role:         &project.Role,
		StartDate:    toTimestamptz(project.StartDate),
		EndDate:      toTimestamptz(project.EndDate),
		Featured:     &project.Featured,
		IsArchived:   &project.IsArchived,
		GithubRepoID: toInt64PtrFromUintPtr(project.GithubRepoID),
	})
	if err != nil {
		return err
	}
	*project = toModelProject(row)
	return nil
}

// Update は既存のプロジェクトを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	row, err := r.q.UpdateProject(ctx, sqlcgen.UpdateProjectParams{
		ID:           int64(project.ID),
		Title:        project.Title,
		Description:  &project.Description,
		TechStack:    &project.TechStack,
		DemoUrl:      &project.DemoURL,
		GithubUrl:    &project.GithubURL,
		ImageUrl:     &project.ImageURL,
		Role:         &project.Role,
		StartDate:    toTimestamptz(project.StartDate),
		EndDate:      toTimestamptz(project.EndDate),
		Featured:     &project.Featured,
		IsArchived:   &project.IsArchived,
		GithubRepoID: toInt64PtrFromUintPtr(project.GithubRepoID),
	})
	if err != nil {
		return err
	}
	*project = toModelProject(row)
	return nil
}

// Delete はプロジェクトを削除する（モデルが論理削除を持つため soft delete になる）。
func (r *projectRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteProject(ctx, int64(id))
}

// FindByID は指定 ID のプロジェクトをユーザー・GitHub リポジトリ付きで取得する。
// 不在の場合は (nil, nil) を返す。
func (r *projectRepository) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	row, err := r.q.GetProjectWithUserAndRepoByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	project := toModelProject(row.Project)
	project.User = toModelUser(row.User)
	attachProjectGithubRepo(&project, row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
		row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
		row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	return &project, nil
}

// FindByUserID はユーザーのプロジェクトを注目優先・作成日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	total, err := r.q.CountProjectsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListProjectsByUserWithRepo(ctx, sqlcgen.ListProjectsByUserWithRepoParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	projects := make([]model.Project, len(rows))
	for i, row := range rows {
		projects[i] = toModelProject(row.Project)
		attachProjectGithubRepo(&projects[i], row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
			row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
			row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	}
	return projects, total, nil
}

// FindArchivedByUserID はアーカイブ済みのプロジェクトを更新日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindArchivedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	total, err := r.q.CountArchivedProjectsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListArchivedProjectsByUserWithRepo(ctx, sqlcgen.ListArchivedProjectsByUserWithRepoParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	projects := make([]model.Project, len(rows))
	for i, row := range rows {
		projects[i] = toModelProject(row.Project)
		attachProjectGithubRepo(&projects[i], row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
			row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
			row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	}
	return projects, total, nil
}

// FindAll は全プロジェクトを作成日の新しい順で取得し、総数も返す。
func (r *projectRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Project, int64, error) {
	total, err := r.q.CountAllProjects(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListAllProjectsWithUserAndRepo(ctx, sqlcgen.ListAllProjectsWithUserAndRepoParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	projects := make([]model.Project, len(rows))
	for i, row := range rows {
		projects[i] = toModelProject(row.Project)
		projects[i].User = toModelUser(row.User)
		attachProjectGithubRepo(&projects[i], row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
			row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
			row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	}
	return projects, total, nil
}

// Search はタイトル・説明・技術スタックへの部分一致で検索する（大文字小文字を区別しない）。
func (r *projectRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error) {
	like := escapeLikePattern(query)

	total, err := r.q.CountSearchProjects(ctx, like)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchProjectsWithUserAndRepo(ctx, sqlcgen.SearchProjectsWithUserAndRepoParams{
		Title:  like,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	projects := make([]model.Project, len(rows))
	for i, row := range rows {
		projects[i] = toModelProject(row.Project)
		projects[i].User = toModelUser(row.User)
		attachProjectGithubRepo(&projects[i], row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
			row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
			row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	}
	return projects, total, nil
}

// FindFeaturedByUserID は注目のプロジェクトを作成日の新しい順で取得する。
func (r *projectRepository) FindFeaturedByUserID(ctx context.Context, userID uint) ([]model.Project, error) {
	rows, err := r.q.ListFeaturedProjectsByUserWithRepo(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	projects := make([]model.Project, len(rows))
	for i, row := range rows {
		projects[i] = toModelProject(row.Project)
		attachProjectGithubRepo(&projects[i], row.RepoID, row.RepoUserID, row.RepoGitHubRepoID,
			row.RepoName, row.RepoFullName, row.RepoDescription, row.RepoLanguage,
			row.RepoStars, row.RepoForks, row.RepoIsPrivate, row.RepoUpdatedAt)
	}
	return projects, nil
}

// CountByUserID はユーザーのプロジェクト総数を返す。
func (r *projectRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountProjectsByUser(ctx, int64(userID))
}

// Archive はプロジェクトをアーカイブする。
func (r *projectRepository) Archive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, true)
}

// Unarchive はプロジェクトのアーカイブを解除する。
func (r *projectRepository) Unarchive(ctx context.Context, id uint) error {
	return r.setArchived(ctx, id, false)
}

// setArchived はアーカイブ状態を更新する共通処理。
func (r *projectRepository) setArchived(ctx context.Context, id uint, archived bool) error {
	return r.q.SetProjectArchived(ctx, sqlcgen.SetProjectArchivedParams{
		ID:         int64(id),
		IsArchived: &archived,
	})
}
