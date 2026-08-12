package usecase

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ConnectZennUseCase は Zenn のユーザー名を設定し、記事を取り込む。
type ConnectZennUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.ZennArticleRepository
	fetcher  repository.ZennArticleFetcher
}

// NewConnectZennUseCase は ConnectZennUseCase を生成する。
func NewConnectZennUseCase(
	users repository.ExternalAccountLinker,
	articles repository.ZennArticleRepository,
	fetcher repository.ZennArticleFetcher,
) *ConnectZennUseCase {
	return &ConnectZennUseCase{users: users, articles: articles, fetcher: fetcher}
}

// Execute はユーザー名を検証して保存し、記事を取り込んだ件数を返す。
func (uc *ConnectZennUseCase) Execute(ctx context.Context, userID uint, username string) (int, error) {
	if err := domain.ValidateExternalUsername(username); err != nil {
		return 0, domain.ErrBadRequest
	}
	exists, err := uc.fetcher.UserExists(ctx, username)
	if err != nil || !exists {
		return 0, domain.ErrBadRequest
	}

	user, err := uc.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return 0, domain.ErrNotFound
	}

	user.ZennUsername = username
	if err := uc.users.Update(ctx, user); err != nil {
		return 0, err
	}

	return syncZennArticles(ctx, uc.fetcher, uc.articles, userID, username)
}

// DisconnectZennUseCase は Zenn 連携を解除し、取り込んだ記事を削除する。
type DisconnectZennUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.ZennArticleRepository
}

// NewDisconnectZennUseCase は DisconnectZennUseCase を生成する。
func NewDisconnectZennUseCase(
	users repository.ExternalAccountLinker,
	articles repository.ZennArticleRepository,
) *DisconnectZennUseCase {
	return &DisconnectZennUseCase{users: users, articles: articles}
}

// Execute は Zenn のユーザー名を空にし、取り込み済みの記事を削除する。
func (uc *DisconnectZennUseCase) Execute(ctx context.Context, userID uint) error {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return domain.ErrNotFound
	}

	user.ZennUsername = ""
	if err := uc.users.Update(ctx, user); err != nil {
		return err
	}
	return uc.articles.DeleteUserArticles(ctx, userID)
}

// SyncZennUseCase は連携済みの Zenn 記事を取り込み直す。
type SyncZennUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.ZennArticleRepository
	fetcher  repository.ZennArticleFetcher
}

// NewSyncZennUseCase は SyncZennUseCase を生成する。
func NewSyncZennUseCase(
	users repository.ExternalAccountLinker,
	articles repository.ZennArticleRepository,
	fetcher repository.ZennArticleFetcher,
) *SyncZennUseCase {
	return &SyncZennUseCase{users: users, articles: articles, fetcher: fetcher}
}

// Execute は連携済みのユーザー名で記事を取り込み直し、件数を返す。
func (uc *SyncZennUseCase) Execute(ctx context.Context, userID uint) (int, error) {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return 0, domain.ErrNotFound
	}
	if user.ZennUsername == "" {
		return 0, domain.ErrBadRequest
	}
	// 保存済みのユーザー名も取得前に検証する（移行前は取得処理の中で検証していた）。
	if err := domain.ValidateExternalUsername(user.ZennUsername); err != nil {
		return 0, err
	}

	return syncZennArticles(ctx, uc.fetcher, uc.articles, userID, user.ZennUsername)
}

// syncZennArticles は記事を取得して取り込み、取り込んだ件数を返す。
// 取り込み時刻は 1 回の同期で揃える。
func syncZennArticles(
	ctx context.Context,
	fetcher repository.ZennArticleFetcher,
	articles repository.ZennArticleRepository,
	userID uint,
	username string,
) (int, error) {
	fetched, err := fetcher.FetchArticles(ctx, username)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	for i := range fetched {
		fetched[i].UpdatedAt = now
	}

	if err := articles.UpsertArticles(ctx, userID, fetched); err != nil {
		return 0, err
	}
	return len(fetched), nil
}

// ListZennArticlesUseCase は取り込み済みの Zenn 記事を一覧する。
type ListZennArticlesUseCase struct {
	articles repository.ZennArticleRepository
}

// NewListZennArticlesUseCase は ListZennArticlesUseCase を生成する。
func NewListZennArticlesUseCase(articles repository.ZennArticleRepository) *ListZennArticlesUseCase {
	return &ListZennArticlesUseCase{articles: articles}
}

// Execute は指定ユーザーの記事を公開日の降順で返す。
func (uc *ListZennArticlesUseCase) Execute(ctx context.Context, userID uint) ([]model.ZennArticle, error) {
	return uc.articles.GetArticles(ctx, userID)
}

// GetZennStatsUseCase は Zenn 記事の統計を取得する。
type GetZennStatsUseCase struct {
	articles repository.ZennArticleRepository
}

// NewGetZennStatsUseCase は GetZennStatsUseCase を生成する。
func NewGetZennStatsUseCase(articles repository.ZennArticleRepository) *GetZennStatsUseCase {
	return &GetZennStatsUseCase{articles: articles}
}

// Execute は指定ユーザーの記事統計を返す。
func (uc *GetZennStatsUseCase) Execute(ctx context.Context, userID uint) (*model.ZennStats, error) {
	return uc.articles.GetStats(ctx, userID)
}
