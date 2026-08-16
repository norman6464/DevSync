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
)

// mockFollowRepo は usecase/repository.FollowRepository の testify/mock 実装。
type mockFollowRepo struct{ mock.Mock }

func (m *mockFollowRepo) Follow(ctx context.Context, followerID, followeeID uint) error {
	return m.Called(ctx, followerID, followeeID).Error(0)
}

func (m *mockFollowRepo) Unfollow(ctx context.Context, followerID, followeeID uint) error {
	return m.Called(ctx, followerID, followeeID).Error(0)
}

func (m *mockFollowRepo) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	users, _ := args.Get(0).([]model.User)
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *mockFollowRepo) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	users, _ := args.Get(0).([]model.User)
	return users, args.Get(1).(int64), args.Error(2)
}

func TestFollowUserUseCase_Execute(t *testing.T) {
	t.Run("フォローすると相手に通知を作る", func(t *testing.T) {
		repo := new(mockFollowRepo)
		notifications := new(mockNotificationCreatorPort)
		repo.On("Follow", mock.Anything, uint(1), uint(2)).Return(nil)
		notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
			// 受信者はフォローされた側、実行者はフォローした側
			return n.UserID == 2 && n.ActorID == 1 && n.Type == model.NotificationTypeFollow
		})).Return(nil)
		uc := usecase.NewFollowUserUseCase(repo, notifications)

		err := uc.Execute(context.Background(), 1, 2)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		notifications.AssertExpectations(t)
	})

	t.Run("自分自身をフォローするとリポジトリを呼ばずエラーを返す", func(t *testing.T) {
		repo := new(mockFollowRepo)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewFollowUserUseCase(repo, notifications)

		err := uc.Execute(context.Background(), 1, 1)

		assert.ErrorIs(t, err, domain.ErrBadRequest)
		repo.AssertNotCalled(t, "Follow")
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("重複フォローは ErrConflict を伝播し通知しない", func(t *testing.T) {
		repo := new(mockFollowRepo)
		notifications := new(mockNotificationCreatorPort)
		repo.On("Follow", mock.Anything, uint(1), uint(2)).Return(domain.ErrConflict)
		uc := usecase.NewFollowUserUseCase(repo, notifications)

		err := uc.Execute(context.Background(), 1, 2)

		assert.ErrorIs(t, err, domain.ErrConflict)
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("フォローに失敗したら通知しない", func(t *testing.T) {
		repo := new(mockFollowRepo)
		notifications := new(mockNotificationCreatorPort)
		repo.On("Follow", mock.Anything, uint(1), uint(2)).Return(errors.New("db error"))
		uc := usecase.NewFollowUserUseCase(repo, notifications)

		err := uc.Execute(context.Background(), 1, 2)

		assert.Error(t, err)
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	// 通知はフォローの付随処理なので、失敗してもフォロー自体は成立させる。
	t.Run("通知の失敗はフォローの成否に影響しない", func(t *testing.T) {
		repo := new(mockFollowRepo)
		notifications := new(mockNotificationCreatorPort)
		repo.On("Follow", mock.Anything, uint(1), uint(2)).Return(nil)
		notifications.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewFollowUserUseCase(repo, notifications)

		assert.NoError(t, uc.Execute(context.Background(), 1, 2))
		notifications.AssertExpectations(t)
	})
}

func TestUnfollowUserUseCase_Execute(t *testing.T) {
	t.Run("正常にフォロー解除できる", func(t *testing.T) {
		repo := new(mockFollowRepo)
		repo.On("Unfollow", mock.Anything, uint(1), uint(2)).Return(nil)
		uc := usecase.NewUnfollowUserUseCase(repo)

		err := uc.Execute(context.Background(), 1, 2)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("自分自身の解除はリポジトリを呼ばずエラーを返す", func(t *testing.T) {
		repo := new(mockFollowRepo)
		uc := usecase.NewUnfollowUserUseCase(repo)

		err := uc.Execute(context.Background(), 1, 1)

		assert.ErrorIs(t, err, domain.ErrBadRequest)
		repo.AssertNotCalled(t, "Unfollow")
	})

	t.Run("リポジトリエラーを伝播する", func(t *testing.T) {
		repo := new(mockFollowRepo)
		repo.On("Unfollow", mock.Anything, uint(1), uint(2)).Return(errors.New("db error"))
		uc := usecase.NewUnfollowUserUseCase(repo)

		err := uc.Execute(context.Background(), 1, 2)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestListFollowersUseCase_Execute(t *testing.T) {
	t.Run("フォロワー一覧を取得する", func(t *testing.T) {
		repo := new(mockFollowRepo)
		expected := []model.User{{Name: "alice"}, {Name: "bob"}}
		repo.On("GetFollowers", mock.Anything, uint(1), 20, 0).Return(expected, int64(2), nil)
		uc := usecase.NewListFollowersUseCase(repo)

		users, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.NoError(t, err)
		assert.Equal(t, expected, users)
		assert.Equal(t, int64(2), total)
		repo.AssertExpectations(t)
	})

	t.Run("空リストを返す", func(t *testing.T) {
		repo := new(mockFollowRepo)
		repo.On("GetFollowers", mock.Anything, uint(1), 20, 0).Return([]model.User{}, int64(0), nil)
		uc := usecase.NewListFollowersUseCase(repo)

		users, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.NoError(t, err)
		assert.Empty(t, users)
		assert.Equal(t, int64(0), total)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラーを伝播する", func(t *testing.T) {
		repo := new(mockFollowRepo)
		repo.On("GetFollowers", mock.Anything, uint(1), 20, 0).Return([]model.User(nil), int64(0), errors.New("db error"))
		uc := usecase.NewListFollowersUseCase(repo)

		users, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, int64(0), total)
		repo.AssertExpectations(t)
	})
}

func TestListFollowingUseCase_Execute(t *testing.T) {
	t.Run("フォロー中一覧を取得する", func(t *testing.T) {
		repo := new(mockFollowRepo)
		expected := []model.User{{Name: "charlie"}}
		repo.On("GetFollowing", mock.Anything, uint(1), 20, 0).Return(expected, int64(1), nil)
		uc := usecase.NewListFollowingUseCase(repo)

		users, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.NoError(t, err)
		assert.Equal(t, expected, users)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラーを伝播する", func(t *testing.T) {
		repo := new(mockFollowRepo)
		repo.On("GetFollowing", mock.Anything, uint(1), 20, 0).Return([]model.User(nil), int64(0), errors.New("db error"))
		uc := usecase.NewListFollowingUseCase(repo)

		users, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, int64(0), total)
		repo.AssertExpectations(t)
	})
}
