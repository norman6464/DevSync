package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// AIAdviceUseCases は AIAdviceHandler が依存する AI 機能の usecase をまとめる。
type AIAdviceUseCases struct {
	Generate           *usecase.GenerateAIAdviceUseCase
	MarkAsRead         *usecase.MarkAIAdviceAsReadUseCase
	Unread             *usecase.GetUnreadAIAdviceUseCase
	DailyChatRemaining *usecase.GetDailyChatRemainingUseCase
	Chat               *usecase.ChatWithAIUseCase
	ListConversations  *usecase.ListAIConversationsUseCase
	GetConversation    *usecase.GetAIConversationUseCase
	DeleteConversation *usecase.DeleteAIConversationUseCase
}

// AIAdviceHandler はAIアドバイス関連のHTTPハンドラ。
type AIAdviceHandler struct {
	uc AIAdviceUseCases
}

// NewAIAdviceHandler は新しいAIAdviceHandlerインスタンスを生成する。
func NewAIAdviceHandler(uc AIAdviceUseCases) *AIAdviceHandler {
	return &AIAdviceHandler{uc: uc}
}

// GetAdvice はルールベースアドバイスを取得する。
// LLM利用可否と本日の残りチャット回数も返す。
func (h *AIAdviceHandler) GetAdvice(c *gin.Context) {
	userID := c.GetUint("userID")

	// ルールエンジンを実行してリアルタイムにアドバイス生成
	advices := h.uc.Generate.Execute(c.Request.Context(), userID)

	// LLM利用可否
	llmAvailable := h.uc.Chat.IsAvailable()

	// 本日の残りチャット回数
	remaining := 0
	if llmAvailable {
		var err error
		remaining, err = h.uc.DailyChatRemaining.Execute(c.Request.Context(), userID)
		if err != nil {
			remaining = 0
		}
	}

	respondOK(c, dto.AIAdviceResponse{
		Advices:            advices,
		LLMAvailable:       llmAvailable,
		DailyChatRemaining: remaining,
	})
}

// MarkAsRead はアドバイスを既読にする。
func (h *AIAdviceHandler) MarkAsRead(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.uc.MarkAsRead.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("marked as read"))
}

// Chat はLLMとのチャットメッセージを送信する。
func (h *AIAdviceHandler) Chat(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.AIChatRequest](c)
	if input == nil {
		return
	}

	conv, err := h.uc.Chat.Execute(c.Request.Context(), userID, input.Message, input.ConversationID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conv)
}

// DeleteConversation は会話を削除する。
func (h *AIAdviceHandler) DeleteConversation(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.uc.DeleteConversation.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("conversation deleted"))
}

// GetConversations は会話履歴一覧を取得する。
func (h *AIAdviceHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")

	limit, offset := parseLimitOffset(c)

	conversations, err := h.uc.ListConversations.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conversations)
}

// GetConversation は会話詳細を取得する。
func (h *AIAdviceHandler) GetConversation(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	conv, err := h.uc.GetConversation.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conv)
}

// GetUnreadAdvice は未読のアドバイスを取得する。
func (h *AIAdviceHandler) GetUnreadAdvice(c *gin.Context) {
	userID := c.GetUint("userID")

	advices, err := h.uc.Unread.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(advices))
}
