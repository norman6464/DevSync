package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// projectMilestoneRepository は [repository.ProjectMilestoneRepository] の sqlc(pgx) 実装。
type projectMilestoneRepository struct {
	q *sqlcgen.Queries
}

// NewProjectMilestoneRepository は ProjectMilestoneRepository の sqlc(pgx) 実装を返す。
func NewProjectMilestoneRepository(q *sqlcgen.Queries) repository.ProjectMilestoneRepository {
	return &projectMilestoneRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectMilestoneRepository = (*projectMilestoneRepository)(nil)

// toModelProjectMilestone は sqlc の生成行を model.ProjectMilestone へ変換する。
func toModelProjectMilestone(row sqlcgen.ProjectMilestone) model.ProjectMilestone {
	return model.ProjectMilestone{
		ID:          uint(row.ID),
		ProjectID:   uint(row.ProjectID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		Status:      model.MilestoneStatus(fromStringPtr(row.Status)),
		DueDate:     fromTimestamptz(row.DueDate),
		CompletedAt: fromTimestamptz(row.CompletedAt),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

// Create はマイルストーンを作成する。
func (r *projectMilestoneRepository) Create(ctx context.Context, milestone *model.ProjectMilestone) error {
	status := string(milestone.Status)
	row, err := r.q.CreateProjectMilestone(ctx, sqlcgen.CreateProjectMilestoneParams{
		ProjectID:   int64(milestone.ProjectID),
		Title:       milestone.Title,
		Description: &milestone.Description,
		Status:      &status,
		DueDate:     toTimestamptz(milestone.DueDate),
	})
	if err != nil {
		return err
	}
	*milestone = toModelProjectMilestone(row)
	return nil
}

// FindByID は指定 ID のマイルストーンを取得する。不在の場合は (nil, nil) を返す。
func (r *projectMilestoneRepository) FindByID(ctx context.Context, id uint) (*model.ProjectMilestone, error) {
	row, err := r.q.GetProjectMilestoneByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	milestone := toModelProjectMilestone(row)
	return &milestone, nil
}

// FindByProjectID は指定プロジェクトのマイルストーン一覧を期日順で取得する。
func (r *projectMilestoneRepository) FindByProjectID(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error) {
	rows, err := r.q.ListProjectMilestonesByProject(ctx, int64(projectID))
	if err != nil {
		return nil, err
	}
	milestones := make([]model.ProjectMilestone, len(rows))
	for i, row := range rows {
		milestones[i] = toModelProjectMilestone(row)
	}
	return milestones, nil
}

// Update はマイルストーンを更新する。
func (r *projectMilestoneRepository) Update(ctx context.Context, milestone *model.ProjectMilestone) error {
	status := string(milestone.Status)
	row, err := r.q.UpdateProjectMilestone(ctx, sqlcgen.UpdateProjectMilestoneParams{
		ID:          int64(milestone.ID),
		Title:       milestone.Title,
		Description: &milestone.Description,
		Status:      &status,
		DueDate:     toTimestamptz(milestone.DueDate),
		CompletedAt: toTimestamptz(milestone.CompletedAt),
	})
	if err != nil {
		return err
	}
	*milestone = toModelProjectMilestone(row)
	return nil
}

// Delete はマイルストーンを削除する。
func (r *projectMilestoneRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteProjectMilestone(ctx, int64(id))
}
