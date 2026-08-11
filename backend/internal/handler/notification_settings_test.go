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

// mockNotificationSettingsRepo は usecase/repository.NotificationSettingsRepository のモック（ctx 付き）。
type mockNotificationSettingsRepo struct{ mock.Mock }

func (m *mockNotificationSettingsRepo) GetOrCreateDefault(ctx context.Context, userID uint) (*model.NotificationSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NotificationSettings)
	return s, args.Error(1)
}

func (m *mockNotificationSettingsRepo) Save(ctx context.Context, settings *model.NotificationSettings) error {
	return m.Called(ctx, settings).Error(0)
}

// setupNotificationSettingsHandler は本物の usecase と port モックで組む。
func setupNotificationSettingsHandler() (*NotificationSettingsHandler, *mockNotificationSettingsRepo) {
	repo := new(mockNotificationSettingsRepo)
	h := NewNotificationSettingsHandler(
		usecase.NewGetNotificationSettingsUseCase(repo),
		usecase.NewUpdateNotificationSettingsUseCase(repo),
	)
	return h, repo
}

// defaultNotificationSettings は 8 項目すべて有効な設定を返す。
func defaultNotificationSettings() *model.NotificationSettings {
	return &model.NotificationSettings{
		ID: 1, UserID: 1,
		EnableLikes: true, EnableComments: true, EnableFollows: true, EnableMessages: true,
		EnableMentions: true, EnableWebPush: true, EnableEmail: true, EnableSound: true,
	}
}

func TestNotificationSettings_GetSettings_Success(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(defaultNotificationSettings(), nil)

	r := newRouter(1)
	r.GET("/notification-settings", h.GetSettings)
	w := doRequest(r, http.MethodGet, "/notification-settings", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, true, body["enable_likes"])
	repo.AssertExpectations(t)
}

func TestNotificationSettings_GetSettings_RepoError(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/notification-settings", h.GetSettings)
	w := doRequest(r, http.MethodGet, "/notification-settings", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// 8 項目はすべてリクエストの値で上書きされる（部分更新ではない）。
func TestNotificationSettings_UpdateSettings_OverwritesAll(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(defaultNotificationSettings(), nil)
	repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.NotificationSettings) bool {
		return !s.EnableLikes && s.EnableComments && !s.EnableFollows && !s.EnableMessages &&
			!s.EnableMentions && !s.EnableWebPush && !s.EnableEmail && !s.EnableSound
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)
	// enable_comments だけ true にし、他は省略（false として扱われる）
	w := doRequest(r, http.MethodPut, "/notification-settings", map[string]interface{}{
		"enable_comments": true,
	})

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestNotificationSettings_UpdateSettings_InvalidJSON(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)
	w := doRequestRaw(r, http.MethodPut, "/notification-settings", "bad json")

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Save")
}

func TestNotificationSettings_UpdateSettings_RepoError(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(defaultNotificationSettings(), nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)
	w := doRequest(r, http.MethodPut, "/notification-settings", map[string]interface{}{"enable_likes": true})

	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// 設定取得に失敗したら保存しない。
func TestNotificationSettings_UpdateSettings_LoadError(t *testing.T) {
	h, repo := setupNotificationSettingsHandler()
	repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.PUT("/notification-settings", h.UpdateSettings)
	w := doRequest(r, http.MethodPut, "/notification-settings", map[string]interface{}{"enable_likes": true})

	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertNotCalled(t, "Save")
}
