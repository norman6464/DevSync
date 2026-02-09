package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// AIAdviceRepository はAIアドバイスデータのDB操作を実装する。
type AIAdviceRepository struct {
	db *gorm.DB
}

// NewAIAdviceRepository は新しいAIAdviceRepositoryインスタンスを生成する。
func NewAIAdviceRepository(db *gorm.DB) *AIAdviceRepository {
	return &AIAdviceRepository{db: db}
}

// Create は新しいAIアドバイスをDBに作成する。
func (r *AIAdviceRepository) Create(advice *model.AIAdvice) error {
	return r.db.Create(advice).Error
}

// CreateBatch は複数のAIアドバイスを一括でDBに作成する。
func (r *AIAdviceRepository) CreateBatch(advices []*model.AIAdvice) error {
	if len(advices) == 0 {
		return nil
	}
	return r.db.Create(&advices).Error
}

// FindByUserID は指定ユーザーIDのアドバイスを優先度順で取得する。
func (r *AIAdviceRepository) FindByUserID(userID uint, limit int) ([]model.AIAdvice, error) {
	var advices []model.AIAdvice
	err := r.db.Where("user_id = ?", userID).
		Order("priority ASC, created_at DESC").
		Limit(limit).
		Find(&advices).Error
	return advices, err
}

// FindUnreadByUserID は指定ユーザーIDの未読アドバイスを取得する。
func (r *AIAdviceRepository) FindUnreadByUserID(userID uint) ([]model.AIAdvice, error) {
	var advices []model.AIAdvice
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("priority ASC, created_at DESC").
		Find(&advices).Error
	return advices, err
}

// MarkAsRead は指定アドバイスを既読にする。
func (r *AIAdviceRepository) MarkAsRead(id, userID uint) error {
	result := r.db.Model(&model.AIAdvice{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// MarkAllAsRead は指定ユーザーの全アドバイスを既読にする。
func (r *AIAdviceRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.AIAdvice{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// DeleteExpired は期限切れのアドバイスを削除する。
func (r *AIAdviceRepository) DeleteExpired() error {
	return r.db.Where("expires_at IS NOT NULL AND expires_at < NOW()").
		Delete(&model.AIAdvice{}).Error
}

// DeleteByUserID は指定ユーザーの全アドバイスを削除する。
func (r *AIAdviceRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.AIAdvice{}).Error
}
