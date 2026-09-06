package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// projectReader は [repository.ProjectReader] の sqlc(pgx) 実装。
// 所有権チェックに必要なプロジェクト読み取りだけを提供する。
type projectReader struct {
	q *sqlcgen.Queries
}

// NewProjectReader は ProjectReader の sqlc(pgx) 実装を返す。
func NewProjectReader(q *sqlcgen.Queries) repository.ProjectReader {
	return &projectReader{q: q}
}

var _ repository.ProjectReader = (*projectReader)(nil)

// toModelProject は sqlc の生成行を model.Project へ変換する（関連の User/GithubRepo は含まない）。
func toModelProject(row sqlcgen.Project) model.Project {
	return model.Project{
		ID:           uint(row.ID),
		UserID:       uint(row.UserID),
		Title:        row.Title,
		Description:  fromStringPtr(row.Description),
		TechStack:    fromStringPtr(row.TechStack),
		DemoURL:      fromStringPtr(row.DemoUrl),
		GithubURL:    fromStringPtr(row.GithubUrl),
		ImageURL:     fromStringPtr(row.ImageUrl),
		Role:         fromStringPtr(row.Role),
		StartDate:    fromTimestamptz(row.StartDate),
		EndDate:      fromTimestamptz(row.EndDate),
		Featured:     row.Featured,
		IsArchived:   row.IsArchived,
		GithubRepoID: fromInt64PtrToUintPtr(row.GithubRepoID),
		CreatedAt:    timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:    timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// FindByID は ID でプロジェクトを取得する（既存 ProjectRepository.FindByID と同じ preload）。
// 不在の場合は (nil, nil) を返す。
func (r *projectReader) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	row, err := r.q.GetProjectWithUserAndRepoByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	project := toModelProject(row.Project)
	project.User = toModelUser(row.User)
	if row.RepoID != nil {
		project.GithubRepo = &model.GitHubRepository{
			ID:           uint(*row.RepoID),
			UserID:       uint(fromInt64PtrValue(row.RepoUserID)),
			GitHubRepoID: fromInt64PtrValue(row.RepoGitHubRepoID),
			Name:         fromStringPtr(row.RepoName),
			FullName:     fromStringPtr(row.RepoFullName),
			Description:  fromStringPtr(row.RepoDescription),
			Language:     fromStringPtr(row.RepoLanguage),
			Stars:        int(fromInt64PtrValue(row.RepoStars)),
			Forks:        int(fromInt64PtrValue(row.RepoForks)),
			IsPrivate:    fromBoolPtr(row.RepoIsPrivate),
			UpdatedAt:    timeValue(fromTimestamptz(row.RepoUpdatedAt)),
		}
	}
	return &project, nil
}
