package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockStudyCircleStatsRepo は usecase/repository.StudyCircleStatsRepository のモック（ctx 付き）。
type mockStudyCircleStatsRepo struct{ mock.Mock }

func (m *mockStudyCircleStatsRepo) GetCircleStats(ctx context.Context, circleID uint) (*model.StudyCircleStats, error) {
	args := m.Called(ctx, circleID)
	s, _ := args.Get(0).(*model.StudyCircleStats)
	return s, args.Error(1)
}

func setupStudyCircleStatsHandler() (*StudyCircleStatsHandler, *mockStudyCircleStatsRepo) {
	stats := new(mockStudyCircleStatsRepo)
	h := NewStudyCircleStatsHandler(usecase.NewGetStudyCircleStatsUseCase(stats))
	return h, stats
}

func TestStudyCircleStats_GetStats_Success(t *testing.T) {
	h, stats := setupStudyCircleStatsHandler()
	stats.On("GetCircleStats", mock.Anything, uint(5)).Return(
		&model.StudyCircleStats{MemberCount: 3, CheckinCount: 5, TotalSteps: 4, CompletedSteps: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/circles/5/stats", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["member_count"].(float64) != 3 {
		t.Errorf("expected member_count 3, got %v", data["member_count"])
	}
	stats.AssertExpectations(t)
}

func TestStudyCircleStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupStudyCircleStatsHandler()

	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/circles/abc/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestStudyCircleStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupStudyCircleStatsHandler()
	stats.On("GetCircleStats", mock.Anything, uint(5)).Return((*model.StudyCircleStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/circles/:id/stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/circles/5/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
