package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// streakFreezeRepository は [repository.StreakFreezeRepository] の GORM 実装。
type streakFreezeRepository struct {
	db *gorm.DB
}

// NewStreakFreezeRepository は StreakFreezeRepository の GORM 実装を返す。
func NewStreakFreezeRepository(db *gorm.DB) repository.StreakFreezeRepository {
	return &streakFreezeRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StreakFreezeRepository = (*streakFreezeRepository)(nil)

// CreateWithinLimits は当日重複・月次上限の判定とフリーズ作成を 1 トランザクションで行う。
// ユーザー行の行ロックで同一ユーザーの同時実行を直列化し、判定と作成の間に
// 他のリクエストが差し込めないようにする（月次上限は件数ベースで一意制約では守れないため）。
func (r *streakFreezeRepository) CreateWithinLimits(ctx context.Context, freeze *model.StreakFreeze, maxPerMonth int) (repository.FreezeUseOutcome, error) {
	outcome := repository.FreezeUseCreated
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").First(&user, freeze.UserID).Error; err != nil {
			return err
		}
		var dayCount int64
		if err := tx.Model(&model.StreakFreeze{}).
			Where("user_id = ? AND used_date = ?", freeze.UserID, freeze.UsedDate).
			Count(&dayCount).Error; err != nil {
			return err
		}
		if dayCount > 0 {
			outcome = repository.FreezeUseDuplicateDay
			return nil
		}
		var monthCount int64
		if err := tx.Model(&model.StreakFreeze{}).
			Where("user_id = ? AND year = ? AND month = ?", freeze.UserID, freeze.Year, freeze.Month).
			Count(&monthCount).Error; err != nil {
			return err
		}
		if monthCount >= int64(maxPerMonth) {
			outcome = repository.FreezeUseMonthlyLimitReached
			return nil
		}
		return tx.Create(freeze).Error
	})
	return outcome, err
}

// GetByUserIDAndMonth は指定ユーザーの指定月のフリーズ一覧を取得する。
func (r *streakFreezeRepository) GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error) {
	var freezes []model.StreakFreeze
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		Order("used_date ASC").
		Find(&freezes).Error
	return freezes, err
}

// HasFreezeOnDate は指定日にフリーズが使用されているかを返す。
func (r *streakFreezeRepository) HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.StreakFreeze{}).
		Where("user_id = ? AND used_date = ?", userID, date).
		Count(&count).Error
	return count > 0, err
}
