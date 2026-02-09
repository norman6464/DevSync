package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// DailyChatLimit は1日あたりのLLMチャット回数制限。
const DailyChatLimit = 5

// AIAdviceService はルールベース推薦エンジンとLLMチャットオーケストレーションを提供する。
type AIAdviceService struct {
	adviceRepo   repository.AIAdviceRepositoryInterface
	convRepo     repository.AIConversationRepositoryInterface
	goalRepo     repository.LearningGoalRepositoryInterface
	logRepo      repository.LearningLogRepositoryInterface
	roadmapRepo  repository.RoadmapRepositoryInterface
	githubRepo   repository.GitHubRepositoryInterface
	resourceRepo repository.LearningResourceRepositoryInterface
	userRepo     repository.UserRepositoryInterface
	llmClient    LLMClientInterface // nil の場合はLLM未設定
}

// NewAIAdviceService は新しいAIAdviceServiceインスタンスを生成する。
func NewAIAdviceService(
	adviceRepo repository.AIAdviceRepositoryInterface,
	convRepo repository.AIConversationRepositoryInterface,
	goalRepo repository.LearningGoalRepositoryInterface,
	logRepo repository.LearningLogRepositoryInterface,
	roadmapRepo repository.RoadmapRepositoryInterface,
	githubRepo repository.GitHubRepositoryInterface,
	resourceRepo repository.LearningResourceRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	llmClient LLMClientInterface,
) *AIAdviceService {
	return &AIAdviceService{
		adviceRepo:   adviceRepo,
		convRepo:     convRepo,
		goalRepo:     goalRepo,
		logRepo:      logRepo,
		roadmapRepo:  roadmapRepo,
		githubRepo:   githubRepo,
		resourceRepo: resourceRepo,
		userRepo:     userRepo,
		llmClient:    llmClient,
	}
}

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

	ctx.goals, err = s.goalRepo.GetByUserID(userID)
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

	ctx.resources, err = s.resourceRepo.FindByUserID(userID, true)
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
	// GitHub TypeScript多用 + React目標なし
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
	// GitHubトップ言語が目標にない
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
	// 直近7日間の平均学習時間が60分以上
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

	// ルール12: Q&A参加（優先度5: Info）
	// ユーザーモデルのフィールドから直接判断はできないので
	// 今後のデータで判断。現在は簡易的にチェック。

	// 優先度順にソート（昇順: 1が最優先）
	sort.Slice(advices, func(i, j int) bool {
		if advices[i].Priority != advices[j].Priority {
			return advices[i].Priority < advices[j].Priority
		}
		return i < j // 同じ優先度なら挿入順維持
	})

	return advices
}

// GetAdvice はキャッシュ済みのアドバイスを取得する。
func (s *AIAdviceService) GetAdvice(userID uint, limit int) ([]model.AIAdvice, error) {
	return s.adviceRepo.FindByUserID(userID, limit)
}

// MarkAsRead はアドバイスを既読にする。
func (s *AIAdviceService) MarkAsRead(id, userID uint) error {
	return s.adviceRepo.MarkAsRead(id, userID)
}

// IsLLMAvailable はLLMクライアントが設定されているかどうかを返す。
func (s *AIAdviceService) IsLLMAvailable() bool {
	return s.llmClient != nil
}

// GetDailyChatRemaining は本日の残りチャット回数を返す。
func (s *AIAdviceService) GetDailyChatRemaining(userID uint) (int, error) {
	count, err := s.convRepo.CountTodayMessages(userID)
	if err != nil {
		return 0, err
	}
	remaining := DailyChatLimit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// Chat はLLMとの会話を行い、結果を保存して返す。
// conversationID が0の場合は新規会話を作成する。
func (s *AIAdviceService) Chat(userID uint, message string, conversationID uint) (*model.AIConversation, error) {
	// LLM利用可否チェック
	if s.llmClient == nil {
		return nil, ErrLLMNotConfigured
	}

	// レート制限チェック
	count, err := s.convRepo.CountTodayMessages(userID)
	if err != nil {
		return nil, err
	}
	if count >= int64(DailyChatLimit) {
		return nil, ErrRateLimitExceeded
	}

	// ユーザーコンテキスト収集（プロンプト用）
	goals, _ := s.goalRepo.GetByUserID(userID)
	streak, _ := s.logRepo.GetStreakInfo(userID)
	roadmaps, _ := s.roadmapRepo.GetByUserID(userID)
	langStats, _ := s.githubRepo.GetLanguageStats(userID)

	// 会話管理
	var conv *model.AIConversation
	if conversationID > 0 {
		conv, err = s.convRepo.FindConversationByID(conversationID)
		if err != nil {
			return nil, ErrNotFound
		}
		if conv.UserID != userID {
			return nil, ErrForbidden
		}
	} else {
		// 新規会話作成
		title := message
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		conv = &model.AIConversation{
			UserID: userID,
			Title:  title,
		}
		if err := s.convRepo.CreateConversation(conv); err != nil {
			return nil, err
		}
	}

	// ユーザーメッセージ保存
	userMsg := &model.AIMessage{
		ConversationID: conv.ID,
		Role:           model.AIMessageRoleUser,
		Content:        message,
	}
	if err := s.convRepo.AddMessage(userMsg); err != nil {
		return nil, err
	}

	// システムプロンプト構築
	systemPrompt := buildSystemPrompt(goals, streak, roadmaps, langStats)

	// 過去のメッセージを取得
	var chatMessages []ChatMessage
	chatMessages = append(chatMessages, ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	if conversationID > 0 && conv.Messages != nil {
		for _, m := range conv.Messages {
			if m.Role == model.AIMessageRoleSystem {
				continue
			}
			chatMessages = append(chatMessages, ChatMessage{
				Role:    string(m.Role),
				Content: m.Content,
			})
		}
	}

	chatMessages = append(chatMessages, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// LLM API呼び出し
	resp, err := s.llmClient.Complete(chatMessages)
	if err != nil {
		return nil, err
	}

	// アシスタントメッセージ保存
	assistantMsg := &model.AIMessage{
		ConversationID: conv.ID,
		Role:           model.AIMessageRoleAssistant,
		Content:        resp.Content,
		TokensUsed:     resp.TokensUsed,
	}
	if err := s.convRepo.AddMessage(assistantMsg); err != nil {
		return nil, err
	}

	// 会話を再取得して返す
	conv, err = s.convRepo.FindConversationByID(conv.ID)
	if err != nil {
		// フォールバック: 手動でメッセージを追加
		conv.Messages = append(conv.Messages, *userMsg, *assistantMsg)
	}

	return conv, nil
}

// DeleteConversation は会話を削除する。所有者チェックを行う。
func (s *AIAdviceService) DeleteConversation(id, userID uint) error {
	conv, err := s.convRepo.FindConversationByID(id)
	if err != nil {
		return ErrNotFound
	}
	if conv.UserID != userID {
		return ErrForbidden
	}
	return s.convRepo.DeleteConversation(id, userID)
}

// GetConversations はユーザーの会話一覧を取得する。
func (s *AIAdviceService) GetConversations(userID uint, limit, offset int) ([]model.AIConversation, error) {
	return s.convRepo.FindConversationsByUserID(userID, limit, offset)
}

// GetConversation は会話詳細を取得する。
func (s *AIAdviceService) GetConversation(id, userID uint) (*model.AIConversation, error) {
	conv, err := s.convRepo.FindConversationByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if conv.UserID != userID {
		return nil, ErrForbidden
	}
	return conv, nil
}

// buildSystemPrompt はユーザーコンテキストからシステムプロンプトを構築する。
func buildSystemPrompt(
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
