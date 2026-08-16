package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"github.com/stretchr/testify/mock"
)

// mockStreakFreezeRepo は usecase/repository.StreakFreezeRepository のモック（ctx 付き）。
type mockStreakFreezeRepo struct{ mock.Mock }

func (m *mockStreakFreezeRepo) CreateWithinLimits(ctx context.Context, freeze *model.StreakFreeze, maxPerMonth int) (repository.FreezeUseOutcome, error) {
	args := m.Called(ctx, freeze, maxPerMonth)
	return args.Get(0).(repository.FreezeUseOutcome), args.Error(1)
}

func (m *mockStreakFreezeRepo) GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error) {
	args := m.Called(ctx, userID, year, month)
	f, _ := args.Get(0).([]model.StreakFreeze)
	return f, args.Error(1)
}

func (m *mockStreakFreezeRepo) HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error) {
	args := m.Called(ctx, userID, date)
	return args.Bool(0), args.Error(1)
}

// setupStreakFreezeHandler は本物の usecase と port モックで StreakFreezeHandler を組む。
func setupStreakFreezeHandler() (*StreakFreezeHandler, *mockStreakFreezeRepo) {
	repo := new(mockStreakFreezeRepo)
	h := NewStreakFreezeHandler(
		usecase.NewUseStreakFreezeUseCase(repo),
		usecase.NewGetStreakFreezeStatusUseCase(repo),
	)
	return h, repo
}

// nowParts は usecase が使う「今日」「今年」「今月」を返す。
func nowParts() (today string, year, month int) {
	now := time.Now()
	return now.Format("2006-01-02"), now.Year(), int(now.Month())
}

// ============================================================
// ストリークフリーズ: 使用ハンドラーテスト
// ============================================================

func TestStreakFreeze_UseFreeze_Success(t *testing.T) {
	h, repo := setupStreakFreezeHandler()
	today, year, month := nowParts()

	repo.On("CreateWithinLimits", mock.Anything, mock.MatchedBy(func(f *model.StreakFreeze) bool {
		return f.UserID == 1 && f.UsedDate == today && f.Year == year && f.Month == month
	}), model.MaxFreezesPerMonth).Return(repository.FreezeUseCreated, nil)

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

// 当日すでに使用済みなら 409 を返す。
func TestStreakFreeze_UseFreeze_AlreadyUsed(t *testing.T) {
	h, repo := setupStreakFreezeHandler()

	repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
		Return(repository.FreezeUseDuplicateDay, nil)

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusConflict)
	repo.AssertExpectations(t)
}

// 今月の上限に達していれば 400 を返す。
func TestStreakFreeze_UseFreeze_MonthlyLimit(t *testing.T) {
	h, repo := setupStreakFreezeHandler()

	repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
		Return(repository.FreezeUseMonthlyLimitReached, nil)

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertExpectations(t)
}

func TestStreakFreeze_UseFreeze_ServiceInternalError(t *testing.T) {
	h, repo := setupStreakFreezeHandler()

	repo.On("CreateWithinLimits", mock.Anything, mock.Anything, model.MaxFreezesPerMonth).
		Return(repository.FreezeUseCreated, errors.New("db error"))

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// ストリークフリーズ: 使用状況ハンドラーテスト
// ============================================================

func TestStreakFreeze_GetStatus_Success(t *testing.T) {
	h, repo := setupStreakFreezeHandler()
	today, year, month := nowParts()

	repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).
		Return([]model.StreakFreeze{{UserID: 1, UsedDate: "2026-02-10"}}, nil)
	repo.On("HasFreezeOnDate", mock.Anything, uint(1), today).Return(false, nil)

	r := gin.New()
	r.GET("/streak-freezes/status", authMiddleware(1), h.GetStatus)

	req, _ := http.NewRequest("GET", "/streak-freezes/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["max_freezes"].(float64) != float64(model.MaxFreezesPerMonth) {
		t.Errorf("expected max_freezes=%d, got %v", model.MaxFreezesPerMonth, result["max_freezes"])
	}
	if result["used_freezes"].(float64) != 1 {
		t.Errorf("expected used_freezes=1, got %v", result["used_freezes"])
	}
	if result["remaining"].(float64) != float64(model.MaxFreezesPerMonth-1) {
		t.Errorf("expected remaining=%d, got %v", model.MaxFreezesPerMonth-1, result["remaining"])
	}
	if result["can_use_today"] != true {
		t.Errorf("expected can_use_today=true, got %v", result["can_use_today"])
	}
	repo.AssertExpectations(t)
}

func TestStreakFreeze_GetStatus_ServiceError(t *testing.T) {
	h, repo := setupStreakFreezeHandler()
	_, year, month := nowParts()

	repo.On("GetByUserIDAndMonth", mock.Anything, uint(1), year, month).
		Return([]model.StreakFreeze(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/streak-freezes/status", authMiddleware(1), h.GetStatus)

	req, _ := http.NewRequest("GET", "/streak-freezes/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}
