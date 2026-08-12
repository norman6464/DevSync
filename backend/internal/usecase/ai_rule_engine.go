package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GenerateAIAdviceUseCase は学習状況からルールベースでアドバイスを生成する。
type GenerateAIAdviceUseCase struct {
	logs      repository.LearningLogRepository
	goals     repository.LearningGoalRepository
	roadmaps  repository.RoadmapRepository
	github    repository.GitHubRepository
	resources repository.LearningResourceRepository
	users     repository.UserSkillsReader
}

// NewGenerateAIAdviceUseCase は GenerateAIAdviceUseCase を生成する。
func NewGenerateAIAdviceUseCase(
	logs repository.LearningLogRepository,
	goals repository.LearningGoalRepository,
	roadmaps repository.RoadmapRepository,
	github repository.GitHubRepository,
	resources repository.LearningResourceRepository,
	users repository.UserSkillsReader,
) *GenerateAIAdviceUseCase {
	return &GenerateAIAdviceUseCase{
		logs: logs, goals: goals, roadmaps: roadmaps,
		github: github, resources: resources, users: users,
	}
}

// aiUserContext はルールエンジンが分析するユーザーデータの集約構造体。
type aiUserContext struct {
	streak    *model.StreakInfo
	goals     []model.LearningGoal
	stats     *model.LearningGoalStats
	roadmaps  []model.Roadmap
	langStats []model.GitHubLanguageStat
	logs      []model.LearningLog
	resources []model.LearningResource
	user      *model.User
}

// collectContext はルールエンジンに必要なユーザーデータを収集する。
func (uc *GenerateAIAdviceUseCase) collectContext(ctx context.Context, userID uint) (*aiUserContext, error) {
	data := &aiUserContext{}
	var err error

	data.streak, err = uc.logs.GetStreakInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	data.goals, _, err = uc.goals.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, err
	}

	data.stats, err = uc.goals.GetStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	data.roadmaps, _, err = uc.roadmaps.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, err
	}

	data.langStats, err = uc.github.GetLanguageStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	data.logs, _, err = uc.logs.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, err
	}

	data.resources, _, err = uc.resources.FindByUserID(ctx, userID, true, 100, 0)
	if err != nil {
		return nil, err
	}

	data.user, err = uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// paramsToJSON はmap[string]stringをJSON文字列に変換する。
func aiParamsToJSON(params map[string]string) string {
	b, _ := json.Marshal(params)
	return string(b)
}

// GenerateAdvice はルールエンジンを実行し、パーソナライズされたアドバイスを優先度順で返す。
func (uc *GenerateAIAdviceUseCase) Execute(ctx context.Context, userID uint) []model.AIAdvice {
	data, err := uc.collectContext(ctx, userID)
	if err != nil {
		return []model.AIAdvice{}
	}

	var advices []model.AIAdvice

	// ルール1: ストリーク途切れ（優先度1: Critical）
	if data.streak.CurrentStreak == 0 && data.streak.TotalDays > 0 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeStreakRecovery,
			Priority:   model.AdvicePriorityCritical,
			TitleKey:   "advice.streakBroken",
			MessageKey: "advice.streakBrokenMsg",
			Params:     aiParamsToJSON(map[string]string{"longestStreak": fmt.Sprintf("%d", data.streak.LongestStreak)}),
			ActionURL:  "/learning-logs",
		})
	}

	// ルール2: ロードマップ停滞（優先度2: High）
	for _, rm := range data.roadmaps {
		if rm.Status != model.RoadmapStatusActive {
			continue
		}
		stalledDays := time.Since(rm.UpdatedAt).Hours() / 24
		if stalledDays >= 7 && rm.CompletedStepCount < rm.StepCount {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeStalledRoadmap,
				Priority:   model.AdvicePriorityHigh,
				TitleKey:   "advice.roadmapStalled",
				MessageKey: "advice.roadmapStalledMsg",
				Params:     aiParamsToJSON(map[string]string{"roadmapTitle": rm.Title, "days": fmt.Sprintf("%d", int(stalledDays))}),
				ActionURL:  fmt.Sprintf("/roadmaps/%d", rm.ID),
			})
			break // 1つのみ表示
		}
	}

	// ルール3: 目標期限超過（優先度2: High）
	for _, goal := range data.goals {
		if goal.Status == model.GoalStatusActive && goal.TargetDate != nil && goal.TargetDate.Before(time.Now()) {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeGoalOverdue,
				Priority:   model.AdvicePriorityHigh,
				TitleKey:   "advice.goalOverdue",
				MessageKey: "advice.goalOverdueMsg",
				Params:     aiParamsToJSON(map[string]string{"goalTitle": goal.Title}),
				ActionURL:  "/goals",
			})
			break // 1つのみ表示
		}
	}

	// ルール4: React提案（優先度3: Medium）
	hasTypeScript := false
	for _, ls := range data.langStats {
		if ls.Language == "TypeScript" && ls.Bytes > 10000 {
			hasTypeScript = true
			break
		}
	}
	if hasTypeScript {
		hasReactGoal := false
		for _, g := range data.goals {
			if strings.Contains(strings.ToLower(g.Title), "react") ||
				strings.Contains(strings.ToLower(g.Title), "next") {
				hasReactGoal = true
				break
			}
		}
		if !hasReactGoal {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeTechSuggestion,
				Priority:   model.AdvicePriorityMedium,
				TitleKey:   "advice.suggestReact",
				MessageKey: "advice.suggestReactMsg",
				ActionURL:  "/goals",
			})
		}
	}

	// ルール5: 技術提案（優先度3: Medium）
	if len(data.langStats) > 0 && !hasTypeScript {
		topLang := data.langStats[0]
		for _, ls := range data.langStats {
			if ls.Bytes > topLang.Bytes {
				topLang = ls
			}
		}
		goalHasLang := false
		for _, g := range data.goals {
			if strings.Contains(strings.ToLower(g.Title), strings.ToLower(topLang.Language)) {
				goalHasLang = true
				break
			}
		}
		if !goalHasLang && topLang.Bytes > 10000 {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeTechSuggestion,
				Priority:   model.AdvicePriorityMedium,
				TitleKey:   "advice.suggestFromGithub",
				MessageKey: "advice.suggestFromGithubMsg",
				Params:     aiParamsToJSON(map[string]string{"language": topLang.Language}),
				ActionURL:  "/goals",
			})
		}
	}

	// ルール6: 目標未設定（優先度3: Medium）
	if data.stats.TotalGoals == 0 && len(data.langStats) > 0 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeGoalSuggestion,
			Priority:   model.AdvicePriorityMedium,
			TitleKey:   "advice.noGoals",
			MessageKey: "advice.noGoalsMsg",
			ActionURL:  "/goals",
		})
	}

	// ルール7: ロードマップ未作成（優先度3: Medium）
	if len(data.goals) > 0 && len(data.roadmaps) == 0 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeGoalSuggestion,
			Priority:   model.AdvicePriorityMedium,
			TitleKey:   "advice.noRoadmap",
			MessageKey: "advice.noRoadmapMsg",
			ActionURL:  "/roadmaps",
		})
	}

	// ルール8: 難易度UP（優先度4: Low）
	if data.stats.CompletedGoals >= 3 && data.stats.AverageProgress > 70 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeDifficultyUp,
			Priority:   model.AdvicePriorityLow,
			TitleKey:   "advice.difficultyUp",
			MessageKey: "advice.difficultyUpMsg",
			Params:     aiParamsToJSON(map[string]string{"completedGoals": fmt.Sprintf("%d", data.stats.CompletedGoals)}),
			ActionURL:  "/goals",
		})
	}

	// ルール9: 学習称賛（優先度4: Low）
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	var recentTotalMinutes int
	var recentDays int
	for _, log := range data.logs {
		if log.CreatedAt.After(sevenDaysAgo) {
			recentTotalMinutes += log.Duration
			recentDays++
		}
	}
	if recentDays > 0 {
		avgMinutes := recentTotalMinutes / 7
		if avgMinutes >= 60 {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypePraise,
				Priority:   model.AdvicePriorityLow,
				TitleKey:   "advice.praise",
				MessageKey: "advice.praiseMsg",
				Params:     aiParamsToJSON(map[string]string{"avgMinutes": fmt.Sprintf("%d", avgMinutes)}),
			})
		}

		// ルール10: 学習時間少（優先度4: Low）
		if avgMinutes < 15 && avgMinutes > 0 {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeGeneral,
				Priority:   model.AdvicePriorityLow,
				TitleKey:   "advice.suggestMoreTime",
				MessageKey: "advice.suggestMoreTimeMsg",
				ActionURL:  "/learning-logs",
			})
		}
	}

	// ルール11: リソース探索（優先度5: Info）
	if len(data.resources) < 3 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeGeneral,
			Priority:   model.AdvicePriorityInfo,
			TitleKey:   "advice.exploreResources",
			MessageKey: "advice.exploreResourcesMsg",
			ActionURL:  "/resources",
		})
	}

	// 優先度順にソート（昇順: 1が最優先）
	sort.Slice(advices, func(i, j int) bool {
		if advices[i].Priority != advices[j].Priority {
			return advices[i].Priority < advices[j].Priority
		}
		return i < j
	})

	return advices
}
