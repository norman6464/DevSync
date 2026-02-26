package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestLearningLogService はLearningLogServiceのテスト用インスタンスを生成するヘルパー。
func newTestLearningLogService() (*LearningLogService, *MockLearningLogRepository) {
	repo := new(MockLearningLogRepository)
	svc := NewLearningLogService(repo, nil)
	return svc, repo
}

// newTestLearningLogServiceWithGoalRepo はゴールリポ付きのテスト用インスタンスを生成するヘルパー。
func newTestLearningLogServiceWithGoalRepo() (*LearningLogService, *MockLearningLogRepository, *MockLearningGoalRepository) {
	repo := new(MockLearningLogRepository)
	goalRepo := new(MockLearningGoalRepository)
	svc := NewLearningLogService(repo, goalRepo)
	return svc, repo, goalRepo
}

// ============================================================
// 学習ログ作成テスト
// ============================================================

func TestLearningLogCreate_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60}
	repo.On("Create", log).Return(nil)

	err := svc.Create(log)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習ログ作成バリデーションテスト
// ============================================================

func TestLearningLogCreate_InvalidCategory(t *testing.T) {
	svc, _ := newTestLearningLogService()

	log := &model.LearningLog{
		Title:    "テスト",
		Content:  "テスト内容",
		UserID:   1,
		Category: model.LogCategory("invalid"),
	}

	err := svc.Create(log)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogCreate_NegativeDuration(t *testing.T) {
	svc, _ := newTestLearningLogService()

	log := &model.LearningLog{
		Title:    "テスト",
		Content:  "テスト内容",
		UserID:   1,
		Duration: -10,
	}

	err := svc.Create(log)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogCreate_ExcessiveDuration(t *testing.T) {
	svc, _ := newTestLearningLogService()

	log := &model.LearningLog{
		Title:    "テスト",
		Content:  "テスト内容",
		UserID:   1,
		Duration: 1500,
	}

	err := svc.Create(log)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogCreate_WithSource(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{
		Title:    "ポモドーロ集中",
		Content:  "25分の集中セッション",
		UserID:   1,
		Duration: 25,
		Source:   model.LogSourcePomodoro,
	}
	repo.On("Create", log).Return(nil)

	err := svc.Create(log)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogCreate_InvalidSource(t *testing.T) {
	svc, _ := newTestLearningLogService()

	log := &model.LearningLog{
		Title:    "テスト",
		Content:  "テスト内容",
		UserID:   1,
		Source:   model.LogSource("unknown"),
	}

	err := svc.Create(log)
	assert.ErrorIs(t, err, ErrBadRequest)
}

// ============================================================
// 学習ログ更新テスト
// ============================================================

func TestLearningLogUpdate_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Old", Content: "Old Content", UserID: 1, Duration: 30}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningLog{Title: "New", Duration: 60}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New", result.Title)
	assert.Equal(t, 60, result.Duration)
	assert.Equal(t, "Old Content", result.Content) // 変更なし
	repo.AssertExpectations(t)
}

func TestLearningLogUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningLog{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningLogUpdate_NotFound(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.LearningLog{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習ログ削除テスト
// ============================================================

func TestLearningLogDelete_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1), uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogDelete_Forbidden(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// ストリーク・カレンダーデータテスト
// ============================================================

func TestLearningLogGetStreakInfo_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	streak := &model.StreakInfo{CurrentStreak: 5, LongestStreak: 10}
	repo.On("GetStreakInfo", uint(1)).Return(streak, nil)

	result, err := svc.GetStreakInfo(1)
	assert.NoError(t, err)
	assert.Equal(t, 5, result.CurrentStreak)
	assert.Equal(t, 10, result.LongestStreak)
	repo.AssertExpectations(t)
}

func TestLearningLogGetCalendarData_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	entries := []model.CalendarEntry{{Date: "2024-01-01", Count: 3}}
	repo.On("GetCalendarData", uint(1)).Return(entries, nil)

	result, err := svc.GetCalendarData(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 3, result[0].Count)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestLearningLogGetByID_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{Title: "Go Study", UserID: 1}
	log.ID = 1
	repo.On("FindByID", uint(1)).Return(log, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Go Study", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("FindByID", uint(999)).Return((*model.LearningLog)(nil), errors.New("not found"))

	_, err := svc.GetByID(999, 1)
	assert.Error(t, err)
}

func TestLearningLogGetByID_Forbidden(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{Title: "Go Study", UserID: 1}
	log.ID = 1
	repo.On("FindByID", uint(1)).Return(log, nil)

	result, err := svc.GetByID(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestLearningLogGetByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{{Title: "Go Study"}, {Title: "React Study"}}
	repo.On("GetByUserID", uint(1), 20, 0).Return(logs, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByUserID", uint(1), 20, 0).Return([]model.LearningLog{}, int64(0), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

func TestLearningLogGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByUserID", uint(1), 10, 10).Return([]model.LearningLog{}, int64(15), nil)

	result, total, err := svc.GetByUserID(1, 10, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(15), total)
	repo.AssertExpectations(t)
}

func TestLearningLogUpdate_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()
	existing := &model.LearningLog{Title: "Old", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))
	updates := &model.LearningLog{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLearningLogDelete_NotFound(t *testing.T) {
	svc, repo := newTestLearningLogService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}

func TestLearningLogUpdate_AllFields(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Old", Content: "Old Content", UserID: 1, Duration: 30}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 全フィールドを更新
	updates := &model.LearningLog{
		Title:    "New Title",
		Content:  "New Content",
		Category: model.LogCategoryCoding,
		Duration: 90,
	}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Content", result.Content)
	assert.Equal(t, model.LogCategoryCoding, result.Category)
	assert.Equal(t, 90, result.Duration)
	repo.AssertExpectations(t)
}

func TestLearningLogUpdate_DurationNegative(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Test", Content: "Content", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningLog{Duration: -10}
	result, err := svc.Update(1, 1, updates)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
}

func TestLearningLogUpdate_DurationExceedsMax(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Test", Content: "Content", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningLog{Duration: 1441}
	result, err := svc.Update(1, 1, updates)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
}

func TestLearningLogUpdate_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Test", Content: "Content", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 空白のみのTitleは更新されず、元の値が保持されるべき
	updates := &model.LearningLog{Title: "   "}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Title) // 変更されない
}

func TestLearningLogUpdate_WhitespaceContent(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Test", Content: "Content", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 空白のみのContentは更新されず、元の値が保持されるべき
	updates := &model.LearningLog{Content: "  \t  "}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Content", result.Content) // 変更されない
}

func TestLearningLogUpdate_WhitespaceCategory(t *testing.T) {
	svc, repo := newTestLearningLogService()

	existing := &model.LearningLog{Title: "Test", Content: "Content", Category: "coding", UserID: 1, Duration: 30}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 空白のみのCategoryは更新されず、元の値が保持されるべき
	updates := &model.LearningLog{Category: "   "}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.LogCategory("coding"), result.Category) // 変更されない
}

// --- ExportCSV ---

func TestLearningLogExportCSV_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "Goの勉強", Content: "基礎を学んだ", Category: model.LogCategoryCoding, Duration: 60, CreatedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)},
		{Title: "設計復習", Content: "DDD読んだ", Category: model.LogCategoryReading, Duration: 30, CreatedAt: time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC)},
	}
	repo.On("GetByPeriod", uint(1), 30).Return(logs, nil)

	data, err := svc.ExportCSV(1, 30)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// BOM付きUTF-8ヘッダー確認
	content := string(data)
	assert.True(t, strings.HasPrefix(content, "\xef\xbb\xbf"))
	assert.Contains(t, content, "Goの勉強")
	assert.Contains(t, content, "設計復習")
	repo.AssertExpectations(t)
}

func TestLearningLogExportCSV_EmptyLogs(t *testing.T) {
	svc, repo := newTestLearningLogService()
	repo.On("GetByPeriod", uint(1), 0).Return([]model.LearningLog{}, nil)

	data, err := svc.ExportCSV(1, 0)
	assert.NoError(t, err)
	// ヘッダー行のみ含まれる
	content := string(data)
	assert.Contains(t, content, "日付")
	assert.Contains(t, content, "タイトル")
	repo.AssertExpectations(t)
}

func TestLearningLogExportCSV_NegativeDays(t *testing.T) {
	svc, _ := newTestLearningLogService()

	_, err := svc.ExportCSV(1, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "期間は0以上の値")
}

func TestLearningLogExportCSV_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()
	repo.On("GetByPeriod", uint(1), 7).Return([]model.LearningLog{}, errors.New("db error"))

	_, err := svc.ExportCSV(1, 7)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByCategory テスト
// ============================================================

func TestLearningLogGetByCategory_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "Go勉強", Category: model.LogCategoryCoding, UserID: 1},
	}
	repo.On("GetByCategory", uint(1), "coding").Return(logs, nil)

	result, err := svc.GetByCategory(1, "coding")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, model.LogCategoryCoding, result[0].Category)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByCategory_EmptyResult(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByCategory", uint(1), "meetup").Return([]model.LearningLog{}, nil)

	result, err := svc.GetByCategory(1, "meetup")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByCategory_InvalidCategory(t *testing.T) {
	svc, _ := newTestLearningLogService()

	result, err := svc.GetByCategory(1, "invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "無効なカテゴリ")
}

func TestLearningLogGetByCategory_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByCategory", uint(1), "coding").Return([]model.LearningLog{}, errors.New("db error"))

	result, err := svc.GetByCategory(1, "coding")
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetBySource テスト
// ============================================================

func TestLearningLogGetBySource_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "ポモドーロ学習1", Source: model.LogSourcePomodoro, UserID: 1},
		{Title: "ポモドーロ学習2", Source: model.LogSourcePomodoro, UserID: 1},
	}
	repo.On("GetBySource", uint(1), "pomodoro").Return(logs, nil)

	result, err := svc.GetBySource(1, "pomodoro")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, model.LogSourcePomodoro, result[0].Source)
	repo.AssertExpectations(t)
}

func TestLearningLogGetBySource_InvalidSource(t *testing.T) {
	svc, _ := newTestLearningLogService()

	result, err := svc.GetBySource(1, "invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "無効なソースです")
}

func TestLearningLogGetBySource_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetBySource", uint(1), "pomodoro").Return([]model.LearningLog{}, nil)

	result, err := svc.GetBySource(1, "pomodoro")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningLogGetBySource_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetBySource", uint(1), "manual").Return([]model.LearningLog{}, errors.New("db error"))

	result, err := svc.GetBySource(1, "manual")
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 週間学習時間合計テスト
// ============================================================

func TestLearningLogGetWeeklyDuration_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("SumDurationByPeriod", uint(1), 7).Return(480, nil)

	duration, err := svc.GetWeeklyDuration(1)
	assert.NoError(t, err)
	assert.Equal(t, 480, duration)
	repo.AssertExpectations(t)
}

func TestLearningLogGetWeeklyDuration_ZeroLogs(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("SumDurationByPeriod", uint(1), 7).Return(0, nil)

	duration, err := svc.GetWeeklyDuration(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, duration)
}

func TestLearningLogGetWeeklyDuration_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("SumDurationByPeriod", uint(1), 7).Return(0, errors.New("db error"))

	duration, err := svc.GetWeeklyDuration(1)
	assert.Error(t, err)
	assert.Equal(t, 0, duration)
}

func TestLearningLogFavorite_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{UserID: 1}
	log.ID = 10
	repo.On("FindByID", uint(10)).Return(log, nil)
	repo.On("Update", mock.MatchedBy(func(l *model.LearningLog) bool {
		return l.ID == 10 && l.IsFavorite
	})).Return(nil)

	err := svc.FavoriteLog(10, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogFavorite_Forbidden(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{UserID: 99}
	log.ID = 10
	repo.On("FindByID", uint(10)).Return(log, nil)

	err := svc.FavoriteLog(10, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestLearningLogFavorite_NotFound(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	err := svc.FavoriteLog(99, 1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestLearningLogUnfavorite_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{UserID: 1, IsFavorite: true}
	log.ID = 10
	repo.On("FindByID", uint(10)).Return(log, nil)
	repo.On("Update", mock.MatchedBy(func(l *model.LearningLog) bool {
		return l.ID == 10 && !l.IsFavorite
	})).Return(nil)

	err := svc.UnfavoriteLog(10, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogUnfavorite_Forbidden(t *testing.T) {
	svc, repo := newTestLearningLogService()

	log := &model.LearningLog{UserID: 99}
	log.ID = 10
	repo.On("FindByID", uint(10)).Return(log, nil)

	err := svc.UnfavoriteLog(10, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

// ============================================================
// GetRecentCategories テスト
// ============================================================

func TestLearningLogGetRecentCategories_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	categories := []string{"Go", "React", "TypeScript"}
	repo.On("GetRecentCategories", uint(1), 5).Return(categories, nil)

	result, err := svc.GetRecentCategories(1)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Go", result[0])
}

func TestLearningLogGetRecentCategories_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetRecentCategories", uint(1), 5).Return([]string{}, nil)

	result, err := svc.GetRecentCategories(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestLearningLogGetRecentCategories_Error(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetRecentCategories", uint(1), 5).Return([]string{}, errors.New("db error"))

	_, err := svc.GetRecentCategories(1)
	assert.Error(t, err)
}

// ============================================================
// ExportCSV 追加エッジケーステスト
// ============================================================

func TestLearningLogExportCSV_SpecificPeriod(t *testing.T) {
	svc, repo := newTestLearningLogService()

	now := time.Now()
	logs := []model.LearningLog{
		{
			Title:    "Go基礎",
			Category: "Go",
			Duration: 60,
			Content:  "変数・型・制御構文",
		},
	}
	logs[0].CreatedAt = now

	repo.On("GetByPeriod", uint(1), 7).Return(logs, nil)

	data, err := svc.ExportCSV(1, 7)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	content := string(data)
	assert.Contains(t, content, "Go基礎")
	assert.Contains(t, content, "60")
	assert.Contains(t, content, "変数・型・制御構文")
}

func TestLearningLogExportCSV_BOMPresent(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByPeriod", uint(1), 30).Return([]model.LearningLog{}, nil)

	data, err := svc.ExportCSV(1, 30)
	assert.NoError(t, err)

	// BOM付きUTF-8を確認
	assert.True(t, len(data) >= 3)
	assert.Equal(t, byte(0xEF), data[0])
	assert.Equal(t, byte(0xBB), data[1])
	assert.Equal(t, byte(0xBF), data[2])
}

// ============================================================
// バッチ作成テスト
// ============================================================

func TestLearningLogBatchCreate_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "Go基礎", Content: "変数を学んだ", UserID: 1, Duration: 60, Category: model.LogCategoryCoding},
		{Title: "React入門", Content: "コンポーネント作成", UserID: 1, Duration: 45, Category: model.LogCategoryCourse},
	}

	repo.On("CreateBatch", mock.MatchedBy(func(l []model.LearningLog) bool {
		return len(l) == 2 && l[0].Title == "Go基礎" && l[1].Title == "React入門"
	})).Return(nil)

	results, err := svc.BatchCreate(1, logs)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, uint(1), results[0].UserID)
	assert.Equal(t, uint(1), results[1].UserID)
	repo.AssertExpectations(t)
}

func TestLearningLogBatchCreate_EmptyList(t *testing.T) {
	svc, _ := newTestLearningLogService()

	results, err := svc.BatchCreate(1, []model.LearningLog{})
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "1件以上")
}

func TestLearningLogBatchCreate_ExceedsMaxCount(t *testing.T) {
	svc, _ := newTestLearningLogService()

	logs := make([]model.LearningLog, 51)
	for i := range logs {
		logs[i] = model.LearningLog{Title: "テスト", Content: "内容", Duration: 10}
	}

	results, err := svc.BatchCreate(1, logs)
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "50件以下")
}

func TestLearningLogBatchCreate_InvalidDuration(t *testing.T) {
	svc, _ := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "正常", Content: "OK", Duration: 60},
		{Title: "異常", Content: "NG", Duration: -10},
	}

	results, err := svc.BatchCreate(1, logs)
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestLearningLogBatchCreate_InvalidCategory(t *testing.T) {
	svc, _ := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "正常", Content: "OK", Category: model.LogCategoryCoding},
		{Title: "異常", Content: "NG", Category: model.LogCategory("invalid")},
	}

	results, err := svc.BatchCreate(1, logs)
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestLearningLogBatchCreate_InvalidSource(t *testing.T) {
	svc, _ := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "正常", Content: "OK", Source: model.LogSourceManual},
		{Title: "異常", Content: "NG", Source: model.LogSource("unknown")},
	}

	results, err := svc.BatchCreate(1, logs)
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestLearningLogBatchCreate_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{
		{Title: "テスト", Content: "内容", Duration: 30},
	}

	repo.On("CreateBatch", mock.AnythingOfType("[]model.LearningLog")).Return(errors.New("db error"))

	results, err := svc.BatchCreate(1, logs)
	assert.Error(t, err)
	assert.Nil(t, results)
	repo.AssertExpectations(t)
}

func TestLearningLogBatchCreate_SetsUserID(t *testing.T) {
	svc, repo := newTestLearningLogService()

	// UserIDが異なる値でも、全て引数のuserIDに上書きされるべき
	logs := []model.LearningLog{
		{Title: "テスト1", Content: "内容1", UserID: 999, Duration: 30},
		{Title: "テスト2", Content: "内容2", UserID: 0, Duration: 45},
	}

	repo.On("CreateBatch", mock.MatchedBy(func(l []model.LearningLog) bool {
		return l[0].UserID == 1 && l[1].UserID == 1
	})).Return(nil)

	results, err := svc.BatchCreate(1, logs)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), results[0].UserID)
	assert.Equal(t, uint(1), results[1].UserID)
	repo.AssertExpectations(t)
}

func TestLearningLogExportCSV_MultipleRows(t *testing.T) {
	svc, repo := newTestLearningLogService()

	now := time.Now()
	logs := []model.LearningLog{
		{Title: "ログ1", Category: "Go", Duration: 30, Content: "メモ1"},
		{Title: "ログ2", Category: "React", Duration: 45, Content: "メモ2"},
		{Title: "ログ3", Category: "Docker", Duration: 15, Content: "メモ3"},
	}
	for i := range logs {
		logs[i].CreatedAt = now
	}

	repo.On("GetByPeriod", uint(1), 0).Return(logs, nil)

	data, err := svc.ExportCSV(1, 0)
	assert.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "ログ1")
	assert.Contains(t, content, "ログ2")
	assert.Contains(t, content, "ログ3")
	assert.Contains(t, content, "Go")
	assert.Contains(t, content, "React")
	assert.Contains(t, content, "Docker")
}

// ============================================================
// ExportJSON テスト
// ============================================================

func TestLearningLogExportJSON_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	now := time.Date(2026, 2, 19, 10, 30, 0, 0, time.UTC)
	logs := []model.LearningLog{
		{Title: "Go基礎", Content: "変数を学んだ", Category: model.LogCategoryCoding, Duration: 60, CreatedAt: now},
		{Title: "設計復習", Content: "DDD読んだ", Category: model.LogCategoryReading, Duration: 30, CreatedAt: now},
	}
	repo.On("GetByPeriod", uint(1), 30).Return(logs, nil)

	data, err := svc.ExportJSON(1, 30)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	content := string(data)
	assert.Contains(t, content, "Go基礎")
	assert.Contains(t, content, "設計復習")
	assert.Contains(t, content, "coding")
	assert.Contains(t, content, "reading")
	repo.AssertExpectations(t)
}

func TestLearningLogExportJSON_EmptyLogs(t *testing.T) {
	svc, repo := newTestLearningLogService()
	repo.On("GetByPeriod", uint(1), 0).Return([]model.LearningLog{}, nil)

	data, err := svc.ExportJSON(1, 0)
	assert.NoError(t, err)
	// 空配列 "[]" が返される
	assert.Equal(t, "[]", string(data))
	repo.AssertExpectations(t)
}

func TestLearningLogExportJSON_NegativeDays(t *testing.T) {
	svc, _ := newTestLearningLogService()

	_, err := svc.ExportJSON(1, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "期間は0以上の値")
}

func TestLearningLogExportJSON_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()
	repo.On("GetByPeriod", uint(1), 7).Return([]model.LearningLog{}, errors.New("db error"))

	_, err := svc.ExportJSON(1, 7)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogExportJSON_ValidJSON(t *testing.T) {
	svc, repo := newTestLearningLogService()

	now := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)
	logs := []model.LearningLog{
		{Title: "テスト", Content: "内容", Category: model.LogCategoryCoding, Duration: 45, CreatedAt: now},
	}
	repo.On("GetByPeriod", uint(1), 30).Return(logs, nil)

	data, err := svc.ExportJSON(1, 30)
	assert.NoError(t, err)

	// 有効なJSONであることを確認
	content := string(data)
	assert.True(t, strings.HasPrefix(content, "["))
	assert.True(t, strings.HasSuffix(content, "]"))
	assert.Contains(t, content, `"title"`)
	assert.Contains(t, content, `"date"`)
	assert.Contains(t, content, `"category"`)
	assert.Contains(t, content, `"duration"`)
	repo.AssertExpectations(t)
}

// ============================================================
// ゴール連携テスト
// ============================================================

func TestLearningLogCreate_WithGoalID_AutoUpdateProgress(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalID := uint(10)
	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60, GoalID: &goalID}

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 10, Status: model.GoalStatusActive, Progress: 0}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("Create", log).Return(nil)
	// 合計120分 = 2時間 / 10時間 = 20%
	repo.On("SumDurationByGoalID", uint(10)).Return(120, nil)
	goalRepo.On("Update", mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.ID == 10 && g.Progress == 20
	})).Return(nil)

	err := svc.Create(log)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogCreate_WithGoalID_AutoComplete(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalID := uint(10)
	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60, GoalID: &goalID}

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 1, Status: model.GoalStatusActive, Progress: 0}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("Create", log).Return(nil)
	// 合計60分 = 1時間 / 1時間 = 100% → 自動完了
	repo.On("SumDurationByGoalID", uint(10)).Return(60, nil)
	goalRepo.On("Update", mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.ID == 10 && g.Progress == 100 && g.Status == model.GoalStatusCompleted && g.CompletedAt != nil
	})).Return(nil)

	err := svc.Create(log)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogCreate_WithGoalID_GoalNotFound(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalID := uint(999)
	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60, GoalID: &goalID}

	goalRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Create(log)
	assert.Error(t, err)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogCreate_WithGoalID_GoalOwnerMismatch(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalID := uint(10)
	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60, GoalID: &goalID}

	goal := &model.LearningGoal{ID: 10, UserID: 999, TargetHours: 10} // 別ユーザーのゴール
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)

	err := svc.Create(log)
	assert.Error(t, err)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogCreate_WithGoalID_NoTargetHours(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalID := uint(10)
	log := &model.LearningLog{Title: "Go勉強", UserID: 1, Duration: 60, GoalID: &goalID}

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 0, Status: model.GoalStatusActive}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("Create", log).Return(nil)
	// TargetHours=0 の場合は進捗更新をスキップ

	err := svc.Create(log)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetLinkedLogs_Success(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 1}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	logs := []model.LearningLog{{ID: 1, Title: "テスト"}}
	repo.On("GetByGoalID", uint(10), 20, 0).Return(logs, int64(1), nil)

	result, total, err := svc.GetLinkedLogs(uint(10), uint(1), 20, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetLinkedLogs_Forbidden(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 999}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)

	_, _, err := svc.GetLinkedLogs(uint(10), uint(1), 20, 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetLinkedLogs_GoalNotFound(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, _, err := svc.GetLinkedLogs(uint(999), uint(1), 20, 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetLinkedLogs_NoGoalRepo(t *testing.T) {
	repo := new(MockLearningLogRepository)
	svc := NewLearningLogService(repo, nil)

	_, _, err := svc.GetLinkedLogs(uint(10), uint(1), 20, 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

// ============================================================
// GetGoalProgress テスト
// ============================================================

func TestLearningLogGetGoalProgress_Success(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 10}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	// 120分 = 2時間 / 10時間 = 20%
	repo.On("SumDurationByGoalID", uint(10)).Return(120, nil)

	progress, err := svc.GetGoalProgress(uint(10), uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(10), progress.GoalID)
	assert.Equal(t, 10, progress.TargetHours)
	assert.Equal(t, 120, progress.ActualMinutes)
	assert.Equal(t, 20, progress.Percentage)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetGoalProgress_ZeroTargetHours(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 0}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("SumDurationByGoalID", uint(10)).Return(60, nil)

	progress, err := svc.GetGoalProgress(uint(10), uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 0, progress.Percentage)
	assert.Equal(t, 60, progress.ActualMinutes)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetGoalProgress_Over100Percent(t *testing.T) {
	svc, repo, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 1, TargetHours: 1}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)
	// 120分 = 2時間 / 1時間 = 200% → キャップ100%
	repo.On("SumDurationByGoalID", uint(10)).Return(120, nil)

	progress, err := svc.GetGoalProgress(uint(10), uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 100, progress.Percentage)
	repo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetGoalProgress_Forbidden(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goal := &model.LearningGoal{ID: 10, UserID: 999}
	goalRepo.On("FindByID", uint(10)).Return(goal, nil)

	_, err := svc.GetGoalProgress(uint(10), uint(1))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetGoalProgress_NotFound(t *testing.T) {
	svc, _, goalRepo := newTestLearningLogServiceWithGoalRepo()

	goalRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.GetGoalProgress(uint(999), uint(1))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	goalRepo.AssertExpectations(t)
}

func TestLearningLogGetGoalProgress_NoGoalRepo(t *testing.T) {
	repo := new(MockLearningLogRepository)
	svc := NewLearningLogService(repo, nil)

	_, err := svc.GetGoalProgress(uint(10), uint(1))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

// ============================================================
// ImportCSV テスト
// ============================================================

func TestLearningLogImportCSV_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	csv := "\xef\xbb\xbf日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,Go学習,60,インターフェースを学んだ\n2025-01-16,reading,書籍読了,30,Clean Architecture\n"

	repo.On("CreateBatch", mock.MatchedBy(func(logs []model.LearningLog) bool {
		return len(logs) == 2 &&
			logs[0].Title == "Go学習" &&
			logs[0].Duration == 60 &&
			logs[0].Category == model.LogCategoryCoding &&
			logs[1].Title == "書籍読了" &&
			logs[1].Duration == 30
	})).Return(nil)

	result, err := svc.ImportCSV(1, []byte(csv))
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningLogImportCSV_NoBOM(t *testing.T) {
	svc, repo := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,Go学習,60,テスト\n"

	repo.On("CreateBatch", mock.Anything).Return(nil)

	result, err := svc.ImportCSV(1, []byte(csv))
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestLearningLogImportCSV_EmptyData(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogImportCSV_InvalidDate(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\nnot-a-date,coding,Test,60,Memo\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Contains(t, err.Error(), "日付形式")
}

func TestLearningLogImportCSV_InvalidCategory(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,invalid_cat,Test,60,Memo\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Contains(t, err.Error(), "カテゴリ")
}

func TestLearningLogImportCSV_InvalidDuration(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,Test,abc,Memo\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Contains(t, err.Error(), "学習時間")
}

func TestLearningLogImportCSV_DurationOutOfRange(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,Test,1500,Memo\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogImportCSV_EmptyTitle(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,,60,Memo\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Contains(t, err.Error(), "タイトル")
}

func TestLearningLogImportCSV_TooManyRows(t *testing.T) {
	svc, _ := newTestLearningLogService()

	var b strings.Builder
	b.WriteString("日付,カテゴリ,タイトル,学習時間(分),メモ\n")
	for i := 0; i < 51; i++ {
		b.WriteString("2025-01-15,coding,Test,60,Memo\n")
	}

	_, err := svc.ImportCSV(1, []byte(b.String()))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Contains(t, err.Error(), "50件")
}

func TestLearningLogImportCSV_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,coding,Test,60,Memo\n"

	repo.On("CreateBatch", mock.Anything).Return(errors.New("db error"))

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
}

func TestLearningLogImportCSV_InsufficientColumns(t *testing.T) {
	svc, _ := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル\n2025-01-15,coding,Test\n"

	_, err := svc.ImportCSV(1, []byte(csv))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestLearningLogImportCSV_DefaultCategory(t *testing.T) {
	svc, repo := newTestLearningLogService()

	csv := "日付,カテゴリ,タイトル,学習時間(分),メモ\n2025-01-15,,Go学習,60,テスト\n"

	repo.On("CreateBatch", mock.MatchedBy(func(logs []model.LearningLog) bool {
		return len(logs) == 1 && logs[0].Category == model.LogCategoryOther
	})).Return(nil)

	result, err := svc.ImportCSV(1, []byte(csv))
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

// ============================================================
// お気に入り一覧テスト
// ============================================================

func TestLearningLogGetFavorites_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	expected := []model.LearningLog{
		{ID: 1, UserID: 1, Title: "お気に入り1", IsFavorite: true},
		{ID: 2, UserID: 1, Title: "お気に入り2", IsFavorite: true},
	}
	repo.On("GetFavorites", uint(1), 20, 0).Return(expected, int64(2), nil)

	logs, total, err := svc.GetFavorites(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, logs, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestLearningLogGetFavorites_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetFavorites", uint(1), 20, 0).Return([]model.LearningLog{}, int64(0), nil)

	logs, total, err := svc.GetFavorites(1, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, logs)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

func TestLearningLogGetFavorites_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetFavorites", uint(1), 20, 0).Return([]model.LearningLog{}, int64(0), errors.New("db error"))

	logs, _, err := svc.GetFavorites(1, 20, 0)
	assert.Error(t, err)
	assert.Empty(t, logs)
	repo.AssertExpectations(t)
}

// ============================================================
// 月別サマリーテスト
// ============================================================

func TestLearningLogGetMonthlySummary_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	expected := []model.MonthlySummary{
		{Month: "2026-01-01", TotalMinutes: 300, LogCount: 10},
		{Month: "2026-02-01", TotalMinutes: 450, LogCount: 15},
	}
	repo.On("GetMonthlySummary", uint(1), 12).Return(expected, nil)

	result, err := svc.GetMonthlySummary(1, 12)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "2026-01-01", result[0].Month)
	assert.Equal(t, 300, result[0].TotalMinutes)
	repo.AssertExpectations(t)
}

func TestLearningLogGetMonthlySummary_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetMonthlySummary", uint(1), 6).Return([]model.MonthlySummary{}, nil)

	result, err := svc.GetMonthlySummary(1, 6)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestLearningLogGetMonthlySummary_InvalidMonths(t *testing.T) {
	svc, _ := newTestLearningLogService()

	result, err := svc.GetMonthlySummary(1, 0)
	assert.Nil(t, result)
	assert.Error(t, err)

	result, err = svc.GetMonthlySummary(1, 25)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestLearningLogGetMonthlySummary_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetMonthlySummary", uint(1), 12).Return([]model.MonthlySummary(nil), errors.New("db error"))

	result, err := svc.GetMonthlySummary(1, 12)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
