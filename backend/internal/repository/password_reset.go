package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PasswordResetRepository はパスワードリセットトークンデータへのアクセスを提供するリポジトリ実装。
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository は新しいPasswordResetRepositoryインスタンスを生成する。
func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create は新しいパスワードリセットトークンをデータベースに作成する。
func (r *PasswordResetRepository) Create(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// FindByToken はトークン文字列でパスワードリセットトークンを検索する。
func (r *PasswordResetRepository) FindByToken(token string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	err := r.db.Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// MarkAsUsed は指定IDのトークンを使用済みにする。
func (r *PasswordResetRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("id = ?", id).Update("used", true).Error
}

// InvalidateUserTokens は指定ユーザーの未使用トークンを全て無効化する。
func (r *PasswordResetRepository) InvalidateUserTokens(userID uint) error {
	return r.db.Model(&model.PasswordResetToken{}).Where("user_id = ? AND used = ?", userID, false).Update("used", true).Error
}

// DeleteExpired は有効期限切れのトークンを全て削除する。
func (r *PasswordResetRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&model.PasswordResetToken{}).Error
}
