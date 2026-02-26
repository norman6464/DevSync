package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupChatMocks はChat用の共通モック設定（ユーザーコンテキスト取得）を行う。
func setupChatMocks(deps *aiAdviceTestDeps) {
	deps.goalRepo.On("GetByUserID", mock.Anything, 100, 0).
		Return([]model.LearningGoal{}, int64(0), nil).Maybe()
	deps.logRepo.On("GetStreakInfo", mock.Anything).
		Return((*model.StreakInfo)(nil), nil).Maybe()
	deps.roadmapRepo.On("GetByUserID", mock.Anything, 100, 0).
		Return([]model.Roadmap{}, int64(0), nil).Maybe()
	deps.githubRepo.On("GetLanguageStats", mock.Anything).
		Return([]model.GitHubLanguageStat{}, nil).Maybe()
}

func TestChat_Success_NewConversation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.convRepo.On("CreateConversation", mock.MatchedBy(func(c *model.AIConversation) bool {
		return c.UserID == 1 && c.Title == "テストメッセージ"
	})).Run(func(args mock.Arguments) {
		conv := args.Get(0).(*model.AIConversation)
		conv.ID = 100
	}).Return(nil)
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleUser && m.Content == "テストメッセージ"
	})).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "AIの回答です",
		TokensUsed: 50,
	}, nil)
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleAssistant && m.Content == "AIの回答です"
	})).Return(nil)
	deps.convRepo.On("FindConversationByID", uint(100)).Return(&model.AIConversation{
		ID:     100,
		UserID: 1,
		Title:  "テストメッセージ",
	}, nil)

	conv, err := svc.Chat(1, "テストメッセージ", 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	assert.Equal(t, uint(100), conv.ID)
}

func TestChat_Success_ExistingConversation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(2), nil)
	deps.convRepo.On("FindConversationByID", uint(50)).Return(&model.AIConversation{
		ID:     50,
		UserID: 1,
		Title:  "既存会話",
		Messages: []model.AIMessage{
			{Role: model.AIMessageRoleUser, Content: "前のメッセージ"},
			{Role: model.AIMessageRoleAssistant, Content: "前の回答"},
		},
	}, nil)
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleUser
	})).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "新しい回答",
		TokensUsed: 30,
	}, nil)
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleAssistant
	})).Return(nil)

	conv, err := svc.Chat(1, "新しい質問", 50)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
}

func TestChat_RateLimitJustUnder(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(DailyChatLimit-1), nil)
	deps.convRepo.On("CreateConversation", mock.Anything).Run(func(args mock.Arguments) {
		conv := args.Get(0).(*model.AIConversation)
		conv.ID = 1
	}).Return(nil)
	deps.convRepo.On("AddMessage", mock.Anything).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{Content: "OK", TokensUsed: 10}, nil)
	deps.convRepo.On("FindConversationByID", mock.Anything).Return(&model.AIConversation{ID: 1, UserID: 1}, nil)

	conv, err := svc.Chat(1, "テスト", 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
}

func TestChat_ConversationNotFound(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.convRepo.On("FindConversationByID", uint(999)).Return(nil, errors.New("not found"))

	conv, err := svc.Chat(1, "テスト", 999)
	assert.Nil(t, conv)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestChat_ConversationForbidden(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.convRepo.On("FindConversationByID", uint(50)).Return(&model.AIConversation{
		ID:     50,
		UserID: 99, // 他のユーザーの会話
	}, nil)

	conv, err := svc.Chat(1, "テスト", 50)
	assert.Nil(t, conv)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestChat_LLMError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.convRepo.On("CreateConversation", mock.Anything).Run(func(args mock.Arguments) {
		conv := args.Get(0).(*model.AIConversation)
		conv.ID = 1
	}).Return(nil)
	deps.convRepo.On("AddMessage", mock.Anything).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(nil, errors.New("LLM API error"))

	conv, err := svc.Chat(1, "テスト", 0)
	assert.Nil(t, conv)
	assert.Error(t, err)
}

func TestChat_TitleTruncation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)
	setupChatMocks(deps)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.convRepo.On("CreateConversation", mock.MatchedBy(func(c *model.AIConversation) bool {
		return len(c.Title) <= 53 && strings.HasSuffix(c.Title, "...")
	})).Run(func(args mock.Arguments) {
		conv := args.Get(0).(*model.AIConversation)
		conv.ID = 1
	}).Return(nil)
	deps.convRepo.On("AddMessage", mock.Anything).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{Content: "OK", TokensUsed: 10}, nil)
	deps.convRepo.On("FindConversationByID", mock.Anything).Return(&model.AIConversation{ID: 1, UserID: 1}, nil)

	longMsg := strings.Repeat("a", 100)
	conv, err := svc.Chat(1, longMsg, 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
}

// ============================================================
// buildSystemPrompt テスト
// ============================================================

func TestBuildSystemPrompt_Empty(t *testing.T) {
	prompt := buildSystemPrompt(nil, nil, nil, nil)
	assert.Contains(t, prompt, "学習アドバイザー")
	assert.NotContains(t, prompt, "学習ストリーク")
	assert.NotContains(t, prompt, "学習目標")
}

func TestBuildSystemPrompt_WithStreak(t *testing.T) {
	streak := &model.StreakInfo{
		CurrentStreak: 10,
		LongestStreak: 30,
		TotalDays:     100,
	}
	prompt := buildSystemPrompt(nil, streak, nil, nil)
	assert.Contains(t, prompt, "10日連続")
	assert.Contains(t, prompt, "最長 30日")
	assert.Contains(t, prompt, "合計 100日")
}

func TestBuildSystemPrompt_WithGoals(t *testing.T) {
	goals := []model.LearningGoal{
		{Title: "Go学習", Status: model.GoalStatusActive, Progress: 50},
		{Title: "React習得", Status: model.GoalStatusCompleted, Progress: 100},
	}
	prompt := buildSystemPrompt(goals, nil, nil, nil)
	assert.Contains(t, prompt, "Go学習")
	assert.Contains(t, prompt, "進行中")
	assert.Contains(t, prompt, "完了")
}

func TestBuildSystemPrompt_WithRoadmaps(t *testing.T) {
	roadmaps := []model.Roadmap{
		{Title: "バックエンド", Progress: 60, CompletedStepCount: 3, StepCount: 5},
	}
	prompt := buildSystemPrompt(nil, nil, roadmaps, nil)
	assert.Contains(t, prompt, "バックエンド")
	assert.Contains(t, prompt, "3/5ステップ")
}

func TestBuildSystemPrompt_WithLangStats(t *testing.T) {
	stats := []model.GitHubLanguageStat{
		{Language: "Go", RepoCount: 10},
		{Language: "TypeScript", RepoCount: 5},
	}
	prompt := buildSystemPrompt(nil, nil, nil, stats)
	assert.Contains(t, prompt, "Go")
	assert.Contains(t, prompt, "10リポジトリ")
}

func TestBuildSystemPrompt_LangStatsMaxFive(t *testing.T) {
	stats := make([]model.GitHubLanguageStat, 8)
	for i := range stats {
		stats[i] = model.GitHubLanguageStat{Language: "Lang" + string(rune('A'+i)), RepoCount: 1}
	}
	prompt := buildSystemPrompt(nil, nil, nil, stats)
	assert.Contains(t, prompt, "LangA")
	assert.Contains(t, prompt, "LangE")
	assert.NotContains(t, prompt, "LangF")
}
