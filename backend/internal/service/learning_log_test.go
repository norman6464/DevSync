package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
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

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Go Study", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("FindByID", uint(999)).Return((*model.LearningLog)(nil), errors.New("not found"))

	_, err := svc.GetByID(999)
	assert.Error(t, err)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestLearningLogGetByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	logs := []model.LearningLog{{Title: "Go Study"}, {Title: "React Study"}}
	repo.On("GetByUserID", uint(1)).Return(logs, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningLogGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetByUserID", uint(1)).Return([]model.LearningLog{}, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
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
