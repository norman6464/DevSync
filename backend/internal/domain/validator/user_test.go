package validator

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserValidator_ValidateRegistration(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  bool
	}{
		{"有効な登録", "testuser", "test@example.com", "password123", false},
		{"無効（ユーザー名が短い）", "u", "test@example.com", "password123", true},
		{"無効（メールが不正）", "testuser", "invalid-email", "password123", true},
		{"無効（パスワードが短い）", "testuser", "test@example.com", "pass", true},
		{"無効（パスワードに数字/記号なし）", "testuser", "test@example.com", "password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateRegistration(tt.username, tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserValidator_ValidateLogin(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"有効なログイン", "test@example.com", "password", false},
		{"無効（メールが不正）", "invalid-email", "password", true},
		{"無効（パスワードが空）", "test@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateLogin(tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserValidator_ValidatePasswordReset(t *testing.T) {
	v := NewUserValidator()

	tests := []struct {
		name        string
		newPassword string
		wantErr     bool
	}{
		{"有効なパスワード", "newpassword123", false},
		{"無効（短すぎる）", "pass", true},
		{"無効（数字/記号なし）", "password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidatePasswordReset(tt.newPassword)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, domain.IsDomainError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
