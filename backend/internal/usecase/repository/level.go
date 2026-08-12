package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// XPStatsReader は XP 計算に必要な統計を読むための最小の契約。
type XPStatsReader interface {
	// GetXPStats は XP 計算に使う各種集計値をまとめて返す。
	GetXPStats(ctx context.Context, userID uint) (*model.XPStats, error)
}
