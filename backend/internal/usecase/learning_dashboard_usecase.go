package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

const (
	// dashboardWeeklyDays は週間学習時間の集計日数。
	dashboardWeeklyDays = 7
	// dashboardTodayDays は当日の学習時間の集計日数。
	dashboardTodayDays = 1
)

// GetLearningDashboardSummaryUseCase は学習ダッシュボードの統合サマリーを取得する。
type GetLearningDashboardSummaryUseCase struct {
	logs      repository.LearningLogSummaryReader
	goals     repository.ActiveLearningGoalReader
	analytics repository.ProductivityStatsReader
}

// NewGetLearningDashboardSummaryUseCase は GetLearningDashboardSummaryUseCase を生成する。
func NewGetLearningDashboardSummaryUseCase(
	logs repository.LearningLogSummaryReader,
	goals repository.ActiveLearningGoalReader,
	analytics repository.ProductivityStatsReader,
) *GetLearningDashboardSummaryUseCase {
	return &GetLearningDashboardSummaryUseCase{logs: logs, goals: goals, analytics: analytics}
}

// Execute はストリーク・週間/当日の学習時間・進行中の目標数・生産性スコアをまとめて返す。
// 途中で取得に失敗した場合はその時点で中断する。
func (uc *GetLearningDashboardSummaryUseCase) Execute(ctx context.Context, userID uint) (*model.LearningDashboardSummary, error) {
	streakInfo, err := uc.logs.GetStreakInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	weeklyMinutes, err := uc.logs.SumDurationByPeriod(ctx, userID, dashboardWeeklyDays)
	if err != nil {
		return nil, err
	}

	todayMinutes, err := uc.logs.SumDurationByPeriod(ctx, userID, dashboardTodayDays)
	if err != nil {
		return nil, err
	}

	activeGoals, err := uc.goals.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats, err := uc.analytics.GetProductivityStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.LearningDashboardSummary{
		StreakInfo:        streakInfo,
		WeeklyMinutes:     weeklyMinutes,
		ActiveGoalCount:   len(activeGoals),
		TodayMinutes:      todayMinutes,
		ProductivityScore: CalculateProductivityScore(stats),
	}, nil
}
