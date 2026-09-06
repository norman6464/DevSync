package persistence

import (
	"context"
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// badgeRepository は [repository.BadgeStatsReader] の sqlc(pgx) 実装。
// 複数テーブルから各種カウントおよびストリークを算出する。
type badgeRepository struct {
	q *sqlcgen.Queries
}

// NewBadgeRepository は BadgeStatsReader の sqlc(pgx) 実装を返す。
func NewBadgeRepository(q *sqlcgen.Queries) repository.BadgeStatsReader {
	return &badgeRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BadgeStatsReader = (*badgeRepository)(nil)

// GetBadgeStats はバッジ判定に必要な全統計を DB から集計する。
// 個々の集計クエリのエラーは移行前と同じく無視し、ストリーク算出の失敗だけを error として返す。
func (r *badgeRepository) GetBadgeStats(ctx context.Context, userID uint) (*model.BadgeStats, error) {
	stats := &model.BadgeStats{}
	uid := int64(userID)

	if v, err := r.q.SumGitHubContributionsByUser(ctx, uid); err == nil {
		stats.TotalContributions = int(v)
	}

	streak, err := r.calculateGitHubStreak(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	if v, err := r.q.CountPostsByUser(ctx, uid); err == nil {
		stats.TotalPosts = int(v)
	}
	if v, err := r.q.SumPostLikesReceivedByUser(ctx, uid); err == nil {
		stats.TotalLikesReceived = int(v)
	}
	if v, err := r.q.CountFollowersByUser(ctx, uid); err == nil {
		stats.FollowerCount = int(v)
	}
	if v, err := r.q.CountFollowingByUser(ctx, uid); err == nil {
		stats.FollowingCount = int(v)
	}
	if v, err := r.q.CountAnswersByUserIncludingDeleted(ctx, uid); err == nil {
		stats.QAAnswerCount = int(v)
	}
	completedStatus := learningGoalCompletedStatus
	if v, err := r.q.CountCompletedLearningGoalsByUser(ctx, sqlcgen.CountCompletedLearningGoalsByUserParams{
		UserID: uid,
		Status: &completedStatus,
	}); err == nil {
		stats.CompletedGoals = int(v)
	}

	// 学習ログ連続記録日数（失敗しても 0 のまま続行する）
	if logStreak, err := r.calculateLearningLogStreak(ctx, userID); err == nil {
		stats.LearningLogStreak = logStreak
	}

	return stats, nil
}

// calculateGitHubStreak は GitHub コントリビューションの連続日数を算出する。
func (r *badgeRepository) calculateGitHubStreak(ctx context.Context, userID uint) (int, error) {
	rows, err := r.q.ListGitHubContributionsByUser(ctx, int64(userID))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ContributedOn.Time.After(rows[j].ContributedOn.Time)
	})

	streak := 0
	today := normalizeToCalendarDay(time.Now())

	for _, row := range rows {
		cDate := normalizeToCalendarDay(row.ContributedOn.Time)
		// 期待日の当日または前日（開始が昨日でも連続とみなす従来仕様）
		expectedDate := today.AddDate(0, 0, -streak)
		if cDate.Equal(expectedDate) || cDate.Equal(expectedDate.AddDate(0, 0, -1)) {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}

// calculateLearningLogStreak は学習ログの連続記録日数を算出する。
func (r *badgeRepository) calculateLearningLogStreak(ctx context.Context, userID uint) (int, error) {
	rows, err := r.q.ListLearningLogDatesByUser(ctx, int64(userID))
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Time.After(rows[j].Time)
	})

	today := normalizeToCalendarDay(time.Now())
	firstDate := normalizeToCalendarDay(rows[0].Time)
	if !isTodayOrYesterday(firstDate, today) {
		return 0, nil
	}

	streak := 1
	for i := 1; i < len(rows); i++ {
		prev := normalizeToCalendarDay(rows[i-1].Time)
		curr := normalizeToCalendarDay(rows[i].Time)
		if isNextCalendarDay(prev, curr) {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}
