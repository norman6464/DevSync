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

// mockXPStatsPort は usecase/repository.XPStatsReader のモック。
type mockXPStatsPort struct{ mock.Mock }

func (m *mockXPStatsPort) GetXPStats(ctx context.Context, userID uint) (*model.XPStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.XPStats)
	return s, args.Error(1)
}

// newTestLevelHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestLevelHandler() (*LevelHandler, *mockXPStatsPort) {
	stats := new(mockXPStatsPort)
	h := NewLevelHandler(
		usecase.NewGetLevelInfoUseCase(stats),
		usecase.NewGetXPBreakdownUseCase(stats),
	)
	return h, stats
}

// ---------- GetMyLevelInfo ----------

func TestLevelGetMyLevelInfo_Success(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/levels/me", h.GetMyLevelInfo)

	// 学習ログ 10 件 * 10 + 600 分 * 0.5 = 400 XP
	stats.On("GetXPStats", mock.Anything, uint(1)).
		Return(&model.XPStats{LearningLogCount: 10, LearningLogTotalDuration: 600}, nil)

	w := doRequest(r, http.MethodGet, "/levels/me", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	// 400 XP はレベル 2（必要累計 XP: Lv1=100, Lv2=300, Lv3=600）
	assertJSONEqual(t, body, "level", float64(2))
	assertJSONEqual(t, body, "total_xp", float64(400))
	assertJSONEqual(t, body, "current_level_xp", float64(300))
	assertJSONEqual(t, body, "next_level_xp", float64(600))
	assertJSONEqual(t, body, "progress_xp", float64(100))
	stats.AssertExpectations(t)
}

func TestLevelGetMyLevelInfo_RepositoryError(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/levels/me", h.GetMyLevelInfo)

	stats.On("GetXPStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/levels/me", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}

// XP が 0 ならレベル 0 で進捗も 0。
func TestLevelGetMyLevelInfo_ZeroXP(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/levels/me", h.GetMyLevelInfo)

	stats.On("GetXPStats", mock.Anything, uint(1)).Return(&model.XPStats{}, nil)

	w := doRequest(r, http.MethodGet, "/levels/me", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "level", float64(0))
	assertJSONEqual(t, body, "total_xp", float64(0))
	assertJSONEqual(t, body, "progress_percent", float64(0))
	stats.AssertExpectations(t)
}

// ---------- GetLevelInfo ----------

func TestLevelGetLevelInfo_Success(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/level", h.GetLevelInfo)

	// 投稿 20 件 * 30 = 600 XP → レベル 3
	stats.On("GetXPStats", mock.Anything, uint(5)).Return(&model.XPStats{PostCount: 20}, nil)

	w := doRequest(r, http.MethodGet, "/users/5/level", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "level", float64(3))
	assertJSONEqual(t, body, "total_xp", float64(600))
	stats.AssertExpectations(t)
}

func TestLevelGetLevelInfo_InvalidID(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/level", h.GetLevelInfo)

	w := doRequest(r, http.MethodGet, "/users/abc/level", nil)
	assertStatus(t, w, http.StatusBadRequest)
	stats.AssertNotCalled(t, "GetXPStats", mock.Anything, mock.Anything)
}

func TestLevelGetLevelInfo_RepositoryError(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/level", h.GetLevelInfo)

	stats.On("GetXPStats", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/5/level", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}

// ---------- GetXPBreakdown ----------

func TestLevelGetXPBreakdown_Success(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)

	stats.On("GetXPStats", mock.Anything, uint(5)).Return(&model.XPStats{
		LearningLogCount:         3,  // 30
		LearningLogTotalDuration: 90, // 45
		PostCount:                2,  // 60
		GitHubContributionDays:   4,  // 20
		CompletedGoals:           1,  // 50
		CommentCount:             2,  // 10
		LikesReceived:            5,  // 15
		CurrentStreak:            15, // (15/7)*20 = 40
	}, nil)

	w := doRequest(r, http.MethodGet, "/users/5/xp", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "learning_logs", float64(75))
	assertJSONEqual(t, body, "posts", float64(60))
	assertJSONEqual(t, body, "github", float64(20))
	assertJSONEqual(t, body, "goals", float64(50))
	assertJSONEqual(t, body, "comments", float64(10))
	assertJSONEqual(t, body, "likes", float64(15))
	assertJSONEqual(t, body, "streak_bonus", float64(40))
	assertJSONEqual(t, body, "total", float64(270))
	stats.AssertExpectations(t)
}

func TestLevelGetXPBreakdown_InvalidID(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)

	w := doRequest(r, http.MethodGet, "/users/abc/xp", nil)
	assertStatus(t, w, http.StatusBadRequest)
	stats.AssertNotCalled(t, "GetXPStats", mock.Anything, mock.Anything)
}

func TestLevelGetXPBreakdown_RepositoryError(t *testing.T) {
	h, stats := newTestLevelHandler()
	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)

	stats.On("GetXPStats", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/5/xp", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}

// ---------- 純粋関数 ----------

func TestCalculateLevelBoundaries(t *testing.T) {
	tests := []struct {
		totalXP int
		want    int
	}{
		{0, 0}, {99, 0}, {100, 1}, {299, 1}, {300, 2}, {599, 2}, {600, 3}, {1000, 4},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, usecase.CalculateLevel(tt.totalXP), "totalXP=%d", tt.totalXP)
	}
}

func TestXPForLevel(t *testing.T) {
	assert.Equal(t, 0, usecase.XPForLevel(0))
	assert.Equal(t, 100, usecase.XPForLevel(1))
	assert.Equal(t, 300, usecase.XPForLevel(2))
	assert.Equal(t, 600, usecase.XPForLevel(3))
}
