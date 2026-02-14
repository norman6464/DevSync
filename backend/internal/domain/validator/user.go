// Package validator provides validation logic for domain entities.
package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// UserValidator handles validation for User entities.
type UserValidator struct{}

// NewUserValidator creates a new UserValidator instance.
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateRegistration validates inputs for user registration.
func (v *UserValidator) ValidateRegistration(name, email, password string) error {
	// ユーザー名のバリデーション
	if err := domain.ValidateUsername(name); err != nil {
		return err
	}

	// メールアドレスのバリデーション
	if err := domain.ValidateEmail(email); err != nil {
		return err
	}

	// パスワードのバリデーション
	if err := domain.ValidatePassword(password); err != nil {
		return err
	}

	return nil
}

// ValidateLogin validates inputs for user login.
func (v *UserValidator) ValidateLogin(email, password string) error {
	// メールアドレスのバリデーション
	if err := domain.ValidateEmail(email); err != nil {
		return err
	}

	// パスワードは空でないことだけチェック（ログイン時は詳細チェック不要）
	if err := domain.ValidateStringLength(password, 1, 0, "パスワード"); err != nil {
		return err
	}

	return nil
}

// ValidatePasswordReset validates new password for password reset.
func (v *UserValidator) ValidatePasswordReset(newPassword string) error {
	return domain.ValidatePassword(newPassword)
}
