package service

import (
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// setupRuleEngineService はテスト用のAIAdviceServiceとモックを生成するヘルパー。
func setupRuleEngineService() (
	*AIAdviceService,
	*MockLearningLogRepository,
	*MockLearningGoalRepository,
	*MockRoadmapRepository,
	*MockGitHubRepository,
	*MockLearningResourceRepository,
	*MockUserRepository,
) {
	adviceRepo := new(MockAIAdviceRepository)
	convRepo := new(MockAIConversationRepository)
	goalRepo := new(MockLearningGoalRepository)
	logRepo := new(MockLearningLogRepository)
	roadmapRepo := new(MockRoadmapRepository)
	githubRepo := new(MockGitHubRepository)
	resourceRepo := new(MockLearningResourceRepository)
	userRepo := new(MockUserRepository)

	svc := NewAIAdviceService(adviceRepo, convRepo, goalRepo, logRepo, roadmapRepo, githubRepo, resourceRepo, userRepo, nil)
	return svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo
}

// setupDefaultMocks は全リポジトリに空/デフォルトの戻り値をセットする。
func setupDefaultMocks(
	logRepo *MockLearningLogRepository,
	goalRepo *MockLearningGoalRepository,
	roadmapRepo *MockRoadmapRepository,
	githubRepo *MockGitHubRepository,
	resourceRepo *MockLearningResourceRepository,
	userRepo *MockUserRepository,
) {
	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
	roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
	githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
	logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
	resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{}, int64(0), nil)
	userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)
}

func TestGenerateAdvice_StreakBroken(t *testing.T) {
	t.Run("ストリークが途切れた場合Criticalアドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 0, TotalDays: 10, LongestStreak: 5}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypeStreakRecovery {
				found = true
				assert.Equal(t, model.AdvicePriorityCritical, a.Priority)
				assert.Equal(t, "advice.streakBroken", a.TitleKey)
				break
			}
		}
		assert.True(t, found, "ストリーク途切れアドバイスが返されるべき")
	})

	t.Run("ストリークが継続中の場合はアドバイスなし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 3, TotalDays: 10}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, model.AdviceTypeStreakRecovery, a.Type)
		}
	})
}

func TestGenerateAdvice_StalledRoadmap(t *testing.T) {
	t.Run("7日以上更新のないアクティブロードマップにHighアドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		staleTime := time.Now().Add(-10 * 24 * time.Hour)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
			{Status: model.RoadmapStatusActive, StepCount: 5, CompletedStepCount: 2, UpdatedAt: staleTime},
		}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypeStalledRoadmap {
				found = true
				assert.Equal(t, model.AdvicePriorityHigh, a.Priority)
				break
			}
		}
		assert.True(t, found, "ロードマップ停滞アドバイスが返されるべき")
	})
}

func TestGenerateAdvice_GoalOverdue(t *testing.T) {
	t.Run("期限超過のアクティブ目標にHighアドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		pastDate := time.Now().Add(-48 * time.Hour)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Status: model.GoalStatusActive, TargetDate: &pastDate, Title: "Go学習"},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1, ActiveGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypeGoalOverdue {
				found = true
				assert.Equal(t, model.AdvicePriorityHigh, a.Priority)
				break
			}
		}
		assert.True(t, found, "目標期限超過アドバイスが返されるべき")
	})
}

func TestGenerateAdvice_ReactSuggestion(t *testing.T) {
	t.Run("TypeScript使用かつReact目標なしの場合React提案を返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Title: "Go学習", Status: model.GoalStatusActive},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "TypeScript", Bytes: 50000},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypeTechSuggestion && a.TitleKey == "advice.suggestReact" {
				found = true
				assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
				break
			}
		}
		assert.True(t, found, "React提案アドバイスが返されるべき")
	})

	t.Run("既にReact目標がある場合は提案なし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Title: "React学習", Status: model.GoalStatusActive},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "TypeScript", Bytes: 50000},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			if a.Type == model.AdviceTypeTechSuggestion {
				assert.NotEqual(t, "advice.suggestReact", a.TitleKey)
			}
		}
	})
}

func TestGenerateAdvice_NoGoals(t *testing.T) {
	t.Run("目標未設定かつGitHubデータありの場合目標提案を返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 0}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Go", Bytes: 20000},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.noGoals" {
				found = true
				assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
				break
			}
		}
		assert.True(t, found, "目標未設定アドバイスが返されるべき")
	})
}

func TestGenerateAdvice_NoRoadmap(t *testing.T) {
	t.Run("目標ありロードマップなしの場合ロードマップ提案を返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Title: "Go学習", Status: model.GoalStatusActive},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.noRoadmap" {
				found = true
				assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
				break
			}
		}
		assert.True(t, found, "ロードマップ未作成アドバイスが返されるべき")
	})
}

func TestGenerateAdvice_DifficultyUp(t *testing.T) {
	t.Run("完了目標3以上かつ平均進捗70超の場合難易度UPアドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{CompletedGoals: 5, AverageProgress: 80}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypeDifficultyUp {
				found = true
				assert.Equal(t, model.AdvicePriorityLow, a.Priority)
				break
			}
		}
		assert.True(t, found, "難易度UPアドバイスが返されるべき")
	})
}

func TestGenerateAdvice_Praise(t *testing.T) {
	t.Run("週平均60分以上の学習で称賛アドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		// 過去7日に毎日90分ずつ学習 → 合計630分 → 平均90分/日
		now := time.Now()
		logs := make([]model.LearningLog, 7)
		for i := 0; i < 7; i++ {
			logs[i] = model.LearningLog{Duration: 90}
			logs[i].CreatedAt = now.Add(-time.Duration(i) * 24 * time.Hour)
		}
		logRepo.On("GetByUserID", uint(1)).Return(logs, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.Type == model.AdviceTypePraise {
				found = true
				assert.Equal(t, model.AdvicePriorityLow, a.Priority)
				break
			}
		}
		assert.True(t, found, "称賛アドバイスが返されるべき")
	})
}

func TestGenerateAdvice_ExploreResources(t *testing.T) {
	t.Run("リソース3未満の場合探索アドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}}, int64(2), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.exploreResources" {
				found = true
				assert.Equal(t, model.AdvicePriorityInfo, a.Priority)
				break
			}
		}
		assert.True(t, found, "リソース探索アドバイスが返されるべき")
	})

	t.Run("リソース3以上の場合探索アドバイスなし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, "advice.exploreResources", a.TitleKey)
		}
	})
}

func TestGenerateAdvice_LowStudyTime(t *testing.T) {
	t.Run("週平均15分未満の学習で学習時間増加アドバイスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		// 過去7日に合計50分の学習（平均7分/日 → 15分未満）
		now := time.Now()
		logs := []model.LearningLog{
			{Duration: 25},
			{Duration: 25},
		}
		logs[0].CreatedAt = now.Add(-1 * 24 * time.Hour)
		logs[1].CreatedAt = now.Add(-3 * 24 * time.Hour)
		logRepo.On("GetByUserID", uint(1)).Return(logs, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.suggestMoreTime" {
				found = true
				assert.Equal(t, model.AdvicePriorityLow, a.Priority)
				assert.Equal(t, "/learning-logs", a.ActionURL)
				break
			}
		}
		assert.True(t, found, "学習時間増加アドバイスが返されるべき")
	})

	t.Run("週平均15分以上60分未満の場合は学習時間アドバイスなし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		// 過去7日に合計210分（平均30分/日 → 15〜60の間）
		now := time.Now()
		logs := []model.LearningLog{
			{Duration: 70},
			{Duration: 70},
			{Duration: 70},
		}
		logs[0].CreatedAt = now.Add(-1 * 24 * time.Hour)
		logs[1].CreatedAt = now.Add(-3 * 24 * time.Hour)
		logs[2].CreatedAt = now.Add(-5 * 24 * time.Hour)
		logRepo.On("GetByUserID", uint(1)).Return(logs, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, "advice.suggestMoreTime", a.TitleKey,
				"週平均15〜60分の場合は学習時間増加アドバイスを返さない")
			assert.NotEqual(t, model.AdviceTypePraise, a.Type,
				"週平均60分未満の場合は称賛アドバイスを返さない")
		}
	})
}

func TestGenerateAdvice_TechSuggestionFromTopLanguage(t *testing.T) {
	t.Run("TypeScript以外のトップ言語で目標なしの場合技術提案を返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Title: "Docker学習", Status: model.GoalStatusActive},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Go", Bytes: 80000, RepoCount: 5},
			{Language: "Python", Bytes: 30000, RepoCount: 2},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.suggestFromGithub" {
				found = true
				assert.Equal(t, model.AdvicePriorityMedium, a.Priority)
				assert.Contains(t, a.Params, "Go")
				break
			}
		}
		assert.True(t, found, "GitHub言語ベースの技術提案が返されるべき")
	})

	t.Run("トップ言語に対応する目標がある場合は提案なし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{
			{Title: "Go言語マスター", Status: model.GoalStatusActive},
		}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 1}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Go", Bytes: 80000, RepoCount: 5},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, "advice.suggestFromGithub", a.TitleKey,
				"既にトップ言語の目標がある場合は提案しない")
		}
	})

	t.Run("トップ言語のバイト数が10000以下の場合は提案なし", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Go", Bytes: 5000, RepoCount: 1},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, "advice.suggestFromGithub", a.TitleKey,
				"バイト数が少ない場合は技術提案しない")
		}
	})
}

func TestGenerateAdvice_TechSuggestionTopLangNotFirst(t *testing.T) {
	t.Run("最初の要素が最大でない場合もトップ言語で技術提案を返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Python", Bytes: 5000, RepoCount: 1},
			{Language: "Go", Bytes: 80000, RepoCount: 5},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		found := false
		for _, a := range advices {
			if a.TitleKey == "advice.suggestFromGithub" {
				found = true
				assert.Contains(t, a.Params, "Go")
				break
			}
		}
		assert.True(t, found, "2番目の要素が最大でも正しくトップ言語として技術提案を返すべき")
	})
}

func TestGenerateAdvice_StalledRoadmapEdgeCases(t *testing.T) {
	t.Run("完了済みロードマップは停滞アドバイスを返さない", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		staleTime := time.Now().Add(-10 * 24 * time.Hour)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
			{Status: model.RoadmapStatusCompleted, StepCount: 5, CompletedStepCount: 5, UpdatedAt: staleTime},
		}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, model.AdviceTypeStalledRoadmap, a.Type,
				"完了済みロードマップには停滞アドバイスを返さない")
		}
	})

	t.Run("全ステップ完了のアクティブロードマップは停滞アドバイスを返さない", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		staleTime := time.Now().Add(-10 * 24 * time.Hour)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
			{Status: model.RoadmapStatusActive, StepCount: 5, CompletedStepCount: 5, UpdatedAt: staleTime},
		}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}, {}, {}}, int64(3), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		for _, a := range advices {
			assert.NotEqual(t, model.AdviceTypeStalledRoadmap, a.Type,
				"全ステップ完了のロードマップには停滞アドバイスを返さない")
		}
	})
}

func TestGenerateAdvice_PrioritySorting(t *testing.T) {
	t.Run("アドバイスは優先度順にソートされる", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		// ストリーク途切れ(Critical=1) + 目標未設定+GitHub有(Medium=3) + リソース少(Info=5)
		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{CurrentStreak: 0, TotalDays: 10, LongestStreak: 5}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{TotalGoals: 0}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{
			{Language: "Go", Bytes: 20000},
		}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{{}}, int64(1), nil)
		userRepo.On("FindByID", uint(1)).Return(&model.User{}, nil)

		advices := svc.GenerateAdvice(1)

		assert.GreaterOrEqual(t, len(advices), 3, "少なくとも3つのアドバイスが返されるべき")

		// 優先度が昇順であることを検証
		for i := 1; i < len(advices); i++ {
			assert.LessOrEqual(t, int(advices[i-1].Priority), int(advices[i].Priority),
				"アドバイスは優先度順にソートされるべき")
		}

		// 最初のアドバイスはCritical（ストリーク途切れ）
		assert.Equal(t, model.AdvicePriorityCritical, advices[0].Priority)
	})
}

func TestGenerateAdvice_EmptyContext(t *testing.T) {
	t.Run("全てデフォルト値の場合リソース探索のみ返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()
		setupDefaultMocks(logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo)

		advices := svc.GenerateAdvice(1)

		// リソース0件 → リソース探索アドバイスのみ
		assert.Len(t, advices, 1)
		assert.Equal(t, "advice.exploreResources", advices[0].TitleKey)
	})
}

func TestGenerateAdvice_ContextCollectionError_StreakInfo(t *testing.T) {
	t.Run("StreakInfo取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, _, _, _, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(nil, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_Goals(t *testing.T) {
	t.Run("Goals取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, _, _, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_GoalStats(t *testing.T) {
	t.Run("GoalStats取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, _, _, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(nil, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_Roadmaps(t *testing.T) {
	t.Run("Roadmap取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, _, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_LanguageStats(t *testing.T) {
	t.Run("LanguageStats取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_Logs(t *testing.T) {
	t.Run("Logs取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, _, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_Resources(t *testing.T) {
	t.Run("Resources取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, _ := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{}, int64(0), assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}

func TestGenerateAdvice_ContextCollectionError_User(t *testing.T) {
	t.Run("User取得エラー時は空スライスを返す", func(t *testing.T) {
		svc, logRepo, goalRepo, roadmapRepo, githubRepo, resourceRepo, userRepo := setupRuleEngineService()

		logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
		goalRepo.On("GetByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
		goalRepo.On("GetStats", uint(1)).Return(&model.LearningGoalStats{}, nil)
		roadmapRepo.On("GetByUserID", uint(1)).Return([]model.Roadmap{}, nil)
		githubRepo.On("GetLanguageStats", uint(1)).Return([]model.GitHubLanguageStat{}, nil)
		logRepo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)
		resourceRepo.On("FindByUserID", uint(1), true, 100, 0).Return([]model.LearningResource{}, int64(0), nil)
		userRepo.On("FindByID", uint(1)).Return(nil, assert.AnError)

		advices := svc.GenerateAdvice(1)
		assert.Empty(t, advices)
	})
}
