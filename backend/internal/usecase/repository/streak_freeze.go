package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// StreakFreezeRepository はストリークフリーズの永続化に対する、usecase 側が要求する契約。
type StreakFreezeRepository interface {
	Create(ctx context.Context, freeze *model.StreakFreeze) error
	GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error)
	HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error)
}
