package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// 学習ダッシュボードは学習ログ・学習目標・学習分析の 3 スライスから値を集めるが、
// それぞれの fat port には依存せず、必要なメソッドだけを最小 port として宣言する。

// LearningLogSummaryReader はダッシュボードが要求する学習ログ側の最小の契約。
type LearningLogSummaryReader interface {
	GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error)
	// SumDurationByPeriod は直近 days 日の学習時間合計（分）を返す。
	SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error)
}

// ActiveLearningGoalReader はダッシュボードが要求する学習目標側の最小の契約。
type ActiveLearningGoalReader interface {
	GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error)
}

// ProductivityStatsReader はダッシュボードが要求する学習分析側の最小の契約。
type ProductivityStatsReader interface {
	GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error)
}
