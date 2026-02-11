package service

import (
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
