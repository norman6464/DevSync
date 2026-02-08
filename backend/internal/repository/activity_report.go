package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ActivityReportRepository はアクティビティレポートの生成・取得を提供するリポジトリ実装。
// GitHubコントリビューション、投稿、コメント、学習目標など複数テーブルからデータを集約する。
type ActivityReportRepository struct {
	db *gorm.DB
}

// NewActivityReportRepository は新しいActivityReportRepositoryインスタンスを生成する。
func NewActivityReportRepository(db *gorm.DB) *ActivityReportRepository {
	return &ActivityReportRepository{db: db}
}

// GetWeeklyReport は指定ユーザーの今週（日曜始まり）のアクティビティレポートを生成する。
func (r *ActivityReportRepository) GetWeeklyReport(userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	// 今週の日曜日を起点として算出
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return r.generateReport(userID, model.ReportPeriodWeekly, startOfWeek, endOfWeek)
}

// GetMonthlyReport は指定ユーザーの今月のアクティビティレポートを生成する。
func (r *ActivityReportRepository) GetMonthlyReport(userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return r.generateReport(userID, model.ReportPeriodMonthly, startOfMonth, endOfMonth)
}

// generateReport は指定期間のアクティビティレポートを各テーブルから集約して生成する。
// コントリビューション、投稿、コメント、いいね、学習目標、フォロワー、メッセージを集計する。
func (r *ActivityReportRepository) generateReport(userID uint, period model.ReportPeriod, startDate, endDate time.Time) (*model.ActivityReport, error) {
	report := &model.ActivityReport{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    userID,
	}

	// 期間内のGitHubコントリビューション合計を取得
	var totalContributions int64
	r.db.Model(&model.GitHubContribution{}).
		Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Select("COALESCE(SUM(count), 0)").
		Scan(&totalContributions)
	report.TotalContributions = int(totalContributions)

	// 期間内の投稿数を取得
	var postsCreated int64
	r.db.Model(&model.Post{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&postsCreated)
	report.PostsCreated = int(postsCreated)

	// 期間内のコメント数を取得
	var commentsCreated int64
	r.db.Model(&model.Comment{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&commentsCreated)
	report.CommentsCreated = int(commentsCreated)

	// 期間内にユーザーの投稿が受け取ったいいね数を取得
	var likesReceived int64
	r.db.Model(&model.Like{}).
		Joins("JOIN posts ON likes.post_id = posts.id").
		Where("posts.user_id = ? AND likes.created_at >= ? AND likes.created_at < ?", userID, startDate, endDate).
		Count(&likesReceived)
	report.LikesReceived = int(likesReceived)

	// 期間内に完了した学習目標数を取得
	var goalsCompleted int64
	r.db.Model(&model.LearningGoal{}).
		Where("user_id = ? AND status = ? AND completed_at >= ? AND completed_at < ?", userID, model.GoalStatusCompleted, startDate, endDate).
		Count(&goalsCompleted)
	report.GoalsCompleted = int(goalsCompleted)

	// アクティブな学習目標の平均進捗率を取得
	var avgProgress float64
	r.db.Model(&model.LearningGoal{}).
		Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).
		Select("COALESCE(AVG(progress), 0)").
		Scan(&avgProgress)
	report.GoalsProgress = int(avgProgress)

	// 期間内の新規フォロワー数を取得
	var newFollowers int64
	r.db.Model(&model.Follow{}).
		Where("followee_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&newFollowers)
	report.NewFollowers = int(newFollowers)

	// 期間内の送受信メッセージ数を合算して取得
	var messagesSent int64
	var messagesReceived int64
	r.db.Model(&model.Message{}).
		Where("sender_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&messagesSent)
	r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&messagesReceived)
	report.MessagesExchanged = int(messagesSent + messagesReceived)

	// 日別のアクティビティデータを生成
	report.DailyContributions = r.getDailyActivity(userID, startDate, endDate)

	// 使用言語トップ5を取得
	report.TopLanguages = r.getTopLanguages(userID)

	return report, nil
}

// getDailyActivity は指定期間の日別アクティビティ（コントリビューション・投稿・コメント）を取得する。
// 期間内の全日付に対してデータを生成し、データがない日は0として返す。
func (r *ActivityReportRepository) getDailyActivity(userID uint, startDate, endDate time.Time) []model.DailyActivity {
	var activities []model.DailyActivity

	// 期間内の全日付をループして各日のデータを収集
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		nextDay := d.AddDate(0, 0, 1)

		activity := model.DailyActivity{
			Date: dateStr,
		}

		// 当日のGitHubコントリビューション数を取得
		var contributions int64
		r.db.Model(&model.GitHubContribution{}).
			Where("user_id = ? AND date = ?", userID, dateStr).
			Select("COALESCE(SUM(count), 0)").
			Scan(&contributions)
		activity.Contributions = int(contributions)

		// 当日の投稿数を取得
		var posts int64
		r.db.Model(&model.Post{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, d, nextDay).
			Count(&posts)
		activity.Posts = int(posts)

		// 当日のコメント数を取得
		var comments int64
		r.db.Model(&model.Comment{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, d, nextDay).
			Count(&comments)
		activity.Comments = int(comments)

		activities = append(activities, activity)
	}

	return activities
}

// getTopLanguages はユーザーのGitHub言語統計から上位5言語を取得する。
// バイト数の降順でソートされる。
func (r *ActivityReportRepository) getTopLanguages(userID uint) []model.LanguageActivity {
	var languages []model.LanguageActivity

	r.db.Model(&model.GitHubLanguageStat{}).
		Where("user_id = ?", userID).
		Select("language, bytes, repo_count as repos").
		Order("bytes DESC").
		Limit(5).
		Scan(&languages)

	return languages
}

// GetComparison は現在の期間と前の期間を比較し、各指標の差分とトレンドを算出する。
// 週次の場合は今週と先週、月次の場合は今月と先月を比較する。
func (r *ActivityReportRepository) GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
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

	currentReport, _ := r.generateReport(userID, period, currentStart, currentEnd)
	prevReport, _ := r.generateReport(userID, period, prevStart, prevEnd)

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
