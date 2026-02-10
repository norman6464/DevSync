package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// === ヘルパー関数 ===

func newTestLevelService() (*LevelService, *MockLevelRepository, *MockNotificationRepository) {
	levelRepo := new(MockLevelRepository)
	notifRepo := new(MockNotificationRepository)
	notifService := NewNotificationService(notifRepo)
	svc := NewLevelService(levelRepo, notifService)
	return svc, levelRepo, notifRepo
}

// === XP計算テスト ===

func TestCalculateXPFromStats_AllZero(t *testing.T) {
	stats := &model.XPStats{}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 0, breakdown.Total)
	assert.Equal(t, 0, breakdown.LearningLogs)
	assert.Equal(t, 0, breakdown.Posts)
	assert.Equal(t, 0, breakdown.GitHub)
	assert.Equal(t, 0, breakdown.Goals)
	assert.Equal(t, 0, breakdown.Comments)
	assert.Equal(t, 0, breakdown.Likes)
	assert.Equal(t, 0, breakdown.StreakBonus)
}

func TestCalculateXPFromStats_LearningLogOnly(t *testing.T) {
	stats := &model.XPStats{
		LearningLogCount:         3,
		LearningLogTotalDuration: 120, // 120分
	}
	breakdown := CalculateXPFromStats(stats)
	// 3 * 10 + int(120 * 0.5) = 30 + 60 = 90
	assert.Equal(t, 90, breakdown.LearningLogs)
	assert.Equal(t, 90, breakdown.Total)
}

func TestCalculateXPFromStats_PostOnly(t *testing.T) {
	stats := &model.XPStats{PostCount: 5}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 150, breakdown.Posts) // 5 * 30
	assert.Equal(t, 150, breakdown.Total)
}

func TestCalculateXPFromStats_GitHubOnly(t *testing.T) {
	stats := &model.XPStats{GitHubContributionDays: 20}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 100, breakdown.GitHub) // 20 * 5
	assert.Equal(t, 100, breakdown.Total)
}

func TestCalculateXPFromStats_GoalOnly(t *testing.T) {
	stats := &model.XPStats{CompletedGoals: 3}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 150, breakdown.Goals) // 3 * 50
	assert.Equal(t, 150, breakdown.Total)
}

func TestCalculateXPFromStats_CommentOnly(t *testing.T) {
	stats := &model.XPStats{CommentCount: 10}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 50, breakdown.Comments) // 10 * 5
	assert.Equal(t, 50, breakdown.Total)
}

func TestCalculateXPFromStats_LikesOnly(t *testing.T) {
	stats := &model.XPStats{LikesReceived: 20}
	breakdown := CalculateXPFromStats(stats)
	assert.Equal(t, 60, breakdown.Likes) // 20 * 3
	assert.Equal(t, 60, breakdown.Total)
}

func TestCalculateXPFromStats_StreakBonus(t *testing.T) {
	// 7日ストリーク → 20 XP
	stats7 := &model.XPStats{CurrentStreak: 7}
	breakdown7 := CalculateXPFromStats(stats7)
	assert.Equal(t, 20, breakdown7.StreakBonus)

	// 14日ストリーク → 40 XP
	stats14 := &model.XPStats{CurrentStreak: 14}
	breakdown14 := CalculateXPFromStats(stats14)
	assert.Equal(t, 40, breakdown14.StreakBonus)
}

func TestCalculateXPFromStats_StreakBonusPartial(t *testing.T) {
	// 6日ストリーク → 0 XP（7日未満）
	stats6 := &model.XPStats{CurrentStreak: 6}
	breakdown6 := CalculateXPFromStats(stats6)
	assert.Equal(t, 0, breakdown6.StreakBonus)

	// 13日ストリーク → 20 XP（7の倍数は1回分）
	stats13 := &model.XPStats{CurrentStreak: 13}
	breakdown13 := CalculateXPFromStats(stats13)
	assert.Equal(t, 20, breakdown13.StreakBonus)
}

func TestCalculateXPFromStats_Combined(t *testing.T) {
	stats := &model.XPStats{
		LearningLogCount:         5,
		LearningLogTotalDuration: 200,
		PostCount:                3,
		GitHubContributionDays:   10,
		CompletedGoals:           2,
		CommentCount:             8,
		LikesReceived:            15,
		CurrentStreak:            14,
	}
	breakdown := CalculateXPFromStats(stats)
	expectedLogs := 5*10 + int(float64(200)*0.5)     // 50 + 100 = 150
	expectedPosts := 3 * 30                           // 90
	expectedGH := 10 * 5                              // 50
	expectedGoals := 2 * 50                           // 100
	expectedComments := 8 * 5                         // 40
	expectedLikes := 15 * 3                           // 45
	expectedStreak := (14 / 7) * 20                   // 40
	expectedTotal := expectedLogs + expectedPosts + expectedGH + expectedGoals + expectedComments + expectedLikes + expectedStreak

	assert.Equal(t, expectedLogs, breakdown.LearningLogs)
	assert.Equal(t, expectedPosts, breakdown.Posts)
	assert.Equal(t, expectedGH, breakdown.GitHub)
	assert.Equal(t, expectedGoals, breakdown.Goals)
	assert.Equal(t, expectedComments, breakdown.Comments)
	assert.Equal(t, expectedLikes, breakdown.Likes)
	assert.Equal(t, expectedStreak, breakdown.StreakBonus)
	assert.Equal(t, expectedTotal, breakdown.Total)
}

// === レベル計算テスト ===

func TestCalculateLevel_Zero(t *testing.T) {
	assert.Equal(t, 0, CalculateLevel(0))
}

func TestCalculateLevel_99(t *testing.T) {
	assert.Equal(t, 0, CalculateLevel(99))
}

func TestCalculateLevel_100(t *testing.T) {
	assert.Equal(t, 1, CalculateLevel(100))
}

func TestCalculateLevel_299(t *testing.T) {
	assert.Equal(t, 1, CalculateLevel(299))
}

func TestCalculateLevel_300(t *testing.T) {
	assert.Equal(t, 2, CalculateLevel(300))
}

func TestCalculateLevel_600(t *testing.T) {
	assert.Equal(t, 3, CalculateLevel(600))
}

func TestCalculateLevel_1000(t *testing.T) {
	assert.Equal(t, 4, CalculateLevel(1000))
}

func TestCalculateLevel_Large(t *testing.T) {
	// Lv10 = 5500 XP
	assert.Equal(t, 10, CalculateLevel(5500))
	assert.Equal(t, 9, CalculateLevel(5499))
}

// === XPForLevel テスト ===

func TestXPForLevel_Boundaries(t *testing.T) {
	assert.Equal(t, 0, XPForLevel(0))
	assert.Equal(t, 100, XPForLevel(1))
	assert.Equal(t, 300, XPForLevel(2))
	assert.Equal(t, 600, XPForLevel(3))
	assert.Equal(t, 1500, XPForLevel(5))
	assert.Equal(t, 5500, XPForLevel(10))
}

// === LevelInfo計算テスト ===

func TestCalculateLevelInfo_Zero(t *testing.T) {
	info := CalculateLevelInfo(0)
	assert.Equal(t, 0, info.Level)
	assert.Equal(t, 0, info.TotalXP)
	assert.Equal(t, 0, info.CurrentLevelXP)
	assert.Equal(t, 100, info.NextLevelXP)
	assert.Equal(t, 0, info.ProgressXP)
	assert.Equal(t, 0.0, info.ProgressPercent)
}

func TestCalculateLevelInfo_Mid(t *testing.T) {
	info := CalculateLevelInfo(150)
	assert.Equal(t, 1, info.Level)
	assert.Equal(t, 150, info.TotalXP)
	assert.Equal(t, 100, info.CurrentLevelXP)
	assert.Equal(t, 300, info.NextLevelXP)
	assert.Equal(t, 50, info.ProgressXP)
	assert.Equal(t, 25.0, info.ProgressPercent)
}

func TestCalculateLevelInfo_Exact(t *testing.T) {
	info := CalculateLevelInfo(300)
	assert.Equal(t, 2, info.Level)
	assert.Equal(t, 300, info.TotalXP)
	assert.Equal(t, 300, info.CurrentLevelXP)
	assert.Equal(t, 600, info.NextLevelXP)
	assert.Equal(t, 0, info.ProgressXP)
	assert.Equal(t, 0.0, info.ProgressPercent)
}

// === Service統合テスト ===

func TestGetLevelInfo_Success(t *testing.T) {
	svc, levelRepo, _ := newTestLevelService()
	stats := &model.XPStats{
		LearningLogCount:         5,
		LearningLogTotalDuration: 100,
		PostCount:                2,
	}
	levelRepo.On("GetXPStats", uint(1)).Return(stats, nil)

	info, err := svc.GetLevelInfo(1)
	assert.NoError(t, err)
	assert.NotNil(t, info)
	// 5*10 + 50 + 2*30 = 50+50+60 = 160 XP → Lv1
	assert.Equal(t, 1, info.Level)
	assert.Equal(t, 160, info.TotalXP)
	levelRepo.AssertExpectations(t)
}

func TestGetLevelInfo_RepoError(t *testing.T) {
	svc, levelRepo, _ := newTestLevelService()
	levelRepo.On("GetXPStats", uint(1)).Return((*model.XPStats)(nil), assert.AnError)

	info, err := svc.GetLevelInfo(1)
	assert.Error(t, err)
	assert.Nil(t, info)
	levelRepo.AssertExpectations(t)
}

func TestGetXPBreakdown_Success(t *testing.T) {
	svc, levelRepo, _ := newTestLevelService()
	stats := &model.XPStats{
		LearningLogCount:         3,
		LearningLogTotalDuration: 60,
		PostCount:                1,
		CommentCount:             4,
	}
	levelRepo.On("GetXPStats", uint(1)).Return(stats, nil)

	breakdown, err := svc.GetXPBreakdown(1)
	assert.NoError(t, err)
	assert.NotNil(t, breakdown)
	assert.Equal(t, 60, breakdown.LearningLogs) // 3*10 + 30
	assert.Equal(t, 30, breakdown.Posts)         // 1*30
	assert.Equal(t, 20, breakdown.Comments)      // 4*5
	levelRepo.AssertExpectations(t)
}

func TestGetXPBreakdown_EmptyStats(t *testing.T) {
	svc, levelRepo, _ := newTestLevelService()
	stats := &model.XPStats{}
	levelRepo.On("GetXPStats", uint(1)).Return(stats, nil)

	breakdown, err := svc.GetXPBreakdown(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, breakdown.Total)
	levelRepo.AssertExpectations(t)
}

func TestCheckAndNotifyLevelUp_LevelUp(t *testing.T) {
	svc, levelRepo, notifRepo := newTestLevelService()
	// previousXP = 90 (Lv0), 現在のXP = 160 (Lv1) → レベルアップ
	stats := &model.XPStats{
		LearningLogCount:         5,
		LearningLogTotalDuration: 100,
		PostCount:                2,
	}
	levelRepo.On("GetXPStats", uint(1)).Return(stats, nil)
	notifRepo.On("Create", mock.AnythingOfType("*model.Notification")).Return(nil)

	err := svc.CheckAndNotifyLevelUp(1, 90)
	assert.NoError(t, err)
	notifRepo.AssertCalled(t, "Create", mock.AnythingOfType("*model.Notification"))
}

func TestCheckAndNotifyLevelUp_NoLevelUp(t *testing.T) {
	svc, levelRepo, notifRepo := newTestLevelService()
	// previousXP = 150 (Lv1), 現在のXP = 160 (Lv1) → レベル変化なし
	stats := &model.XPStats{
		LearningLogCount:         5,
		LearningLogTotalDuration: 100,
		PostCount:                2,
	}
	levelRepo.On("GetXPStats", uint(1)).Return(stats, nil)

	err := svc.CheckAndNotifyLevelUp(1, 150)
	assert.NoError(t, err)
	notifRepo.AssertNotCalled(t, "Create", mock.Anything)
}
