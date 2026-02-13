package dto_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.RegisterRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantErr: false,
		},
		{
			name: "名前が空",
			request: dto.RegisterRequest{
				Name:            "",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantErr: true,
		},
		{
			name: "無効なメールアドレス",
			request: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "invalid-email",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantErr: true,
		},
		{
			name: "パスワードが短すぎる",
			request: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "short",
				ConfirmPassword: "short",
			},
			wantErr: true,
		},
		{
			name: "確認パスワードが空",
			request: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.LoginRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "無効なメールアドレス",
			request: dto.LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "パスワードが空",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPasswordResetRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.PasswordResetRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.PasswordResetRequest{
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "無効なメールアドレス",
			request: dto.PasswordResetRequest{
				Email: "invalid-email",
			},
			wantErr: true,
		},
		{
			name: "メールアドレスが空",
			request: dto.PasswordResetRequest{
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResetPasswordRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.ResetPasswordRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.ResetPasswordRequest{
				Token:       "valid-token",
				NewPassword: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "トークンが空",
			request: dto.ResetPasswordRequest{
				Token:       "",
				NewPassword: "newpassword123",
			},
			wantErr: true,
		},
		{
			name: "新しいパスワードが短すぎる",
			request: dto.ResetPasswordRequest{
				Token:       "valid-token",
				NewPassword: "short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteAccountRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request dto.DeleteAccountRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト",
			request: dto.DeleteAccountRequest{
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "パスワードが空",
			request: dto.DeleteAccountRequest{
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
