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
