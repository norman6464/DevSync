package handler

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLearningLogRepo は usecase/repository.LearningLogRepository のモック（ctx 付き）。
type mockLearningLogRepo struct{ mock.Mock }

func (m *mockLearningLogRepo) Create(ctx context.Context, log *model.LearningLog) error {
	return m.Called(ctx, log).Error(0)
}
func (m *mockLearningLogRepo) CreateBatch(ctx context.Context, logs []model.LearningLog) error {
	return m.Called(ctx, logs).Error(0)
}
func (m *mockLearningLogRepo) Update(ctx context.Context, log *model.LearningLog) error {
	return m.Called(ctx, log).Error(0)
}
func (m *mockLearningLogRepo) Delete(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *mockLearningLogRepo) FindByID(ctx context.Context, id uint) (*model.LearningLog, error) {
	args := m.Called(ctx, id)
	l, _ := args.Get(0).(*model.LearningLog)
	return l, args.Error(1)
}
func (m *mockLearningLogRepo) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningLogRepo) GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningLog, error) {
	args := m.Called(ctx, userID, category)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Error(1)
}
func (m *mockLearningLogRepo) GetBySource(ctx context.Context, userID uint, source string) ([]model.LearningLog, error) {
	args := m.Called(ctx, userID, source)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Error(1)
}
func (m *mockLearningLogRepo) GetByPeriod(ctx context.Context, userID uint, days int) ([]model.LearningLog, error) {
	args := m.Called(ctx, userID, days)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Error(1)
}
func (m *mockLearningLogRepo) GetFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningLogRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockLearningLogRepo) SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error) {
	args := m.Called(ctx, userID, days)
	return args.Int(0), args.Error(1)
}
func (m *mockLearningLogRepo) SumDurationByGoalID(ctx context.Context, goalID uint) (int, error) {
	args := m.Called(ctx, goalID)
	return args.Int(0), args.Error(1)
}
func (m *mockLearningLogRepo) GetByGoalID(ctx context.Context, goalID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(ctx, goalID, limit, offset)
	l, _ := args.Get(0).([]model.LearningLog)
	return l, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningLogRepo) GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.StreakInfo)
	return s, args.Error(1)
}
func (m *mockLearningLogRepo) GetCalendarData(ctx context.Context, userID uint) ([]model.CalendarEntry, error) {
	args := m.Called(ctx, userID)
	e, _ := args.Get(0).([]model.CalendarEntry)
	return e, args.Error(1)
}
func (m *mockLearningLogRepo) GetRecentCategories(ctx context.Context, userID uint, limit int) ([]string, error) {
	args := m.Called(ctx, userID, limit)
	cats, _ := args.Get(0).([]string)
	return cats, args.Error(1)
}
func (m *mockLearningLogRepo) GetMonthlySummary(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error) {
	args := m.Called(ctx, userID, months)
	s, _ := args.Get(0).([]model.MonthlySummary)
	return s, args.Error(1)
}

// mockLearningGoalLinker は usecase/repository.LearningGoalLinker のモック（ctx 付き）。
type mockLearningGoalLinker struct{ mock.Mock }

func (m *mockLearningGoalLinker) FindByID(ctx context.Context, id uint) (*model.LearningGoal, error) {
	args := m.Called(ctx, id)
	g, _ := args.Get(0).(*model.LearningGoal)
	return g, args.Error(1)
}
func (m *mockLearningGoalLinker) Update(ctx context.Context, goal *model.LearningGoal) error {
	return m.Called(ctx, goal).Error(0)
}

// learningLogPorts は LearningLogHandler が使う port モックの束。
type learningLogPorts struct {
	Logs  *mockLearningLogRepo
	Goals *mockLearningGoalLinker
}

// newTestLearningLogHandler は本物の usecase に port モックを注入した LearningLogHandler を生成する。
func newTestLearningLogHandler() (*LearningLogHandler, learningLogPorts) {
	logs := new(mockLearningLogRepo)
	goals := new(mockLearningGoalLinker)
	h := NewLearningLogHandler(
		usecase.NewCreateLearningLogUseCase(logs, goals),
		usecase.NewBatchCreateLearningLogsUseCase(logs),
		usecase.NewImportLearningLogsCSVUseCase(logs),
		usecase.NewGetLearningLogUseCase(logs),
		usecase.NewListLearningLogsUseCase(logs),
		usecase.NewUpdateLearningLogUseCase(logs),
		usecase.NewDeleteLearningLogUseCase(logs),
		usecase.NewGetLearningStreakUseCase(logs),
		usecase.NewGetLearningCalendarUseCase(logs),
		usecase.NewExportLearningLogsCSVUseCase(logs),
		usecase.NewExportLearningLogsJSONUseCase(logs),
		usecase.NewListLearningLogsByCategoryUseCase(logs),
		usecase.NewListLearningLogsBySourceUseCase(logs),
		usecase.NewGetWeeklyLearningDurationUseCase(logs),
		usecase.NewFavoriteLearningLogUseCase(logs),
		usecase.NewUnfavoriteLearningLogUseCase(logs),
		usecase.NewListRecentLearningCategoriesUseCase(logs),
		usecase.NewListGoalLinkedLogsUseCase(logs, goals),
		usecase.NewGetGoalProgressUseCase(logs, goals),
		usecase.NewListFavoriteLearningLogsUseCase(logs),
		usecase.NewGetLearningLogMonthlySummaryUseCase(logs),
		usecase.NewCountLearningLogsUseCase(logs),
	)
	return h, learningLogPorts{Logs: logs, Goals: goals}
}

// logOwnedBy は指定ユーザーが所有する学習ログを返すテスト用ヘルパー。
func logOwnedBy(id, userID uint) *model.LearningLog {
	return &model.LearningLog{
		ID: id, UserID: userID, Title: "既存ログ", Content: "本文",
		Category: model.LogCategoryCoding, Duration: 60, Source: model.LogSourceManual,
	}
}

// ============================================================
// Create
// ============================================================

func TestLearningLogHandler_Create(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	p.Logs.On("Create", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		// カテゴリ未指定はその他で補われる。
		return l.UserID == 1 && l.Title == "Go学習" && l.Category == model.LogCategoryOther
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "duration": 60,
	})
	assertStatus(t, w, http.StatusCreated)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Create_InvalidJSON(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/learning-logs", "not json")
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// タイトル欠落は DTO の binding で 400 になる。
func TestLearningLogHandler_Create_ValidationError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{"content": "本文"})
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 無効なカテゴリは usecase の検証で 400 になる。
func TestLearningLogHandler_Create_InvalidCategory(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "category": "unknown",
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "無効なカテゴリです")
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_Create_InvalidSource(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "source": "unknown",
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "無効なソースです")
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ゴール紐付けで、ゴールが見つからなければ 404。
func TestLearningLogHandler_Create_GoalNotFound(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	p.Goals.On("FindByID", mock.Anything, uint(9)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "goal_id": 9,
	})
	assertStatus(t, w, http.StatusNotFound)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 他人のゴールへは紐付けできない。
func TestLearningLogHandler_Create_GoalForbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	p.Goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 2}, nil)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "goal_id": 9,
	})
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ゴール紐付けが成功すると、そのゴールの進捗も更新される。
func TestLearningLogHandler_Create_UpdatesGoalProgress(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	goal := &model.LearningGoal{ID: 9, UserID: 1, TargetHours: 10, Status: model.GoalStatusActive}
	p.Goals.On("FindByID", mock.Anything, uint(9)).Return(goal, nil)
	p.Logs.On("Create", mock.Anything, mock.AnythingOfType("*model.LearningLog")).Return(nil)
	p.Logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(300, nil)
	p.Goals.On("Update", mock.Anything, mock.MatchedBy(func(g *model.LearningGoal) bool {
		return g.Progress == 50
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文", "duration": 60, "goal_id": 9,
	})
	assertStatus(t, w, http.StatusCreated)
	p.Goals.AssertExpectations(t)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Create_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	p.Logs.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": "Go学習", "content": "本文",
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// BatchCreate
// ============================================================

func TestLearningLogHandler_BatchCreate(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/batch", h.BatchCreate)

	p.Logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(logs []model.LearningLog) bool {
		return len(logs) == 2 && logs[0].UserID == 1 && logs[1].UserID == 1
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/learning-logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{
			{"title": "1件目", "content": "本文", "duration": 30},
			{"title": "2件目", "content": "本文", "duration": 30},
		},
	})
	assertStatus(t, w, http.StatusCreated)
	p.Logs.AssertExpectations(t)
}

// 空配列は DTO の binding（min=1）で 400 になる。
func TestLearningLogHandler_BatchCreate_Empty(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/batch", h.BatchCreate)

	w := doRequest(r, http.MethodPost, "/learning-logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{},
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_BatchCreate_InvalidJSON(t *testing.T) {
	h, _ := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/batch", h.BatchCreate)

	w := doRequestRaw(r, http.MethodPost, "/learning-logs/batch", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLogHandler_BatchCreate_MissingTitle(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/batch", h.BatchCreate)

	w := doRequest(r, http.MethodPost, "/learning-logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{{"content": "本文", "duration": 30}},
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_BatchCreate_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/batch", h.BatchCreate)

	p.Logs.On("CreateBatch", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/learning-logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{{"title": "1件目", "content": "本文", "duration": 30}},
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByID / 一覧
// ============================================================

func TestLearningLogHandler_GetByID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/:id", h.GetByID)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/1", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存ログ")
}

func TestLearningLogHandler_GetByID_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/:id", h.GetByID)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/1", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// 不在のログは 500 になる（移行前から変わらない挙動）。
func TestLearningLogHandler_GetByID_NotFoundIs500(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/:id", h.GetByID)

	p.Logs.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/99", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestLearningLogHandler_GetMyLogs(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs", h.GetMyLogs)

	p.Logs.On("GetByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total":1`)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetMyLogs_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs", h.GetMyLogs)

	p.Logs.On("GetByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.LearningLog(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/learning-logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogHandler_GetByUserID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/learning-logs", h.GetByUserID)

	p.Logs.On("GetByUserID", mock.Anything, uint(7), 20, 0).
		Return([]model.LearningLog{*logOwnedBy(1, 7)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/users/7/learning-logs", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetByUserID_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/learning-logs", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/learning-logs", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetByUserID_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/learning-logs", h.GetByUserID)

	p.Logs.On("GetByUserID", mock.Anything, uint(7), 20, 0).
		Return([]model.LearningLog(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/learning-logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Update / Delete
// ============================================================

func TestLearningLogHandler_Update(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.PUT("/learning-logs/:id", h.Update)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 1), nil)
	p.Logs.On("Update", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		// 指定しなかったフィールドは据え置かれる。
		return l.Title == "更新後" && l.Content == "本文" && l.Duration == 60
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/learning-logs/1", map[string]interface{}{"title": "更新後"})
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Update_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.PUT("/learning-logs/:id", h.Update)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPut, "/learning-logs/1", map[string]interface{}{"title": "更新後"})
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 不在のログの更新は 500 になる（移行前から変わらない挙動）。
func TestLearningLogHandler_Update_NotFoundIs500(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.PUT("/learning-logs/:id", h.Update)

	p.Logs.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/learning-logs/99", map[string]interface{}{"title": "更新後"})
	assertStatus(t, w, http.StatusNotFound)
	p.Logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_Update_InvalidJSON(t *testing.T) {
	h, _ := newTestLearningLogHandler()
	r := newRouter(1)
	r.PUT("/learning-logs/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/learning-logs/1", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLogHandler_Update_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.PUT("/learning-logs/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/learning-logs/abc", map[string]interface{}{"title": "更新後"})
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_Delete(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.DELETE("/learning-logs/:id", h.Delete)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 1), nil)
	p.Logs.On("Delete", mock.Anything, uint(1), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/learning-logs/1", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Delete_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.DELETE("/learning-logs/:id", h.Delete)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodDelete, "/learning-logs/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

// ============================================================
// ストリーク / カレンダー / 週間学習時間 / 月別サマリー
// ============================================================

func TestLearningLogHandler_GetStreakInfo(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/streak", h.GetStreakInfo)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(7)).
		Return(&model.StreakInfo{CurrentStreak: 3, LongestStreak: 5, TotalDays: 10}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/streak", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"current_streak":3`)
}

func TestLearningLogHandler_GetStreakInfo_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/streak", h.GetStreakInfo)

	w := doRequest(r, http.MethodGet, "/users/abc/streak", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetStreakInfo", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetStreakInfo_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/streak", h.GetStreakInfo)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(7)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/streak", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogHandler_GetMyStreakInfo(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/my/streak", h.GetMyStreakInfo)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{CurrentStreak: 1}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/my/streak", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetCalendarData(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/calendar", h.GetCalendarData)

	p.Logs.On("GetCalendarData", mock.Anything, uint(7)).
		Return([]model.CalendarEntry{{Date: "2026-01-01", Count: 2}}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/calendar", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "2026-01-01")
}

// 0 件でも null ではなく空配列を返す。
func TestLearningLogHandler_GetCalendarData_Empty(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/calendar", h.GetCalendarData)

	p.Logs.On("GetCalendarData", mock.Anything, uint(7)).Return([]model.CalendarEntry(nil), nil)

	w := doRequest(r, http.MethodGet, "/users/7/calendar", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestLearningLogHandler_GetCalendarData_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/calendar", h.GetCalendarData)

	p.Logs.On("GetCalendarData", mock.Anything, uint(7)).
		Return([]model.CalendarEntry(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/calendar", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogHandler_GetMyCalendarData(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/my/calendar", h.GetMyCalendarData)

	p.Logs.On("GetCalendarData", mock.Anything, uint(1)).Return([]model.CalendarEntry{}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/my/calendar", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetWeeklyDuration(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/weekly-duration", h.GetWeeklyDuration)

	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(7), 7).Return(180, nil)

	w := doRequest(r, http.MethodGet, "/users/7/weekly-duration", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"duration":180`)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetWeeklyDuration_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/weekly-duration", h.GetWeeklyDuration)

	w := doRequest(r, http.MethodGet, "/users/abc/weekly-duration", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "SumDurationByPeriod", mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetWeeklyDuration_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/weekly-duration", h.GetWeeklyDuration)

	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(7), 7).Return(0, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/weekly-duration", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogHandler_GetMyWeeklyDuration(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/my/weekly-duration", h.GetMyWeeklyDuration)

	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(60, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/my/weekly-duration", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetMonthlySummary(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/monthly-summary", h.GetMonthlySummary)

	p.Logs.On("GetMonthlySummary", mock.Anything, uint(7), 12).
		Return([]model.MonthlySummary{{Month: "2026-01-01", TotalMinutes: 120, LogCount: 3}}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/monthly-summary", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

// months が範囲外なら usecase の検証で 400 になる。
func TestLearningLogHandler_GetMonthlySummary_OutOfRange(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/monthly-summary", h.GetMonthlySummary)

	w := doRequest(r, http.MethodGet, "/users/7/monthly-summary?months=25", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetMonthlySummary", mock.Anything, mock.Anything, mock.Anything)
}

// ============================================================
// カテゴリ / ソース / 最近のカテゴリ
// ============================================================

func TestLearningLogHandler_GetByCategory(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	p.Logs.On("GetByCategory", mock.Anything, uint(1), "coding").
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/category/coding", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

// 0 件でも null ではなく空配列を返す。
func TestLearningLogHandler_GetByCategory_Empty(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	p.Logs.On("GetByCategory", mock.Anything, uint(1), "coding").Return([]model.LearningLog(nil), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/category/coding", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestLearningLogHandler_GetByCategory_Invalid(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	w := doRequest(r, http.MethodGet, "/learning-logs/category/unknown", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetByCategory", mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetBySource(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/source/:source", h.GetBySource)

	p.Logs.On("GetBySource", mock.Anything, uint(1), "manual").
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/source/manual", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetBySource_Invalid(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/source/:source", h.GetBySource)

	w := doRequest(r, http.MethodGet, "/learning-logs/source/unknown", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetBySource", mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetRecentCategories(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/recent-categories", h.GetRecentCategories)

	p.Logs.On("GetRecentCategories", mock.Anything, uint(1), 5).
		Return([]string{"coding", "reading"}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/recent-categories", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "coding")
	p.Logs.AssertExpectations(t)
}

// 0 件でも null ではなく空配列を返す。
func TestLearningLogHandler_GetRecentCategories_Empty(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/recent-categories", h.GetRecentCategories)

	p.Logs.On("GetRecentCategories", mock.Anything, uint(1), 5).Return([]string(nil), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/recent-categories", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestLearningLogHandler_GetRecentCategories_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/recent-categories", h.GetRecentCategories)

	p.Logs.On("GetRecentCategories", mock.Anything, uint(1), 5).Return([]string(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/recent-categories", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// お気に入り
// ============================================================

func TestLearningLogHandler_Favorite(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/:id/favorite", h.Favorite)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 1), nil)
	p.Logs.On("Update", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		return l.IsFavorite
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/learning-logs/1/favorite", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Favorite_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/:id/favorite", h.Favorite)

	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(logOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPost, "/learning-logs/1/favorite", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_Favorite_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs/:id/favorite", h.Favorite)

	w := doRequest(r, http.MethodPost, "/learning-logs/abc/favorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_Unfavorite(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.DELETE("/learning-logs/:id/favorite", h.Unfavorite)

	favorited := logOwnedBy(1, 1)
	favorited.IsFavorite = true
	p.Logs.On("FindByID", mock.Anything, uint(1)).Return(favorited, nil)
	p.Logs.On("Update", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		return !l.IsFavorite
	})).Return(nil)

	w := doRequest(r, http.MethodDelete, "/learning-logs/1/favorite", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_Unfavorite_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.DELETE("/learning-logs/:id/favorite", h.Unfavorite)

	w := doRequest(r, http.MethodDelete, "/learning-logs/abc/favorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetFavorites(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/favorites", h.GetFavorites)

	p.Logs.On("GetFavorites", mock.Anything, uint(1), 20, 0).
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/favorites", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
}

// ============================================================
// エクスポート
// ============================================================

func TestLearningLogHandler_ExportLogs_CSV(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	p.Logs.On("GetByPeriod", mock.Anything, uint(1), 30).
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, w.Body.String(), "日付,カテゴリ,タイトル,学習時間(分),メモ")
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_ExportLogs_AllPeriod(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	p.Logs.On("GetByPeriod", mock.Anything, uint(1), 0).Return([]model.LearningLog{}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?period=all", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "learning-logs-all-")
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_ExportLogs_InvalidPeriod(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?period=abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetByPeriod", mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_ExportLogs_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	p.Logs.On("GetByPeriod", mock.Anything, uint(1), 30).
		Return([]model.LearningLog(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/export", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestLearningLogHandler_ExportLogs_JSON(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	p.Logs.On("GetByPeriod", mock.Anything, uint(1), 30).
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=json", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	assert.Contains(t, w.Body.String(), `"title":"既存ログ"`)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_ExportLogs_JSONAllPeriod(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	p.Logs.On("GetByPeriod", mock.Anything, uint(1), 0).Return([]model.LearningLog{}, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?period=all&format=json", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "learning-logs-all-")
}

func TestLearningLogHandler_ExportLogs_InvalidFormat(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=xml", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "GetByPeriod", mock.Anything, mock.Anything, mock.Anything)
}

// ============================================================
// CSV インポート
// ============================================================

// createCSVMultipartRequest はCSVファイル付きのマルチパートリクエストを組み立てる。
func createCSVMultipartRequest(csvContent string) (*http.Request, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "import.csv")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write([]byte(csvContent)); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/logs/import", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestLearningLogHandler_ImportCSV(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/import", h.ImportCSV)

	p.Logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(logs []model.LearningLog) bool {
		return len(logs) == 1 && logs[0].Title == "Go学習" && logs[0].Duration == 60 && logs[0].UserID == 1
	})).Return(nil)

	req, err := createCSVMultipartRequest("日付,カテゴリ,タイトル,学習時間(分),メモ\n2026-01-15,coding,Go学習,60,テスト\n")
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["imported"])
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_ImportCSV_NoFile(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/import", h.ImportCSV)

	w := doRequest(r, http.MethodPost, "/logs/import", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

// 行の内容が壊れている CSV は、行番号つきのメッセージで 400 になる。
func TestLearningLogHandler_ImportCSV_InvalidRow(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/import", h.ImportCSV)

	req, err := createCSVMultipartRequest("日付,カテゴリ,タイトル,学習時間(分),メモ\nxxxx,coding,Go学習,60,テスト\n")
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "2行目")
	p.Logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_ImportCSV_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/import", h.ImportCSV)

	p.Logs.On("CreateBatch", mock.Anything, mock.Anything).Return(errors.New("db error"))

	req, err := createCSVMultipartRequest("日付,カテゴリ,タイトル,学習時間(分),メモ\n2026-01-15,coding,Go学習,60,テスト\n")
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// ゴール連携
// ============================================================

func TestLearningLogHandler_GetLinkedLogs(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)

	p.Goals.On("FindByID", mock.Anything, uint(5)).Return(&model.LearningGoal{ID: 5, UserID: 1}, nil)
	p.Logs.On("GetByGoalID", mock.Anything, uint(5), 20, 0).
		Return([]model.LearningLog{*logOwnedBy(1, 1)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/goals/5/linked-logs", nil)
	assertStatus(t, w, http.StatusOK)
	p.Logs.AssertExpectations(t)
	p.Goals.AssertExpectations(t)
}

func TestLearningLogHandler_GetLinkedLogs_GoalNotFound(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)

	p.Goals.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/goals/5/linked-logs", nil)
	assertStatus(t, w, http.StatusNotFound)
	p.Logs.AssertNotCalled(t, "GetByGoalID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetLinkedLogs_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)

	p.Goals.On("FindByID", mock.Anything, uint(5)).Return(&model.LearningGoal{ID: 5, UserID: 2}, nil)

	w := doRequest(r, http.MethodGet, "/goals/5/linked-logs", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "GetByGoalID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetLinkedLogs_InvalidID(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)

	w := doRequest(r, http.MethodGet, "/goals/abc/linked-logs", nil)
	assertStatus(t, w, http.StatusBadRequest)
	p.Goals.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetGoalProgress(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/progress", h.GetGoalProgress)

	p.Goals.On("FindByID", mock.Anything, uint(5)).
		Return(&model.LearningGoal{ID: 5, UserID: 1, TargetHours: 10}, nil)
	p.Logs.On("SumDurationByGoalID", mock.Anything, uint(5)).Return(300, nil)

	w := doRequest(r, http.MethodGet, "/goals/5/progress", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"percentage":50`)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetGoalProgress_Forbidden(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/progress", h.GetGoalProgress)

	p.Goals.On("FindByID", mock.Anything, uint(5)).Return(&model.LearningGoal{ID: 5, UserID: 2}, nil)

	w := doRequest(r, http.MethodGet, "/goals/5/progress", nil)
	assertStatus(t, w, http.StatusForbidden)
	p.Logs.AssertNotCalled(t, "SumDurationByGoalID", mock.Anything, mock.Anything)
}

func TestLearningLogHandler_GetGoalProgress_GoalNotFound(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/progress", h.GetGoalProgress)

	p.Goals.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/goals/5/progress", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ============================================================
// GetMyCount
// ============================================================

func TestLearningLogHandler_GetMyCount(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/my/count", h.GetMyCount)

	p.Logs.On("CountByUserID", mock.Anything, uint(1)).Return(int64(15), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":15`)
	p.Logs.AssertExpectations(t)
}

func TestLearningLogHandler_GetMyCount_RepositoryError(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/my/count", h.GetMyCount)

	p.Logs.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Logs.AssertExpectations(t)
}

// 長すぎるタイトルは DTO の binding で弾かれ、usecase まで届かない。
func TestLearningLogHandler_Create_TitleTooLong(t *testing.T) {
	h, p := newTestLearningLogHandler()
	r := newRouter(1)
	r.POST("/learning-logs", h.Create)

	w := doRequest(r, http.MethodPost, "/learning-logs", map[string]interface{}{
		"title": strings.Repeat("あ", 201), "content": "本文",
	})
	assertStatus(t, w, http.StatusBadRequest)
	p.Logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
