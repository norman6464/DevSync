package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmailPreferences_GetPreferences_Success(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	svc.On("GetByID", uint(1)).Return(&model.User{
		EmailWeeklyReport: true,
		EmailLanguage:     "ja",
	}, nil)

	r := newRouter(1)
	r.GET("/email-preferences", h.GetPreferences)

	w := doRequest(r, http.MethodGet, "/email-preferences", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, true, body["email_weekly_report"])
	assert.Equal(t, "ja", body["email_language"])
	svc.AssertExpectations(t)
}

func TestEmailPreferences_GetPreferences_ServiceError(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	svc.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/email-preferences", h.GetPreferences)

	w := doRequest(r, http.MethodGet, "/email-preferences", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestEmailPreferences_UpdatePreferences_Success(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	user := &model.User{EmailWeeklyReport: true, EmailLanguage: "ja"}
	svc.On("GetByID", uint(1)).Return(user, nil)
	svc.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_weekly_report": false,
		"email_language":      "en",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestEmailPreferences_UpdatePreferences_InvalidJSON(t *testing.T) {
	h, _ := setupEmailPreferencesHandler()

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequestRaw(r, http.MethodPut, "/email-preferences", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestEmailPreferences_UpdatePreferences_InvalidLanguage(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	user := &model.User{EmailWeeklyReport: true, EmailLanguage: "ja"}
	svc.On("GetByID", uint(1)).Return(user, nil)

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_language": "invalid",
	})
	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}

func TestEmailPreferences_UpdatePreferences_ServiceError(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	user := &model.User{EmailWeeklyReport: true, EmailLanguage: "ja"}
	svc.On("GetByID", uint(1)).Return(user, nil)
	svc.On("Update", mock.AnythingOfType("*model.User")).Return(errors.New("update failed"))

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_weekly_report": false,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
