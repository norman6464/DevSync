package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningLogStatsRepository は [repository.LearningLogStatsRepository] の sqlc(pgx) 実装。
type learningLogStatsRepository struct {
	q *sqlcgen.Queries
}

// NewLearningLogStatsRepository は LearningLogStatsRepository の sqlc(pgx) 実装を返す。
func NewLearningLogStatsRepository(q *sqlcgen.Queries) repository.LearningLogStatsRepository {
	return &learningLogStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogStatsRepository = (*learningLogStatsRepository)(nil)

// GetLearningLogStats は指定ユーザーの学習ログ集計統計を返す。
func (r *learningLogStatsRepository) GetLearningLogStats(ctx context.Context, userID uint) (*model.LearningLogStats, error) {
	total, err := r.q.CountLearningLogsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	duration, err := r.q.SumLearningLogDurationByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	categories, err := r.q.CountLearningLogCategoriesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountLearningLogsByUserSince(ctx, sqlcgen.CountLearningLogsByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.LearningLogStats{
		TotalLogs:     total,
		TotalDuration: duration,
		CategoryCount: categories,
		LogsThisMonth: thisMonth,
	}, nil
}
