package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
)

// userContext はルールエンジンが分析するユーザーデータの集約構造体。
type userContext struct {
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
func (s *AIAdviceService) collectContext(userID uint) (*userContext, error) {
	ctx := &userContext{}
	var err error

	ctx.streak, err = s.logRepo.GetStreakInfo(userID)
	if err != nil {
		return nil, err
	}

	ctx.goals, _, err = s.goalRepo.GetByUserID(userID, 100, 0)
	if err != nil {
		return nil, err
	}

	ctx.stats, err = s.goalRepo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	ctx.roadmaps, err = s.roadmapRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	ctx.langStats, err = s.githubRepo.GetLanguageStats(userID)
	if err != nil {
		return nil, err
	}

	ctx.logs, err = s.logRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	ctx.resources, _, err = s.resourceRepo.FindByUserID(userID, true, 100, 0)
	if err != nil {
		return nil, err
	}

	ctx.user, err = s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

// paramsToJSON はmap[string]stringをJSON文字列に変換する。
func paramsToJSON(params map[string]string) string {
	b, _ := json.Marshal(params)
	return string(b)
}

// GenerateAdvice はルールエンジンを実行し、パーソナライズされたアドバイスを優先度順で返す。
func (s *AIAdviceService) GenerateAdvice(userID uint) []model.AIAdvice {
	ctx, err := s.collectContext(userID)
	if err != nil {
		return []model.AIAdvice{}
	}

	var advices []model.AIAdvice

	// ルール1: ストリーク途切れ（優先度1: Critical）
	if ctx.streak.CurrentStreak == 0 && ctx.streak.TotalDays > 0 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeStreakRecovery,
			Priority:   model.AdvicePriorityCritical,
			TitleKey:   "advice.streakBroken",
			MessageKey: "advice.streakBrokenMsg",
			Params:     paramsToJSON(map[string]string{"longestStreak": fmt.Sprintf("%d", ctx.streak.LongestStreak)}),
			ActionURL:  "/learning-logs",
		})
	}

	// ルール2: ロードマップ停滞（優先度2: High）
	for _, rm := range ctx.roadmaps {
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
				Params:     paramsToJSON(map[string]string{"roadmapTitle": rm.Title, "days": fmt.Sprintf("%d", int(stalledDays))}),
				ActionURL:  fmt.Sprintf("/roadmaps/%d", rm.ID),
			})
			break // 1つのみ表示
		}
	}

	// ルール3: 目標期限超過（優先度2: High）
	for _, goal := range ctx.goals {
		if goal.Status == model.GoalStatusActive && goal.TargetDate != nil && goal.TargetDate.Before(time.Now()) {
			advices = append(advices, model.AIAdvice{
				UserID:     userID,
				Type:       model.AdviceTypeGoalOverdue,
				Priority:   model.AdvicePriorityHigh,
				TitleKey:   "advice.goalOverdue",
				MessageKey: "advice.goalOverdueMsg",
				Params:     paramsToJSON(map[string]string{"goalTitle": goal.Title}),
				ActionURL:  "/goals",
			})
			break // 1つのみ表示
		}
	}

	// ルール4: React提案（優先度3: Medium）
	hasTypeScript := false
	for _, ls := range ctx.langStats {
		if ls.Language == "TypeScript" && ls.Bytes > 10000 {
			hasTypeScript = true
			break
		}
	}
	if hasTypeScript {
		hasReactGoal := false
		for _, g := range ctx.goals {
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
	if len(ctx.langStats) > 0 && !hasTypeScript {
		topLang := ctx.langStats[0]
		for _, ls := range ctx.langStats {
			if ls.Bytes > topLang.Bytes {
				topLang = ls
			}
		}
		goalHasLang := false
		for _, g := range ctx.goals {
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
				Params:     paramsToJSON(map[string]string{"language": topLang.Language}),
				ActionURL:  "/goals",
			})
		}
	}

	// ルール6: 目標未設定（優先度3: Medium）
	if ctx.stats.TotalGoals == 0 && len(ctx.langStats) > 0 {
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
	if len(ctx.goals) > 0 && len(ctx.roadmaps) == 0 {
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
	if ctx.stats.CompletedGoals >= 3 && ctx.stats.AverageProgress > 70 {
		advices = append(advices, model.AIAdvice{
			UserID:     userID,
			Type:       model.AdviceTypeDifficultyUp,
			Priority:   model.AdvicePriorityLow,
			TitleKey:   "advice.difficultyUp",
			MessageKey: "advice.difficultyUpMsg",
			Params:     paramsToJSON(map[string]string{"completedGoals": fmt.Sprintf("%d", ctx.stats.CompletedGoals)}),
			ActionURL:  "/goals",
		})
	}

	// ルール9: 学習称賛（優先度4: Low）
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	var recentTotalMinutes int
	var recentDays int
	for _, log := range ctx.logs {
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
				Params:     paramsToJSON(map[string]string{"avgMinutes": fmt.Sprintf("%d", avgMinutes)}),
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
	if len(ctx.resources) < 3 {
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
