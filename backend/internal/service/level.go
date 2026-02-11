package service

import (
	"fmt"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LevelService はレベルシステムのビジネスロジックを提供する。
// XP計算、レベル計算、レベルアップ通知を担当する。
type LevelService struct {
	repo                repository.LevelRepositoryInterface
	notificationService NotificationServiceInterface
}

// NewLevelService は新しいLevelServiceインスタンスを生成する。
func NewLevelService(repo repository.LevelRepositoryInterface, ns NotificationServiceInterface) *LevelService {
	return &LevelService{repo: repo, notificationService: ns}
}

// GetLevelInfo は指定ユーザーのレベル情報を取得する。
func (s *LevelService) GetLevelInfo(userID uint) (*model.LevelInfo, error) {
	stats, err := s.repo.GetXPStats(userID)
	if err != nil {
		return nil, err
	}
	breakdown := CalculateXPFromStats(stats)
	info := CalculateLevelInfo(breakdown.Total)
	return info, nil
}

// GetXPBreakdown は指定ユーザーのXP内訳を取得する。
func (s *LevelService) GetXPBreakdown(userID uint) (*model.XPBreakdown, error) {
	stats, err := s.repo.GetXPStats(userID)
	if err != nil {
		return nil, err
	}
	return CalculateXPFromStats(stats), nil
}

// CheckAndNotifyLevelUp は前回のXPからレベルアップしたかを確認し、通知を送信する。
func (s *LevelService) CheckAndNotifyLevelUp(userID uint, previousXP int) error {
	stats, err := s.repo.GetXPStats(userID)
	if err != nil {
		return err
	}
	breakdown := CalculateXPFromStats(stats)
	currentLevel := CalculateLevel(breakdown.Total)
	previousLevel := CalculateLevel(previousXP)

	if currentLevel > previousLevel {
		notification := &model.Notification{
			UserID:  userID,
			Type:    model.NotificationTypeLevelUp,
			ActorID: userID,
		}
		return s.notificationService.CreateNotification(notification)
	}
	return nil
}

// CalculateXPFromStats はXPStatsからXP内訳を計算する純粋関数。
func CalculateXPFromStats(stats *model.XPStats) *model.XPBreakdown {
	breakdown := &model.XPBreakdown{}

	// 学習ログ: 10 * count + int(duration * 0.5)
	breakdown.LearningLogs = stats.LearningLogCount*10 + int(float64(stats.LearningLogTotalDuration)*0.5)

	// 投稿: 30 * count
	breakdown.Posts = stats.PostCount * 30

	// GitHub: 5 * days
	breakdown.GitHub = stats.GitHubContributionDays * 5

	// 目標完了: 50 * count
	breakdown.Goals = stats.CompletedGoals * 50

	// コメント: 5 * count
	breakdown.Comments = stats.CommentCount * 5

	// いいね受取: 3 * count
	breakdown.Likes = stats.LikesReceived * 3

	// ストリークボーナス: (streak / 7) * 20
	breakdown.StreakBonus = (stats.CurrentStreak / 7) * 20

	// 合計
	breakdown.Total = breakdown.LearningLogs + breakdown.Posts + breakdown.GitHub +
		breakdown.Goals + breakdown.Comments + breakdown.Likes + breakdown.StreakBonus

	return breakdown
}

// CalculateLevel は累計XPからレベルを算出する純粋関数。
// レベルNに必要な累計XP = 100 * N * (N+1) / 2
func CalculateLevel(totalXP int) int {
	level := 0
	for {
		requiredXP := 100 * (level + 1) * (level + 2) / 2
		if totalXP < requiredXP {
			break
		}
		level++
	}
	return level
}

// XPForLevel は指定レベルに必要な累計XPを返す純粋関数。
func XPForLevel(level int) int {
	return 100 * level * (level + 1) / 2
}

// CalculateLevelInfo は累計XPからLevelInfo全体を算出する純粋関数。
func CalculateLevelInfo(totalXP int) *model.LevelInfo {
	level := CalculateLevel(totalXP)
	currentLevelXP := XPForLevel(level)
	nextLevelXP := XPForLevel(level + 1)
	progressXP := totalXP - currentLevelXP
	xpRange := nextLevelXP - currentLevelXP

	var progressPercent float64
	if xpRange > 0 {
		progressPercent = float64(progressXP) / float64(xpRange) * 100
	}

	return &model.LevelInfo{
		Level:           level,
		TotalXP:         totalXP,
		CurrentLevelXP:  currentLevelXP,
		NextLevelXP:     nextLevelXP,
		ProgressXP:      progressXP,
		ProgressPercent: progressPercent,
	}
}

// GetLevelTitle はレベルに応じた称号を返す（将来拡張用）。
func GetLevelTitle(level int) string {
	switch {
	case level == 0:
		return "Newcomer"
	case level <= 5:
		return "Beginner"
	case level <= 10:
		return "Intermediate"
	case level <= 20:
		return "Advanced"
	case level <= 30:
		return "Expert"
	case level <= 40:
		return "Master"
	default:
		return fmt.Sprintf("Legend Lv.%d", level)
	}
}
