package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningAnalyticsRepository は学習分析の集計に対する、usecase 側が要求する契約。
// いずれも集計結果を返す読み取り専用の操作。
type LearningAnalyticsRepository interface {
	// GetHeatmapData は曜日×時間帯ごとの学習時間を返す。
	GetHeatmapData(ctx context.Context, userID uint) ([]model.HeatmapEntry, error)
	// GetCategoryBreakdown はカテゴリ別の学習時間とログ件数を返す（割合は usecase 側で計算する）。
	GetCategoryBreakdown(ctx context.Context, userID uint) ([]model.CategoryBreakdown, error)
	// GetWeeklyTrends は直近 weeks 週の週別集計を返す。
	GetWeeklyTrends(ctx context.Context, userID uint, weeks int) ([]model.WeeklyTrend, error)
	// GetProductivityStats は生産性スコアの算出に必要な統計を返す。
	GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error)
}
