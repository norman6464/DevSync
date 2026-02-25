package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWeeklyChallengeService は WeeklyChallengeServiceInterface のモック実装。
type MockWeeklyChallengeService struct{ mock.Mock }

func (m *MockWeeklyChallengeService) GetCurrentChallenge(userID uint) (*model.WeeklyChallenge, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WeeklyChallenge), args.Error(1)
}

func (m *MockWeeklyChallengeService) UpdateProgress(userID uint, value int) (*model.WeeklyChallenge, error) {
	args := m.Called(userID, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WeeklyChallenge), args.Error(1)
}

func setupWeeklyChallengeHandler() (*WeeklyChallengeHandler, *MockWeeklyChallengeService) {
	svc := new(MockWeeklyChallengeService)
	h := NewWeeklyChallengeHandler(svc)
	return h, svc
}

// ============================================================
// ウィークリーチャレンジ: 取得ハンドラーテスト
// ============================================================

func TestWeeklyChallenge_GetCurrent_Success(t *testing.T) {
	h, svc := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	challenge := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
	}
	svc.On("GetCurrentChallenge", uint(1)).Return(challenge, nil)

	r := gin.New()
	r.GET("/weekly-challenges/current", authMiddleware(1), h.GetCurrent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weekly-challenges/current", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

// ============================================================
// ウィークリーチャレンジ: 進捗更新ハンドラーテスト
// ============================================================

func TestWeeklyChallenge_UpdateProgress_Success(t *testing.T) {
	h, svc := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	challenge := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
		CurrentValue:  200,
	}
	svc.On("UpdateProgress", uint(1), 200).Return(challenge, nil)

	r := gin.New()
	r.PUT("/weekly-challenges/progress", authMiddleware(1), h.UpdateProgress)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/weekly-challenges/progress", jsonBody(map[string]interface{}{
		"value": 200,
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
