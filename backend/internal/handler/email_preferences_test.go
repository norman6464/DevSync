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
	svc.On("UpdateEmailPreferences", uint(1), mock.AnythingOfType("*bool"), mock.AnythingOfType("*string")).
		Return(&model.User{EmailWeeklyReport: false, EmailLanguage: "en"}, nil)

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
	svc.On("UpdateEmailPreferences", uint(1), (*bool)(nil), mock.AnythingOfType("*string")).
		Return(nil, errors.New("無効なメール言語設定です"))

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_language": "invalid",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestEmailPreferences_UpdatePreferences_ServiceError(t *testing.T) {
	h, svc := setupEmailPreferencesHandler()
	svc.On("UpdateEmailPreferences", uint(1), mock.AnythingOfType("*bool"), (*string)(nil)).
		Return(nil, errors.New("update failed"))

	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_weekly_report": false,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
