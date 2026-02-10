package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LevelRepository はレベルシステムの統計集計を提供するリポジトリ実装。
// 複数テーブルからRaw SQLでXP計算に必要なデータを集計する。
type LevelRepository struct {
	db *gorm.DB
}

// NewLevelRepository は新しいLevelRepositoryインスタンスを生成する。
func NewLevelRepository(db *gorm.DB) *LevelRepository {
	return &LevelRepository{db: db}
}

// GetXPStats はXP計算に必要な全統計をDBから集計する。
func (r *LevelRepository) GetXPStats(userID uint) (*model.XPStats, error) {
	stats := &model.XPStats{}

	// 学習ログ: 件数と合計時間
	r.db.Raw("SELECT COUNT(*) FROM learning_logs WHERE user_id = ?", userID).Scan(&stats.LearningLogCount)
	r.db.Raw("SELECT COALESCE(SUM(duration), 0) FROM learning_logs WHERE user_id = ?", userID).Scan(&stats.LearningLogTotalDuration)

	// 投稿数
	r.db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.PostCount)

	// GitHubコントリビューション日数（count > 0 の日数）
	r.db.Raw("SELECT COUNT(DISTINCT date) FROM git_hub_contributions WHERE user_id = ? AND count > 0", userID).Scan(&stats.GitHubContributionDays)

	// 完了した学習目標数
	r.db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// コメント数
	r.db.Raw("SELECT COUNT(*) FROM comments WHERE user_id = ?", userID).Scan(&stats.CommentCount)

	// 受け取ったいいね総数（自分の投稿のlike_countの合計）
	r.db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.LikesReceived)

	// 学習ログストリーク（連続記録日数）
	streak, err := r.calculateLearningLogStreak(userID)
	if err == nil {
		stats.CurrentStreak = streak
	}

	return stats, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
// BadgeServiceのcalculateLearningLogStreakと同じロジック。
func (r *LevelRepository) calculateLearningLogStreak(userID uint) (int, error) {
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

	today := time.Now().UTC().Truncate(24 * time.Hour)
	streak := 0

	for i, d := range dates {
		dDate := d.Date.UTC().Truncate(24 * time.Hour)
		if i == 0 {
			// 最新の日付が今日か昨日でなければストリークは0
			diffToToday := today.Sub(dDate)
			if diffToToday >= 48*time.Hour {
				return 0, nil
			}
			streak = 1
			continue
		}

		prevDate := dates[i-1].Date.UTC().Truncate(24 * time.Hour)
		diff := prevDate.Sub(dDate)
		if diff >= 24*time.Hour && diff < 48*time.Hour {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}
