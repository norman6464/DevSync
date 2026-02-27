package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// WeeklyGoalService はカテゴリ別週間学習目標のビジネスロジックを提供する。
type WeeklyGoalService struct {
	repo repository.WeeklyGoalRepositoryInterface
}

// NewWeeklyGoalService は新しいWeeklyGoalServiceインスタンスを生成する。
func NewWeeklyGoalService(repo repository.WeeklyGoalRepositoryInterface) *WeeklyGoalService {
	return &WeeklyGoalService{repo: repo}
}

// SetGoal はカテゴリ別の週間学習目標を設定する（既存があれば更新）。
func (s *WeeklyGoalService) SetGoal(userID uint, category string, targetMinutes int) (*model.WeeklyGoal, error) {
	if !model.ValidCategories[model.LogCategory(category)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidCategory, nil)
	}
	if targetMinutes < 0 || targetMinutes > 10080 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "目標時間は0〜10080分（1週間）で指定してください", nil)
	}

	goal := &model.WeeklyGoal{
		UserID:        userID,
		Category:      model.LogCategory(category),
		TargetMinutes: targetMinutes,
	}
	if err := s.repo.Upsert(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// GetGoals は指定ユーザーの全カテゴリ週間目標を取得する。
func (s *WeeklyGoalService) GetGoals(userID uint) ([]model.WeeklyGoal, error) {
	return s.repo.GetByUserID(userID)
}

// GetProgress は指定ユーザーの全カテゴリ週間目標の達成状況を返す。
func (s *WeeklyGoalService) GetProgress(userID uint) ([]model.WeeklyGoalProgress, error) {
	goals, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var progress []model.WeeklyGoalProgress
	for _, g := range goals {
		actual, err := s.repo.SumDurationByUserCategoryThisWeek(userID, string(g.Category))
		if err != nil {
			return nil, err
		}
		pct := 0
		if g.TargetMinutes > 0 {
			pct = actual * 100 / g.TargetMinutes
		}
		progress = append(progress, model.WeeklyGoalProgress{
			Category:        g.Category,
			TargetMinutes:   g.TargetMinutes,
			ActualMinutes:   actual,
			ProgressPercent: pct,
		})
	}
	return progress, nil
}
