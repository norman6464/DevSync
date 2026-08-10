package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// challengeTemplate はチャレンジ生成用のテンプレート。
type challengeTemplate struct {
	challengeType model.ChallengeType
	description   string
	targetValue   int
}

// challengeTemplates は自動生成時に選ばれる候補。
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

// GetCurrentWeeklyChallengeUseCase は今週のチャレンジを取得する。未登録なら自動生成する。
type GetCurrentWeeklyChallengeUseCase struct {
	challenges repository.WeeklyChallengeRepository
}

// NewGetCurrentWeeklyChallengeUseCase は GetCurrentWeeklyChallengeUseCase を生成する。
func NewGetCurrentWeeklyChallengeUseCase(challenges repository.WeeklyChallengeRepository) *GetCurrentWeeklyChallengeUseCase {
	return &GetCurrentWeeklyChallengeUseCase{challenges: challenges}
}

// Execute は今週のチャレンジを返す。存在しない場合はテンプレートから 1 つ選んで作成する。
func (uc *GetCurrentWeeklyChallengeUseCase) Execute(ctx context.Context, userID uint) (*model.WeeklyChallenge, error) {
	year, week := time.Now().ISOWeek()

	challenge, err := uc.challenges.FindByUserAndWeek(ctx, userID, year, week)
	if err != nil {
		return nil, err
	}
	if challenge != nil {
		return challenge, nil
	}

	// 未登録なので新しいチャレンジを自動生成する
	tmpl := challengeTemplates[rand.Intn(len(challengeTemplates))]
	challenge = &model.WeeklyChallenge{
		UserID:        userID,
		Year:          year,
		Week:          week,
		ChallengeType: tmpl.challengeType,
		Description:   tmpl.description,
		TargetValue:   tmpl.targetValue,
	}

	if err := uc.challenges.Create(ctx, challenge); err != nil {
		return nil, err
	}
	return challenge, nil
}

// UpdateWeeklyChallengeProgressUseCase は今週のチャレンジの進捗を更新する。
type UpdateWeeklyChallengeProgressUseCase struct {
	challenges repository.WeeklyChallengeRepository
}

// NewUpdateWeeklyChallengeProgressUseCase は UpdateWeeklyChallengeProgressUseCase を生成する。
func NewUpdateWeeklyChallengeProgressUseCase(challenges repository.WeeklyChallengeRepository) *UpdateWeeklyChallengeProgressUseCase {
	return &UpdateWeeklyChallengeProgressUseCase{challenges: challenges}
}

// Execute は進捗値を反映し、目標到達時に完了フラグを立てる。
// 対象週のチャレンジが無い場合は移行前と同じく内部エラー扱いにする（HTTP 500）。
func (uc *UpdateWeeklyChallengeProgressUseCase) Execute(ctx context.Context, userID uint, value int) (*model.WeeklyChallenge, error) {
	year, week := time.Now().ISOWeek()

	challenge, err := uc.challenges.FindByUserAndWeek(ctx, userID, year, week)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, errors.New("今週のチャレンジが見つかりません")
	}

	challenge.CurrentValue = value

	if !challenge.IsCompleted && value >= challenge.TargetValue {
		challenge.IsCompleted = true
		now := time.Now()
		challenge.CompletedAt = &now
	}

	if err := uc.challenges.Update(ctx, challenge); err != nil {
		return nil, err
	}
	return challenge, nil
}
