package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *PasswordResetRepository) FindByToken(token string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	err := r.db.Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *PasswordResetRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("id = ?", id).Update("used", true).Error
}

func (r *PasswordResetRepository) InvalidateUserTokens(userID uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("user_id = ? AND used = ?", userID, false).Update("used", true).Error
}

func (r *PasswordResetRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&model.PasswordResetToken{}).Error
}
