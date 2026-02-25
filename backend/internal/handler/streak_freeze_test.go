package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockStreakFreezeService
// ============================================================

type MockStreakFreezeService struct{ mock.Mock }

func (m *MockStreakFreezeService) UseFreeze(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockStreakFreezeService) GetFreezeStatus(userID uint) (*model.StreakFreezeStatus, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StreakFreezeStatus), args.Error(1)
}

func (m *MockStreakFreezeService) GetFreezeDates(userID uint) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

// ============================================================
// StreakFreezeHandler テスト
// ============================================================

func TestStreakFreeze_UseFreeze_Success(t *testing.T) {
	mockSvc := new(MockStreakFreezeService)
	h := NewStreakFreezeHandler(mockSvc)

	mockSvc.On("UseFreeze", uint(1)).Return(nil)

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
	mockSvc.AssertExpectations(t)
}

func TestStreakFreeze_UseFreeze_AlreadyUsed(t *testing.T) {
	mockSvc := new(MockStreakFreezeService)
	h := NewStreakFreezeHandler(mockSvc)

	mockSvc.On("UseFreeze", uint(1)).Return(domain.ErrConflict)

	r := gin.New()
	r.POST("/streak-freezes", authMiddleware(1), h.UseFreeze)

	req, _ := http.NewRequest("POST", "/streak-freezes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusConflict)
}

func TestStreakFreeze_GetStatus_Success(t *testing.T) {
	mockSvc := new(MockStreakFreezeService)
	h := NewStreakFreezeHandler(mockSvc)

	mockSvc.On("GetFreezeStatus", uint(1)).Return(&model.StreakFreezeStatus{
		MaxFreezes:  2,
		UsedFreezes: 1,
		Remaining:   1,
		UsedDates:   []string{"2026-02-10"},
		TodayUsed:   false,
		CanUseToday: true,
	}, nil)

	r := gin.New()
	r.GET("/streak-freezes/status", authMiddleware(1), h.GetStatus)

	req, _ := http.NewRequest("GET", "/streak-freezes/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["max_freezes"].(float64) != 2 {
		t.Errorf("expected max_freezes=2, got %v", result["max_freezes"])
	}
	if result["remaining"].(float64) != 1 {
		t.Errorf("expected remaining=1, got %v", result["remaining"])
	}
}
