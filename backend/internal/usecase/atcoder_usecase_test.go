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

// mockAtCoderRatingFetcher は usecase/repository.AtCoderRatingFetcher のモック。
type mockAtCoderRatingFetcher struct{ mock.Mock }

func (m *mockAtCoderRatingFetcher) FetchRatingHistory(ctx context.Context, username string) ([]model.AtCoderRatingEntry, error) {
	args := m.Called(ctx, username)
	h, _ := args.Get(0).([]model.AtCoderRatingEntry)
	return h, args.Error(1)
}

func (m *mockAtCoderRatingFetcher) UserExists(ctx context.Context, username string) bool {
	return m.Called(ctx, username).Bool(0)
}

// mockAccountLinker は usecase/repository.ExternalAccountLinker のモック。
type mockAccountLinker struct{ mock.Mock }

func (m *mockAccountLinker) FindByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *mockAccountLinker) Update(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

// ============================================================
// レーティング取得
// ============================================================

func TestGetAtCoderRatingUseCase(t *testing.T) {
	t.Run("最新のレーティングから色とランクを求める", func(t *testing.T) {
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewGetAtCoderRatingUseCase(ratings)

		ratings.On("FetchRatingHistory", mock.Anything, "testuser").Return([]model.AtCoderRatingEntry{
			{NewRating: 800}, {NewRating: 1500},
		}, nil)

		info, err := uc.Execute(context.Background(), "testuser")
		require.NoError(t, err)
		assert.Equal(t, "testuser", info.Username)
		assert.Equal(t, 1500, info.Rating)
		assert.Equal(t, "cyan", info.Color)
		assert.Equal(t, "水色", info.Rank)
	})

	t.Run("履歴が空ならレーティング 0（灰）", func(t *testing.T) {
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewGetAtCoderRatingUseCase(ratings)

		ratings.On("FetchRatingHistory", mock.Anything, "newuser").Return([]model.AtCoderRatingEntry{}, nil)

		info, err := uc.Execute(context.Background(), "newuser")
		require.NoError(t, err)
		assert.Equal(t, 0, info.Rating)
		assert.Equal(t, "gray", info.Color)
		assert.Equal(t, "灰", info.Rank)
	})

	t.Run("ユーザー名が不正なら取得しない", func(t *testing.T) {
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewGetAtCoderRatingUseCase(ratings)

		_, err := uc.Execute(context.Background(), "bad!user")
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeValidation, domain.GetDomainError(err).Code)
		ratings.AssertNotCalled(t, "FetchRatingHistory", mock.Anything, mock.Anything)
	})

	t.Run("取得エラーはそのまま返す", func(t *testing.T) {
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewGetAtCoderRatingUseCase(ratings)

		fetchErr := domain.NewError(domain.ErrCodeServiceUnavailable, "AtCoder APIリクエストに失敗", nil)
		ratings.On("FetchRatingHistory", mock.Anything, "testuser").Return(nil, fetchErr)

		_, err := uc.Execute(context.Background(), "testuser")
		assert.ErrorIs(t, err, fetchErr)
	})
}

func TestAtCoderRatingColor(t *testing.T) {
	tests := []struct {
		name   string
		rating int
		want   string
	}{
		{"赤（2800以上）", 2800, "red"},
		{"橙（2400-2799）", 2400, "orange"},
		{"黄（2000-2399）", 2000, "yellow"},
		{"青（1600-1999）", 1600, "blue"},
		{"水色（1200-1599）", 1200, "cyan"},
		{"緑（800-1199）", 800, "green"},
		{"茶（400-799）", 400, "brown"},
		{"灰（0-399）", 399, "gray"},
		{"灰（0）", 0, "gray"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, usecase.AtCoderRatingColor(tt.rating))
		})
	}
}

func TestAtCoderRatingRank(t *testing.T) {
	tests := []struct {
		name   string
		rating int
		want   string
	}{
		{"赤（2800以上）", 3000, "赤"},
		{"橙（2400-2799）", 2500, "橙"},
		{"黄（2000-2399）", 2100, "黄"},
		{"青（1600-1999）", 1700, "青"},
		{"水色（1200-1599）", 1300, "水色"},
		{"緑（800-1199）", 900, "緑"},
		{"茶（400-799）", 500, "茶"},
		{"灰（0-399）", 100, "灰"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, usecase.AtCoderRatingRank(tt.rating))
		})
	}
}

// ============================================================
// 連携・連携解除
// ============================================================

func TestConnectAtCoderUseCase(t *testing.T) {
	t.Run("存在するユーザー名を保存する", func(t *testing.T) {
		users := new(mockAccountLinker)
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewConnectAtCoderUseCase(users, ratings)

		ratings.On("UserExists", mock.Anything, "myuser").Return(true)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.AtCoderUsername == "myuser"
		})).Return(nil)

		user, err := uc.Execute(context.Background(), 1, "myuser")
		require.NoError(t, err)
		assert.Equal(t, "myuser", user.AtCoderUsername)
		users.AssertExpectations(t)
	})

	t.Run("形式が不正なら外部 API を呼ばずに 400", func(t *testing.T) {
		users := new(mockAccountLinker)
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewConnectAtCoderUseCase(users, ratings)

		_, err := uc.Execute(context.Background(), 1, "bad!user")
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeBadRequest, domain.GetDomainError(err).Code)
		ratings.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
		users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("AtCoder 上に存在しなければ 400", func(t *testing.T) {
		users := new(mockAccountLinker)
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewConnectAtCoderUseCase(users, ratings)

		ratings.On("UserExists", mock.Anything, "ghost").Return(false)

		_, err := uc.Execute(context.Background(), 1, "ghost")
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeBadRequest, domain.GetDomainError(err).Code)
		users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("ユーザーが不在なら 404", func(t *testing.T) {
		users := new(mockAccountLinker)
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewConnectAtCoderUseCase(users, ratings)

		ratings.On("UserExists", mock.Anything, "myuser").Return(true)
		users.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := uc.Execute(context.Background(), 1, "myuser")
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeNotFound, domain.GetDomainError(err).Code)
		users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("保存エラーはそのまま返す", func(t *testing.T) {
		users := new(mockAccountLinker)
		ratings := new(mockAtCoderRatingFetcher)
		uc := usecase.NewConnectAtCoderUseCase(users, ratings)

		dbErr := errors.New("db error")
		ratings.On("UserExists", mock.Anything, "myuser").Return(true)
		users.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		users.On("Update", mock.Anything, mock.Anything).Return(dbErr)

		_, err := uc.Execute(context.Background(), 1, "myuser")
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestDisconnectAtCoderUseCase(t *testing.T) {
	t.Run("ユーザー名を空にして保存する", func(t *testing.T) {
		users := new(mockAccountLinker)
		uc := usecase.NewDisconnectAtCoderUseCase(users)

		users.On("FindByID", mock.Anything, uint(1)).
			Return(&model.User{ID: 1, AtCoderUsername: "myuser"}, nil)
		users.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.AtCoderUsername == ""
		})).Return(nil)

		user, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Empty(t, user.AtCoderUsername)
		users.AssertExpectations(t)
	})

	t.Run("取得に失敗したら 404 に潰す", func(t *testing.T) {
		users := new(mockAccountLinker)
		uc := usecase.NewDisconnectAtCoderUseCase(users)

		users.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		_, err := uc.Execute(context.Background(), 1)
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeNotFound, domain.GetDomainError(err).Code)
		users.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}
