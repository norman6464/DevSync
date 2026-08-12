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

// mockBadgeStatsPort は usecase/repository.BadgeStatsReader のモック。
type mockBadgeStatsPort struct{ mock.Mock }

func (m *mockBadgeStatsPort) GetBadgeStats(ctx context.Context, userID uint) (*model.BadgeStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.BadgeStats)
	return s, args.Error(1)
}

// mockNotificationCreator は usecase/repository.NotificationCreator のモック。
type mockNotificationCreator struct{ mock.Mock }

func (m *mockNotificationCreator) Create(ctx context.Context, notification *model.Notification) error {
	return m.Called(ctx, notification).Error(0)
}

// newTestBadgeHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestBadgeHandler() (*BadgeHandler, *mockBadgeStatsPort, *mockNotificationCreator) {
	stats := new(mockBadgeStatsPort)
	notifications := new(mockNotificationCreator)
	h := NewBadgeHandler(
		usecase.NewGetUserBadgesUseCase(stats),
		usecase.NewNotifyBadgeEarnedUseCase(notifications),
	)
	return h, stats, notifications
}

// ---------- GetUserBadges ----------

func TestBadgeGetUserBadges_Success(t *testing.T) {
	h, stats, _ := newTestBadgeHandler()
	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)

	stats.On("GetBadgeStats", mock.Anything, uint(1)).
		Return(&model.BadgeStats{TotalContributions: 60, TotalPosts: 1}, nil)

	w := doRequest(r, http.MethodGet, "/users/1/badges", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	// 閾値を超えたバッジは獲得済み、超えないものは未獲得で返る。
	assert.Contains(t, body, `{"id":"first-commit","name":"badges.firstCommit","description":"badges.firstCommitDesc","category":"contribution","earned":true}`)
	assert.Contains(t, body, `{"id":"code-warrior","name":"badges.codeWarrior","description":"badges.codeWarriorDesc","category":"contribution","earned":false}`)
	assert.Contains(t, body, `{"id":"first-post","name":"badges.firstPost","description":"badges.firstPostDesc","category":"post","earned":true}`)
	stats.AssertExpectations(t)
}

// 学習ログのストリークが GitHub より長ければそちらで判定する。
func TestBadgeGetUserBadges_UsesLongerStreak(t *testing.T) {
	h, stats, _ := newTestBadgeHandler()
	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)

	stats.On("GetBadgeStats", mock.Anything, uint(1)).
		Return(&model.BadgeStats{CurrentStreak: 3, LearningLogStreak: 8}, nil)

	w := doRequest(r, http.MethodGet, "/users/1/badges", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"id":"week-streak","name":"badges.weekStreak","description":"badges.weekStreakDesc","category":"streak","earned":true`)
	stats.AssertExpectations(t)
}

func TestBadgeGetUserBadges_InvalidID(t *testing.T) {
	h, stats, _ := newTestBadgeHandler()
	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)

	w := doRequest(r, http.MethodGet, "/users/abc/badges", nil)
	assertStatus(t, w, http.StatusBadRequest)
	stats.AssertNotCalled(t, "GetBadgeStats", mock.Anything, mock.Anything)
}

func TestBadgeGetUserBadges_RepositoryError(t *testing.T) {
	h, stats, _ := newTestBadgeHandler()
	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)

	stats.On("GetBadgeStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/1/badges", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}

// ---------- NotifyBadgeEarned ----------

func TestBadgeNotifyBadgeEarned_Success(t *testing.T) {
	h, _, notifications := newTestBadgeHandler()
	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)

	notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
		// 通知は認証ユーザー本人宛に、本人を実行者として作られる。
		return n.UserID == 1 && n.ActorID == 1 && n.Type == model.NotificationTypeBadge &&
			n.BadgeID != nil && *n.BadgeID == "first-commit"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/badges/notify", map[string]string{"badge_id": "first-commit"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "badge notification created")
	notifications.AssertExpectations(t)
}

func TestBadgeNotifyBadgeEarned_MissingBadgeID(t *testing.T) {
	h, _, notifications := newTestBadgeHandler()
	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)

	w := doRequest(r, http.MethodPost, "/badges/notify", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
	notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestBadgeNotifyBadgeEarned_RepositoryError(t *testing.T) {
	h, _, notifications := newTestBadgeHandler()
	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)

	notifications.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/badges/notify", map[string]string{"badge_id": "first-commit"})
	assertStatus(t, w, http.StatusInternalServerError)
	notifications.AssertExpectations(t)
}
