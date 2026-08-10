package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// StudyCircleStatsRepository はスタディサークル集計統計の取得に対する、usecase 側が要求する契約。
type StudyCircleStatsRepository interface {
	GetCircleStats(ctx context.Context, circleID uint) (*model.StudyCircleStats, error)
}
