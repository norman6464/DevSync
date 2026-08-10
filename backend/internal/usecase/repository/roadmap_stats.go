package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// RoadmapStatsRepository はユーザーロードマップ統計の取得に対する、usecase 側が要求する契約。
type RoadmapStatsRepository interface {
	GetRoadmapStats(ctx context.Context, userID uint) (*model.RoadmapStats, error)
}
