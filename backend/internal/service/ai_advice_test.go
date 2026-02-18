package service

import (
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// テストヘルパー
// ============================================================

// newTestAIAdviceService はAIAdviceServiceのテスト用インスタンスを生成する。
type aiAdviceTestDeps struct {
	adviceRepo   *MockAIAdviceRepository
	convRepo     *MockAIConversationRepository
	goalRepo     *MockLearningGoalRepository
	logRepo      *MockLearningLogRepository
	roadmapRepo  *MockRoadmapRepository
	githubRepo   *MockGitHubRepository
	resourceRepo *MockLearningResourceRepository
	userRepo     *MockUserRepository
	llmClient    *MockLLMClient
}

func newTestAIAdviceService(withLLM bool) (*AIAdviceService, *aiAdviceTestDeps) {
	deps := &aiAdviceTestDeps{
		adviceRepo:   new(MockAIAdviceRepository),
		convRepo:     new(MockAIConversationRepository),
		goalRepo:     new(MockLearningGoalRepository),
		logRepo:      new(MockLearningLogRepository),
		roadmapRepo:  new(MockRoadmapRepository),
		githubRepo:   new(MockGitHubRepository),
		resourceRepo: new(MockLearningResourceRepository),
		userRepo:     new(MockUserRepository),
		llmClient:    new(MockLLMClient),
	}

	var llm LLMClientInterface
	if withLLM {
		llm = deps.llmClient
	}

	svc := NewAIAdviceService(
		deps.adviceRepo,
		deps.convRepo,
		deps.goalRepo,
		deps.logRepo,
		deps.roadmapRepo,
		deps.githubRepo,
		deps.resourceRepo,
		deps.userRepo,
		llm,
	)
	return svc, deps
}

// ============================================================
// ルールエンジン テスト（10ケース）
// ============================================================

func TestRuleEngine_StreakBroken(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// ストリーク0 + 学習歴あり → streak_recovery アドバイス
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{
		CurrentStreak: 0,
		TotalDays:     10,
	}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeStreakRecovery {
			found = true
			assert.Equal(t, model.AdvicePriorityCritical, a.Priority)
		}
	}
	assert.True(t, found, "streak_recovery アドバイスが生成されるべき")
}

func TestRuleEngine_RoadmapStalled(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// ロードマップのステップが7日以上未更新 → stalled_roadmap
	stalledTime := time.Now().Add(-8 * 24 * time.Hour)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
		{
			ID: 1, Title: "Go学習", Status: model.RoadmapStatusActive,
			StepCount: 5, CompletedStepCount: 2, Progress: 40,
			Steps: []model.RoadmapStep{
				{ID: 1, IsCompleted: false, UpdatedAt: stalledTime},
			},
			UpdatedAt: stalledTime,
		},
	}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeStalledRoadmap {
			found = true
			assert.Equal(t, model.AdvicePriorityHigh, a.Priority)
		}
	}
	assert.True(t, found, "stalled_roadmap アドバイスが生成されるべき")
}

func TestRuleEngine_GoalOverdue(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 目標期限超過 → goal_overdue
	pastDate := time.Now().Add(-3 * 24 * time.Hour)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "React学習", Status: model.GoalStatusActive, TargetDate: &pastDate, Progress: 50},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{ActiveGoals: 1}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeGoalOverdue {
			found = true
			assert.Equal(t, model.AdvicePriorityHigh, a.Priority)
		}
	}
	assert.True(t, found, "goal_overdue アドバイスが生成されるべき")
}

func TestRuleEngine_TechGapReact(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// GitHub TypeScript多用 + React目標なし → React提案
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "TypeScript基礎", Category: model.GoalCategoryLanguage, Status: model.GoalStatusActive},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{ActiveGoals: 1}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
		{Language: "TypeScript", Bytes: 100000, RepoCount: 5},
		{Language: "JavaScript", Bytes: 50000, RepoCount: 3},
	}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeTechSuggestion {
			found = true
			assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
		}
	}
	assert.True(t, found, "tech_suggestion アドバイスが生成されるべき")
}

func TestRuleEngine_NoGoalsWithGitHub(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 目標0件 + GitHub活動あり → 目標設定提案
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 0}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
		{Language: "Go", Bytes: 80000, RepoCount: 3},
	}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeGoalSuggestion && a.TitleKey == "advice.noGoals" {
			found = true
			assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
		}
	}
	assert.True(t, found, "goal_suggestion (noGoals) アドバイスが生成されるべき")
}

func TestRuleEngine_HighCompletionRate(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 達成率高い → 難易度UP提案
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 5, TotalDays: 30}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Status: model.GoalStatusActive, Progress: 80},
		{Status: model.GoalStatusCompleted, Progress: 100},
		{Status: model.GoalStatusCompleted, Progress: 100},
		{Status: model.GoalStatusCompleted, Progress: 100},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{
		CompletedGoals: 3, ActiveGoals: 1, AverageProgress: 80,
	}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
		{ID: 1, Status: model.RoadmapStatusActive},
	}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{
		{}, {}, {},
	}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypeDifficultyUp {
			found = true
			assert.Equal(t, model.AdvicePriorityLow, a.Priority)
		}
	}
	assert.True(t, found, "difficulty_up アドバイスが生成されるべき")
}

func TestRuleEngine_ConsistentLearner(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 7日平均60分/日以上 → 称賛
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 10, TotalDays: 30}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Status: model.GoalStatusActive, Progress: 50},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{ActiveGoals: 1}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
		{ID: 1, Status: model.RoadmapStatusActive},
	}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)

	// 直近7日間に毎日90分の学習ログ
	logs := make([]model.LearningLog, 7)
	for i := 0; i < 7; i++ {
		logs[i] = model.LearningLog{
			Duration:  90,
			CreatedAt: time.Now().Add(-time.Duration(i) * 24 * time.Hour),
		}
	}
	deps.logRepo.On("GetByUserID", uint(1)).Return(logs, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{
		{}, {}, {},
	}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.Type == model.AdviceTypePraise {
			found = true
			assert.Equal(t, model.AdvicePriorityLow, a.Priority)
		}
	}
	assert.True(t, found, "praise アドバイスが生成されるべき")
}

func TestRuleEngine_NoRoadmap(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 目標あり + ロードマップなし → ロードマップ提案
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "Go学習", Status: model.GoalStatusActive},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{ActiveGoals: 1}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	found := false
	for _, a := range advices {
		if a.TitleKey == "advice.noRoadmap" {
			found = true
			assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
		}
	}
	assert.True(t, found, "noRoadmap アドバイスが生成されるべき")
}

func TestRuleEngine_PriorityOrdering(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 複数ルールに該当 → 優先度順に並ぶこと
	pastDate := time.Now().Add(-3 * 24 * time.Hour)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{
		CurrentStreak: 0, TotalDays: 10,
	}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "React", Status: model.GoalStatusActive, TargetDate: &pastDate, Progress: 30},
	}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{ActiveGoals: 1}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	assert.True(t, len(advices) >= 2, "2つ以上のアドバイスが生成されるべき")
	// 優先度は昇順（1が最優先）
	for i := 1; i < len(advices); i++ {
		assert.True(t, advices[i-1].Priority <= advices[i].Priority,
			"アドバイスが優先度順に並んでいるべき: %d <= %d",
			advices[i-1].Priority, advices[i].Priority)
	}
}

func TestRuleEngine_NewUserNoData(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// データなし → 空または初期ウェルカム
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	deps.resourceRepo.On("FindByUserID", uint(1), true).Return([]model.LearningResource{}, nil)
	deps.userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

	advices := svc.GenerateAdvice(1)

	// 新規ユーザーはstreakBrokenにはならない（TotalDays=0）
	for _, a := range advices {
		assert.NotEqual(t, model.AdviceTypeStreakRecovery, a.Type,
			"新規ユーザーにstreak_recoveryは不要")
	}
}

// ============================================================
// LLMチャット テスト（4ケース）
// ============================================================

func TestChat_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	// レート制限OK
	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(2), nil)
	// ユーザーコンテキスト収集
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "Go学習", Status: model.GoalStatusActive, Progress: 50},
	}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 5}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
		{Language: "Go", Bytes: 50000},
	}, nil)
	// 会話作成
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(nil)
	// ユーザーメッセージ保存
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	// LLM呼び出し
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "Go学習の次のステップとしてWeb開発をおすすめします。",
		TokensUsed: 150,
	}, nil)
	// 会話再取得
	deps.convRepo.On("FindConversationByID", mock.AnythingOfType("uint")).Return(&model.AIConversation{
		UserID: 1,
		Title:  "Goの次に何を学ぶべきですか？",
	}, nil)

	conv, err := svc.Chat(1, "Goの次に何を学ぶべきですか？", 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	assert.Equal(t, uint(1), conv.UserID)
	deps.convRepo.AssertExpectations(t)
	deps.llmClient.AssertExpectations(t)
}

func TestChat_RateLimitExceeded(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	// 本日5回超 → エラー
	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(5), nil)

	conv, err := svc.Chat(1, "質問", 0)
	assert.ErrorIs(t, err, ErrRateLimitExceeded)
	assert.Nil(t, conv)
}

func TestChat_LLMNotConfigured(t *testing.T) {
	svc, _ := newTestAIAdviceService(false) // LLMなし

	conv, err := svc.Chat(1, "質問", 0)
	assert.ErrorIs(t, err, ErrLLMNotConfigured)
	assert.Nil(t, conv)
}

func TestIsLLMAvailable(t *testing.T) {
	// LLMあり
	svcWithLLM, _ := newTestAIAdviceService(true)
	assert.True(t, svcWithLLM.IsLLMAvailable())

	// LLMなし
	svcNoLLM, _ := newTestAIAdviceService(false)
	assert.False(t, svcNoLLM.IsLLMAvailable())
}

// ============================================================
// 会話削除 テスト（3ケース）
// ============================================================

func TestDeleteConversation_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 会話が存在し、所有者が一致 → 正常削除
	deps.convRepo.On("FindConversationByID", uint(1)).Return(&model.AIConversation{
		UserID: 1,
		Title:  "テスト会話",
	}, nil)
	deps.convRepo.On("DeleteConversation", uint(1), uint(1)).Return(nil)

	err := svc.DeleteConversation(1, 1)
	assert.NoError(t, err)
	deps.convRepo.AssertExpectations(t)
}

func TestDeleteConversation_NotFound(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 会話が存在しない → ErrNotFound
	deps.convRepo.On("FindConversationByID", uint(999)).Return(nil, ErrNotFound)

	err := svc.DeleteConversation(999, 1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDeleteConversation_Forbidden(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	// 他人の会話を削除しようとした → ErrForbidden
	deps.convRepo.On("FindConversationByID", uint(1)).Return(&model.AIConversation{
		UserID: 2, // 別ユーザー
		Title:  "他人の会話",
	}, nil)

	err := svc.DeleteConversation(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

// ============================================================
// GetAdvice テスト
// ============================================================

func TestGetAdvice_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.adviceRepo.On("FindByUserID", uint(1), 10).Return([]model.AIAdvice{
		{UserID: 1, TitleKey: "advice.test"},
	}, nil)

	advices, err := svc.GetAdvice(1, 10)
	assert.NoError(t, err)
	assert.Len(t, advices, 1)
}

func TestGetAdvice_Error(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.adviceRepo.On("FindByUserID", uint(1), 10).Return([]model.AIAdvice{}, ErrNotFound)

	_, err := svc.GetAdvice(1, 10)
	assert.Error(t, err)
}

// ============================================================
// MarkAsRead テスト
// ============================================================

func TestMarkAsRead_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.adviceRepo.On("MarkAsRead", uint(5), uint(1)).Return(nil)

	err := svc.MarkAsRead(5, 1)
	assert.NoError(t, err)
}

func TestMarkAsRead_Error(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.adviceRepo.On("MarkAsRead", uint(999), uint(1)).Return(ErrNotFound)

	err := svc.MarkAsRead(999, 1)
	assert.ErrorIs(t, err, ErrNotFound)
}

// ============================================================
// GetDailyChatRemaining テスト
// ============================================================

func TestGetDailyChatRemaining_HasRemaining(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(2), nil)

	remaining, err := svc.GetDailyChatRemaining(1)
	assert.NoError(t, err)
	assert.Equal(t, 3, remaining) // DailyChatLimit(5) - 2 = 3
}

func TestGetDailyChatRemaining_AtLimit(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(5), nil)

	remaining, err := svc.GetDailyChatRemaining(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, remaining)
}

func TestGetDailyChatRemaining_OverLimit(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(10), nil)

	remaining, err := svc.GetDailyChatRemaining(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, remaining) // 負にならない
}

func TestGetDailyChatRemaining_Error(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), ErrNotFound)

	_, err := svc.GetDailyChatRemaining(1)
	assert.Error(t, err)
}

// ============================================================
// GetConversations テスト
// ============================================================

func TestGetConversations_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.convRepo.On("FindConversationsByUserID", uint(1), 20, 0).Return([]model.AIConversation{
		{UserID: 1, Title: "会話1"},
		{UserID: 1, Title: "会話2"},
	}, nil)

	convs, err := svc.GetConversations(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, convs, 2)
}

func TestGetConversations_Error(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.convRepo.On("FindConversationsByUserID", uint(1), 20, 0).Return([]model.AIConversation{}, ErrNotFound)

	_, err := svc.GetConversations(1, 20, 0)
	assert.Error(t, err)
}

// ============================================================
// GetConversation テスト
// ============================================================

func TestGetConversation_Success(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.convRepo.On("FindConversationByID", uint(1)).Return(&model.AIConversation{
		UserID: 1,
		Title:  "テスト会話",
	}, nil)

	conv, err := svc.GetConversation(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "テスト会話", conv.Title)
}

func TestGetConversation_NotFound(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.convRepo.On("FindConversationByID", uint(999)).Return(nil, ErrNotFound)

	conv, err := svc.GetConversation(999, 1)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, conv)
}

func TestGetConversation_Forbidden(t *testing.T) {
	svc, deps := newTestAIAdviceService(false)

	deps.convRepo.On("FindConversationByID", uint(1)).Return(&model.AIConversation{
		UserID: 2, // 別ユーザー
		Title:  "他人の会話",
	}, nil)

	conv, err := svc.GetConversation(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, conv)
}

// ============================================================
// Chat 追加テスト（既存会話・Forbidden）
// ============================================================

func TestChat_ExistingConversation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(1), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)

	// 既存会話を返す（所有者一致）
	deps.convRepo.On("FindConversationByID", uint(10)).Return(&model.AIConversation{
		UserID:   1,
		Title:    "既存会話",
		Messages: []model.AIMessage{},
	}, nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "回答です",
		TokensUsed: 100,
	}, nil)
	// 会話再取得
	deps.convRepo.On("FindConversationByID", mock.AnythingOfType("uint")).Return(&model.AIConversation{
		UserID: 1,
		Title:  "既存会話",
	}, nil)

	conv, err := svc.Chat(1, "追加の質問", 10)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
}

func TestChat_ForbiddenConversation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(1), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)

	// 他人の会話
	deps.convRepo.On("FindConversationByID", uint(10)).Return(&model.AIConversation{
		UserID: 2,
		Title:  "他人の会話",
	}, nil)

	conv, err := svc.Chat(1, "不正アクセス", 10)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, conv)
}

// ============================================================
// buildSystemPrompt テスト
// ============================================================

func TestBuildSystemPrompt(t *testing.T) {
	t.Run("全データありの場合プロンプトにコンテキストが含まれる", func(t *testing.T) {
		goals := []model.LearningGoal{
			{Title: "Go学習", Status: model.GoalStatusActive, Progress: 50},
			{Title: "React学習", Status: model.GoalStatusCompleted, Progress: 100},
		}
		streak := &model.StreakInfo{CurrentStreak: 5, LongestStreak: 10, TotalDays: 30}
		roadmaps := []model.Roadmap{
			{Title: "Goロードマップ", Progress: 60, CompletedStepCount: 3, StepCount: 5},
		}
		langStats := []model.GitHubLanguageStat{
			{Language: "Go", RepoCount: 5},
			{Language: "TypeScript", RepoCount: 3},
		}

		prompt := buildSystemPrompt(goals, streak, roadmaps, langStats)

		assert.Contains(t, prompt, "5日連続")
		assert.Contains(t, prompt, "Go学習")
		assert.Contains(t, prompt, "完了")
		assert.Contains(t, prompt, "Goロードマップ")
		assert.Contains(t, prompt, "Go (5リポジトリ)")
	})

	t.Run("データなしの場合基本プロンプトのみ", func(t *testing.T) {
		prompt := buildSystemPrompt(nil, nil, nil, nil)

		assert.Contains(t, prompt, "学習アドバイザー")
		assert.NotContains(t, prompt, "【学習ストリーク】")
		assert.NotContains(t, prompt, "【学習目標】")
	})

	t.Run("GitHub言語統計が5件超の場合上位5件のみ表示", func(t *testing.T) {
		langStats := []model.GitHubLanguageStat{
			{Language: "Go", RepoCount: 10},
			{Language: "TypeScript", RepoCount: 8},
			{Language: "Python", RepoCount: 6},
			{Language: "Rust", RepoCount: 4},
			{Language: "Java", RepoCount: 3},
			{Language: "C++", RepoCount: 2},
			{Language: "Ruby", RepoCount: 1},
		}

		prompt := buildSystemPrompt(nil, nil, nil, langStats)

		assert.Contains(t, prompt, "Go (10リポジトリ)")
		assert.Contains(t, prompt, "Java (3リポジトリ)")
		assert.NotContains(t, prompt, "C++")
		assert.NotContains(t, prompt, "Ruby")
	})
}

// ============================================================
// Chat エラーパス テスト
// ============================================================

func TestChat_CountTodayMessagesError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), ErrNotFound)

	conv, err := svc.Chat(1, "質問", 0)
	assert.Error(t, err)
	assert.Nil(t, conv)
}

func TestChat_CreateConversationError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(ErrNotFound)

	conv, err := svc.Chat(1, "質問", 0)
	assert.Error(t, err)
	assert.Nil(t, conv)
}

func TestChat_AddUserMessageError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(ErrNotFound).Once()

	conv, err := svc.Chat(1, "質問", 0)
	assert.Error(t, err)
	assert.Nil(t, conv)
}

func TestChat_LLMCompleteError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return((*ChatResponse)(nil), ErrNotFound)

	conv, err := svc.Chat(1, "質問", 0)
	assert.Error(t, err)
	assert.Nil(t, conv)
}

func TestChat_AddAssistantMessageError(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(nil)
	// ユーザーメッセージ成功、アシスタントメッセージ失敗
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleUser
	})).Return(nil)
	deps.convRepo.On("AddMessage", mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleAssistant
	})).Return(ErrNotFound)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "回答",
		TokensUsed: 50,
	}, nil)

	conv, err := svc.Chat(1, "質問", 0)
	assert.Error(t, err)
	assert.Nil(t, conv)
}

func TestChat_ExistingConversationWithMessages(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	// 既存会話にメッセージ履歴あり（systemメッセージはスキップされること）
	deps.convRepo.On("FindConversationByID", uint(10)).Return(&model.AIConversation{
		UserID: 1,
		Title:  "既存会話",
		Messages: []model.AIMessage{
			{Role: model.AIMessageRoleSystem, Content: "system prompt"},
			{Role: model.AIMessageRoleUser, Content: "前の質問"},
			{Role: model.AIMessageRoleAssistant, Content: "前の回答"},
		},
	}, nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	deps.llmClient.On("Complete", mock.MatchedBy(func(msgs []ChatMessage) bool {
		// systemメッセージが先頭、過去のuserとassistantが含まれ、systemはスキップ
		return len(msgs) >= 4 && msgs[0].Role == "system"
	})).Return(&ChatResponse{
		Content:    "回答",
		TokensUsed: 50,
	}, nil)
	deps.convRepo.On("FindConversationByID", mock.AnythingOfType("uint")).Return(&model.AIConversation{
		UserID: 1,
		Title:  "既存会話",
	}, nil)

	conv, err := svc.Chat(1, "追加の質問", 10)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
}

func TestChat_LongTitleTruncation(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	// 会話作成時にタイトルが50文字+...に切り詰められることを検証
	deps.convRepo.On("CreateConversation", mock.MatchedBy(func(c *model.AIConversation) bool {
		return len(c.Title) == 53 && c.Title[50:] == "..."
	})).Return(nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "回答",
		TokensUsed: 50,
	}, nil)
	deps.convRepo.On("FindConversationByID", mock.AnythingOfType("uint")).Return(&model.AIConversation{
		UserID: 1,
	}, nil)

	longMessage := "これは50文字を超える非常に長いメッセージです。テストのために長い文章を書いています。さらに追加のテキスト。"
	conv, err := svc.Chat(1, longMessage, 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	deps.convRepo.AssertExpectations(t)
}

func TestChat_ExistingConversationNotFound(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	// 既存会話が見つからない
	deps.convRepo.On("FindConversationByID", uint(999)).Return((*model.AIConversation)(nil), ErrNotFound)

	conv, err := svc.Chat(1, "質問", 999)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, conv)
}

func TestChat_RefetchFallback(t *testing.T) {
	svc, deps := newTestAIAdviceService(true)

	deps.convRepo.On("CountTodayMessages", uint(1)).Return(int64(0), nil)
	deps.goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	deps.logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	deps.roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	deps.githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	deps.convRepo.On("CreateConversation", mock.AnythingOfType("*model.AIConversation")).Return(nil)
	deps.convRepo.On("AddMessage", mock.AnythingOfType("*model.AIMessage")).Return(nil)
	deps.llmClient.On("Complete", mock.Anything).Return(&ChatResponse{
		Content:    "回答です",
		TokensUsed: 100,
	}, nil)
	// 再取得が失敗 → フォールバックでメッセージを手動追加
	deps.convRepo.On("FindConversationByID", mock.AnythingOfType("uint")).Return(&model.AIConversation{
		UserID: 1,
	}, ErrNotFound)

	conv, err := svc.Chat(1, "質問", 0)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	// フォールバックでユーザーメッセージとアシスタントメッセージが手動追加される
	assert.Len(t, conv.Messages, 2)
	assert.Equal(t, model.AIMessageRoleUser, conv.Messages[0].Role)
	assert.Equal(t, model.AIMessageRoleAssistant, conv.Messages[1].Role)
	assert.Equal(t, "回答です", conv.Messages[1].Content)
}
