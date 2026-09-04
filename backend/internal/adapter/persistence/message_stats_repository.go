package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// messageStatsRepository は [repository.MessageStatsRepository] の sqlc(pgx) 実装。
type messageStatsRepository struct {
	q *sqlcgen.Queries
}

// NewMessageStatsRepository は MessageStatsRepository の sqlc(pgx) 実装を返す。
func NewMessageStatsRepository(q *sqlcgen.Queries) repository.MessageStatsRepository {
	return &messageStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MessageStatsRepository = (*messageStatsRepository)(nil)

// GetMessageStats は指定ユーザーのメッセージ集計統計を返す。
func (r *messageStatsRepository) GetMessageStats(ctx context.Context, userID uint) (*model.MessageStats, error) {
	sent, err := r.q.CountMessagesSentByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	received, err := r.q.CountMessagesReceivedByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	conversations, err := r.q.CountConversationsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountMessagesSentByUserSince(ctx, sqlcgen.CountMessagesSentByUserSinceParams{
		SenderID:  int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.MessageStats{
		TotalSent:          sent,
		TotalReceived:      received,
		ConversationsCount: conversations,
		MessagesThisMonth:  thisMonth,
	}, nil
}
