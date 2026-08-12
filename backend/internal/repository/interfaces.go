// Package repository はDevSyncアプリケーションのデータアクセス層を提供する。
// 各リポジトリはGORMを使用してPostgreSQLに対するCRUD操作を実装する。
package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
)

// UserRepositoryInterface はユーザーデータ操作の契約を定義する。
type UserRepositoryInterface interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Search(query string) ([]model.User, error)
	FindByGitHubID(githubID int64) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	DeleteWithRelatedData(id uint) error
	UpdatePassword(userID uint, hashedPassword string) error
}

// PasswordResetRepositoryInterface はパスワードリセットトークンデータ操作の契約を定義する。
type PasswordResetRepositoryInterface interface {
	Create(token *model.PasswordResetToken) error
	FindByToken(token string) (*model.PasswordResetToken, error)
	MarkAsUsed(id uint) error
	InvalidateUserTokens(userID uint) error
	DeleteExpired() error
}
