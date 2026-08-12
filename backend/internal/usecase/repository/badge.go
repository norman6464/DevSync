package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// BadgeStatsReader はバッジ判定に必要な統計を読むための最小の契約。
type BadgeStatsReader interface {
	// GetBadgeStats はバッジ判定に使う各種集計値をまとめて返す。
	GetBadgeStats(ctx context.Context, userID uint) (*model.BadgeStats, error)
}
