package service

import (
	"errors"
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

// ============================================================
// エラーハンドリングテスト
// ============================================================

func TestGetCurrentChallenge_DBError(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	dbErr := errors.New("connection refused")
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(nil, dbErr)

	result, err := svc.GetCurrentChallenge(uint(1))
	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	repo.AssertNotCalled(t, "Create")
	repo.AssertExpectations(t)
}

func TestGetCurrentChallenge_CreateError(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(nil, gorm.ErrRecordNotFound)
	repo.On("Create", mock.AnythingOfType("*model.WeeklyChallenge")).Return(errors.New("db write error"))

	result, err := svc.GetCurrentChallenge(uint(1))
	assert.Nil(t, result)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateProgress_UpdateError(t *testing.T) {
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
	repo.On("Update", mock.AnythingOfType("*model.WeeklyChallenge")).Return(errors.New("db update error"))

	result, err := svc.UpdateProgress(uint(1), 200)
	assert.Nil(t, result)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateProgress_AlreadyCompleted(t *testing.T) {
	svc, repo := newTestWeeklyChallengeService()

	year, week := time.Now().ISOWeek()
	completedAt := time.Now().Add(-1 * time.Hour)
	challenge := &model.WeeklyChallenge{
		ID:            1,
		UserID:        1,
		Year:          year,
		Week:          week,
		ChallengeType: model.ChallengeDurationTotal,
		TargetValue:   300,
		CurrentValue:  300,
		IsCompleted:   true,
		CompletedAt:   &completedAt,
	}
	repo.On("FindByUserAndWeek", uint(1), year, week).Return(challenge, nil)
	repo.On("Update", mock.MatchedBy(func(c *model.WeeklyChallenge) bool {
		// 既に完了済みの場合、CompletedAtが上書きされないことを検証
		return c.IsCompleted && c.CompletedAt.Equal(completedAt)
	})).Return(nil)

	result, err := svc.UpdateProgress(uint(1), 500)
	assert.NoError(t, err)
	assert.True(t, result.IsCompleted)
	assert.Equal(t, 500, result.CurrentValue)
	// CompletedAtは元の値のままであること
	assert.Equal(t, completedAt, *result.CompletedAt)
	repo.AssertExpectations(t)
}
