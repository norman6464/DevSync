package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UseStreakFreezeUseCase は今日のストリークフリーズを使用する。
type UseStreakFreezeUseCase struct {
	freezes repository.StreakFreezeRepository
}

// NewUseStreakFreezeUseCase は UseStreakFreezeUseCase を生成する。
func NewUseStreakFreezeUseCase(freezes repository.StreakFreezeRepository) *UseStreakFreezeUseCase {
	return &UseStreakFreezeUseCase{freezes: freezes}
}

// Execute は当日の重複と月次上限を検査したうえでフリーズを記録する。
// 判定と作成は repository 側で不可分に行われ、同時実行でも重複・上限超過は起きない。
func (uc *UseStreakFreezeUseCase) Execute(ctx context.Context, userID uint) error {
	now := time.Now()

	outcome, err := uc.freezes.CreateWithinLimits(ctx, &model.StreakFreeze{
		UserID:   userID,
		UsedDate: now.Format("2006-01-02"),
		Month:    int(now.Month()),
		Year:     now.Year(),
	}, model.MaxFreezesPerMonth)
	if err != nil {
		return err
	}
	switch outcome {
	case repository.FreezeUseDuplicateDay:
		return domain.NewError(domain.ErrCodeConflict, "今日は既にフリーズを使用済みです", nil)
	case repository.FreezeUseMonthlyLimitReached:
		return domain.NewError(domain.ErrCodeBadRequest,
			fmt.Sprintf("今月のフリーズ回数上限（%d回）に達しています", model.MaxFreezesPerMonth), nil)
	}
	return nil
}

// GetStreakFreezeStatusUseCase は今月のフリーズ使用状況を取得する。
type GetStreakFreezeStatusUseCase struct {
	freezes repository.StreakFreezeRepository
}

// NewGetStreakFreezeStatusUseCase は GetStreakFreezeStatusUseCase を生成する。
func NewGetStreakFreezeStatusUseCase(freezes repository.StreakFreezeRepository) *GetStreakFreezeStatusUseCase {
	return &GetStreakFreezeStatusUseCase{freezes: freezes}
}

// Execute は今月の使用回数・残り回数・当日の使用有無を組み立てて返す。
func (uc *GetStreakFreezeStatusUseCase) Execute(ctx context.Context, userID uint) (*model.StreakFreezeStatus, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	freezes, err := uc.freezes.GetByUserIDAndMonth(ctx, userID, now.Year(), int(now.Month()))
	if err != nil {
		return nil, err
	}

	todayUsed, err := uc.freezes.HasFreezeOnDate(ctx, userID, today)
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
