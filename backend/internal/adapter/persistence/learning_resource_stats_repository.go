package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningResourceStatsRepository は [repository.LearningResourceStatsRepository] の sqlc(pgx) 実装。
type learningResourceStatsRepository struct {
	q *sqlcgen.Queries
}

// NewLearningResourceStatsRepository は LearningResourceStatsRepository の sqlc(pgx) 実装を返す。
func NewLearningResourceStatsRepository(q *sqlcgen.Queries) repository.LearningResourceStatsRepository {
	return &learningResourceStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningResourceStatsRepository = (*learningResourceStatsRepository)(nil)

// GetLearningResourceStats は指定ユーザーの学習リソース活動集計統計を返す。
func (r *learningResourceStatsRepository) GetLearningResourceStats(ctx context.Context, userID uint) (*model.LearningResourceStats, error) {
	total, err := r.q.CountLearningResourcesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	likes, err := r.q.SumLearningResourceLikeCountByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	saves, err := r.q.SumLearningResourceSaveCountByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	categories, err := r.q.CountLearningResourceCategoriesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.LearningResourceStats{
		TotalResources: total,
		TotalLikes:     likes,
		TotalSaves:     saves,
		CategoryCount:  categories,
	}, nil
}
