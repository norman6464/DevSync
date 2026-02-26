package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// StreakFreezeService テスト
// ============================================================

func TestUseFreeze_Success(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	// 今月の使用済みフリーズ: 0件
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{}, nil)
	// 今日のフリーズ未使用
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)
	// フリーズ作成成功
	mockRepo.On("Create", mock.MatchedBy(func(f *model.StreakFreeze) bool {
		return f.UserID == 1 && f.UsedDate == today && f.Month == int(now.Month()) && f.Year == now.Year()
	})).Return(nil)

	err := svc.UseFreeze(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUseFreeze_AlreadyUsedToday(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{{ID: 1, UserID: 1, UsedDate: today}}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(true, nil)

	err := svc.UseFreeze(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "今日は既にフリーズを使用済み")
}

func TestUseFreeze_MonthlyLimitReached(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	// 今月2件使用済み
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{
			{ID: 1, UserID: 1, UsedDate: "2026-02-01"},
			{ID: 2, UserID: 1, UsedDate: "2026-02-10"},
		}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)

	err := svc.UseFreeze(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "今月のフリーズ回数上限")
}

func TestGetFreezeStatus_Success(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{
			{ID: 1, UserID: 1, UsedDate: "2026-02-05"},
		}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)

	status, err := svc.GetFreezeStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, model.MaxFreezesPerMonth, status.MaxFreezes)
	assert.Equal(t, 1, status.UsedFreezes)
	assert.Equal(t, 1, status.Remaining)
	assert.False(t, status.TodayUsed)
	assert.True(t, status.CanUseToday)
	assert.Equal(t, []string{"2026-02-05"}, status.UsedDates)
}

func TestGetFreezeStatus_AllUsed(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{
			{ID: 1, UserID: 1, UsedDate: "2026-02-05"},
			{ID: 2, UserID: 1, UsedDate: "2026-02-15"},
		}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)

	status, err := svc.GetFreezeStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, status.Remaining)
	assert.False(t, status.CanUseToday)
}

func TestGetFreezeStatus_TodayAlreadyUsed(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")

	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{
			{ID: 1, UserID: 1, UsedDate: today},
		}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(true, nil)

	status, err := svc.GetFreezeStatus(1)
	assert.NoError(t, err)
	assert.True(t, status.TodayUsed)
	assert.False(t, status.CanUseToday)
}

// ============================================================
// GetFreezeDates テスト（カバレッジ0%の修正）
// ============================================================

func TestGetFreezeDates_Success(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	dates := []string{"2026-02-01", "2026-02-15"}
	mockRepo.On("GetFreezeDates", uint(1)).Return(dates, nil)

	result, err := svc.GetFreezeDates(1)
	assert.NoError(t, err)
	assert.Equal(t, dates, result)
	mockRepo.AssertExpectations(t)
}

func TestGetFreezeDates_Empty(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	mockRepo.On("GetFreezeDates", uint(1)).Return([]string{}, nil)

	result, err := svc.GetFreezeDates(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetFreezeDates_RepoError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	mockRepo.On("GetFreezeDates", uint(1)).Return([]string(nil), errors.New("db error"))

	result, err := svc.GetFreezeDates(1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// UseFreeze エラー伝播テスト
// ============================================================

func TestUseFreeze_HasFreezeOnDateError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	today := time.Now().Format("2006-01-02")
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, errors.New("db error"))

	err := svc.UseFreeze(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestUseFreeze_GetByUserIDAndMonthError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze(nil), errors.New("db error"))

	err := svc.UseFreeze(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestUseFreeze_CreateError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{}, nil)
	mockRepo.On("Create", mock.AnythingOfType("*model.StreakFreeze")).Return(errors.New("db error"))

	err := svc.UseFreeze(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// ============================================================
// GetFreezeStatus エラー伝播テスト
// ============================================================

func TestGetFreezeStatus_GetByUserIDAndMonthError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze(nil), errors.New("db error"))

	status, err := svc.GetFreezeStatus(1)
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestGetFreezeStatus_HasFreezeOnDateError(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, errors.New("db error"))

	status, err := svc.GetFreezeStatus(1)
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestGetFreezeStatus_NoFreezes(t *testing.T) {
	mockRepo := new(MockStreakFreezeRepository)
	svc := NewStreakFreezeService(mockRepo)

	now := time.Now()
	today := now.Format("2006-01-02")
	mockRepo.On("GetByUserIDAndMonth", uint(1), now.Year(), int(now.Month())).
		Return([]model.StreakFreeze{}, nil)
	mockRepo.On("HasFreezeOnDate", uint(1), today).Return(false, nil)

	status, err := svc.GetFreezeStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, model.MaxFreezesPerMonth, status.MaxFreezes)
	assert.Equal(t, 0, status.UsedFreezes)
	assert.Equal(t, model.MaxFreezesPerMonth, status.Remaining)
	assert.False(t, status.TodayUsed)
	assert.True(t, status.CanUseToday)
	assert.Empty(t, status.UsedDates)
}
