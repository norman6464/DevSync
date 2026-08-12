package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// BuildAIChatPromptUseCase は LLM へ渡すシステムプロンプトを組み立てる。
// ユーザーの学習状況を読み取るだけで、取得に失敗した項目は黙って省く（チャット自体は継続させる）。
type BuildAIChatPromptUseCase struct {
	goals    repository.LearningGoalRepository
	logs     repository.LearningLogRepository
	roadmaps repository.RoadmapRepository
	github   repository.GitHubRepository
}

// NewBuildAIChatPromptUseCase は BuildAIChatPromptUseCase を生成する。
func NewBuildAIChatPromptUseCase(
	goals repository.LearningGoalRepository,
	logs repository.LearningLogRepository,
	roadmaps repository.RoadmapRepository,
	github repository.GitHubRepository,
) *BuildAIChatPromptUseCase {
	return &BuildAIChatPromptUseCase{goals: goals, logs: logs, roadmaps: roadmaps, github: github}
}

// Execute はユーザーの学習状況を集めてシステムプロンプトを返す。
func (uc *BuildAIChatPromptUseCase) Execute(ctx context.Context, userID uint) string {
	goals, _, _ := uc.goals.GetByUserID(ctx, userID, 100, 0)
	streak, _ := uc.logs.GetStreakInfo(ctx, userID)
	roadmaps, _, _ := uc.roadmaps.GetByUserID(ctx, userID, 100, 0)
	langStats, _ := uc.github.GetLanguageStats(ctx, userID)

	return buildAIChatSystemPrompt(goals, streak, roadmaps, langStats)
}

// buildSystemPrompt はユーザーコンテキストからシステムプロンプトを構築する。
func buildAIChatSystemPrompt(
	goals []model.LearningGoal,
	streak *model.StreakInfo,
	roadmaps []model.Roadmap,
	langStats []model.GitHubLanguageStat,
) string {
	var sb strings.Builder

	sb.WriteString("あなたは駆け出しエンジニア向けの学習アドバイザーです。")
	sb.WriteString("ユーザーの学習状況を踏まえて、具体的で実践的なアドバイスを日本語で提供してください。\n\n")

	// ストリーク情報
	if streak != nil {
		sb.WriteString(fmt.Sprintf("【学習ストリーク】現在 %d日連続 / 最長 %d日 / 合計 %d日\n",
			streak.CurrentStreak, streak.LongestStreak, streak.TotalDays))
	}

	// 学習目標
	if len(goals) > 0 {
		sb.WriteString("【学習目標】\n")
		for _, g := range goals {
			status := "進行中"
			if g.Status == model.GoalStatusCompleted {
				status = "完了"
			}
			sb.WriteString(fmt.Sprintf("- %s (%s, 進捗 %d%%)\n", g.Title, status, g.Progress))
		}
	}

	// ロードマップ
	if len(roadmaps) > 0 {
		sb.WriteString("【ロードマップ】\n")
		for _, rm := range roadmaps {
			sb.WriteString(fmt.Sprintf("- %s (進捗 %d%%, %d/%dステップ完了)\n",
				rm.Title, rm.Progress, rm.CompletedStepCount, rm.StepCount))
		}
	}

	// GitHub言語統計
	if len(langStats) > 0 {
		sb.WriteString("【GitHub使用言語】\n")
		for i, ls := range langStats {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("- %s (%dリポジトリ)\n", ls.Language, ls.RepoCount))
		}
	}

	sb.WriteString("\n回答は簡潔に、次のアクションが明確になるようにしてください。")
	return sb.String()
}
