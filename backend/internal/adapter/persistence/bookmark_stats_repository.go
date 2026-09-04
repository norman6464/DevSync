package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// bookmarkStatsRepository は [repository.BookmarkStatsRepository] の sqlc(pgx) 実装。
type bookmarkStatsRepository struct {
	q *sqlcgen.Queries
}

// NewBookmarkStatsRepository は BookmarkStatsRepository の sqlc(pgx) 実装を返す。
func NewBookmarkStatsRepository(q *sqlcgen.Queries) repository.BookmarkStatsRepository {
	return &bookmarkStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookmarkStatsRepository = (*bookmarkStatsRepository)(nil)

// GetBookmarkStats は指定ユーザーのブックマーク集計統計を返す。
func (r *bookmarkStatsRepository) GetBookmarkStats(ctx context.Context, userID uint) (*model.BookmarkStats, error) {
	made, err := r.q.CountBookmarksMadeByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	received, err := r.q.CountBookmarksReceivedByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountBookmarksMadeByUserSince(ctx, sqlcgen.CountBookmarksMadeByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.BookmarkStats{
		TotalBookmarksMade:     made,
		TotalBookmarksReceived: received,
		BookmarksThisMonth:     thisMonth,
	}, nil
}
