// Package dto provides Data Transfer Objects for API requests and responses.
package dto

import "github.com/norman6464/devsync/backend/internal/model"

// UserResponse はユーザー情報レスポンス（認証後のレスポンスで使用）
type UserResponse struct {
	User model.User `json:"user"`
}

// RegisterRequest は新規ユーザー登録リクエスト
type RegisterRequest struct {
	Name            string `json:"name" binding:"required" validate:"required"`
	Email           string `json:"email" binding:"required,email" validate:"required,email"`
	Password        string `json:"password" binding:"required,min=6" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// PasswordResetRequest はパスワードリセットリクエスト
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

// ResetPasswordRequest はトークンによるパスワードリセットリクエスト
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
}

// DeleteAccountRequest はアカウント削除リクエスト
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required" validate:"required"`
}
