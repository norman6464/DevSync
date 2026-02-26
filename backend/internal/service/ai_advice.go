package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// AIAdviceService はルールベース推薦エンジンとLLMチャットオーケストレーションを提供する。
// ルールエンジンは ai_rule_engine.go、LLMチャットは ai_chat.go に分離されている。
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

// GetAdvice はキャッシュ済みのアドバイスを取得する。
func (s *AIAdviceService) GetAdvice(userID uint, limit int) ([]model.AIAdvice, error) {
	return s.adviceRepo.FindByUserID(userID, limit)
}

// GetUnreadAdvice は未読のアドバイスを優先度順で取得する。
func (s *AIAdviceService) GetUnreadAdvice(userID uint) ([]model.AIAdvice, error) {
	return s.adviceRepo.FindUnreadByUserID(userID)
}

// MarkAsRead はアドバイスを既読にする。
func (s *AIAdviceService) MarkAsRead(id, userID uint) error {
	if err := s.adviceRepo.MarkAsRead(id, userID); err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "アドバイスが見つかりません", err)
	}
	return nil
}

// IsLLMAvailable はLLMクライアントが設定されているかどうかを返す。
func (s *AIAdviceService) IsLLMAvailable() bool {
	return s.llmClient != nil
}

// GetDailyChatRemaining は本日の残りチャット回数を返す。
func (s *AIAdviceService) GetDailyChatRemaining(userID uint) (int, error) {
	count, err := s.convRepo.CountTodayMessages(userID)
	if err != nil {
		return 0, domain.NewError(domain.ErrCodeDatabase, "チャット回数の取得に失敗しました", err)
	}
	remaining := DailyChatLimit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// findAndCheckConversationOwnership は会話を取得し、指定ユーザーが所有者かを検証する。
func (s *AIAdviceService) findAndCheckConversationOwnership(id, userID uint) (*model.AIConversation, error) {
	conv, err := s.convRepo.FindConversationByID(id)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "会話が見つかりません", err)
	}
	if conv.UserID != userID {
		return nil, domain.NewError(domain.ErrCodeForbidden, "この会話にアクセスする権限がありません", nil)
	}
	return conv, nil
}

// DeleteConversation は会話を削除する。所有者チェックを行う。
func (s *AIAdviceService) DeleteConversation(id, userID uint) error {
	if _, err := s.findAndCheckConversationOwnership(id, userID); err != nil {
		return err
	}
	return s.convRepo.DeleteConversation(id, userID)
}

// GetConversations はユーザーの会話一覧を取得する。
func (s *AIAdviceService) GetConversations(userID uint, limit, offset int) ([]model.AIConversation, error) {
	return s.convRepo.FindConversationsByUserID(userID, limit, offset)
}

// GetConversation は会話詳細を取得する。
func (s *AIAdviceService) GetConversation(id, userID uint) (*model.AIConversation, error) {
	return s.findAndCheckConversationOwnership(id, userID)
}
