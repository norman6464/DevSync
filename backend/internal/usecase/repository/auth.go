package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// AuthUserRepository は認証がユーザーに対して必要とする永続化の契約。
// 参照系の UserRepository とは別に、作成・削除・パスワード更新まで含む。
type AuthUserRepository interface {
	UserSkillsReader

	// FindByEmail はメールアドレスでユーザーを返す。パスワードハッシュも含めて返す
	// （ログイン処理での照合に使うため）。不在の場合は (nil, nil) を返す。
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// FindByUsername はユーザー名でユーザーを返す。不在の場合は (nil, nil) を返す。
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// FindByGitHubID は GitHub の ID でユーザーを返す。パスワードハッシュも含めて返す。
	// 不在の場合は (nil, nil) を返す。
	FindByGitHubID(ctx context.Context, githubID int64) (*model.User, error)
	// FindByIDWithPassword はIDでユーザーを返す。パスワードハッシュも含めて返す
	// （退会時の本人確認に使うため）。不在の場合は (nil, nil) を返す。
	FindByIDWithPassword(ctx context.Context, id uint) (*model.User, error)

	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	// UpdatePassword はパスワードハッシュだけを更新する。
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error
	// DeleteWithRelatedData はユーザーと関連データをトランザクション内で削除する。
	DeleteWithRelatedData(ctx context.Context, id uint) error
}

// PasswordResetTokenRepository はパスワードリセットトークンの永続化に対する契約。
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *model.PasswordResetToken) error
	// FindByToken はハッシュ済みトークンで検索する。不在の場合は (nil, nil) を返す。
	FindByToken(ctx context.Context, hashedToken string) (*model.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, id uint) error
	// InvalidateUserTokens は指定ユーザーの未使用トークンを無効化する。
	InvalidateUserTokens(ctx context.Context, userID uint) error
}
