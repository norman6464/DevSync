package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockWeeklyChallengeRepo は usecase/repository.WeeklyChallengeRepository のモック（ctx 付き）。
// FindByUserAndWeek は「不在」を (nil, nil) で表す契約。
type mockWeeklyChallengeRepo struct{ mock.Mock }

func (m *mockWeeklyChallengeRepo) Create(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return m.Called(ctx, challenge).Error(0)
}

func (m *mockWeeklyChallengeRepo) FindByUserAndWeek(ctx context.Context, userID uint, year, week int) (*model.WeeklyChallenge, error) {
	args := m.Called(ctx, userID, year, week)
	c, _ := args.Get(0).(*model.WeeklyChallenge)
	return c, args.Error(1)
}

func (m *mockWeeklyChallengeRepo) Update(ctx context.Context, challenge *model.WeeklyChallenge) error {
	return m.Called(ctx, challenge).Error(0)
}

// setupWeeklyChallengeHandler は本物の usecase と port モックで WeeklyChallengeHandler を組む。
func setupWeeklyChallengeHandler() (*WeeklyChallengeHandler, *mockWeeklyChallengeRepo) {
	repo := new(mockWeeklyChallengeRepo)
	h := NewWeeklyChallengeHandler(
		usecase.NewGetCurrentWeeklyChallengeUseCase(repo),
		usecase.NewUpdateWeeklyChallengeProgressUseCase(repo),
	)
	return h, repo
}

// currentWeekChallenge は今週分のチャレンジを返す。
func currentWeekChallenge() *model.WeeklyChallenge {
	year, week := time.Now().ISOWeek()
	return &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
	}
}

// ============================================================
// ウィークリーチャレンジ: 取得ハンドラーテスト
// ============================================================

func TestWeeklyChallenge_GetCurrent_Success(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(currentWeekChallenge(), nil)

	r := gin.New()
	r.GET("/weekly-challenges/current", authMiddleware(1), h.GetCurrent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weekly-challenges/current", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertNotCalled(t, "Create")
	repo.AssertExpectations(t)
}

// 今週分が未登録なら自動生成して 200 を返す。
// 旧テストは service モックが domain.ErrNotFound を返す 404 ケースを持っていたが、
// 実際の実装は未登録時に自動生成するため到達しない経路だった。
func TestWeeklyChallenge_GetCurrent_GeneratesWhenAbsent(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
		Return((*model.WeeklyChallenge)(nil), nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(c *model.WeeklyChallenge) bool {
		return c.UserID == 1 && c.Year == year && c.Week == week && c.TargetValue > 0
	})).Return(nil)

	r := gin.New()
	r.GET("/weekly-challenges/current", authMiddleware(1), h.GetCurrent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weekly-challenges/current", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

// ============================================================
// ウィークリーチャレンジ: 進捗更新ハンドラーテスト
// ============================================================

func TestWeeklyChallenge_UpdateProgress_Success(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).Return(currentWeekChallenge(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.WeeklyChallenge")).Return(nil)

	r := gin.New()
	r.PUT("/weekly-challenges/progress", authMiddleware(1), h.UpdateProgress)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/weekly-challenges/progress", jsonBody(map[string]interface{}{
		"value": 200,
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

// ============================================================
// ウィークリーチャレンジ: エラーパステスト
// ============================================================

func TestWeeklyChallenge_GetCurrent_ServiceError(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
		Return((*model.WeeklyChallenge)(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/weekly-challenges/current", authMiddleware(1), h.GetCurrent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weekly-challenges/current", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertNotCalled(t, "Create")
}

// 今週分が未登録のまま進捗更新を呼ぶと 500（移行前と同じ扱い）。
func TestWeeklyChallenge_UpdateProgress_AbsentChallenge(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
		Return((*model.WeeklyChallenge)(nil), nil)

	r := gin.New()
	r.PUT("/weekly-challenges/progress", authMiddleware(1), h.UpdateProgress)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/weekly-challenges/progress", jsonBody(map[string]interface{}{
		"value": 100,
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertNotCalled(t, "Update")
}

func TestWeeklyChallenge_UpdateProgress_ServiceError(t *testing.T) {
	h, repo := setupWeeklyChallengeHandler()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", mock.Anything, uint(1), year, week).
		Return((*model.WeeklyChallenge)(nil), errors.New("db error"))

	r := gin.New()
	r.PUT("/weekly-challenges/progress", authMiddleware(1), h.UpdateProgress)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/weekly-challenges/progress", jsonBody(map[string]interface{}{
		"value": 100,
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWeeklyChallenge_UpdateProgress_InvalidJSON(t *testing.T) {
	h, _ := setupWeeklyChallengeHandler()

	r := gin.New()
	r.PUT("/weekly-challenges/progress", authMiddleware(1), h.UpdateProgress)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/weekly-challenges/progress", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
