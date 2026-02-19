package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestStudyCircleStatsService はテスト用のStudyCircleStatsServiceを生成する。
func newTestStudyCircleStatsService() (*StudyCircleStatsService, *MockStudyCircleStatsRepository) {
	repo := new(MockStudyCircleStatsRepository)
	svc := NewStudyCircleStatsService(repo)
	return svc, repo
}

// ============================================================
// GetCircleStats テスト
// ============================================================

func TestStudyCircleStatsService_GetCircleStats_Success(t *testing.T) {
	svc, repo := newTestStudyCircleStatsService()

	expected := &model.StudyCircleStats{
		MemberCount:    5,
		CheckinCount:   30,
		TotalSteps:     10,
		CompletedSteps: 7,
	}
	repo.On("GetCircleStats", uint(1)).Return(expected, nil)

	result, err := svc.GetCircleStats(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), result.MemberCount)
	assert.Equal(t, int64(30), result.CheckinCount)
	assert.Equal(t, int64(10), result.TotalSteps)
	assert.Equal(t, int64(7), result.CompletedSteps)
	repo.AssertExpectations(t)
}

func TestStudyCircleStatsService_GetCircleStats_EmptyCircle(t *testing.T) {
	svc, repo := newTestStudyCircleStatsService()

	empty := &model.StudyCircleStats{}
	repo.On("GetCircleStats", uint(2)).Return(empty, nil)

	result, err := svc.GetCircleStats(2)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.MemberCount)
	assert.Equal(t, int64(0), result.CheckinCount)
	repo.AssertExpectations(t)
}

func TestStudyCircleStatsService_GetCircleStats_InvalidCircleID(t *testing.T) {
	svc, _ := newTestStudyCircleStatsService()

	result, err := svc.GetCircleStats(0)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "circleIDは必須です")
}

func TestStudyCircleStatsService_GetCircleStats_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleStatsService()

	repo.On("GetCircleStats", uint(1)).Return((*model.StudyCircleStats)(nil), errors.New("db error"))

	result, err := svc.GetCircleStats(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
