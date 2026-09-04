package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// projectStatsRepository は [repository.ProjectStatsRepository] の sqlc(pgx) 実装。
type projectStatsRepository struct {
	q *sqlcgen.Queries
}

// NewProjectStatsRepository は ProjectStatsRepository の sqlc(pgx) 実装を返す。
func NewProjectStatsRepository(q *sqlcgen.Queries) repository.ProjectStatsRepository {
	return &projectStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectStatsRepository = (*projectStatsRepository)(nil)

// GetProjectStats は指定ユーザーのプロジェクト活動集計統計を返す。
func (r *projectStatsRepository) GetProjectStats(ctx context.Context, userID uint) (*model.ProjectStats, error) {
	row, err := r.q.GetProjectStats(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return &model.ProjectStats{
		TotalProjects:     row.TotalProjects,
		FeaturedProjects:  row.FeaturedProjects,
		OngoingProjects:   row.OngoingProjects,
		CompletedProjects: row.CompletedProjects,
	}, nil
}
