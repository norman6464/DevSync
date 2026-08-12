package persistence

import (
	"context"
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// badgeRepository は [repository.BadgeStatsReader] の GORM 実装。
// 複数テーブルから Raw SQL で各種カウントおよびストリークを算出する。
type badgeRepository struct {
	db *gorm.DB
}

// NewBadgeRepository は BadgeStatsReader の GORM 実装を返す。
func NewBadgeRepository(db *gorm.DB) repository.BadgeStatsReader {
	return &badgeRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BadgeStatsReader = (*badgeRepository)(nil)

// GetBadgeStats はバッジ判定に必要な全統計を DB から集計する。
// 個々の集計クエリのエラーは移行前と同じく無視し、ストリーク算出の失敗だけを error として返す。
func (r *badgeRepository) GetBadgeStats(ctx context.Context, userID uint) (*model.BadgeStats, error) {
	db := r.db.WithContext(ctx)
	stats := &model.BadgeStats{}

	// GitHubコントリビューション総数
	db.Raw("SELECT COALESCE(SUM(count), 0) FROM git_hub_contributions WHERE user_id = ?", userID).Scan(&stats.TotalContributions)

	// GitHub連続コントリビューション日数
	streak, err := r.calculateGitHubStreak(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	// 投稿総数
	db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)

	// 受け取ったいいね総数
	db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalLikesReceived)

	// フォロワー数
	db.Raw("SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID).Scan(&stats.FollowerCount)

	// フォロー中の数
	db.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID).Scan(&stats.FollowingCount)

	// Q&A回答数
	db.Raw("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&stats.QAAnswerCount)

	// 完了した学習目標数
	db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// 学習ログ連続記録日数（失敗しても 0 のまま続行する）
	if logStreak, err := r.calculateLearningLogStreak(ctx, userID); err == nil {
		stats.LearningLogStreak = logStreak
	}

	return stats, nil
}

// calculateGitHubStreak は GitHub コントリビューションの連続日数を算出する。
func (r *badgeRepository) calculateGitHubStreak(ctx context.Context, userID uint) (int, error) {
	type dateCount struct {
		Date  time.Time
		Count int
	}
	var contributions []dateCount
	err := r.db.WithContext(ctx).
		Raw("SELECT date, count FROM git_hub_contributions WHERE user_id = ? AND count > 0 ORDER BY date DESC", userID).
		Scan(&contributions).Error
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
func (r *badgeRepository) calculateLearningLogStreak(ctx context.Context, userID uint) (int, error) {
	var dates []struct {
		Date time.Time
	}
	err := r.db.WithContext(ctx).
		Raw("SELECT DISTINCT DATE(created_at) as date FROM learning_logs WHERE user_id = ? ORDER BY date DESC", userID).
		Scan(&dates).Error
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
	if today.Sub(firstDate) >= 48*time.Hour {
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
