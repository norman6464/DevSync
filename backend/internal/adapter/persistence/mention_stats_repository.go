package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// mentionStatsRepository は [repository.MentionStatsRepository] の sqlc(pgx) 実装。
type mentionStatsRepository struct {
	q *sqlcgen.Queries
}

// NewMentionStatsRepository は MentionStatsRepository の sqlc(pgx) 実装を返す。
func NewMentionStatsRepository(q *sqlcgen.Queries) repository.MentionStatsRepository {
	return &mentionStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MentionStatsRepository = (*mentionStatsRepository)(nil)

// GetMentionStats は指定ユーザーのメンション集計統計を返す。
func (r *mentionStatsRepository) GetMentionStats(ctx context.Context, userID uint) (*model.MentionStats, error) {
	received, err := r.q.CountMentionsReceivedByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	made, err := r.q.CountMentionsMadeByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountMentionsReceivedByUserSince(ctx, sqlcgen.CountMentionsReceivedByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.MentionStats{
		MentionsReceived:  received,
		MentionsMade:      made,
		MentionsThisMonth: thisMonth,
	}, nil
}
