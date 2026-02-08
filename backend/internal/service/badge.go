package service

import (
	"sort"
	"time"

	"gorm.io/gorm"
)

// BadgeStats holds the aggregated statistics needed for badge evaluation.
type BadgeStats struct {
	TotalContributions int
	CurrentStreak      int
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

// GetBadgeStats collects all statistics from the database needed for badge evaluation.
func GetBadgeStats(db *gorm.DB, userID uint) (*BadgeStats, error) {
	stats := &BadgeStats{}

	// Total contributions
	db.Raw("SELECT COALESCE(SUM(count), 0) FROM git_hub_contributions WHERE user_id = ?", userID).Scan(&stats.TotalContributions)

	// Current streak
	streak, err := calculateStreak(db, userID)
	if err != nil {
		return nil, err
	}
	stats.CurrentStreak = streak

	// Total posts
	db.Raw("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalPosts)

	// Total likes received
	db.Raw("SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ?", userID).Scan(&stats.TotalLikesReceived)

	// Follower count
	db.Raw("SELECT COUNT(*) FROM follows WHERE followee_id = ?", userID).Scan(&stats.FollowerCount)

	// Following count
	db.Raw("SELECT COUNT(*) FROM follows WHERE follower_id = ?", userID).Scan(&stats.FollowingCount)

	// QA answer count
	db.Raw("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&stats.QAAnswerCount)

	// Completed goals
	db.Raw("SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = ?", userID, "completed").Scan(&stats.CompletedGoals)

	return stats, nil
}

// calculateStreak calculates the current contribution streak (consecutive days).
func calculateStreak(db *gorm.DB, userID uint) (int, error) {
	type DateCount struct {
		Date  time.Time
		Count int
	}
	var contributions []DateCount
	err := db.Raw("SELECT date, count FROM git_hub_contributions WHERE user_id = ? AND count > 0 ORDER BY date DESC", userID).Scan(&contributions).Error
	if err != nil {
		return 0, err
	}
	if len(contributions) == 0 {
		return 0, nil
	}

	// Sort descending by date (already ordered by query, but ensure)
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

// EvaluateBadges evaluates all badges based on the given stats.
func EvaluateBadges(stats *BadgeStats) []BadgeResult {
	return []BadgeResult{
		// Contribution badges
		{ID: "first-commit", Name: "badges.firstCommit", Description: "badges.firstCommitDesc", Category: "contribution", Earned: stats.TotalContributions >= 1},
		{ID: "contributor", Name: "badges.contributor", Description: "badges.contributorDesc", Category: "contribution", Earned: stats.TotalContributions >= 50},
		{ID: "code-warrior", Name: "badges.codeWarrior", Description: "badges.codeWarriorDesc", Category: "contribution", Earned: stats.TotalContributions >= 200},
		{ID: "commit-master", Name: "badges.commitMaster", Description: "badges.commitMasterDesc", Category: "contribution", Earned: stats.TotalContributions >= 500},
		{ID: "legend", Name: "badges.legend", Description: "badges.legendDesc", Category: "contribution", Earned: stats.TotalContributions >= 1000},

		// Streak badges
		{ID: "week-streak", Name: "badges.weekStreak", Description: "badges.weekStreakDesc", Category: "streak", Earned: stats.CurrentStreak >= 7},
		{ID: "month-streak", Name: "badges.monthStreak", Description: "badges.monthStreakDesc", Category: "streak", Earned: stats.CurrentStreak >= 30},

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

		// Q&A badges (new)
		{ID: "qa-first-answer", Name: "badges.qaFirstAnswer", Description: "badges.qaFirstAnswerDesc", Category: "qa", Earned: stats.QAAnswerCount >= 1},
		{ID: "qa-helper", Name: "badges.qaHelper", Description: "badges.qaHelperDesc", Category: "qa", Earned: stats.QAAnswerCount >= 10},

		// Goal badges (new)
		{ID: "goal-achiever", Name: "badges.goalAchiever", Description: "badges.goalAchieverDesc", Category: "goal", Earned: stats.CompletedGoals >= 5},
		{ID: "goal-master", Name: "badges.goalMaster", Description: "badges.goalMasterDesc", Category: "goal", Earned: stats.CompletedGoals >= 20},
	}
}
