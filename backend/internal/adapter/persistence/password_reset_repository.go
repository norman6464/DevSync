package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// passwordResetRepository は [repository.PasswordResetTokenRepository] の GORM 実装。
type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository は PasswordResetTokenRepository の GORM 実装を返す。
func NewPasswordResetTokenRepository(db *gorm.DB) repository.PasswordResetTokenRepository {
	return &passwordResetRepository{db: db}
}

var _ repository.PasswordResetTokenRepository = (*passwordResetRepository)(nil)

// Create はリセットトークンを保存する。
func (r *passwordResetRepository) Create(ctx context.Context, token *model.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByToken はハッシュ済みトークンで検索する。不在の場合は (nil, nil) を返す。
func (r *passwordResetRepository) FindByToken(ctx context.Context, hashedToken string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	if err := r.db.WithContext(ctx).Where("token = ?", hashedToken).First(&resetToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &resetToken, nil
}

// MarkAsUsed はトークンを使用済みにする。
func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.PasswordResetToken{}).
		Where("id = ?", id).Update("used", true).Error
}

// InvalidateUserTokens は指定ユーザーの未使用トークンをすべて無効化する。
func (r *passwordResetRepository) InvalidateUserTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.PasswordResetToken{}).
		Where("user_id = ? AND used = ?", userID, false).Update("used", true).Error
}
