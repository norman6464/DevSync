package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// topLanguagesLimit はトップ言語一覧の最大件数。
const topLanguagesLimit = 5

// activityReportRepository は [repository.ActivityReportRepository] の sqlc(pgx) 実装。
// GitHub コントリビューション、投稿、コメント、学習目標など複数テーブルからデータを集約する。
type activityReportRepository struct {
	q *sqlcgen.Queries
}

// NewActivityReportRepository は ActivityReportRepository の sqlc(pgx) 実装を返す。
func NewActivityReportRepository(q *sqlcgen.Queries) repository.ActivityReportRepository {
	return &activityReportRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ActivityReportRepository = (*activityReportRepository)(nil)

// GetWeeklyReport は指定ユーザーの今週（日曜始まり）のアクティビティレポートを生成する。
func (r *activityReportRepository) GetWeeklyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	// 今週の日曜日を起点として算出
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return r.generateReport(ctx, userID, model.ReportPeriodWeekly, startOfWeek, endOfWeek)
}

// GetMonthlyReport は指定ユーザーの今月のアクティビティレポートを生成する。
func (r *activityReportRepository) GetMonthlyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return r.generateReport(ctx, userID, model.ReportPeriodMonthly, startOfMonth, endOfMonth)
}

// generateReport は指定期間のアクティビティレポートを各テーブルから集約して生成する。
// コントリビューション、投稿、コメント、いいね、学習目標、フォロワー、メッセージを集計する。
func (r *activityReportRepository) generateReport(ctx context.Context, userID uint, period model.ReportPeriod, startDate, endDate time.Time) (*model.ActivityReport, error) {
	report := &model.ActivityReport{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    userID,
	}

	start := toTimestamptzNotNull(startDate)
	end := toTimestamptzNotNull(endDate)

	// 期間内のGitHubコントリビューション合計を取得
	totalContributions, err := r.q.SumContributionsInRange(ctx, sqlcgen.SumContributionsInRangeParams{
		UserID: int64(userID), ContributedOn: toDateNotNull(startDate), ContributedOn_2: toDateNotNull(endDate),
	})
	if err != nil {
		return nil, err
	}
	report.TotalContributions = int(totalContributions)

	// 期間内の投稿数を取得
	postsCreated, err := r.q.CountPostsInRange(ctx, sqlcgen.CountPostsInRangeParams{
		UserID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.PostsCreated = int(postsCreated)

	// 期間内のコメント数を取得
	commentsCreated, err := r.q.CountCommentsInRange(ctx, sqlcgen.CountCommentsInRangeParams{
		UserID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.CommentsCreated = int(commentsCreated)

	// 期間内にユーザーの投稿が受け取ったいいね数を取得
	likesReceived, err := r.q.CountLikesReceivedInRange(ctx, sqlcgen.CountLikesReceivedInRangeParams{
		UserID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.LikesReceived = int(likesReceived)

	// 期間内に完了した学習目標数を取得
	goalsCompleted, err := r.q.CountCompletedGoalsInRange(ctx, sqlcgen.CountCompletedGoalsInRangeParams{
		UserID: int64(userID), CompletedAt: start, CompletedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.GoalsCompleted = int(goalsCompleted)

	// アクティブな学習目標の平均進捗率を取得
	avgProgress, err := r.q.GetAverageActiveProgressByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	report.GoalsProgress = int(avgProgress)

	// 期間内の新規フォロワー数を取得
	newFollowers, err := r.q.CountNewFollowersInRange(ctx, sqlcgen.CountNewFollowersInRangeParams{
		FolloweeID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.NewFollowers = int(newFollowers)

	// 期間内の送受信メッセージ数を合算して取得
	messagesSent, err := r.q.CountMessagesSentInRange(ctx, sqlcgen.CountMessagesSentInRangeParams{
		SenderID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	messagesReceived, err := r.q.CountMessagesReceivedInRange(ctx, sqlcgen.CountMessagesReceivedInRangeParams{
		ReceiverID: int64(userID), CreatedAt: start, CreatedAt_2: end,
	})
	if err != nil {
		return nil, err
	}
	report.MessagesExchanged = int(messagesSent + messagesReceived)

	// 日別のアクティビティデータを生成
	dailyContributions, err := r.getDailyActivity(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	report.DailyContributions = dailyContributions

	// 使用言語トップ5を取得
	topLanguages, err := r.getTopLanguages(ctx, userID)
	if err != nil {
		return nil, err
	}
	report.TopLanguages = topLanguages

	return report, nil
}

// getDailyActivity は指定期間の日別アクティビティ（コントリビューション・投稿・コメント）を取得する。
// 期間内の全日付に対してデータを生成し、データがない日は0として返す。
func (r *activityReportRepository) getDailyActivity(ctx context.Context, userID uint, startDate, endDate time.Time) ([]model.DailyActivity, error) {
	var activities []model.DailyActivity

	// 期間内の全日付をループして各日のデータを収集
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		nextDay := d.AddDate(0, 0, 1)
		dStart := toTimestamptzNotNull(d)
		dEnd := toTimestamptzNotNull(nextDay)

		activity := model.DailyActivity{
			Date: dateStr,
		}

		// 当日のGitHubコントリビューション数を取得（contributed_onはdate型なのでタイムゾーンの影響を受けない）
		contributions, err := r.q.SumContributionsInRange(ctx, sqlcgen.SumContributionsInRangeParams{
			UserID: int64(userID), ContributedOn: toDateNotNull(d), ContributedOn_2: toDateNotNull(nextDay),
		})
		if err != nil {
			return nil, err
		}
		activity.Contributions = int(contributions)

		// 当日の投稿数を取得
		posts, err := r.q.CountPostsInRange(ctx, sqlcgen.CountPostsInRangeParams{
			UserID: int64(userID), CreatedAt: dStart, CreatedAt_2: dEnd,
		})
		if err != nil {
			return nil, err
		}
		activity.Posts = int(posts)

		// 当日のコメント数を取得
		comments, err := r.q.CountCommentsInRange(ctx, sqlcgen.CountCommentsInRangeParams{
			UserID: int64(userID), CreatedAt: dStart, CreatedAt_2: dEnd,
		})
		if err != nil {
			return nil, err
		}
		activity.Comments = int(comments)

		activities = append(activities, activity)
	}

	return activities, nil
}

// getTopLanguages はユーザーのGitHub言語統計から上位5言語を取得する。
// バイト数の降順でソートされる。
func (r *activityReportRepository) getTopLanguages(ctx context.Context, userID uint) ([]model.LanguageActivity, error) {
	rows, err := r.q.ListGitHubLanguageStatsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if len(rows) > topLanguagesLimit {
		rows = rows[:topLanguagesLimit]
	}
	languages := make([]model.LanguageActivity, len(rows))
	for i, row := range rows {
		languages[i] = model.LanguageActivity{
			Language: row.Language,
			Bytes:    row.Bytes,
			Repos:    int(row.RepoCount),
		}
	}
	return languages, nil
}

// GetComparison は現在の期間と前の期間を比較し、各指標の差分とトレンドを算出する。
// 週次の場合は今週と先週、月次の場合は今月と先月を比較する。
func (r *activityReportRepository) GetComparison(ctx context.Context, userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	now := time.Now()
	var currentStart, currentEnd, prevStart, prevEnd time.Time

	if period == model.ReportPeriodWeekly {
		// 週次: 今週と先週を比較
		currentStart = now.AddDate(0, 0, -int(now.Weekday()))
		currentStart = time.Date(currentStart.Year(), currentStart.Month(), currentStart.Day(), 0, 0, 0, 0, now.Location())
		currentEnd = currentStart.AddDate(0, 0, 7)
		prevStart = currentStart.AddDate(0, 0, -7)
		prevEnd = currentStart
	} else {
		// 月次: 今月と先月を比較
		currentStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		currentEnd = currentStart.AddDate(0, 1, 0)
		prevStart = currentStart.AddDate(0, -1, 0)
		prevEnd = currentStart
	}

	currentReport, _ := r.generateReport(ctx, userID, period, currentStart, currentEnd)
	prevReport, _ := r.generateReport(ctx, userID, period, prevStart, prevEnd)

	comparison := &model.ReportComparison{
		ContributionsDiff: currentReport.TotalContributions - prevReport.TotalContributions,
		PostsDiff:         currentReport.PostsCreated - prevReport.PostsCreated,
		FollowersDiff:     currentReport.NewFollowers - prevReport.NewFollowers,
		GoalsDiff:         currentReport.GoalsCompleted - prevReport.GoalsCompleted,
	}

	// トレンド率を算出（コントリビューション+投稿×10+目標完了×20の重み付きスコアで比較）
	prevTotal := float64(prevReport.TotalContributions + prevReport.PostsCreated*10 + prevReport.GoalsCompleted*20)
	currTotal := float64(currentReport.TotalContributions + currentReport.PostsCreated*10 + currentReport.GoalsCompleted*20)
	if prevTotal > 0 {
		comparison.TrendPercentage = ((currTotal - prevTotal) / prevTotal) * 100
	} else if currTotal > 0 {
		comparison.TrendPercentage = 100
	}

	return comparison, nil
}
