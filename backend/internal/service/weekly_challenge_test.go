package service

import (
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func newTestWeeklyChallengeService() (*WeeklyChallengeService, *MockWeeklyChallengeRepository) {
	repo := new(MockWeeklyChallengeRepository)
	svc := NewWeeklyChallengeService(repo)
	return svc, repo
}

// ============================================================
// 今週のチャレンジ取得テスト
// ============================================================

func TestGetCurrentChallenge_ExistingChallenge(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	existing := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
		CurrentValue:  120,
	}
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(existing, nil)

	result, err := svc.GetCurrentChallenge(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, existing, result)
	repo.AssertExpectations(t)
}

func TestGetCurrentChallenge_CreatesNew(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(nil, gorm.ErrRecordNotFound)
	repo.On("Create", mock.AnythingOfType("*model.WeeklyChallenge")).Return(nil)

	result, err := svc.GetCurrentChallenge(uint(1))
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, year, result.Year)
	assert.Equal(t, week, result.Week)
	assert.True(t, result.TargetValue > 0)
	repo.AssertExpectations(t)
}

// ============================================================
// 進捗更新テスト
// ============================================================

func TestUpdateProgress_Success(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	challenge := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
		CurrentValue:  100,
	}
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(challenge, nil)
	repo.On("Update", mock.AnythingOfType("*model.WeeklyChallenge")).Return(nil)

	result, err := svc.UpdateProgress(uint(1), 200)
	assert.NoError(t, err)
	assert.Equal(t, 200, result.CurrentValue)
	repo.AssertExpectations(t)
}

func TestUpdateProgress_Completes(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	challenge := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
		CurrentValue:  100,
	}
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(challenge, nil)
	repo.On("Update", mock.MatchedBy(func(c *model.WeeklyChallenge) bool {
		return c.IsCompleted && c.CompletedAt != nil
	})).Return(nil)

	result, err := svc.UpdateProgress(uint(1), 350)
	assert.NoError(t, err)
	assert.True(t, result.IsCompleted)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestUpdateProgress_NoChallengeFound(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(nil, gorm.ErrRecordNotFound)

	_, err := svc.UpdateProgress(uint(1), 100)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
