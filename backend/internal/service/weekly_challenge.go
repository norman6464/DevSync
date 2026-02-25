package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"gorm.io/gorm"
)

// WeeklyChallengeService はウィークリーチャレンジのビジネスロジックを提供する。
type WeeklyChallengeService struct {
	repo repository.WeeklyChallengeRepositoryInterface
}

// NewWeeklyChallengeService は新しいWeeklyChallengeServiceを生成する。
func NewWeeklyChallengeService(repo repository.WeeklyChallengeRepositoryInterface) *WeeklyChallengeService {
	return &WeeklyChallengeService{repo: repo}
}

// challengeTemplate はチャレンジ生成用のテンプレート。
type challengeTemplate struct {
	challengeType model.ChallengeType
	description   string
	targetValue   int
}

var challengeTemplates = []challengeTemplate{
	{model.ChallengeDurationTotal, "weekly_duration_180", 180},
	{model.ChallengeDurationTotal, "weekly_duration_300", 300},
	{model.ChallengeDurationTotal, "weekly_duration_420", 420},
	{model.ChallengeStreakDays, "weekly_streak_3", 3},
	{model.ChallengeStreakDays, "weekly_streak_5", 5},
	{model.ChallengeStreakDays, "weekly_streak_7", 7},
	{model.ChallengeCategoryCount, "weekly_categories_3", 3},
	{model.ChallengeLogCount, "weekly_logs_5", 5},
	{model.ChallengeLogCount, "weekly_logs_7", 7},
	{model.ChallengeLogCount, "weekly_logs_10", 10},
}

// GetCurrentChallenge は今週のチャレンジを返す。存在しない場合は自動生成する。
func (s *WeeklyChallengeService) GetCurrentChallenge(userID uint) (*model.WeeklyChallenge, error) {
	year, week := time.Now().ISOWeek()

	challenge, err := s.repo.FindByUserAndWeek(userID, year, week)
	if err == nil {
		return challenge, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 新しいチャレンジを自動生成
	tmpl := challengeTemplates[rand.Intn(len(challengeTemplates))]
	challenge = &model.WeeklyChallenge{
		UserID:        userID,
		Year:          year,
		Week:          week,
		ChallengeType: tmpl.challengeType,
		Description:   tmpl.description,
		TargetValue:   tmpl.targetValue,
	}

	if err := s.repo.Create(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

// UpdateProgress はチャレンジの進捗を更新する。
func (s *WeeklyChallengeService) UpdateProgress(userID uint, value int) (*model.WeeklyChallenge, error) {
	year, week := time.Now().ISOWeek()

	challenge, err := s.repo.FindByUserAndWeek(userID, year, week)
	if err != nil {
		return nil, err
	}

	challenge.CurrentValue = value

	if !challenge.IsCompleted && value >= challenge.TargetValue {
		challenge.IsCompleted = true
		now := time.Now()
		challenge.CompletedAt = &now
	}

	if err := s.repo.Update(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
