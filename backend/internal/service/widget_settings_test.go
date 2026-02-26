package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestWidgetSettingsService() (*WidgetSettingsService, *MockWidgetSettingsRepository) {
	repo := new(MockWidgetSettingsRepository)
	svc := NewWidgetSettingsService(repo)
	return svc, repo
}

// ============================================================
// GetSettings テスト
// ============================================================

func TestWidgetSettingsGetSettings_Success(t *testing.T) {
	svc, repo := newTestWidgetSettingsService()

	settings := &model.WidgetSettings{
		UserID:   1,
		Settings: `[{"key":"userProfile","visible":true,"order":0}]`,
	}
	repo.On("FindByUserID", uint(1)).Return(settings, nil)

	result, err := svc.GetSettings(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
	assert.Contains(t, result.Settings, "userProfile")
	repo.AssertExpectations(t)
}

func TestWidgetSettingsGetSettings_NotFound_ReturnsDefault(t *testing.T) {
	svc, repo := newTestWidgetSettingsService()

	repo.On("FindByUserID", uint(1)).Return(nil, errors.New("record not found"))

	result, err := svc.GetSettings(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
	// デフォルト設定が返される
	assert.Contains(t, result.Settings, "userProfile")
	assert.Contains(t, result.Settings, "level")
	assert.Contains(t, result.Settings, "streak")
}

// ============================================================
// UpdateSettings テスト
// ============================================================

func TestWidgetSettingsUpdateSettings_Success(t *testing.T) {
	svc, repo := newTestWidgetSettingsService()

	newSettings := `[{"key":"userProfile","visible":true,"order":0},{"key":"level","visible":false,"order":1}]`
	repo.On("Upsert", mock.MatchedBy(func(s *model.WidgetSettings) bool {
		return s.UserID == 1 && s.Settings == newSettings
	})).Return(nil)

	err := svc.UpdateSettings(1, newSettings)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestWidgetSettingsUpdateSettings_EmptySettings(t *testing.T) {
	svc, _ := newTestWidgetSettingsService()

	err := svc.UpdateSettings(1, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定は必須です")
}

func TestWidgetSettingsUpdateSettings_InvalidJSON(t *testing.T) {
	svc, _ := newTestWidgetSettingsService()

	err := svc.UpdateSettings(1, "not valid json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定は有効なJSONである必要があります")
}

func TestWidgetSettingsUpdateSettings_RepoError(t *testing.T) {
	svc, repo := newTestWidgetSettingsService()

	validSettings := `[{"key":"userProfile","visible":true,"order":0}]`
	repo.On("Upsert", mock.AnythingOfType("*model.WidgetSettings")).Return(errors.New("db error"))

	err := svc.UpdateSettings(1, validSettings)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestWidgetSettingsUpdateSettings_NotArray(t *testing.T) {
	svc, _ := newTestWidgetSettingsService()

	err := svc.UpdateSettings(1, `{"key":"value"}`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "設定はJSON配列である必要があります")
}
