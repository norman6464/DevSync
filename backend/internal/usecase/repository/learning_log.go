package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogRepository は学習ログの永続化に対する、usecase 側が要求する契約。
type LearningLogRepository interface {
	Create(ctx context.Context, log *model.LearningLog) error
	CreateBatch(ctx context.Context, logs []model.LearningLog) error
	Update(ctx context.Context, log *model.LearningLog) error
	// Delete は所有者本人のログだけを削除する。
	Delete(ctx context.Context, id, userID uint) error
	// FindByID は指定 ID のログを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.LearningLog, error)

	// GetByUserID はユーザーのログを作成日の新しい順でページ取得し、総数も返す。
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningLog, error)
	GetBySource(ctx context.Context, userID uint, source string) ([]model.LearningLog, error)
	// GetByPeriod は直近 days 日のログを返す。days が 0 以下なら全期間を返す。
	GetByPeriod(ctx context.Context, userID uint, days int) ([]model.LearningLog, error)
	GetFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// SumDurationByPeriod は直近 days 日の学習時間合計（分）を返す。
	SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error)
	// SumDurationByGoalID は指定ゴールに紐付いたログの学習時間合計（分）を返す。
	SumDurationByGoalID(ctx context.Context, goalID uint) (int, error)
	GetByGoalID(ctx context.Context, goalID uint, limit, offset int) ([]model.LearningLog, int64, error)

	GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error)
	GetCalendarData(ctx context.Context, userID uint) ([]model.CalendarEntry, error)
	// GetRecentCategories は使用回数の多い順にカテゴリを limit 件返す。
	GetRecentCategories(ctx context.Context, userID uint, limit int) ([]string, error)
	GetMonthlySummary(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error)
}

// LearningGoalLinker は学習ログのゴール紐付けと進捗の自動更新に必要な最小の契約。
// 目標スライスの fat port には依存しない。
type LearningGoalLinker interface {
	// FindByID は指定 ID の目標を返す。不在の場合は (nil, nil) を返す。
	FindByID(ctx context.Context, id uint) (*model.LearningGoal, error)
	Update(ctx context.Context, goal *model.LearningGoal) error
}
