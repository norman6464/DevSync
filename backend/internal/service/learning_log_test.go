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
	svc := NewLearningLogService(repo)
	return svc, repo
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
