package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetLevelInfoUseCase は指定ユーザーのレベル情報を返す。
type GetLevelInfoUseCase struct {
	stats repository.XPStatsReader
}

// NewGetLevelInfoUseCase は GetLevelInfoUseCase を生成する。
func NewGetLevelInfoUseCase(stats repository.XPStatsReader) *GetLevelInfoUseCase {
	return &GetLevelInfoUseCase{stats: stats}
}

// Execute は統計から XP を算出し、レベル情報にまとめて返す。
func (uc *GetLevelInfoUseCase) Execute(ctx context.Context, userID uint) (*model.LevelInfo, error) {
	stats, err := uc.stats.GetXPStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return CalculateLevelInfo(CalculateXPFromStats(stats).Total), nil
}

// GetXPBreakdownUseCase は指定ユーザーの XP 内訳を返す。
type GetXPBreakdownUseCase struct {
	stats repository.XPStatsReader
}

// NewGetXPBreakdownUseCase は GetXPBreakdownUseCase を生成する。
func NewGetXPBreakdownUseCase(stats repository.XPStatsReader) *GetXPBreakdownUseCase {
	return &GetXPBreakdownUseCase{stats: stats}
}

// Execute は統計から XP の内訳を算出して返す。
func (uc *GetXPBreakdownUseCase) Execute(ctx context.Context, userID uint) (*model.XPBreakdown, error) {
	stats, err := uc.stats.GetXPStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return CalculateXPFromStats(stats), nil
}

// CalculateXPFromStats は XPStats から XP 内訳を計算する純粋関数。
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

	breakdown.Total = breakdown.LearningLogs + breakdown.Posts + breakdown.GitHub +
		breakdown.Goals + breakdown.Comments + breakdown.Likes + breakdown.StreakBonus

	return breakdown
}

// CalculateLevel は累計 XP からレベルを算出する純粋関数。
// レベル N に必要な累計 XP = 100 * N * (N+1) / 2
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

// XPForLevel は指定レベルに必要な累計 XP を返す純粋関数。
func XPForLevel(level int) int {
	return 100 * level * (level + 1) / 2
}

// CalculateLevelInfo は累計 XP から LevelInfo 全体を算出する純粋関数。
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
