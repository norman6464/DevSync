package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// AIAdviceRepository はルールエンジンが生成したアドバイスの永続化に対する、usecase 側が要求する契約。
type AIAdviceRepository interface {
	CreateBatch(ctx context.Context, advices []*model.AIAdvice) error
	// FindByUserID は優先度の高い順・作成の新しい順にアドバイスを返す。
	FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AIAdvice, error)
	// FindUnreadByUserID は未読のアドバイスを優先度順で返す。
	FindUnreadByUserID(ctx context.Context, userID uint) ([]model.AIAdvice, error)
	// MarkAsRead は本人のアドバイス 1 件を既読にする。
	MarkAsRead(ctx context.Context, id, userID uint) error
	// DeleteByUserID は指定ユーザーのアドバイスをすべて削除する。
	DeleteByUserID(ctx context.Context, userID uint) error
}

// AIConversationRepository は AI チャットの会話・メッセージの永続化に対する契約。
type AIConversationRepository interface {
	CreateConversation(ctx context.Context, conv *model.AIConversation) error
	// FindConversationsByUserID は会話を更新の新しい順に返す。
	FindConversationsByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.AIConversation, error)
	// FindConversationByID は会話をメッセージ付きで返す。存在しなければ (nil, nil) を返す。
	FindConversationByID(ctx context.Context, id uint) (*model.AIConversation, error)
	AddMessage(ctx context.Context, msg *model.AIMessage) error
	// CountTodayMessages は当日のユーザー発言数を返す（レート制限に使う）。
	CountTodayMessages(ctx context.Context, userID uint) (int64, error)
	// DeleteConversation は本人の会話をメッセージごと削除する。
	DeleteConversation(ctx context.Context, id, userID uint) error
}

// LLMClient は LLM への問い合わせに対する契約。実装は adapter/external に置く。
type LLMClient interface {
	// Complete は会話履歴を渡して応答を得る。
	Complete(ctx context.Context, messages []model.ChatMessage) (*model.ChatResponse, error)
}
