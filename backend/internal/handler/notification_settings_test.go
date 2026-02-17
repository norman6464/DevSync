package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestNotificationSettings_GetSettings_Success(t *testing.T) {
	h, svc := setupNotificationSettingsHandler()
	svc.On("GetSettings", uint(1)).Return(&model.NotificationSettings{
		EnableLikes:    true,
		EnableComments: true,
	}, nil)

	r := newRouter(1)
	r.GET("/notification-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/notification-settings", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNotificationSettings_GetSettings_ServiceError(t *testing.T) {
	h, svc := setupNotificationSettingsHandler()
	svc.On("GetSettings", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/notification-settings", h.GetSettings)

	w := doRequest(r, http.MethodGet, "/notification-settings", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestNotificationSettings_UpdateSettings_Success(t *testing.T) {
	h, svc := setupNotificationSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("*model.NotificationSettings")).Return(&model.NotificationSettings{
		EnableLikes:    false,
		EnableComments: true,
	}, nil)

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/notification-settings", map[string]interface{}{
		"enable_likes":    false,
		"enable_comments": true,
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestNotificationSettings_UpdateSettings_InvalidJSON(t *testing.T) {
	h, _ := setupNotificationSettingsHandler()

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)

	w := doRequestRaw(r, http.MethodPut, "/notification-settings", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNotificationSettings_UpdateSettings_ServiceError(t *testing.T) {
	h, svc := setupNotificationSettingsHandler()
	svc.On("UpdateSettings", uint(1), mock.AnythingOfType("*model.NotificationSettings")).Return(nil, errors.New("update failed"))

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)

	w := doRequest(r, http.MethodPut, "/notification-settings", map[string]interface{}{
		"enable_likes": false,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
