package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// FreezeUseOutcome はフリーズ使用の判定結果を表す。
type FreezeUseOutcome int

const (
	// FreezeUseCreated はフリーズを記録できたことを表す。
	FreezeUseCreated FreezeUseOutcome = iota
	// FreezeUseDuplicateDay は同じ日に使用済みで記録しなかったことを表す。
	FreezeUseDuplicateDay
	// FreezeUseMonthlyLimitReached は月次上限に達していて記録しなかったことを表す。
	FreezeUseMonthlyLimitReached
)

// StreakFreezeRepository はストリークフリーズの永続化に対する、usecase 側が要求する契約。
type StreakFreezeRepository interface {
	// CreateWithinLimits は当日重複・月次上限の判定とフリーズ作成を不可分に行う。
	// 同一ユーザーの同時実行でも重複行や上限超過が起きないことを実装が保証する。
	CreateWithinLimits(ctx context.Context, freeze *model.StreakFreeze, maxPerMonth int) (FreezeUseOutcome, error)
	GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error)
	HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error)
}
