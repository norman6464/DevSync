package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// GetSettings テスト
// ============================================================

func TestWidgetSettings_GetSettings_Success(t *testing.T) {
	h, svc := setupWidgetSettingsHandler()
	settings := &model.WidgetSettings{
		UserID:   1,
		Settings: `[{"key":"userProfile","visible":true,"order":0}]`,
	}
	svc.On("GetSettings", uint(1)).Return(settings, nil)

	r := newRouter(1)
	r.GET("/widget-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/widget-settings", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "userProfile")
	svc.AssertExpectations(t)
}

func TestWidgetSettings_GetSettings_ServiceError(t *testing.T) {
	h, svc := setupWidgetSettingsHandler()
	svc.On("GetSettings", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/widget-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/widget-settings", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// UpdateSettings テスト
// ============================================================

func TestWidgetSettings_UpdateSettings_Success(t *testing.T) {
	h, svc := setupWidgetSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("string")).Return(nil)

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{
		"settings": []map[string]interface{}{
			{"key": "userProfile", "visible": true, "order": 0},
			{"key": "level", "visible": false, "order": 1},
		},
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
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

func TestWidgetSettings_UpdateSettings_ServiceError(t *testing.T) {
	h, svc := setupWidgetSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("string")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/widget-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/widget-settings", map[string]interface{}{
		"settings": []map[string]interface{}{
			{"key": "userProfile", "visible": true, "order": 0},
		},
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
