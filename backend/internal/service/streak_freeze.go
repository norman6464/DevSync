package service

import (
	"fmt"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StreakFreezeService はストリークフリーズのビジネスロジックを提供する。
type StreakFreezeService struct {
	repo repository.StreakFreezeRepositoryInterface
}

// NewStreakFreezeService は新しいStreakFreezeServiceインスタンスを生成する。
func NewStreakFreezeService(repo repository.StreakFreezeRepositoryInterface) *StreakFreezeService {
	return &StreakFreezeService{repo: repo}
}

// UseFreeze は今日のストリークフリーズを使用する。
// 月2回の上限チェックと、当日の重複チェックを行う。
func (s *StreakFreezeService) UseFreeze(userID uint) error {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 今日既に使用済みか
	used, err := s.repo.HasFreezeOnDate(userID, today)
	if err != nil {
		return err
	}
	if used {
		return domain.NewError(domain.ErrCodeConflict, "今日は既にフリーズを使用済みです", nil)
	}

	// 今月の使用回数チェック
	freezes, err := s.repo.GetByUserIDAndMonth(userID, now.Year(), int(now.Month()))
	if err != nil {
		return err
	}
	if len(freezes) >= model.MaxFreezesPerMonth {
		return domain.NewError(domain.ErrCodeBadRequest,
			fmt.Sprintf("今月のフリーズ回数上限（%d回）に達しています", model.MaxFreezesPerMonth), nil)
	}

	freeze := &model.StreakFreeze{
		UserID:   userID,
		UsedDate: today,
		Month:    int(now.Month()),
		Year:     now.Year(),
	}

	return s.repo.Create(freeze)
}

// GetFreezeStatus は今月のフリーズ使用状況を返す。
func (s *StreakFreezeService) GetFreezeStatus(userID uint) (*model.StreakFreezeStatus, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	freezes, err := s.repo.GetByUserIDAndMonth(userID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, err
	}

	todayUsed, err := s.repo.HasFreezeOnDate(userID, today)
	if err != nil {
		return nil, err
	}

	usedDates := make([]string, 0, len(freezes))
	for _, f := range freezes {
		usedDates = append(usedDates, f.UsedDate)
	}

	remaining := model.MaxFreezesPerMonth - len(freezes)
	if remaining < 0 {
		remaining = 0
	}

	return &model.StreakFreezeStatus{
		MaxFreezes:  model.MaxFreezesPerMonth,
		UsedFreezes: len(freezes),
		Remaining:   remaining,
		UsedDates:   usedDates,
		TodayUsed:   todayUsed,
		CanUseToday: !todayUsed && remaining > 0,
	}, nil
}

// GetFreezeDates はユーザーの全フリーズ使用日を返す（ストリーク計算用）。
func (s *StreakFreezeService) GetFreezeDates(userID uint) ([]string, error) {
	return s.repo.GetFreezeDates(userID)
}
