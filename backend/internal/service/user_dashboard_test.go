package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// TestUserDashboardGetStats_Success は正常にダッシュボード統計が取得できることを確認する。
func TestUserDashboardGetStats_Success(t *testing.T) {
	mockRepo := &MockUserDashboardRepository{}
	svc := NewUserDashboardService(mockRepo)

	expected := &model.UserDashboardStats{
		PostCount:        10,
		LikesReceived:    50,
		CommentsReceived: 20,
		ViewsReceived:    300,
		FollowerCount:    15,
		FollowingCount:   8,
	}
	mockRepo.On("GetDashboardStats", uint(1)).Return(expected, nil)

	result, err := svc.GetStats(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

// TestUserDashboardGetStats_ZeroValues は全カウントが0のユーザーでも正常に取得できることを確認する。
func TestUserDashboardGetStats_ZeroValues(t *testing.T) {
	mockRepo := &MockUserDashboardRepository{}
	svc := NewUserDashboardService(mockRepo)

	expected := &model.UserDashboardStats{}
	mockRepo.On("GetDashboardStats", uint(2)).Return(expected, nil)

	result, err := svc.GetStats(2)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.PostCount)
	assert.Equal(t, int64(0), result.LikesReceived)
	mockRepo.AssertExpectations(t)
}

// TestUserDashboardGetStats_InvalidUserID はuserID=0のときBadRequestエラーを返すことを確認する。
func TestUserDashboardGetStats_InvalidUserID(t *testing.T) {
	mockRepo := &MockUserDashboardRepository{}
	svc := NewUserDashboardService(mockRepo)

	result, err := svc.GetStats(0)

	assert.Error(t, err)
	assert.Nil(t, result)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	mockRepo.AssertNotCalled(t, "GetDashboardStats")
}

// TestUserDashboardGetStats_RepoError はリポジトリエラーが正しく伝播することを確認する。
func TestUserDashboardGetStats_RepoError(t *testing.T) {
	mockRepo := &MockUserDashboardRepository{}
	svc := NewUserDashboardService(mockRepo)

	mockRepo.On("GetDashboardStats", uint(999)).Return(nil, errors.New("DB error"))

	result, err := svc.GetStats(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "DB error")
	mockRepo.AssertExpectations(t)
}
