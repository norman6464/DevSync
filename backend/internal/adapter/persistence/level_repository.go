package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningGoalCompletedStatus は完了扱いとする learning_goals.status の値（移行前の GORM 実装と同じ）。
const learningGoalCompletedStatus = "completed"

// levelRepository は [repository.XPStatsReader] の sqlc(pgx) 実装。
// 複数テーブルから XP 計算に必要なデータを集計する。
type levelRepository struct {
	q *sqlcgen.Queries
}

// NewLevelRepository は XPStatsReader の sqlc(pgx) 実装を返す。
func NewLevelRepository(q *sqlcgen.Queries) repository.XPStatsReader {
	return &levelRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.XPStatsReader = (*levelRepository)(nil)

// GetXPStats は XP 計算に必要な全統計を DB から集計する。
// 個々の集計クエリのエラーは移行前と同じく無視する。
func (r *levelRepository) GetXPStats(ctx context.Context, userID uint) (*model.XPStats, error) {
	stats := &model.XPStats{}
	uid := int64(userID)

	if v, err := r.q.CountLearningLogsByUser(ctx, uid); err == nil {
		stats.LearningLogCount = int(v)
	}
	if v, err := r.q.SumLearningLogDurationByUser(ctx, uid); err == nil {
		stats.LearningLogTotalDuration = int(v)
	}
	if v, err := r.q.CountPostsByUser(ctx, uid); err == nil {
		stats.PostCount = int(v)
	}
	if v, err := r.q.CountGitHubContributionDaysByUser(ctx, uid); err == nil {
		stats.GitHubContributionDays = int(v)
	}
	completedStatus := learningGoalCompletedStatus
	if v, err := r.q.CountCompletedLearningGoalsByUser(ctx, sqlcgen.CountCompletedLearningGoalsByUserParams{
		UserID: uid,
		Status: &completedStatus,
	}); err == nil {
		stats.CompletedGoals = int(v)
	}
	if v, err := r.q.CountCommentsByUser(ctx, uid); err == nil {
		stats.CommentCount = int(v)
	}
	if v, err := r.q.SumPostLikesReceivedByUser(ctx, uid); err == nil {
		stats.LikesReceived = int(v)
	}

	// 学習ログストリーク（失敗しても 0 のまま続行する）
	if streak, err := r.calculateLearningLogStreak(ctx, userID); err == nil {
		stats.CurrentStreak = streak
	}

	return stats, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
func (r *levelRepository) calculateLearningLogStreak(ctx context.Context, userID uint) (int, error) {
	rows, err := r.q.ListLearningLogDatesByUser(ctx, int64(userID))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}

	today := normalizeToCalendarDay(time.Now())
	streak := 0

	for i, d := range rows {
		dDate := normalizeToCalendarDay(d.Time)
		if i == 0 {
			// 最新の日付が今日か昨日でなければストリークは 0
			if !isTodayOrYesterday(dDate, today) {
				return 0, nil
			}
			streak = 1
			continue
		}

		prevDate := normalizeToCalendarDay(rows[i-1].Time)
		if isNextCalendarDay(prevDate, dDate) {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}
