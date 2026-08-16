package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockReminderSettingsRepo は usecase/repository.ReminderSettingsRepository のモック（ctx 付き）。
type mockReminderSettingsRepo struct{ mock.Mock }

func (m *mockReminderSettingsRepo) GetOrCreateDefault(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ReminderSettings)
	return s, args.Error(1)
}

func (m *mockReminderSettingsRepo) FindByUserID(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ReminderSettings)
	return s, args.Error(1)
}

func (m *mockReminderSettingsRepo) Save(ctx context.Context, settings *model.ReminderSettings) error {
	return m.Called(ctx, settings).Error(0)
}

// setupReminderSettingsHandler は本物の usecase と port モックで ReminderSettingsHandler を組む。
func setupReminderSettingsHandler() (*ReminderSettingsHandler, *mockReminderSettingsRepo) {
	repo := new(mockReminderSettingsRepo)
	h := NewReminderSettingsHandler(
		usecase.NewGetReminderSettingsUseCase(repo),
		usecase.NewUpdateReminderSettingsUseCase(repo),
	)
	return h, repo
}

// storedReminderSettings はハンドラーテストの起点となる既存設定を返す。
func storedReminderSettings() *model.ReminderSettings {
	return &model.ReminderSettings{
		ID:               10,
		UserID:           1,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
		EnableWeb:        true,
	}
}

// ============================================================
// リマインダー設定: 取得ハンドラーテスト
// ============================================================

func TestReminderSettings_GetSettings_Success(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)

	r := newRouter(1)
	r.GET("/reminder-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/reminder-settings", nil)
	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["frequency"] != "daily" {
		t.Errorf("expected frequency=daily, got %v", result["frequency"])
	}
	if result["notification_time"] != "09:00" {
		t.Errorf("expected notification_time=09:00, got %v", result["notification_time"])
	}
	repo.AssertExpectations(t)
}

func TestReminderSettings_GetSettings_ServiceError(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/reminder-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/reminder-settings", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// リマインダー設定: 更新ハンドラーテスト
// ============================================================

func TestReminderSettings_UpdateSettings_Success(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)
	repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.ReminderSettings) bool {
		return s.Frequency == model.ReminderFrequencyWeekly && s.Enabled
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"enabled":   true,
		"frequency": "weekly",
	})
	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["frequency"] != "weekly" {
		t.Errorf("expected frequency=weekly, got %v", result["frequency"])
	}
	repo.AssertExpectations(t)
}

func TestReminderSettings_UpdateSettings_InvalidJSON(t *testing.T) {
	h, _ := setupReminderSettingsHandler()

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequestRaw(r, http.MethodPut, "/reminder-settings", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReminderSettings_UpdateSettings_ServiceError(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("update failed"))

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"enabled": false,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// 未知の頻度は 400 を返し、保存しない。
func TestReminderSettings_UpdateSettings_ValidationError(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"frequency": "hourly",
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

// 通知時間の形式違反は 400 を返し、保存しない。
func TestReminderSettings_UpdateSettings_InvalidNotificationTime(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"notification_time": "9:00",
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

// 非活動日数の上限超過は 400 を返し、保存しない。
func TestReminderSettings_UpdateSettings_InactiveDaysTooLarge(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"inactive_days": 31,
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

// 設定が未登録のまま更新されてもデフォルト設定を作成して反映する（従来は 500 だった）。
func TestReminderSettings_UpdateSettings_FirstTime(t *testing.T) {
	h, repo := setupReminderSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(storedReminderSettings(), nil)
	repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.ReminderSettings) bool {
		return s.Enabled
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"enabled": true,
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}
