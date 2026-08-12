package handler

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// mockAIAdviceRepo は usecase/repository.AIAdviceRepository のモック。
type mockAIAdviceRepo struct{ mock.Mock }

func (m *mockAIAdviceRepo) CreateBatch(ctx context.Context, advices []*model.AIAdvice) error {
	return m.Called(ctx, advices).Error(0)
}

func (m *mockAIAdviceRepo) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AIAdvice, error) {
	args := m.Called(ctx, userID, limit)
	a, _ := args.Get(0).([]model.AIAdvice)
	return a, args.Error(1)
}

func (m *mockAIAdviceRepo) FindUnreadByUserID(ctx context.Context, userID uint) ([]model.AIAdvice, error) {
	args := m.Called(ctx, userID)
	a, _ := args.Get(0).([]model.AIAdvice)
	return a, args.Error(1)
}

func (m *mockAIAdviceRepo) MarkAsRead(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockAIAdviceRepo) DeleteByUserID(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

// mockAIConversationRepo は usecase/repository.AIConversationRepository のモック。
type mockAIConversationRepo struct{ mock.Mock }

func (m *mockAIConversationRepo) CreateConversation(ctx context.Context, conv *model.AIConversation) error {
	return m.Called(ctx, conv).Error(0)
}

func (m *mockAIConversationRepo) FindConversationsByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.AIConversation, error) {
	args := m.Called(ctx, userID, limit, offset)
	c, _ := args.Get(0).([]model.AIConversation)
	return c, args.Error(1)
}

func (m *mockAIConversationRepo) FindConversationByID(ctx context.Context, id uint) (*model.AIConversation, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.AIConversation)
	return c, args.Error(1)
}

func (m *mockAIConversationRepo) AddMessage(ctx context.Context, msg *model.AIMessage) error {
	return m.Called(ctx, msg).Error(0)
}

func (m *mockAIConversationRepo) CountTodayMessages(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockAIConversationRepo) DeleteConversation(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}

// mockLLMClient は usecase/repository.LLMClient のモック。
type mockLLMClient struct{ mock.Mock }

func (m *mockLLMClient) Complete(ctx context.Context, messages []model.ChatMessage) (*model.ChatResponse, error) {
	args := m.Called(ctx, messages)
	r, _ := args.Get(0).(*model.ChatResponse)
	return r, args.Error(1)
}

// aiPorts は AI 機能の usecase に注入した port モックをまとめる。
type aiPorts struct {
	Advices       *mockAIAdviceRepo
	Conversations *mockAIConversationRepo
	LLM           *mockLLMClient
	Goals         *MockLearningGoalRepository
	Logs          *mockLearningLogRepo
	Roadmaps      *mockRoadmapRepo
	GitHub        *mockGitHubRepo
	Resources     *mockLearningResourceRepo
	Users         *mockUserPort
}
