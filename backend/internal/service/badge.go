package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BadgeResult は個別バッジの獲得状況を表す。
type BadgeResult struct {
	ID          string `json:"id"`          // バッジ識別子
	Name        string `json:"name"`        // バッジ名（i18nキー）
	Description string `json:"description"` // バッジ説明（i18nキー）
	Category    string `json:"category"`    // バッジカテゴリ
	Earned      bool   `json:"earned"`      // 獲得済みフラグ
}

// BadgeService はバッジ判定のビジネスロジックを提供する。
// 各種統計を集計し、閾値ベースでバッジの獲得状況を評価する。
type BadgeService struct {
	repo                repository.BadgeRepositoryInterface
	notificationService NotificationServiceInterface
}

// NewBadgeService は新しいBadgeServiceインスタンスを生成する。
func NewBadgeService(repo repository.BadgeRepositoryInterface, notificationService NotificationServiceInterface) *BadgeService {
	return &BadgeService{repo: repo, notificationService: notificationService}
}

// GetUserBadges は指定ユーザーの全バッジと獲得状況を返す。
// 統計を集計した後、全18バッジを閾値で評価する。
func (s *BadgeService) GetUserBadges(userID uint) ([]BadgeResult, error) {
	stats, err := s.repo.GetBadgeStats(userID)
	if err != nil {
		return nil, err
	}
	return evaluateBadges(stats), nil
}

// NotifyBadgeEarned は新しいバッジ獲得の通知を作成する。
func (s *BadgeService) NotifyBadgeEarned(userID uint, badgeID string) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBadge,
		ActorID: userID,
		BadgeID: &badgeID,
	}
	return s.notificationService.CreateNotification(notification)
}

// evaluateBadges は統計データに基づいて全18バッジの獲得状況を評価する。
func evaluateBadges(stats *model.BadgeStats) []BadgeResult {
	// GitHubストリークと学習ログストリークの大きい方を使用
	combinedStreak := stats.CurrentStreak
	if stats.LearningLogStreak > combinedStreak {
		combinedStreak = stats.LearningLogStreak
	}

	return []BadgeResult{
		// コントリビューションバッジ（1, 50, 200, 500, 1000回）
		{ID: "first-commit", Name: "badges.firstCommit", Description: "badges.firstCommitDesc", Category: "contribution", Earned: stats.TotalContributions >= 1},
		{ID: "contributor", Name: "badges.contributor", Description: "badges.contributorDesc", Category: "contribution", Earned: stats.TotalContributions >= 50},
		{ID: "code-warrior", Name: "badges.codeWarrior", Description: "badges.codeWarriorDesc", Category: "contribution", Earned: stats.TotalContributions >= 200},
		{ID: "commit-master", Name: "badges.commitMaster", Description: "badges.commitMasterDesc", Category: "contribution", Earned: stats.TotalContributions >= 500},
		{ID: "legend", Name: "badges.legend", Description: "badges.legendDesc", Category: "contribution", Earned: stats.TotalContributions >= 1000},

		// ストリークバッジ（7日, 30日連続）
		{ID: "week-streak", Name: "badges.weekStreak", Description: "badges.weekStreakDesc", Category: "streak", Earned: combinedStreak >= 7},
		{ID: "month-streak", Name: "badges.monthStreak", Description: "badges.monthStreakDesc", Category: "streak", Earned: combinedStreak >= 30},

		// 投稿バッジ（1件, 10件）
		{ID: "first-post", Name: "badges.firstPost", Description: "badges.firstPostDesc", Category: "post", Earned: stats.TotalPosts >= 1},
		{ID: "blogger", Name: "badges.blogger", Description: "badges.bloggerDesc", Category: "post", Earned: stats.TotalPosts >= 10},

		// エンゲージメントバッジ（いいね10, 50件）
		{ID: "liked", Name: "badges.liked", Description: "badges.likedDesc", Category: "engagement", Earned: stats.TotalLikesReceived >= 10},
		{ID: "popular", Name: "badges.popular", Description: "badges.popularDesc", Category: "engagement", Earned: stats.TotalLikesReceived >= 50},

		// ソーシャルバッジ（フォロー5人, フォロワー10人, 50人）
		{ID: "friendly", Name: "badges.friendly", Description: "badges.friendlyDesc", Category: "social", Earned: stats.FollowingCount >= 5},
		{ID: "influencer", Name: "badges.influencer", Description: "badges.influencerDesc", Category: "social", Earned: stats.FollowerCount >= 10},
		{ID: "star", Name: "badges.star", Description: "badges.starDesc", Category: "social", Earned: stats.FollowerCount >= 50},

		// Q&Aバッジ（回答1件, 10件）
		{ID: "qa-first-answer", Name: "badges.qaFirstAnswer", Description: "badges.qaFirstAnswerDesc", Category: "qa", Earned: stats.QAAnswerCount >= 1},
		{ID: "qa-helper", Name: "badges.qaHelper", Description: "badges.qaHelperDesc", Category: "qa", Earned: stats.QAAnswerCount >= 10},

		// 目標達成バッジ（5件, 20件完了）
		{ID: "goal-achiever", Name: "badges.goalAchiever", Description: "badges.goalAchieverDesc", Category: "goal", Earned: stats.CompletedGoals >= 5},
		{ID: "goal-master", Name: "badges.goalMaster", Description: "badges.goalMasterDesc", Category: "goal", Earned: stats.CompletedGoals >= 20},
	}
}
