package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// studyCircleStatsRepository は [repository.StudyCircleStatsRepository] の sqlc(pgx) 実装。
type studyCircleStatsRepository struct {
	q *sqlcgen.Queries
}

// NewStudyCircleStatsRepository は StudyCircleStatsRepository の sqlc(pgx) 実装を返す。
func NewStudyCircleStatsRepository(q *sqlcgen.Queries) repository.StudyCircleStatsRepository {
	return &studyCircleStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StudyCircleStatsRepository = (*studyCircleStatsRepository)(nil)

// GetCircleStats は指定サークルの集計統計を返す。
func (r *studyCircleStatsRepository) GetCircleStats(ctx context.Context, circleID uint) (*model.StudyCircleStats, error) {
	members, err := r.q.CountStudyCircleMembersByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	checkins, err := r.q.CountStudyCircleCheckinsByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	steps, err := r.q.CountStudyCircleStepsByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	completedSteps, err := r.q.CountStudyCircleCompletedStepsByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	return &model.StudyCircleStats{
		MemberCount:    members,
		CheckinCount:   checkins,
		TotalSteps:     steps,
		CompletedSteps: completedSteps,
	}, nil
}
