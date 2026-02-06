package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type ActivityReportRepository struct {
	db *gorm.DB
}

func NewActivityReportRepository(db *gorm.DB) *ActivityReportRepository {
	return &ActivityReportRepository{db: db}
}

// GetWeeklyReport generates a weekly report for a user
func (r *ActivityReportRepository) GetWeeklyReport(userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	// Get the start of this week (Sunday)
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return r.generateReport(userID, model.ReportPeriodWeekly, startOfWeek, endOfWeek)
}

// GetMonthlyReport generates a monthly report for a user
func (r *ActivityReportRepository) GetMonthlyReport(userID uint) (*model.ActivityReport, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return r.generateReport(userID, model.ReportPeriodMonthly, startOfMonth, endOfMonth)
}

func (r *ActivityReportRepository) generateReport(userID uint, period model.ReportPeriod, startDate, endDate time.Time) (*model.ActivityReport, error) {
	report := &model.ActivityReport{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		UserID:    userID,
	}

	// Get GitHub contributions in the period
	var totalContributions int64
	r.db.Model(&model.GitHubContribution{}).
		Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Select("COALESCE(SUM(count), 0)").
		Scan(&totalContributions)
	report.TotalContributions = int(totalContributions)

	// Get posts created
	var postsCreated int64
	r.db.Model(&model.Post{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&postsCreated)
	report.PostsCreated = int(postsCreated)

	// Get comments created
	var commentsCreated int64
	r.db.Model(&model.Comment{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&commentsCreated)
	report.CommentsCreated = int(commentsCreated)

	// Get likes received on user's posts
	var likesReceived int64
	r.db.Model(&model.Like{}).
		Joins("JOIN posts ON likes.post_id = posts.id").
		Where("posts.user_id = ? AND likes.created_at >= ? AND likes.created_at < ?", userID, startDate, endDate).
		Count(&likesReceived)
	report.LikesReceived = int(likesReceived)

	// Get goals completed
	var goalsCompleted int64
	r.db.Model(&model.LearningGoal{}).
		Where("user_id = ? AND status = ? AND completed_at >= ? AND completed_at < ?", userID, model.GoalStatusCompleted, startDate, endDate).
		Count(&goalsCompleted)
	report.GoalsCompleted = int(goalsCompleted)

	// Get average progress of active goals
	var avgProgress float64
	r.db.Model(&model.LearningGoal{}).
		Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).
		Select("COALESCE(AVG(progress), 0)").
		Scan(&avgProgress)
	report.GoalsProgress = int(avgProgress)

	// Get new followers
	var newFollowers int64
	r.db.Model(&model.Follow{}).
		Where("followee_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&newFollowers)
	report.NewFollowers = int(newFollowers)

	// Get messages exchanged
	var messagesSent int64
	var messagesReceived int64
	r.db.Model(&model.Message{}).
		Where("sender_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&messagesSent)
	r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).
		Count(&messagesReceived)
	report.MessagesExchanged = int(messagesSent + messagesReceived)

	// Get daily contributions
	report.DailyContributions = r.getDailyActivity(userID, startDate, endDate)

	// Get top languages
	report.TopLanguages = r.getTopLanguages(userID)

	return report, nil
}

func (r *ActivityReportRepository) getDailyActivity(userID uint, startDate, endDate time.Time) []model.DailyActivity {
	var activities []model.DailyActivity

	// Generate all dates in the range
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		nextDay := d.AddDate(0, 0, 1)

		activity := model.DailyActivity{
			Date: dateStr,
		}

		// Get contributions for this day
		var contributions int64
		r.db.Model(&model.GitHubContribution{}).
			Where("user_id = ? AND date = ?", userID, dateStr).
			Select("COALESCE(SUM(count), 0)").
			Scan(&contributions)
		activity.Contributions = int(contributions)

		// Get posts for this day
		var posts int64
		r.db.Model(&model.Post{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, d, nextDay).
			Count(&posts)
		activity.Posts = int(posts)

		// Get comments for this day
		var comments int64
		r.db.Model(&model.Comment{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, d, nextDay).
			Count(&comments)
		activity.Comments = int(comments)

		activities = append(activities, activity)
	}

	return activities
}

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

// GetComparison compares current period with previous period
func (r *ActivityReportRepository) GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	now := time.Now()
	var currentStart, currentEnd, prevStart, prevEnd time.Time

	if period == model.ReportPeriodWeekly {
		currentStart = now.AddDate(0, 0, -int(now.Weekday()))
		currentStart = time.Date(currentStart.Year(), currentStart.Month(), currentStart.Day(), 0, 0, 0, 0, now.Location())
		currentEnd = currentStart.AddDate(0, 0, 7)
		prevStart = currentStart.AddDate(0, 0, -7)
		prevEnd = currentStart
	} else {
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

	// Calculate trend percentage
	prevTotal := float64(prevReport.TotalContributions + prevReport.PostsCreated*10 + prevReport.GoalsCompleted*20)
	currTotal := float64(currentReport.TotalContributions + currentReport.PostsCreated*10 + currentReport.GoalsCompleted*20)
	if prevTotal > 0 {
		comparison.TrendPercentage = ((currTotal - prevTotal) / prevTotal) * 100
	} else if currTotal > 0 {
		comparison.TrendPercentage = 100
	}

	return comparison, nil
}
