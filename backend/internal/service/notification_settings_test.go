package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationSettingsRepository は NotificationSettingsRepository のモック実装。
type MockNotificationSettingsRepository struct{ mock.Mock }

func (m *MockNotificationSettingsRepository) CreateOrUpdate(settings *model.NotificationSettings) error {
	return m.Called(settings).Error(0)
}

func (m *MockNotificationSettingsRepository) GetByUserID(userID uint) (*model.NotificationSettings, error) {
	args := m.Called(userID)
	if settings := args.Get(0); settings != nil {
		return settings.(*model.NotificationSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNotificationSettingsRepository) GetOrCreateDefault(userID uint) (*model.NotificationSettings, error) {
	args := m.Called(userID)
	if settings := args.Get(0); settings != nil {
		return settings.(*model.NotificationSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupNotificationSettingsService() (*NotificationSettingsService, *MockNotificationSettingsRepository) {
	repo := new(MockNotificationSettingsRepository)
	svc := NewNotificationSettingsService(repo)
	return svc, repo
}

func TestNotificationSettingsService_GetSettings(t *testing.T) {
	svc, repo := setupNotificationSettingsService()

	settings := &model.NotificationSettings{
		ID:             1,
		UserID:         1,
		EnableLikes:    true,
		EnableComments: false,
		EnableSound:    true,
	}

	repo.On("GetOrCreateDefault", uint(1)).Return(settings, nil)

	result, err := svc.GetSettings(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, true, result.EnableLikes)
	assert.Equal(t, false, result.EnableComments)
	repo.AssertExpectations(t)
}

func TestNotificationSettingsService_UpdateSettings(t *testing.T) {
	svc, repo := setupNotificationSettingsService()

	existingSettings := &model.NotificationSettings{
		ID:             1,
		UserID:         1,
		EnableLikes:    true,
		EnableComments: true,
		EnableFollows:  true,
		EnableMessages: true,
		EnableMentions: true,
		EnableWebPush:  true,
		EnableEmail:    true,
		EnableSound:    true,
	}

	repo.On("GetOrCreateDefault", uint(1)).Return(existingSettings, nil)
	repo.On("CreateOrUpdate", mock.AnythingOfType("*model.NotificationSettings")).Return(nil)

	updates := &model.NotificationSettings{
		EnableLikes:    false,
		EnableComments: true,
		EnableFollows:  false,
		EnableMessages: true,
		EnableMentions: true,
		EnableWebPush:  true,
		EnableEmail:    true,
		EnableSound:    false,
	}

	result, err := svc.UpdateSettings(1, updates)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, false, result.EnableLikes)    // 更新された
	assert.Equal(t, true, result.EnableComments)  // 更新された
	assert.Equal(t, false, result.EnableSound)    // 更新された
	assert.Equal(t, false, result.EnableFollows)  // 更新された
	repo.AssertExpectations(t)
}

func TestNotificationSettingsService_ShouldNotify(t *testing.T) {
	svc, repo := setupNotificationSettingsService()

	settings := &model.NotificationSettings{
		UserID:         1,
		EnableLikes:    true,
		EnableComments: false,
		EnableFollows:  true,
		EnableMessages: true,
	}

	repo.On("GetOrCreateDefault", uint(1)).Return(settings, nil)

	// いいね通知: 有効
	should, err := svc.ShouldNotify(1, model.NotificationTypeLike)
	assert.NoError(t, err)
	assert.True(t, should)

	// コメント通知: 無効
	should, err = svc.ShouldNotify(1, model.NotificationTypeComment)
	assert.NoError(t, err)
	assert.False(t, should)

	// フォロー通知: 有効
	should, err = svc.ShouldNotify(1, model.NotificationTypeFollow)
	assert.NoError(t, err)
	assert.True(t, should)

	// メッセージ通知: 有効
	should, err = svc.ShouldNotify(1, model.NotificationTypeMessage)
	assert.NoError(t, err)
	assert.True(t, should)

	repo.AssertExpectations(t)
}
