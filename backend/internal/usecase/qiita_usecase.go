package usecase

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ConnectQiitaUseCase は Qiita のユーザー名を設定し、記事を取り込む。
type ConnectQiitaUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.QiitaArticleRepository
	fetcher  repository.QiitaArticleFetcher
}

// NewConnectQiitaUseCase は ConnectQiitaUseCase を生成する。
func NewConnectQiitaUseCase(
	users repository.ExternalAccountLinker,
	articles repository.QiitaArticleRepository,
	fetcher repository.QiitaArticleFetcher,
) *ConnectQiitaUseCase {
	return &ConnectQiitaUseCase{users: users, articles: articles, fetcher: fetcher}
}

// Execute はユーザー名を検証して保存し、記事を取り込んだ件数を返す。
func (uc *ConnectQiitaUseCase) Execute(ctx context.Context, userID uint, username string) (int, error) {
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

	user.QiitaUsername = username
	if err := uc.users.Update(ctx, user); err != nil {
		return 0, err
	}

	return syncQiitaArticles(ctx, uc.fetcher, uc.articles, userID, username)
}

// DisconnectQiitaUseCase は Qiita 連携を解除し、取り込んだ記事を削除する。
type DisconnectQiitaUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.QiitaArticleRepository
}

// NewDisconnectQiitaUseCase は DisconnectQiitaUseCase を生成する。
func NewDisconnectQiitaUseCase(
	users repository.ExternalAccountLinker,
	articles repository.QiitaArticleRepository,
) *DisconnectQiitaUseCase {
	return &DisconnectQiitaUseCase{users: users, articles: articles}
}

// Execute は Qiita のユーザー名を空にし、取り込み済みの記事を削除する。
func (uc *DisconnectQiitaUseCase) Execute(ctx context.Context, userID uint) error {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return domain.ErrNotFound
	}

	user.QiitaUsername = ""
	if err := uc.users.Update(ctx, user); err != nil {
		return err
	}
	return uc.articles.DeleteUserArticles(ctx, userID)
}

// SyncQiitaUseCase は連携済みの Qiita 記事を取り込み直す。
type SyncQiitaUseCase struct {
	users    repository.ExternalAccountLinker
	articles repository.QiitaArticleRepository
	fetcher  repository.QiitaArticleFetcher
}

// NewSyncQiitaUseCase は SyncQiitaUseCase を生成する。
func NewSyncQiitaUseCase(
	users repository.ExternalAccountLinker,
	articles repository.QiitaArticleRepository,
	fetcher repository.QiitaArticleFetcher,
) *SyncQiitaUseCase {
	return &SyncQiitaUseCase{users: users, articles: articles, fetcher: fetcher}
}

// Execute は連携済みのユーザー名で記事を取り込み直し、件数を返す。
func (uc *SyncQiitaUseCase) Execute(ctx context.Context, userID uint) (int, error) {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return 0, domain.ErrNotFound
	}
	if user.QiitaUsername == "" {
		return 0, domain.ErrBadRequest
	}
	// 保存済みのユーザー名も取得前に検証する（移行前は取得処理の中で検証していた）。
	if err := domain.ValidateExternalUsername(user.QiitaUsername); err != nil {
		return 0, err
	}

	return syncQiitaArticles(ctx, uc.fetcher, uc.articles, userID, user.QiitaUsername)
}

// syncQiitaArticles は記事を取得して取り込み、取り込んだ件数を返す。
// 取り込み時刻は 1 回の同期で揃える。
func syncQiitaArticles(
	ctx context.Context,
	fetcher repository.QiitaArticleFetcher,
	articles repository.QiitaArticleRepository,
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

// ListQiitaArticlesUseCase は取り込み済みの Qiita 記事を一覧する。
type ListQiitaArticlesUseCase struct {
	articles repository.QiitaArticleRepository
}

// NewListQiitaArticlesUseCase は ListQiitaArticlesUseCase を生成する。
func NewListQiitaArticlesUseCase(articles repository.QiitaArticleRepository) *ListQiitaArticlesUseCase {
	return &ListQiitaArticlesUseCase{articles: articles}
}

// Execute は指定ユーザーの記事を公開日の降順で返す。
func (uc *ListQiitaArticlesUseCase) Execute(ctx context.Context, userID uint) ([]model.QiitaArticle, error) {
	return uc.articles.GetArticles(ctx, userID)
}

// GetQiitaStatsUseCase は Qiita 記事の統計を取得する。
type GetQiitaStatsUseCase struct {
	articles repository.QiitaArticleRepository
}

// NewGetQiitaStatsUseCase は GetQiitaStatsUseCase を生成する。
func NewGetQiitaStatsUseCase(articles repository.QiitaArticleRepository) *GetQiitaStatsUseCase {
	return &GetQiitaStatsUseCase{articles: articles}
}

// Execute は指定ユーザーの記事統計を返す。
func (uc *GetQiitaStatsUseCase) Execute(ctx context.Context, userID uint) (*model.QiitaStats, error) {
	return uc.articles.GetStats(ctx, userID)
}
