package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// UserSkillsReader は ID からユーザーを 1 件読むだけの最小の契約。
// おすすめユーザーの算出のようにプロフィールを参照するだけの usecase はこちらに依存する。
type UserSkillsReader interface {
	// FindByID は指定 ID のユーザーを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.User, error)
}

// UserRepository はユーザー情報の永続化に対する、usecase 側が要求する契約。
// 認証や外部サービス連携が使う操作（作成・削除・パスワード更新など）は含まない。
type UserRepository interface {
	UserSkillsReader

	FindAll(ctx context.Context) ([]model.User, error)
	// FindByUsername は指定ユーザー名のユーザーを返す。不在の場合は (nil, nil) を返す。
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	// Search は名前またはメールアドレスへの部分一致で検索する（大文字小文字を区別しない）。
	Search(ctx context.Context, query string) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
}
