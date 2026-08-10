package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// MessageStatsRepository はユーザーメッセージ集計統計の取得に対する、usecase 側が要求する契約。
type MessageStatsRepository interface {
	GetMessageStats(ctx context.Context, userID uint) (*model.MessageStats, error)
}
