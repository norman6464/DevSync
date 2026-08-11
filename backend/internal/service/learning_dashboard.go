package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningDashboardService は学習ダッシュボードの統合サマリーを提供する。
type LearningDashboardService struct {
	logRepo       repository.LearningLogRepositoryInterface
	goalRepo      repository.LearningGoalRepositoryInterface
	analyticsRepo repository.LearningAnalyticsRepositoryInterface
}

// NewLearningDashboardService は新しいLearningDashboardServiceインスタンスを生成する。
func NewLearningDashboardService(
	logRepo repository.LearningLogRepositoryInterface,
	goalRepo repository.LearningGoalRepositoryInterface,
	analyticsRepo repository.LearningAnalyticsRepositoryInterface,
) *LearningDashboardService {
	return &LearningDashboardService{
		logRepo:       logRepo,
		goalRepo:      goalRepo,
		analyticsRepo: analyticsRepo,
	}
}

// GetSummary は学習ダッシュボードの統合サマリーを返す。
func (s *LearningDashboardService) GetSummary(userID uint) (*model.LearningDashboardSummary, error) {
	streakInfo, err := s.logRepo.GetStreakInfo(userID)
	if err != nil {
		return nil, err
	}

	weeklyMinutes, err := s.logRepo.SumDurationByPeriod(userID, 7)
	if err != nil {
		return nil, err
	}

	todayMinutes, err := s.logRepo.SumDurationByPeriod(userID, 1)
	if err != nil {
		return nil, err
	}

	activeGoals, err := s.goalRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	stats, err := s.analyticsRepo.GetProductivityStats(userID)
	if err != nil {
		return nil, err
	}
	// 生産性スコアの算出は学習分析スライスの移行で usecase 側へ移った。
	// 複製せずそちらを参照する（本スライスの移行時にこの依存ごと解消する）。
	productivityScore := usecase.CalculateProductivityScore(stats)

	return &model.LearningDashboardSummary{
		StreakInfo:         streakInfo,
		WeeklyMinutes:     weeklyMinutes,
		ActiveGoalCount:   len(activeGoals),
		TodayMinutes:      todayMinutes,
		ProductivityScore: productivityScore,
	}, nil
}
