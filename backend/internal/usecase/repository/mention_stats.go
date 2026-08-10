package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// MentionStatsRepository はユーザーメンション集計統計の取得に対する、usecase 側が要求する契約。
type MentionStatsRepository interface {
	GetMentionStats(ctx context.Context, userID uint) (*model.MentionStats, error)
}
