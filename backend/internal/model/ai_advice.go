package model

import "time"

// AdviceType はAIアドバイスの種別を表す型。
type AdviceType string

// AIアドバイスの種別定数群。
const (
	AdviceTypeStreakRecovery AdviceType = "streak_recovery" // ストリーク途切れ回復
	AdviceTypeStalledRoadmap AdviceType = "stalled_roadmap" // ロードマップ停滞
	AdviceTypeGoalOverdue    AdviceType = "goal_overdue"    // 目標期限超過
	AdviceTypeTechSuggestion AdviceType = "tech_suggestion" // 技術提案
	AdviceTypeGoalSuggestion AdviceType = "goal_suggestion" // 目標提案
	AdviceTypeDifficultyUp   AdviceType = "difficulty_up"   // 難易度アップ
	AdviceTypePraise         AdviceType = "praise"          // 称賛
	AdviceTypeGeneral        AdviceType = "general"         // 一般
)

// AdvicePriority はAIアドバイスの優先度を表す型（1=Critical〜5=Info）。
type AdvicePriority int

// AIアドバイスの優先度定数群。
const (
	AdvicePriorityCritical AdvicePriority = 1 // 最重要
	AdvicePriorityHigh     AdvicePriority = 2 // 高
	AdvicePriorityMedium   AdvicePriority = 3 // 中
	AdvicePriorityLow      AdvicePriority = 4 // 低
	AdvicePriorityInfo     AdvicePriority = 5 // 情報
)

// AIAdvice はルールエンジンが生成するパーソナライズされた学習アドバイスを表す。
// TitleKey/MessageKeyはi18nキーで、フロントエンドで翻訳される。
// Paramsはテンプレート補間用のJSONパラメータ。
type AIAdvice struct {
	ID         uint           `json:"id"`
	UserID     uint           `json:"user_id"`
	Type       AdviceType     `json:"type"`
	Priority   AdvicePriority `json:"priority"`
	TitleKey   string         `json:"title_key"`
	MessageKey string         `json:"message_key"`
	Params     string         `json:"params"`
	ActionURL  string         `json:"action_url"`
	IsRead     bool           `json:"is_read"`
	ExpiresAt  *time.Time     `json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// AIMessageRole はAIメッセージの送信者を表す型。
type AIMessageRole string

// AIメッセージロール定数群。
const (
	AIMessageRoleUser      AIMessageRole = "user"
	AIMessageRoleAssistant AIMessageRole = "assistant"
	AIMessageRoleSystem    AIMessageRole = "system"
)

// AIConversation はユーザーとLLMの会話セッションを表す。
type AIConversation struct {
	ID        uint        `json:"id"`
	UserID    uint        `json:"user_id"`
	Title     string      `json:"title"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Messages  []AIMessage `json:"messages,omitempty"`
}

// AIMessage はAI会話内の個別メッセージを表す。
type AIMessage struct {
	ID             uint          `json:"id"`
	ConversationID uint          `json:"conversation_id"`
	Role           AIMessageRole `json:"role"`
	Content        string        `json:"content"`
	TokensUsed     int           `json:"tokens_used"`
	CreatedAt      time.Time     `json:"created_at"`
}

// ChatMessage は LLM に送信する 1 メッセージを表す。
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse は LLM からの応答を表す。
type ChatResponse struct {
	Content    string `json:"content"`
	TokensUsed int    `json:"tokens_used"`
}
