package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReminderSettingsRepository はReminderSettingsRepositoryのモック実装。
type MockReminderSettingsRepository struct {
	mock.Mock
}

func (m *MockReminderSettingsRepository) CreateOrUpdate(settings *model.ReminderSettings) error {
	return m.Called(settings).Error(0)
}

func (m *MockReminderSettingsRepository) GetByUserID(userID uint) (*model.ReminderSettings, error) {
	args := m.Called(userID)
	if settings := args.Get(0); settings != nil {
		return settings.(*model.ReminderSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReminderSettingsRepository) GetOrCreateDefault(userID uint) (*model.ReminderSettings, error) {
	args := m.Called(userID)
	if settings := args.Get(0); settings != nil {
		return settings.(*model.ReminderSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReminderSettingsRepository) GetEnabledSettings() ([]model.ReminderSettings, error) {
	args := m.Called()
	return args.Get(0).([]model.ReminderSettings), args.Error(1)
}

func (m *MockReminderSettingsRepository) UpdateLastRemindedAt(userID uint) error {
	return m.Called(userID).Error(0)
}

// newTestReminderSettingsService はテスト用のReminderSettingsServiceを生成する。
func newTestReminderSettingsService() (*ReminderSettingsService, *MockReminderSettingsRepository) {
	repo := new(MockReminderSettingsRepository)
	svc := NewReminderSettingsService(repo)
	return svc, repo
}

// ============================================================
// GetSettings テスト
// ============================================================

func TestReminderSettingsService_GetSettings(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	settings := &model.ReminderSettings{
		UserID:           1,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
	}

	repo.On("GetOrCreateDefault", uint(1)).Return(settings, nil)

	result, err := svc.GetSettings(1)
	assert.NoError(t, err)
	assert.Equal(t, settings, result)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateSettings テスト
// ============================================================

func TestReminderSettingsService_UpdateSettings(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	existingSettings := &model.ReminderSettings{
		ID:               1,
		UserID:           1,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
		EnableWeb:        true,
		EnableEmail:      false,
	}

	updates := &model.ReminderSettings{
		Frequency:   model.ReminderFrequencyWeekly,
		EnableEmail: true,
	}

	repo.On("GetByUserID", uint(1)).Return(existingSettings, nil)
	repo.On("CreateOrUpdate", mock.MatchedBy(func(s *model.ReminderSettings) bool {
		return s.UserID == 1 && s.Frequency == model.ReminderFrequencyWeekly && s.EnableEmail == true
	})).Return(nil)

	result, err := svc.UpdateSettings(1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.ReminderFrequencyWeekly, result.Frequency)
	assert.Equal(t, true, result.EnableEmail)
	repo.AssertExpectations(t)
}

// ============================================================
// ShouldRemind テスト
// ============================================================

func TestReminderSettingsService_ShouldRemind(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	tests := []struct {
		name           string
		settings       *model.ReminderSettings
		lastActivity   time.Time
		expectedResult bool
	}{
		{
			name: "リマインダー無効 → false",
			settings: &model.ReminderSettings{
				UserID:       1,
				Enabled:      false,
				InactiveDays: 3,
			},
			lastActivity:   time.Now().Add(-5 * 24 * time.Hour),
			expectedResult: false,
		},
		{
			name: "非アクティブ期間未満 → false",
			settings: &model.ReminderSettings{
				UserID:       1,
				Enabled:      true,
				InactiveDays: 3,
			},
			lastActivity:   time.Now().Add(-2 * 24 * time.Hour),
			expectedResult: false,
		},
		{
			name: "非アクティブ期間超過 → true",
			settings: &model.ReminderSettings{
				UserID:       1,
				Enabled:      true,
				InactiveDays: 3,
			},
			lastActivity:   time.Now().Add(-4 * 24 * time.Hour),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByUserID", tt.settings.UserID).Return(tt.settings, nil).Once()

			result := svc.ShouldRemind(tt.settings.UserID, tt.lastActivity)
			assert.Equal(t, tt.expectedResult, result)
		})
	}

	repo.AssertExpectations(t)
}

// ============================================================
// エラーケーステスト
// ============================================================

func TestReminderSettingsService_GetSettings_Error(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	repo.On("GetOrCreateDefault", uint(1)).Return(nil, errors.New("db error"))

	result, err := svc.GetSettings(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestReminderSettingsService_UpdateSettings_GetError(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	repo.On("GetByUserID", uint(1)).Return(nil, errors.New("db error"))

	updates := &model.ReminderSettings{Frequency: model.ReminderFrequencyDaily}
	result, err := svc.UpdateSettings(1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestReminderSettingsService_UpdateSettings_SaveError(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	existing := &model.ReminderSettings{
		ID:               1,
		UserID:           1,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
	}

	repo.On("GetByUserID", uint(1)).Return(existing, nil)
	repo.On("CreateOrUpdate", mock.AnythingOfType("*model.ReminderSettings")).Return(errors.New("save error"))

	updates := &model.ReminderSettings{Frequency: model.ReminderFrequencyWeekly}
	result, err := svc.UpdateSettings(1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestReminderSettingsService_ShouldRemind_RepoError(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	repo.On("GetByUserID", uint(1)).Return(nil, errors.New("db error"))

	result := svc.ShouldRemind(1, time.Now().Add(-5*24*time.Hour))
	assert.False(t, result)
	repo.AssertExpectations(t)
}

func TestReminderSettingsService_SendReminder_Success(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	repo.On("UpdateLastRemindedAt", uint(1)).Return(nil)

	err := svc.SendReminder(1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestReminderSettingsService_SendReminder_Error(t *testing.T) {
	svc, repo := newTestReminderSettingsService()

	repo.On("UpdateLastRemindedAt", uint(1)).Return(errors.New("db error"))

	err := svc.SendReminder(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
