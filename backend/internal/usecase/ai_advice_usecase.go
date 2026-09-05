package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

const (
	// DailyChatLimit は 1 日あたりの LLM チャット回数の上限。
	DailyChatLimit = 5
	// MaxChatMessageLength はチャットメッセージの最大文字数。
	MaxChatMessageLength = 5000
	// aiConversationTitleMaxRunes は会話タイトルに使う先頭文字数。
	aiConversationTitleMaxRunes = 50
)

// ErrLLMNotConfigured / ErrChatRateLimitExceeded は移行前と同じ domain の sentinel を使う。
var (
	// ErrLLMNotConfigured は LLM クライアントが未設定のときに返す（503）。
	ErrLLMNotConfigured = domain.ErrServiceUnavailable
	// ErrChatRateLimitExceeded は 1 日のチャット上限に達したときに返す（429）。
	ErrChatRateLimitExceeded = domain.ErrRateLimitExceeded
)

// GetAIAdviceUseCase はキャッシュ済みのアドバイスを取得する。
type GetAIAdviceUseCase struct {
	advices repository.AIAdviceRepository
}

// NewGetAIAdviceUseCase は GetAIAdviceUseCase を生成する。
func NewGetAIAdviceUseCase(advices repository.AIAdviceRepository) *GetAIAdviceUseCase {
	return &GetAIAdviceUseCase{advices: advices}
}

// Execute は優先度順にアドバイスを返す。
func (uc *GetAIAdviceUseCase) Execute(ctx context.Context, userID uint, limit int) ([]model.AIAdvice, error) {
	return uc.advices.FindByUserID(ctx, userID, limit)
}

// GetUnreadAIAdviceUseCase は未読のアドバイスを取得する。
type GetUnreadAIAdviceUseCase struct {
	advices repository.AIAdviceRepository
}

// NewGetUnreadAIAdviceUseCase は GetUnreadAIAdviceUseCase を生成する。
func NewGetUnreadAIAdviceUseCase(advices repository.AIAdviceRepository) *GetUnreadAIAdviceUseCase {
	return &GetUnreadAIAdviceUseCase{advices: advices}
}

// Execute は未読のアドバイスを優先度順で返す。
func (uc *GetUnreadAIAdviceUseCase) Execute(ctx context.Context, userID uint) ([]model.AIAdvice, error) {
	return uc.advices.FindUnreadByUserID(ctx, userID)
}

// MarkAIAdviceAsReadUseCase はアドバイスを既読にする。
type MarkAIAdviceAsReadUseCase struct {
	advices repository.AIAdviceRepository
}

// NewMarkAIAdviceAsReadUseCase は MarkAIAdviceAsReadUseCase を生成する。
func NewMarkAIAdviceAsReadUseCase(advices repository.AIAdviceRepository) *MarkAIAdviceAsReadUseCase {
	return &MarkAIAdviceAsReadUseCase{advices: advices}
}

// Execute は本人のアドバイスを既読にする。対象が無ければ 404 を返す。
func (uc *MarkAIAdviceAsReadUseCase) Execute(ctx context.Context, id, userID uint) error {
	if err := uc.advices.MarkAsRead(ctx, id, userID); err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "アドバイスが見つかりません", err)
	}
	return nil
}

// RefreshAIAdviceUseCase はルールエンジンでアドバイスを生成し直して保存する。
type RefreshAIAdviceUseCase struct {
	advices  repository.AIAdviceRepository
	generate *GenerateAIAdviceUseCase
}

// NewRefreshAIAdviceUseCase は RefreshAIAdviceUseCase を生成する。
func NewRefreshAIAdviceUseCase(
	advices repository.AIAdviceRepository,
	generate *GenerateAIAdviceUseCase,
) *RefreshAIAdviceUseCase {
	return &RefreshAIAdviceUseCase{advices: advices, generate: generate}
}

// Execute は既存のアドバイスを削除してから生成し直し、保存した結果を返す。
func (uc *RefreshAIAdviceUseCase) Execute(ctx context.Context, userID uint) ([]model.AIAdvice, error) {
	if err := uc.advices.DeleteByUserID(ctx, userID); err != nil {
		return nil, domain.NewError(domain.ErrCodeDatabase, "既存アドバイスの削除に失敗しました", err)
	}

	generated := uc.generate.Execute(ctx, userID)
	if len(generated) == 0 {
		return []model.AIAdvice{}, nil
	}

	pointers := make([]*model.AIAdvice, 0, len(generated))
	for i := range generated {
		pointers = append(pointers, &generated[i])
	}
	if err := uc.advices.CreateBatch(ctx, pointers); err != nil {
		return nil, domain.NewError(domain.ErrCodeDatabase, "アドバイスの保存に失敗しました", err)
	}
	return generated, nil
}

// GetDailyChatRemainingUseCase は本日の残りチャット回数を返す。
type GetDailyChatRemainingUseCase struct {
	conversations repository.AIConversationRepository
}

// NewGetDailyChatRemainingUseCase は GetDailyChatRemainingUseCase を生成する。
func NewGetDailyChatRemainingUseCase(conversations repository.AIConversationRepository) *GetDailyChatRemainingUseCase {
	return &GetDailyChatRemainingUseCase{conversations: conversations}
}

// Execute は残り回数を返す（0 未満にはならない）。
func (uc *GetDailyChatRemainingUseCase) Execute(ctx context.Context, userID uint) (int, error) {
	count, err := uc.conversations.CountTodayMessages(ctx, userID)
	if err != nil {
		return 0, domain.NewError(domain.ErrCodeDatabase, "チャット回数の取得に失敗しました", err)
	}
	remaining := DailyChatLimit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// ListAIConversationsUseCase は会話一覧を取得する。
type ListAIConversationsUseCase struct {
	conversations repository.AIConversationRepository
}

// NewListAIConversationsUseCase は ListAIConversationsUseCase を生成する。
func NewListAIConversationsUseCase(conversations repository.AIConversationRepository) *ListAIConversationsUseCase {
	return &ListAIConversationsUseCase{conversations: conversations}
}

// Execute は会話を更新の新しい順に返す。
func (uc *ListAIConversationsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.AIConversation, error) {
	return uc.conversations.FindConversationsByUserID(ctx, userID, limit, offset)
}

// GetAIConversationUseCase は会話の詳細を取得する。
type GetAIConversationUseCase struct {
	conversations repository.AIConversationRepository
}

// NewGetAIConversationUseCase は GetAIConversationUseCase を生成する。
func NewGetAIConversationUseCase(conversations repository.AIConversationRepository) *GetAIConversationUseCase {
	return &GetAIConversationUseCase{conversations: conversations}
}

// Execute は本人の会話をメッセージ付きで返す。
func (uc *GetAIConversationUseCase) Execute(ctx context.Context, id, userID uint) (*model.AIConversation, error) {
	return findOwnedConversation(ctx, uc.conversations, id, userID)
}

// DeleteAIConversationUseCase は会話を削除する。
type DeleteAIConversationUseCase struct {
	conversations repository.AIConversationRepository
}

// NewDeleteAIConversationUseCase は DeleteAIConversationUseCase を生成する。
func NewDeleteAIConversationUseCase(conversations repository.AIConversationRepository) *DeleteAIConversationUseCase {
	return &DeleteAIConversationUseCase{conversations: conversations}
}

// Execute は所有権を検証したうえで会話を削除する。
func (uc *DeleteAIConversationUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedConversation(ctx, uc.conversations, id, userID); err != nil {
		return err
	}
	return uc.conversations.DeleteConversation(ctx, id, userID)
}

// findOwnedConversation は会話を取得し、userID が所有者であることを検証する。
// 会話が無ければ 404、他人の会話なら 403 を返す。
func findOwnedConversation(
	ctx context.Context,
	conversations repository.AIConversationRepository,
	id, userID uint,
) (*model.AIConversation, error) {
	conv, err := conversations.FindConversationByID(ctx, id)
	if err != nil || conv == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "会話が見つかりません", err)
	}
	if conv.UserID != userID {
		return nil, domain.NewError(domain.ErrCodeForbidden, "この会話にアクセスする権限がありません", nil)
	}
	return conv, nil
}

// ChatWithAIUseCase は LLM と会話し、やり取りを保存する。
type ChatWithAIUseCase struct {
	conversations repository.AIConversationRepository
	llm           repository.LLMClient
	prompt        *BuildAIChatPromptUseCase
}

// NewChatWithAIUseCase は ChatWithAIUseCase を生成する。
// llm が nil の場合、チャットは利用不可として扱う。
func NewChatWithAIUseCase(
	conversations repository.AIConversationRepository,
	llm repository.LLMClient,
	prompt *BuildAIChatPromptUseCase,
) *ChatWithAIUseCase {
	return &ChatWithAIUseCase{conversations: conversations, llm: llm, prompt: prompt}
}

// IsAvailable は LLM が設定されているかどうかを返す。
func (uc *ChatWithAIUseCase) IsAvailable() bool {
	return uc.llm != nil
}

// Execute はメッセージを検証し、LLM の応答を保存して会話を返す。
// conversationID が 0 の場合は新規会話を作成する。
func (uc *ChatWithAIUseCase) Execute(ctx context.Context, userID uint, message string, conversationID uint) (*model.AIConversation, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "メッセージを入力してください", nil)
	}
	if len(message) > MaxChatMessageLength {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "メッセージは5000文字以内で入力してください", nil)
	}
	if uc.llm == nil {
		return nil, ErrLLMNotConfigured
	}

	count, err := uc.conversations.CountTodayMessages(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= int64(DailyChatLimit) {
		return nil, ErrChatRateLimitExceeded
	}

	systemPrompt := uc.prompt.Execute(ctx, userID)

	conv, err := uc.resolveConversation(ctx, userID, message, conversationID)
	if err != nil {
		return nil, err
	}

	userMsg := &model.AIMessage{
		ConversationID: conv.ID,
		Role:           model.AIMessageRoleUser,
		Content:        message,
	}
	if err := uc.conversations.AddMessage(ctx, userMsg); err != nil {
		return nil, err
	}

	messages := []model.ChatMessage{{Role: "system", Content: systemPrompt}}
	if conversationID > 0 {
		for _, m := range conv.Messages {
			if m.Role == model.AIMessageRoleSystem {
				continue
			}
			messages = append(messages, model.ChatMessage{Role: string(m.Role), Content: m.Content})
		}
	}
	messages = append(messages, model.ChatMessage{Role: "user", Content: message})

	resp, err := uc.llm.Complete(ctx, messages)
	if err != nil {
		return nil, err
	}

	assistantMsg := &model.AIMessage{
		ConversationID: conv.ID,
		Role:           model.AIMessageRoleAssistant,
		Content:        resp.Content,
		TokensUsed:     resp.TokensUsed,
	}
	if err := uc.conversations.AddMessage(ctx, assistantMsg); err != nil {
		return nil, err
	}

	// 保存後の状態で返す。取得できない場合は今回のやり取りを手元で足して返す。
	// 会話の保存自体は既に成功しているため、再取得の失敗だけで全体を失敗にはしない。
	refreshed, err := uc.conversations.FindConversationByID(ctx, conv.ID)
	if err != nil || refreshed == nil {
		conv.Messages = append(conv.Messages, *userMsg, *assistantMsg)
		return conv, nil //nolint:nilerr // 再取得失敗時は手元の値へ意図的にフォールバックする
	}
	return refreshed, nil
}

// resolveConversation は既存の会話を取得するか、新しい会話を作成する。
func (uc *ChatWithAIUseCase) resolveConversation(ctx context.Context, userID uint, message string, conversationID uint) (*model.AIConversation, error) {
	if conversationID > 0 {
		conv, err := uc.conversations.FindConversationByID(ctx, conversationID)
		if err != nil || conv == nil {
			return nil, domain.ErrNotFound
		}
		if conv.UserID != userID {
			return nil, domain.ErrForbidden
		}
		return conv, nil
	}

	conv := &model.AIConversation{UserID: userID, Title: conversationTitle(message)}
	if err := uc.conversations.CreateConversation(ctx, conv); err != nil {
		return nil, err
	}
	return conv, nil
}

// conversationTitle は最初のメッセージから会話タイトルを作る。
func conversationTitle(message string) string {
	runes := []rune(message)
	if len(runes) > aiConversationTitleMaxRunes {
		return string(runes[:aiConversationTitleMaxRunes]) + "..."
	}
	return message
}
