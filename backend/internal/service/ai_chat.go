package service

import (
	"fmt"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// DailyChatLimit は1日あたりのLLMチャット回数制限。
const DailyChatLimit = 5

// MaxChatMessageLength はチャットメッセージの最大文字数。
const MaxChatMessageLength = 5000

// Chat はLLMとの会話を行い、結果を保存して返す。
// conversationID が0の場合は新規会話を作成する。
func (s *AIAdviceService) Chat(userID uint, message string, conversationID uint) (*model.AIConversation, error) {
	// メッセージバリデーション
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "メッセージを入力してください", nil)
	}
	if len(trimmed) > MaxChatMessageLength {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "メッセージは5000文字以内で入力してください", nil)
	}
	message = trimmed

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
	goals, _, _ := s.goalRepo.GetByUserID(userID, 100, 0)
	streak, _ := s.logRepo.GetStreakInfo(userID)
	roadmaps, _, _ := s.roadmapRepo.GetByUserID(userID, 100, 0)
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
		if runes := []rune(title); len(runes) > 50 {
			title = string(runes[:50]) + "..."
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
