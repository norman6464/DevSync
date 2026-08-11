package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestEmailPreferencesHandler は本物の usecase に port モックを注入したハンドラーを生成する。
// port モック（mockUserPort）は user スライスのテストと共用する。
func newTestEmailPreferencesHandler() (*EmailPreferencesHandler, *mockUserPort) {
	repo := new(mockUserPort)
	h := NewEmailPreferencesHandler(
		usecase.NewGetEmailPreferencesUseCase(repo),
		usecase.NewUpdateEmailPreferencesUseCase(repo),
	)
	return h, repo
}

// ============================================================
// 取得
// ============================================================

func TestEmailPreferencesHandler_GetPreferences(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.GET("/email-preferences", h.GetPreferences)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)

	w := doRequest(r, http.MethodGet, "/email-preferences", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"email_weekly_report":true`)
	assert.Contains(t, body, `"email_language":"ja"`)
	repo.AssertExpectations(t)
}

func TestEmailPreferencesHandler_GetPreferences_NotFound(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.GET("/email-preferences", h.GetPreferences)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/email-preferences", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// DB 障害も 404 に潰す（移行前から変わらない挙動）。
func TestEmailPreferencesHandler_GetPreferences_RepositoryError(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.GET("/email-preferences", h.GetPreferences)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/email-preferences", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ============================================================
// 更新
// ============================================================

func TestEmailPreferencesHandler_UpdatePreferences(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return !u.EmailWeeklyReport && u.EmailLanguage == "en"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_weekly_report": false, "email_language": "en",
	})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"email_language":"en"`)
	repo.AssertExpectations(t)
}

// 指定しなかった項目は据え置く。
func TestEmailPreferencesHandler_UpdatePreferences_Partial(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	repo.On("FindByID", mock.Anything, uint(1)).
		Return(&model.User{ID: 1, EmailWeeklyReport: true, EmailLanguage: "ja"}, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
		return u.EmailWeeklyReport && u.EmailLanguage == "ko"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{"email_language": "ko"})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestEmailPreferencesHandler_UpdatePreferences_InvalidLanguage(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{
		"email_language": "xx",
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "無効なメール言語設定です")
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestEmailPreferencesHandler_UpdatePreferences_InvalidJSON(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	w := doRequestRaw(r, http.MethodPut, "/email-preferences", "invalid json")
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestEmailPreferencesHandler_UpdatePreferences_NotFound(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{"email_language": "ja"})
	assertStatus(t, w, http.StatusNotFound)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestEmailPreferencesHandler_UpdatePreferences_RepositoryError(t *testing.T) {
	h, repo := newTestEmailPreferencesHandler()
	r := newRouter(1)
	r.PUT("/email-preferences", h.UpdatePreferences)

	repo.On("FindByID", mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
	repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/email-preferences", map[string]interface{}{"email_language": "ja"})
	assertStatus(t, w, http.StatusInternalServerError)
}
