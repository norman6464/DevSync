package service

import (
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BadgeStats はバッジ判定に必要な集計統計を保持する構造体。
type BadgeStats struct {
	TotalContributions int // GitHub総コントリビューション数
	CurrentStreak      int // GitHub連続コントリビューション日数
	LearningLogStreak  int // 学習ログ連続記録日数
	TotalPosts         int // 投稿総数
	TotalLikesReceived int // 受け取ったいいね総数
	FollowerCount      int // フォロワー数
	FollowingCount     int // フォロー中の数
	QAAnswerCount      int // Q&A回答数
	CompletedGoals     int // 完了した学習目標数
}

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
	db                  *gorm.DB             // 統計集計用DBコネクション
	notificationService *NotificationService // バッジ獲得通知用
}

// NewBadgeService は新しいBadgeServiceインスタンスを生成する。
func NewBadgeService(db *gorm.DB, notificationService *NotificationService) *BadgeService {
	return &BadgeService{db: db, notificationService: notificationService}
}

// GetUserBadges は指定ユーザーの全バッジと獲得状況を返す。
// 統計を集計した後、全18バッジを閾値で評価する。
func (s *BadgeService) GetUserBadges(userID uint) ([]BadgeResult, error) {
	stats, err := s.getBadgeStats(userID)
	if err != nil {
		return nil, err
	}
	return s.evaluateBadges(stats), nil
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

// getBadgeStats はバッジ判定に必要な全統計をDBから集計する。
// 複数テーブルからRaw SQLで各種カウントを取得する。
func (s *BadgeService) getBadgeStats(userID uint) (*BadgeStats, error) {
	stats := &BadgeStats{}

	// GitHubコントリビューション総数
	s.db.Raw("SELECT COALESCE(SUM(count), 0) FROM git_hub_contributions WHERE user_id = ?", userID).Scan(&stats.TotalContributions)

	// GitHub連続コントリビューション日数
	streak, err := s.calculateStreak(userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	// 投稿総数
	s.db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)

	// 受け取ったいいね総数
	s.db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalLikesReceived)

	// フォロワー数
	s.db.Raw("SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID).Scan(&stats.FollowerCount)

	// フォロー中の数
	s.db.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID).Scan(&stats.FollowingCount)

	// Q&A回答数
	s.db.Raw("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&stats.QAAnswerCount)

	// 完了した学習目標数
	s.db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// 学習ログ連続記録日数
	logStreak, err := s.calculateLearningLogStreak(userID)
	if err == nil {
		stats.LearningLogStreak = logStreak
	}

	return stats, nil
}

// calculateStreak はGitHubコントリビューションの連続日数（ストリーク）を算出する。
// 直近の日付から遡り、連続してコントリビューションがある日数をカウントする。
func (s *BadgeService) calculateStreak(userID uint) (int, error) {
	type DateCount struct {
		Date  time.Time
		Count int
	}
	var contributions []DateCount
	err := s.db.Raw("SELECT date, count FROM git_hub_contributions WHERE user_id = ? AND count > 0 ORDER BY date DESC", userID).Scan(&contributions).Error
	if err != nil {
		return 0, err
	}
	if len(contributions) == 0 {
		return 0, nil
	}

	// 日付の降順でソート
	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].Date.After(contributions[j].Date)
	})

	streak := 0
	today := time.Now().UTC().Truncate(24 * time.Hour)

	for _, c := range contributions {
		cDate := c.Date.UTC().Truncate(24 * time.Hour)
		expectedDate := today.AddDate(0, 0, -streak)
		diff := expectedDate.Sub(cDate)

		// 48時間以内の差異を許容（タイムゾーンの考慮）
		if diff >= 0 && diff < 48*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
// 直近の記録日から遡り、1日ごとに連続している日数をカウントする。
func (s *BadgeService) calculateLearningLogStreak(userID uint) (int, error) {
	var dates []struct {
		Date time.Time
	}
	err := s.db.Raw("SELECT DISTINCT DATE(created_at) as date FROM learning_logs WHERE user_id = ? ORDER BY date DESC", userID).Scan(&dates).Error
	if err != nil {
		return 0, err
	}
	if len(dates) == 0 {
		return 0, nil
	}

	// 日付の降順でソート
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Date.After(dates[j].Date)
	})

	today := time.Now().UTC().Truncate(24 * time.Hour)
	firstDate := dates[0].Date.UTC().Truncate(24 * time.Hour)
	diffToToday := today.Sub(firstDate)

	// 最新の記録が48時間以上前ならストリークなし
	if diffToToday >= 48*time.Hour {
		return 0, nil
	}

	streak := 1
	for i := 1; i < len(dates); i++ {
		prev := dates[i-1].Date.UTC().Truncate(24 * time.Hour)
		curr := dates[i].Date.UTC().Truncate(24 * time.Hour)
		diff := prev.Sub(curr)
		// 1日差（24〜48時間）なら連続とみなす
		if diff >= 24*time.Hour && diff < 48*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}

// evaluateBadges は統計データに基づいて全18バッジの獲得状況を評価する。
// 7カテゴリ（contribution, streak, post, engagement, social, qa, goal）のバッジを返す。
func (s *BadgeService) evaluateBadges(stats *BadgeStats) []BadgeResult {
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
