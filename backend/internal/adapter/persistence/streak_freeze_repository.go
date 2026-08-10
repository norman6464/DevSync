package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
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

// Create は新しいストリークフリーズをデータベースに作成する。
func (r *streakFreezeRepository) Create(ctx context.Context, freeze *model.StreakFreeze) error {
	return r.db.WithContext(ctx).Create(freeze).Error
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
