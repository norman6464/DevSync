package repository

import (
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BadgeRepository はバッジ判定に必要な統計データを集計するリポジトリ実装。
// 複数テーブルからRaw SQLで各種カウントおよびストリークを算出する。
type BadgeRepository struct {
	db *gorm.DB
}

// NewBadgeRepository は新しいBadgeRepositoryインスタンスを生成する。
func NewBadgeRepository(db *gorm.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}

// GetBadgeStats はバッジ判定に必要な全統計をDBから集計する。
func (r *BadgeRepository) GetBadgeStats(userID uint) (*model.BadgeStats, error) {
	stats := &model.BadgeStats{}

	// GitHubコントリビューション総数
	r.db.Raw("SELECT COALESCE(SUM(count), 0) FROM git_hub_contributions WHERE user_id = ?", userID).Scan(&stats.TotalContributions)

	// GitHub連続コントリビューション日数
	streak, err := r.calculateGitHubStreak(userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	// 投稿総数
	r.db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)

	// 受け取ったいいね総数
	r.db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalLikesReceived)

	// フォロワー数
	r.db.Raw("SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID).Scan(&stats.FollowerCount)

	// フォロー中の数
	r.db.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID).Scan(&stats.FollowingCount)

	// Q&A回答数
	r.db.Raw("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&stats.QAAnswerCount)

	// 完了した学習目標数
	r.db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// 学習ログ連続記録日数
	logStreak, err := r.calculateLearningLogStreak(userID)
	if err == nil {
		stats.LearningLogStreak = logStreak
	}

	return stats, nil
}

// calculateGitHubStreak はGitHubコントリビューションの連続日数（ストリーク）を算出する。
func (r *BadgeRepository) calculateGitHubStreak(userID uint) (int, error) {
	type DateCount struct {
		Date  time.Time
		Count int
	}
	var contributions []DateCount
	err := r.db.Raw("SELECT date, count FROM git_hub_contributions WHERE user_id = ? AND count > 0 ORDER BY date DESC", userID).Scan(&contributions).Error
	if err != nil {
		return 0, err
	}
	if len(contributions) == 0 {
		return 0, nil
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].Date.After(contributions[j].Date)
	})

	streak := 0
	today := time.Now().UTC().Truncate(24 * time.Hour)

	for _, c := range contributions {
		cDate := c.Date.UTC().Truncate(24 * time.Hour)
		expectedDate := today.AddDate(0, 0, -streak)
		diff := expectedDate.Sub(cDate)

		if diff >= 0 && diff < 48*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
func (r *BadgeRepository) calculateLearningLogStreak(userID uint) (int, error) {
	var dates []struct {
		Date time.Time
	}
	err := r.db.Raw("SELECT DISTINCT DATE(created_at) as date FROM learning_logs WHERE user_id = ? ORDER BY date DESC", userID).Scan(&dates).Error
	if err != nil {
		return 0, err
	}
	if len(dates) == 0 {
		return 0, nil
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Date.After(dates[j].Date)
	})

	today := time.Now().UTC().Truncate(24 * time.Hour)
	firstDate := dates[0].Date.UTC().Truncate(24 * time.Hour)
	diffToToday := today.Sub(firstDate)

	if diffToToday >= 48*time.Hour {
		return 0, nil
	}

	streak := 1
	for i := 1; i < len(dates); i++ {
		prev := dates[i-1].Date.UTC().Truncate(24 * time.Hour)
		curr := dates[i].Date.UTC().Truncate(24 * time.Hour)
		diff := prev.Sub(curr)
		if diff >= 24*time.Hour && diff < 48*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}
