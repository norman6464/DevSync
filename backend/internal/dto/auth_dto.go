// Package dto provides Data Transfer Objects for API requests and responses.
package dto

// RegisterRequest は新規ユーザー登録リクエスト
type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// PasswordResetRequest はパスワードリセットリクエスト
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest はトークンによるパスワードリセットリクエスト
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// DeleteAccountRequest はアカウント削除リクエスト
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required"`
}
