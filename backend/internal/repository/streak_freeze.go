package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// StreakFreezeRepository はストリークフリーズデータへのアクセスを提供するリポジトリ実装。
type StreakFreezeRepository struct {
	db *gorm.DB
}

// NewStreakFreezeRepository は新しいStreakFreezeRepositoryインスタンスを生成する。
func NewStreakFreezeRepository(db *gorm.DB) *StreakFreezeRepository {
	return &StreakFreezeRepository{db: db}
}

// Create は新しいストリークフリーズをデータベースに作成する。
func (r *StreakFreezeRepository) Create(freeze *model.StreakFreeze) error {
	return r.db.Create(freeze).Error
}

// GetByUserIDAndMonth は指定ユーザーの指定月のフリーズ一覧を取得する。
func (r *StreakFreezeRepository) GetByUserIDAndMonth(userID uint, year, month int) ([]model.StreakFreeze, error) {
	var freezes []model.StreakFreeze
	err := r.db.Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		Order("used_date ASC").
		Find(&freezes).Error
	return freezes, err
}

// GetFreezeDates はユーザーの全フリーズ使用日を返す。
func (r *StreakFreezeRepository) GetFreezeDates(userID uint) ([]string, error) {
	var dates []string
	err := r.db.Model(&model.StreakFreeze{}).
		Where("user_id = ?", userID).
		Pluck("used_date", &dates).Error
	return dates, err
}

// HasFreezeOnDate は指定日にフリーズが使用されているかを返す。
func (r *StreakFreezeRepository) HasFreezeOnDate(userID uint, date string) (bool, error) {
	var count int64
	err := r.db.Model(&model.StreakFreeze{}).
		Where("user_id = ? AND used_date = ?", userID, date).
		Count(&count).Error
	return count > 0, err
}
