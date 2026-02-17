package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestBadgeService() (*BadgeService, *MockBadgeRepository, *MockNotificationService) {
	repo := new(MockBadgeRepository)
	notifSvc := new(MockNotificationService)
	svc := NewBadgeService(repo, notifSvc)
	return svc, repo, notifSvc
}

// === GetUserBadges ===

func TestGetUserBadges_Success(t *testing.T) {
	svc, repo, _ := newTestBadgeService()

	stats := &model.BadgeStats{
		TotalContributions: 100,
		CurrentStreak:      10,
		TotalPosts:         5,
		TotalLikesReceived: 20,
		FollowerCount:      15,
		FollowingCount:     8,
		QAAnswerCount:      3,
		CompletedGoals:     7,
	}
	repo.On("GetBadgeStats", uint(1)).Return(stats, nil)

	badges, err := svc.GetUserBadges(1)

	assert.NoError(t, err)
	assert.Len(t, badges, 18)
	repo.AssertExpectations(t)
}

func TestGetUserBadges_RepoError(t *testing.T) {
	svc, repo, _ := newTestBadgeService()

	repo.On("GetBadgeStats", uint(1)).Return(nil, errors.New("db error"))

	badges, err := svc.GetUserBadges(1)

	assert.Error(t, err)
	assert.Nil(t, badges)
	repo.AssertExpectations(t)
}

// === NotifyBadgeEarned ===

func TestNotifyBadgeEarned_Success(t *testing.T) {
	svc, _, notifSvc := newTestBadgeService()

	notifSvc.On("CreateNotification", mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 1 &&
			n.Type == model.NotificationTypeBadge &&
			n.ActorID == 1 &&
			n.BadgeID != nil && *n.BadgeID == "first-commit"
	})).Return(nil)

	err := svc.NotifyBadgeEarned(1, "first-commit")

	assert.NoError(t, err)
	notifSvc.AssertExpectations(t)
}

func TestNotifyBadgeEarned_Error(t *testing.T) {
	svc, _, notifSvc := newTestBadgeService()

	notifSvc.On("CreateNotification", mock.Anything).Return(errors.New("notification error"))

	err := svc.NotifyBadgeEarned(1, "first-commit")

	assert.Error(t, err)
	assert.Equal(t, "notification error", err.Error())
	notifSvc.AssertExpectations(t)
}

// === evaluateBadges: コントリビューションバッジ ===

func TestEvaluateBadges_ContributionBadges(t *testing.T) {
	tests := []struct {
		name          string
		contributions int
		expectedIDs   []string
	}{
		{"0件 → バッジなし", 0, nil},
		{"1件 → first-commit", 1, []string{"first-commit"}},
		{"50件 → +contributor", 50, []string{"first-commit", "contributor"}},
		{"200件 → +code-warrior", 200, []string{"first-commit", "contributor", "code-warrior"}},
		{"500件 → +commit-master", 500, []string{"first-commit", "contributor", "code-warrior", "commit-master"}},
		{"1000件 → +legend", 1000, []string{"first-commit", "contributor", "code-warrior", "commit-master", "legend"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{TotalContributions: tt.contributions}
			badges := evaluateBadges(stats)

			earnedContrib := filterByCategory(badges, "contribution")
			earnedIDs := extractEarnedIDs(earnedContrib)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: ストリークバッジ ===

func TestEvaluateBadges_StreakBadges(t *testing.T) {
	tests := []struct {
		name              string
		currentStreak     int
		learningLogStreak int
		expectedIDs       []string
	}{
		{"0日 → バッジなし", 0, 0, nil},
		{"GitHub 7日 → week-streak", 7, 0, []string{"week-streak"}},
		{"学習ログ 7日 → week-streak", 0, 7, []string{"week-streak"}},
		{"GitHub 30日 → +month-streak", 30, 0, []string{"week-streak", "month-streak"}},
		{"学習ログ 30日 → +month-streak", 5, 30, []string{"week-streak", "month-streak"}},
		{"両方の大きい方を使用", 3, 10, []string{"week-streak"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{
				CurrentStreak:     tt.currentStreak,
				LearningLogStreak: tt.learningLogStreak,
			}
			badges := evaluateBadges(stats)

			earnedStreak := filterByCategory(badges, "streak")
			earnedIDs := extractEarnedIDs(earnedStreak)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: 投稿バッジ ===

func TestEvaluateBadges_PostBadges(t *testing.T) {
	tests := []struct {
		name        string
		totalPosts  int
		expectedIDs []string
	}{
		{"0件 → バッジなし", 0, nil},
		{"1件 → first-post", 1, []string{"first-post"}},
		{"10件 → +blogger", 10, []string{"first-post", "blogger"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{TotalPosts: tt.totalPosts}
			badges := evaluateBadges(stats)

			earned := filterByCategory(badges, "post")
			earnedIDs := extractEarnedIDs(earned)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: エンゲージメントバッジ ===

func TestEvaluateBadges_EngagementBadges(t *testing.T) {
	tests := []struct {
		name        string
		likes       int
		expectedIDs []string
	}{
		{"0件 → バッジなし", 0, nil},
		{"10件 → liked", 10, []string{"liked"}},
		{"50件 → +popular", 50, []string{"liked", "popular"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{TotalLikesReceived: tt.likes}
			badges := evaluateBadges(stats)

			earned := filterByCategory(badges, "engagement")
			earnedIDs := extractEarnedIDs(earned)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: ソーシャルバッジ ===

func TestEvaluateBadges_SocialBadges(t *testing.T) {
	tests := []struct {
		name           string
		followingCount int
		followerCount  int
		expectedIDs    []string
	}{
		{"0/0 → バッジなし", 0, 0, nil},
		{"5フォロー → friendly", 5, 0, []string{"friendly"}},
		{"10フォロワー → influencer", 0, 10, []string{"influencer"}},
		{"50フォロワー → +star", 5, 50, []string{"friendly", "influencer", "star"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{
				FollowingCount: tt.followingCount,
				FollowerCount:  tt.followerCount,
			}
			badges := evaluateBadges(stats)

			earned := filterByCategory(badges, "social")
			earnedIDs := extractEarnedIDs(earned)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: Q&Aバッジ ===

func TestEvaluateBadges_QABadges(t *testing.T) {
	tests := []struct {
		name        string
		answers     int
		expectedIDs []string
	}{
		{"0件 → バッジなし", 0, nil},
		{"1件 → qa-first-answer", 1, []string{"qa-first-answer"}},
		{"10件 → +qa-helper", 10, []string{"qa-first-answer", "qa-helper"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{QAAnswerCount: tt.answers}
			badges := evaluateBadges(stats)

			earned := filterByCategory(badges, "qa")
			earnedIDs := extractEarnedIDs(earned)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: 目標達成バッジ ===

func TestEvaluateBadges_GoalBadges(t *testing.T) {
	tests := []struct {
		name        string
		goals       int
		expectedIDs []string
	}{
		{"0件 → バッジなし", 0, nil},
		{"5件 → goal-achiever", 5, []string{"goal-achiever"}},
		{"20件 → +goal-master", 20, []string{"goal-achiever", "goal-master"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &model.BadgeStats{CompletedGoals: tt.goals}
			badges := evaluateBadges(stats)

			earned := filterByCategory(badges, "goal")
			earnedIDs := extractEarnedIDs(earned)
			assert.Equal(t, tt.expectedIDs, earnedIDs)
		})
	}
}

// === evaluateBadges: 全バッジ獲得 ===

func TestEvaluateBadges_AllBadgesEarned(t *testing.T) {
	stats := &model.BadgeStats{
		TotalContributions: 1000,
		CurrentStreak:      30,
		LearningLogStreak:  30,
		TotalPosts:         10,
		TotalLikesReceived: 50,
		FollowerCount:      50,
		FollowingCount:     5,
		QAAnswerCount:      10,
		CompletedGoals:     20,
	}
	badges := evaluateBadges(stats)

	for _, b := range badges {
		assert.True(t, b.Earned, "バッジ %s が未獲得", b.ID)
	}
}

// === evaluateBadges: 全バッジ未獲得 ===

func TestEvaluateBadges_NoBadgesEarned(t *testing.T) {
	stats := &model.BadgeStats{}
	badges := evaluateBadges(stats)

	assert.Len(t, badges, 18)
	for _, b := range badges {
		assert.False(t, b.Earned, "バッジ %s が誤って獲得判定", b.ID)
	}
}

// === evaluateBadges: バッジの構造検証 ===

func TestEvaluateBadges_BadgeStructure(t *testing.T) {
	stats := &model.BadgeStats{}
	badges := evaluateBadges(stats)

	for _, b := range badges {
		assert.NotEmpty(t, b.ID, "IDが空")
		assert.NotEmpty(t, b.Name, "Nameが空")
		assert.NotEmpty(t, b.Description, "Descriptionが空")
		assert.NotEmpty(t, b.Category, "Categoryが空")
	}

	// カテゴリの種類を検証
	categories := make(map[string]bool)
	for _, b := range badges {
		categories[b.Category] = true
	}
	assert.Equal(t, 7, len(categories), "7カテゴリ存在すべき")
	assert.True(t, categories["contribution"])
	assert.True(t, categories["streak"])
	assert.True(t, categories["post"])
	assert.True(t, categories["engagement"])
	assert.True(t, categories["social"])
	assert.True(t, categories["qa"])
	assert.True(t, categories["goal"])
}

// === ヘルパー関数 ===

func filterByCategory(badges []BadgeResult, category string) []BadgeResult {
	var result []BadgeResult
	for _, b := range badges {
		if b.Category == category {
			result = append(result, b)
		}
	}
	return result
}

func extractEarnedIDs(badges []BadgeResult) []string {
	var ids []string
	for _, b := range badges {
		if b.Earned {
			ids = append(ids, b.ID)
		}
	}
	return ids
}
