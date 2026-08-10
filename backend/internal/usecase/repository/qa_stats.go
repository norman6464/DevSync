package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// QAStatsRepository はユーザー Q&A 活動集計統計の取得に対する、usecase 側が要求する契約。
type QAStatsRepository interface {
	GetQAStats(ctx context.Context, userID uint) (*model.QAStats, error)
}
