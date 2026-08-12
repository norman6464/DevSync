package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ZennArticleFetcher は Zenn の記事を外部から取得するための最小の契約。
// HTTP など取得手段の詳細は adapter 側に閉じ込める。
type ZennArticleFetcher interface {
	// FetchArticles は指定ユーザーの記事をすべて取得する（ページングは実装側で解決する）。
	FetchArticles(ctx context.Context, username string) ([]model.ZennArticle, error)
	// UserExists は指定ユーザーの記事一覧を取得できるかどうかを返す。
	UserExists(ctx context.Context, username string) (bool, error)
}

// ZennArticleRepository は取得済み Zenn 記事の永続化に対する、usecase 側が要求する契約。
type ZennArticleRepository interface {
	// UpsertArticles は記事を Zenn 側の ID で重複判定して保存する。
	UpsertArticles(ctx context.Context, userID uint, articles []model.ZennArticle) error
	// GetArticles は指定ユーザーの記事を公開日の降順で返す。
	GetArticles(ctx context.Context, userID uint) ([]model.ZennArticle, error)
	// GetStats は指定ユーザーの記事統計を返す。
	GetStats(ctx context.Context, userID uint) (*model.ZennStats, error)
	// DeleteUserArticles は指定ユーザーの記事をすべて削除する。
	DeleteUserArticles(ctx context.Context, userID uint) error
}
