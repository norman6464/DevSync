package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFollowService_Follow(t *testing.T) {
	t.Run("正常にフォローできる", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("Follow", uint(1), uint(2)).Return(nil)

		err := svc.Follow(1, 2)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("自分自身をフォローするとエラーを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		err := svc.Follow(1, 1)
		assert.ErrorIs(t, err, ErrBadRequest)
		repo.AssertNotCalled(t, "Follow")
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("Follow", uint(1), uint(2)).Return(errors.New("db error"))

		err := svc.Follow(1, 2)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestFollowService_Unfollow(t *testing.T) {
	t.Run("正常にフォロー解除できる", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("Unfollow", uint(1), uint(2)).Return(nil)

		err := svc.Unfollow(1, 2)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("Unfollow", uint(1), uint(2)).Return(errors.New("db error"))

		err := svc.Unfollow(1, 2)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestFollowService_IsFollowing(t *testing.T) {
	t.Run("フォロー中の場合trueを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("IsFollowing", uint(1), uint(2)).Return(true)

		result := svc.IsFollowing(1, 2)
		assert.True(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("フォローしていない場合falseを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("IsFollowing", uint(1), uint(2)).Return(false)

		result := svc.IsFollowing(1, 2)
		assert.False(t, result)
		repo.AssertExpectations(t)
	})
}

func TestFollowService_GetFollowers(t *testing.T) {
	t.Run("正常にフォロワー一覧を取得", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		expected := []model.User{
			{Name: "alice"},
			{Name: "bob"},
		}
		repo.On("GetFollowers", uint(1)).Return(expected, nil)

		result, err := svc.GetFollowers(1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("GetFollowers", uint(1)).Return([]model.User(nil), errors.New("db error"))

		result, err := svc.GetFollowers(1)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestFollowService_GetFollowing(t *testing.T) {
	t.Run("正常にフォロー中ユーザー一覧を取得", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		expected := []model.User{
			{Name: "charlie"},
		}
		repo.On("GetFollowing", uint(1)).Return(expected, nil)

		result, err := svc.GetFollowing(1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("リポジトリエラー時にエラーを返す", func(t *testing.T) {
		repo := new(MockFollowRepository)
		svc := NewFollowService(repo)

		repo.On("GetFollowing", uint(1)).Return([]model.User(nil), errors.New("db error"))

		result, err := svc.GetFollowing(1)
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}
