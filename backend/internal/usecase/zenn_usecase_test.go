package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockZennArticleRepo は usecase/repository.ZennArticleRepository のモック。
type mockZennArticleRepo struct{ mock.Mock }

func (m *mockZennArticleRepo) UpsertArticles(ctx context.Context, userID uint, articles []model.ZennArticle) error {
	return m.Called(ctx, userID, articles).Error(0)
}

func (m *mockZennArticleRepo) GetArticles(ctx context.Context, userID uint) ([]model.ZennArticle, error) {
	args := m.Called(ctx, userID)
	a, _ := args.Get(0).([]model.ZennArticle)
	return a, args.Error(1)
}

func (m *mockZennArticleRepo) GetStats(ctx context.Context, userID uint) (*model.ZennStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ZennStats)
	return s, args.Error(1)
}

func (m *mockZennArticleRepo) DeleteUserArticles(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockZennArticleFetcher は usecase/repository.ZennArticleFetcher のモック。
type mockZennArticleFetcher struct{ mock.Mock }

func (m *mockZennArticleFetcher) FetchArticles(ctx context.Context, username string) ([]model.ZennArticle, error) {
	args := m.Called(ctx, username)
	a, _ := args.Get(0).([]model.ZennArticle)
	return a, args.Error(1)
}

func (m *mockZennArticleFetcher) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

// ============================================================
// 連携
// ============================================================

func TestConnectZennUseCase(t *testing.T) {
	t.Run("ユーザー名を保存して記事を取り込む", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewConnectZennUseCase(users, articles, fetcher)

		fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ZennUsername == "testuser"
		})).Return(nil)
		fetcher.On("FetchArticles", mock.Anything, "testuser").
			Return([]model.ZennArticle{{ZennID: 1}, {ZennID: 2}}, nil)
		articles.On("UpsertArticles", mock.Anything, uint(1), mock.MatchedBy(func(a []model.ZennArticle) bool {
			// 取り込み時刻を揃えて入れる。
			return len(a) == 2 && !a[0].UpdatedAt.IsZero() && a[0].UpdatedAt.Equal(a[1].UpdatedAt)
		})).Return(nil)

		count, err := uc.Execute(context.Background(), 1, "testuser")
		require.NoError(t, err)
		assert.Equal(t, 2, count)
		users.AssertExpectations(t)
		articles.AssertExpectations(t)
		fetcher.AssertExpectations(t)
	})

	t.Run("形式が不正なら 400 で外部 API を呼ばない", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewConnectZennUseCase(users, articles, fetcher)

		_, err := uc.Execute(context.Background(), 1, "bad!user")
		assert.ErrorIs(t, err, domain.ErrBadRequest)
		fetcher.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
		users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("存在確認に失敗しても 400", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewConnectZennUseCase(users, articles, fetcher)

		fetcher.On("UserExists", mock.Anything, "testuser").Return(false, errors.New("api error"))

		_, err := uc.Execute(context.Background(), 1, "testuser")
		assert.ErrorIs(t, err, domain.ErrBadRequest)
		users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
		fetcher.AssertExpectations(t)
	})

	t.Run("ユーザーが不在なら 404", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewConnectZennUseCase(users, articles, fetcher)

		fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := uc.Execute(context.Background(), 1, "testuser")
		assert.ErrorIs(t, err, domain.ErrNotFound)
		users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		users.AssertExpectations(t)
	})

	t.Run("記事の取り込みに失敗したらエラーを返す", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewConnectZennUseCase(users, articles, fetcher)

		upsertErr := errors.New("db error")
		fetcher.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		users.On("Update", mock.Anything, mock.Anything).Return(nil)
		fetcher.On("FetchArticles", mock.Anything, "testuser").Return([]model.ZennArticle{{ZennID: 1}}, nil)
		articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(upsertErr)

		count, err := uc.Execute(context.Background(), 1, "testuser")
		assert.ErrorIs(t, err, upsertErr)
		assert.Zero(t, count)
	})
}

// ============================================================
// 連携解除
// ============================================================

func TestDisconnectZennUseCase(t *testing.T) {
	t.Run("ユーザー名を空にして記事を削除する", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		uc := usecase.NewDisconnectZennUseCase(users, articles)

		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "testuser"}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.ZennUsername == ""
		})).Return(nil)
		articles.On("DeleteUserArticles", mock.Anything, uint(1)).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 1))
		users.AssertExpectations(t)
		articles.AssertExpectations(t)
	})

	t.Run("ユーザーが不在なら 404 で記事を消さない", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		uc := usecase.NewDisconnectZennUseCase(users, articles)

		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1), domain.ErrNotFound)
		articles.AssertNotCalled(t, "DeleteUserArticles", mock.Anything, mock.Anything)
		users.AssertExpectations(t)
	})
}

// ============================================================
// 同期
// ============================================================

func TestSyncZennUseCase(t *testing.T) {
	t.Run("連携済みのユーザー名で取り込み直す", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewSyncZennUseCase(users, articles, fetcher)

		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "testuser"}, nil)
		fetcher.On("FetchArticles", mock.Anything, "testuser").
			Return([]model.ZennArticle{{ZennID: 1}, {ZennID: 2}, {ZennID: 3}}, nil)
		articles.On("UpsertArticles", mock.Anything, uint(1), mock.Anything).Return(nil)

		count, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, 3, count)
		users.AssertExpectations(t)
		articles.AssertExpectations(t)
		fetcher.AssertExpectations(t)
	})

	t.Run("未連携なら 400", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewSyncZennUseCase(users, articles, fetcher)

		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

		_, err := uc.Execute(context.Background(), 1)
		assert.ErrorIs(t, err, domain.ErrBadRequest)
		fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
	})

	t.Run("保存済みのユーザー名が不正なら検証エラー", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewSyncZennUseCase(users, articles, fetcher)

		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1, ZennUsername: "bad!user"}, nil)

		_, err := uc.Execute(context.Background(), 1)
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeValidation, domain.GetDomainError(err).Code)
		fetcher.AssertNotCalled(t, "FetchArticles", mock.Anything, mock.Anything)
	})

	t.Run("ユーザーが不在なら 404", func(t *testing.T) {
		users := new(mockAccountLinker)
		articles := new(mockZennArticleRepo)
		fetcher := new(mockZennArticleFetcher)
		uc := usecase.NewSyncZennUseCase(users, articles, fetcher)

		users.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, err := uc.Execute(context.Background(), 1)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

// ============================================================
// 参照
// ============================================================

func TestListZennArticlesUseCase(t *testing.T) {
	articles := new(mockZennArticleRepo)
	uc := usecase.NewListZennArticlesUseCase(articles)

	articles.On("GetArticles", mock.Anything, uint(5)).Return([]model.ZennArticle{{Title: "A"}}, nil)

	got, err := uc.Execute(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	articles.AssertExpectations(t)
}

func TestGetZennStatsUseCase(t *testing.T) {
	articles := new(mockZennArticleRepo)
	uc := usecase.NewGetZennStatsUseCase(articles)

	articles.On("GetStats", mock.Anything, uint(5)).Return(&model.ZennStats{TotalArticles: 8}, nil)

	stats, err := uc.Execute(context.Background(), 5)
	require.NoError(t, err)
	assert.Equal(t, 8, stats.TotalArticles)
	articles.AssertExpectations(t)
}
