package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// roadmapStatsRepository は [repository.RoadmapStatsRepository] の sqlc(pgx) 実装。
type roadmapStatsRepository struct {
	q *sqlcgen.Queries
}

// NewRoadmapStatsRepository は RoadmapStatsRepository の sqlc(pgx) 実装を返す。
func NewRoadmapStatsRepository(q *sqlcgen.Queries) repository.RoadmapStatsRepository {
	return &roadmapStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RoadmapStatsRepository = (*roadmapStatsRepository)(nil)

// GetRoadmapStats は指定ユーザーのロードマップ統計を返す。
func (r *roadmapStatsRepository) GetRoadmapStats(ctx context.Context, userID uint) (*model.RoadmapStats, error) {
	total, err := r.q.CountRoadmapsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	activeStatus := string(model.RoadmapStatusActive)
	active, err := r.q.CountRoadmapsByUserAndStatus(ctx, sqlcgen.CountRoadmapsByUserAndStatusParams{
		UserID: int64(userID),
		Status: &activeStatus,
	})
	if err != nil {
		return nil, err
	}

	completedStatus := string(model.RoadmapStatusCompleted)
	completed, err := r.q.CountRoadmapsByUserAndStatus(ctx, sqlcgen.CountRoadmapsByUserAndStatusParams{
		UserID: int64(userID),
		Status: &completedStatus,
	})
	if err != nil {
		return nil, err
	}

	totalSteps, err := r.q.SumRoadmapStepCountByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	completedSteps, err := r.q.SumRoadmapCompletedStepCountByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.RoadmapStats{
		TotalRoadmaps:     int(total),
		ActiveRoadmaps:    int(active),
		CompletedRoadmaps: int(completed),
		TotalSteps:        int(totalSteps),
		CompletedSteps:    int(completedSteps),
	}, nil
}
