package service

import (
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BadgeStats holds the aggregated statistics needed for badge evaluation.
type BadgeStats struct {
	TotalContributions int
	CurrentStreak      int
	LearningLogStreak  int
	TotalPosts         int
	TotalLikesReceived int
	FollowerCount      int
	FollowingCount     int
	QAAnswerCount      int
	CompletedGoals     int
}

// BadgeResult represents a single badge with its earned status.
type BadgeResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Earned      bool   `json:"earned"`
}

// BadgeService handles badge evaluation business logic.
type BadgeService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

// NewBadgeService creates a new BadgeService.
func NewBadgeService(db *gorm.DB, notificationService *NotificationService) *BadgeService {
	return &BadgeService{db: db, notificationService: notificationService}
}

// GetUserBadges returns all badges with earned status for the given user.
func (s *BadgeService) GetUserBadges(userID uint) ([]BadgeResult, error) {
	stats, err := s.getBadgeStats(userID)
	if err != nil {
		return nil, err
	}
	return s.evaluateBadges(stats), nil
}

// NotifyBadgeEarned creates a notification for a newly earned badge.
func (s *BadgeService) NotifyBadgeEarned(userID uint, badgeID string) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBadge,
		ActorID: userID,
		BadgeID: &badgeID,
	}
	return s.notificationService.CreateNotification(notification)
}

// getBadgeStats collects all statistics from the database needed for badge evaluation.
func (s *BadgeService) getBadgeStats(userID uint) (*BadgeStats, error) {
	stats := &BadgeStats{}

	// Total contributions
	s.db.Raw("SELECT COALESCE(SUM(count), 0) FROM git_hub_contributions WHERE user_id = ?", userID).Scan(&stats.TotalContributions)

	// Current streak
	streak, err := s.calculateStreak(userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	// Total posts
	s.db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)

	// Total likes received
	s.db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalLikesReceived)

	// Follower count
	s.db.Raw("SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID).Scan(&stats.FollowerCount)

	// Following count
	s.db.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID).Scan(&stats.FollowingCount)

	// QA answer count
	s.db.Raw("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&stats.QAAnswerCount)

	// Completed goals
	s.db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	// Learning log streak
	logStreak, err := s.calculateLearningLogStreak(userID)
	if err == nil {
		stats.LearningLogStreak = logStreak
	}

	return stats, nil
}

// calculateStreak calculates the current contribution streak (consecutive days).
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

// calculateLearningLogStreak calculates the current streak from learning logs.
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

// evaluateBadges evaluates all badges based on the given stats.
func (s *BadgeService) evaluateBadges(stats *BadgeStats) []BadgeResult {
	combinedStreak := stats.CurrentStreak
	if stats.LearningLogStreak > combinedStreak {
		combinedStreak = stats.LearningLogStreak
	}

	return []BadgeResult{
		// Contribution badges
		{ID: "first-commit", Name: "badges.firstCommit", Description: "badges.firstCommitDesc", Category: "contribution", Earned: stats.TotalContributions >= 1},
		{ID: "contributor", Name: "badges.contributor", Description: "badges.contributorDesc", Category: "contribution", Earned: stats.TotalContributions >= 50},
		{ID: "code-warrior", Name: "badges.codeWarrior", Description: "badges.codeWarriorDesc", Category: "contribution", Earned: stats.TotalContributions >= 200},
		{ID: "commit-master", Name: "badges.commitMaster", Description: "badges.commitMasterDesc", Category: "contribution", Earned: stats.TotalContributions >= 500},
		{ID: "legend", Name: "badges.legend", Description: "badges.legendDesc", Category: "contribution", Earned: stats.TotalContributions >= 1000},

		// Streak badges
		{ID: "week-streak", Name: "badges.weekStreak", Description: "badges.weekStreakDesc", Category: "streak", Earned: combinedStreak >= 7},
		{ID: "month-streak", Name: "badges.monthStreak", Description: "badges.monthStreakDesc", Category: "streak", Earned: combinedStreak >= 30},

		// Post badges
		{ID: "first-post", Name: "badges.firstPost", Description: "badges.firstPostDesc", Category: "post", Earned: stats.TotalPosts >= 1},
		{ID: "blogger", Name: "badges.blogger", Description: "badges.bloggerDesc", Category: "post", Earned: stats.TotalPosts >= 10},

		// Engagement badges
		{ID: "liked", Name: "badges.liked", Description: "badges.likedDesc", Category: "engagement", Earned: stats.TotalLikesReceived >= 10},
		{ID: "popular", Name: "badges.popular", Description: "badges.popularDesc", Category: "engagement", Earned: stats.TotalLikesReceived >= 50},

		// Social badges
		{ID: "friendly", Name: "badges.friendly", Description: "badges.friendlyDesc", Category: "social", Earned: stats.FollowingCount >= 5},
		{ID: "influencer", Name: "badges.influencer", Description: "badges.influencerDesc", Category: "social", Earned: stats.FollowerCount >= 10},
		{ID: "star", Name: "badges.star", Description: "badges.starDesc", Category: "social", Earned: stats.FollowerCount >= 50},

		// Q&A badges
		{ID: "qa-first-answer", Name: "badges.qaFirstAnswer", Description: "badges.qaFirstAnswerDesc", Category: "qa", Earned: stats.QAAnswerCount >= 1},
		{ID: "qa-helper", Name: "badges.qaHelper", Description: "badges.qaHelperDesc", Category: "qa", Earned: stats.QAAnswerCount >= 10},

		// Goal badges
		{ID: "goal-achiever", Name: "badges.goalAchiever", Description: "badges.goalAchieverDesc", Category: "goal", Earned: stats.CompletedGoals >= 5},
		{ID: "goal-master", Name: "badges.goalMaster", Description: "badges.goalMasterDesc", Category: "goal", Earned: stats.CompletedGoals >= 20},
	}
}
