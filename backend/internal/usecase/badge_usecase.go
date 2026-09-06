package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// validBadgeIDs はEvaluateBadgesが返しうるバッジIDの集合。統計値には依存しないため
// ゼロ値のBadgeStatsで評価して抽出する。クライアントから届くbadge_idがこの集合に
// 含まれない場合は不正な値として拒否する（無検証保存の防止）。
var validBadgeIDs = func() map[string]bool {
	ids := make(map[string]bool)
	for _, badge := range EvaluateBadges(&model.BadgeStats{}) {
		ids[badge.ID] = true
	}
	return ids
}()

// GetUserBadgesUseCase は指定ユーザーの全バッジと獲得状況を返す。
type GetUserBadgesUseCase struct {
	stats repository.BadgeStatsReader
}

// NewGetUserBadgesUseCase は GetUserBadgesUseCase を生成する。
func NewGetUserBadgesUseCase(stats repository.BadgeStatsReader) *GetUserBadgesUseCase {
	return &GetUserBadgesUseCase{stats: stats}
}

// Execute は統計を集計し、全バッジの獲得状況を評価して返す。
func (uc *GetUserBadgesUseCase) Execute(ctx context.Context, userID uint) ([]model.BadgeResult, error) {
	stats, err := uc.stats.GetBadgeStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return EvaluateBadges(stats), nil
}

// NotifyBadgeEarnedUseCase はバッジ獲得の通知を作成する。
type NotifyBadgeEarnedUseCase struct {
	notifications repository.NotificationCreator
}

// NewNotifyBadgeEarnedUseCase は NotifyBadgeEarnedUseCase を生成する。
func NewNotifyBadgeEarnedUseCase(notifications repository.NotificationCreator) *NotifyBadgeEarnedUseCase {
	return &NotifyBadgeEarnedUseCase{notifications: notifications}
}

// Execute は本人宛のバッジ獲得通知を作成する。badgeIDはクライアントから届く値のため、
// EvaluateBadgesが定義する既知のバッジID集合に含まれるものだけを受け付ける。
func (uc *NotifyBadgeEarnedUseCase) Execute(ctx context.Context, userID uint, badgeID string) error {
	if !validBadgeIDs[badgeID] {
		return domain.NewError(domain.ErrCodeBadRequest, "不明なバッジIDです", nil)
	}
	return uc.notifications.Create(ctx, &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBadge,
		ActorID: userID,
		BadgeID: &badgeID,
	})
}

// EvaluateBadges は統計データに基づいて全 18 バッジの獲得状況を評価する。
func EvaluateBadges(stats *model.BadgeStats) []model.BadgeResult {
	// GitHubストリークと学習ログストリークの大きい方を使用
	combinedStreak := stats.CurrentStreak
	if stats.LearningLogStreak > combinedStreak {
		combinedStreak = stats.LearningLogStreak
	}

	return []model.BadgeResult{
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
