package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestReminderSettings_GetSettings_Success(t *testing.T) {
	h, svc := setupReminderSettingsHandler()
	svc.On("GetSettings", uint(1)).Return(&model.ReminderSettings{
		Enabled:   true,
		Frequency: "daily",
	}, nil)

	r := newRouter(1)
	r.GET("/reminder-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/reminder-settings", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestReminderSettings_GetSettings_ServiceError(t *testing.T) {
	h, svc := setupReminderSettingsHandler()
	svc.On("GetSettings", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/reminder-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/reminder-settings", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestReminderSettings_UpdateSettings_Success(t *testing.T) {
	h, svc := setupReminderSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("*model.ReminderSettings")).Return(&model.ReminderSettings{
		Enabled:   true,
		Frequency: "weekly",
	}, nil)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"enabled":   true,
		"frequency": "weekly",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestReminderSettings_UpdateSettings_InvalidJSON(t *testing.T) {
	h, _ := setupReminderSettingsHandler()

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequestRaw(r, http.MethodPut, "/reminder-settings", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestReminderSettings_UpdateSettings_ServiceError(t *testing.T) {
	h, svc := setupReminderSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("*model.ReminderSettings")).Return(nil, errors.New("update failed"))

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"enabled": false,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestReminderSettings_UpdateSettings_ValidationError(t *testing.T) {
	h, svc := setupReminderSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("*model.ReminderSettings")).Return(
		nil, domain.NewError(domain.ErrCodeBadRequest, "頻度はdailyまたはweeklyのみ有効です", nil),
	)

	r := newRouter(1)
	r.PUT("/reminder-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/reminder-settings", map[string]interface{}{
		"frequency": "hourly",
	})
	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}
