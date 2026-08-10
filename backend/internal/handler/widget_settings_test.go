package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockWidgetSettingsRepo は usecase/repository.WidgetSettingsRepository のモック（ctx 付き）。
type mockWidgetSettingsRepo struct{ mock.Mock }

func (m *mockWidgetSettingsRepo) FindByUserID(ctx context.Context, userID uint) (*model.WidgetSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.WidgetSettings)
	return s, args.Error(1)
}

func (m *mockWidgetSettingsRepo) Upsert(ctx context.Context, settings *model.WidgetSettings) error {
	return m.Called(ctx, settings).Error(0)
}

// setupWidgetSettingsHandler は本物の usecase と port モックで WidgetSettingsHandler を組む。
func setupWidgetSettingsHandler() (*WidgetSettingsHandler, *mockWidgetSettingsRepo) {
	repo := new(mockWidgetSettingsRepo)
	h := NewWidgetSettingsHandler(
		usecase.NewGetWidgetSettingsUseCase(repo),
		usecase.NewUpdateWidgetSettingsUseCase(repo),
	)
	return h, repo
}

// ============================================================
// GetSettings テスト
// ============================================================

func TestWidgetSettings_GetSettings_Success(t *testing.T) {
	h, repo := setupWidgetSettingsHandler()
	settings := &model.WidgetSettings{
		UserID:   1,
		Settings: `[{"key":"userProfile","visible":true,"order":0}]`,
	}
	repo.On("FindByUserID", mock.Anything, uint(1)).Return(settings, nil)

	r := newRouter(1)
	r.GET("/widget-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/widget-settings", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "userProfile")
	repo.AssertExpectations(t)
}

// 未登録（および repo エラー）のときはデフォルト設定を 200 で返す既存挙動を確認する。
func TestWidgetSettings_GetSettings_NotFoundReturnsDefaults(t *testing.T) {
	h, repo := setupWidgetSettingsHandler()
	repo.On("FindByUserID", mock.Anything, uint(1)).Return((*model.WidgetSettings)(nil), errors.New("record not found"))

	r := newRouter(1)
	r.GET("/widget-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/widget-settings", nil)
	assertStatus(t, w, http.StatusOK)
	// デフォルト配置の 14 項目が返る
	assert.Contains(t, w.Body.String(), "quickStats")
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateSettings テスト
// ============================================================

func TestWidgetSettings_UpdateSettings_Success(t *testing.T) {
	h, repo := setupWidgetSettingsHandler()
	repo.On("Upsert", mock.Anything, mock.AnythingOfType("*model.WidgetSettings")).Return(nil)

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{
		"settings": []map[string]interface{}{
			{"key": "userProfile", "visible": true, "order": 0},
			{"key": "level", "visible": false, "order": 1},
		},
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestWidgetSettings_UpdateSettings_InvalidJSON(t *testing.T) {
	h, _ := setupWidgetSettingsHandler()

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequestRaw(r, http.MethodPut, "/widget-settings", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestWidgetSettings_UpdateSettings_MissingSettings(t *testing.T) {
	h, _ := setupWidgetSettingsHandler()

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
}

// JSON ではあるが配列でない場合は Upsert せず 400 を返す。
func TestWidgetSettings_UpdateSettings_NotArray(t *testing.T) {
	h, repo := setupWidgetSettingsHandler()

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{
		"settings": map[string]interface{}{"key": "userProfile"},
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Upsert")
}

func TestWidgetSettings_UpdateSettings_ServiceError(t *testing.T) {
	h, repo := setupWidgetSettingsHandler()
	repo.On("Upsert", mock.Anything, mock.AnythingOfType("*model.WidgetSettings")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{
		"settings": []map[string]interface{}{
			{"key": "userProfile", "visible": true, "order": 0},
		},
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
