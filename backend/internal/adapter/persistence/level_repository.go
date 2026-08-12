package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// levelRepository は [repository.XPStatsReader] の GORM 実装。
// 複数テーブルから Raw SQL で XP 計算に必要なデータを集計する。
type levelRepository struct {
	db *gorm.DB
}

// NewLevelRepository は XPStatsReader の GORM 実装を返す。
func NewLevelRepository(db *gorm.DB) repository.XPStatsReader {
	return &levelRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.XPStatsReader = (*levelRepository)(nil)

// GetXPStats は XP 計算に必要な全統計を DB から集計する。
// 個々の集計クエリのエラーは移行前と同じく無視する。
func (r *levelRepository) GetXPStats(ctx context.Context, userID uint) (*model.XPStats, error) {
	db := r.db.WithContext(ctx)
	stats := &model.XPStats{}

	// 学習ログ: 件数と合計時間
	db.Raw("SELECT COUNT(*) FROM learning_logs WHERE user_id = ?", userID).Scan(&stats.LearningLogCount)
	db.Raw("SELECT COALESCE(SUM(duration), 0) FROM learning_logs WHERE user_id = ?", userID).Scan(&stats.LearningLogTotalDuration)

	// 投稿数
	db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.PostCount)

	// GitHubコントリビューション日数（count > 0 の日数）
	db.Raw("SELECT COUNT(DISTINCT date) FROM git_hub_contributions WHERE user_id = ? AND count > 0", userID).Scan(&stats.GitHubContributionDays)

	// 完了した学習目標数
	db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// コメント数
	db.Raw("SELECT COUNT(*) FROM comments WHERE user_id = ?", userID).Scan(&stats.CommentCount)

	// 受け取ったいいね総数（自分の投稿の like_count の合計）
	db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.LikesReceived)

	// 学習ログストリーク（失敗しても 0 のまま続行する）
	if streak, err := r.calculateLearningLogStreak(ctx, userID); err == nil {
		stats.CurrentStreak = streak
	}

	return stats, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
func (r *levelRepository) calculateLearningLogStreak(ctx context.Context, userID uint) (int, error) {
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

	today := time.Now().UTC().Truncate(24 * time.Hour)
	streak := 0

	for i, d := range dates {
		dDate := d.Date.UTC().Truncate(24 * time.Hour)
		if i == 0 {
			// 最新の日付が今日か昨日でなければストリークは 0
			if today.Sub(dDate) >= 48*time.Hour {
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
