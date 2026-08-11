package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockLearningLogRepo は usecase/repository.LearningLogRepository のモック。
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
	c, _ := args.Get(0).([]string)
	return c, args.Error(1)
}
func (m *mockLearningLogRepo) GetMonthlySummary(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error) {
	args := m.Called(ctx, userID, months)
	s, _ := args.Get(0).([]model.MonthlySummary)
	return s, args.Error(1)
}

// mockLearningGoalLinker は usecase/repository.LearningGoalLinker のモック。
type mockLearningGoalLinker struct{ mock.Mock }

func (m *mockLearningGoalLinker) FindByID(ctx context.Context, id uint) (*model.LearningGoal, error) {
	args := m.Called(ctx, id)
	g, _ := args.Get(0).(*model.LearningGoal)
	return g, args.Error(1)
}
func (m *mockLearningGoalLinker) Update(ctx context.Context, goal *model.LearningGoal) error {
	return m.Called(ctx, goal).Error(0)
}

// ownedLearningLog は指定ユーザーが所有する学習ログを返すテスト用ヘルパー。
func ownedLearningLog(id, userID uint) *model.LearningLog {
	return &model.LearningLog{
		ID: id, UserID: userID, Title: "既存ログ", Content: "本文",
		Category: model.LogCategoryCoding, Duration: 60, Source: model.LogSourceManual,
	}
}

// ============================================================
// Create
// ============================================================

func TestCreateLearningLogUseCase_Execute(t *testing.T) {
	t.Run("カテゴリ未指定はその他で補われる", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("Create", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
			return l.Category == model.LogCategoryOther
		})).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, nil)

		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", Duration: 30})

		assert.NoError(t, err)
		logs.AssertExpectations(t)
	})

	t.Run("内容は空でも作成できる", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, nil)

		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題"})

		assert.NoError(t, err)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name string
			log  *model.LearningLog
			msg  string
		}{
			{"タイトルが空", &model.LearningLog{Title: ""}, "タイトル"},
			{"タイトルが上限超過", &model.LearningLog{Title: strings.Repeat("あ", 201)}, "タイトル"},
			{"内容が上限超過", &model.LearningLog{Title: "題", Content: strings.Repeat("あ", 10001)}, "内容"},
			{"学習時間が負", &model.LearningLog{Title: "題", Duration: -1}, "学習時間"},
			{"学習時間が上限超過", &model.LearningLog{Title: "題", Duration: 1441}, "学習時間"},
			{"無効なカテゴリ", &model.LearningLog{Title: "題", Category: "unknown"}, "無効なカテゴリです"},
			{"無効なソース", &model.LearningLog{Title: "題", Source: "unknown"}, "無効なソースです"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				logs := new(mockLearningLogRepo)
				uc := usecase.NewCreateLearningLogUseCase(logs, nil)

				err := uc.Execute(context.Background(), c.log)

				require.Error(t, err)
				domainErr := domain.GetDomainError(err)
				require.NotNil(t, domainErr)
				assert.Contains(t, domainErr.Message, c.msg)
				logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("書き込み失敗は内部エラーに包む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateLearningLogUseCase(logs, nil)

		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題"})

		require.Error(t, err)
		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "学習ログの作成に失敗しました", domainErr.Message)
	})

	t.Run("ゴールが取得できなければ 404", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(nil, nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
		logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("ゴール取得が DB 障害でも 404 に潰す", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(nil, errors.New("db error"))
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	})

	t.Run("他人のゴールには紐付けられない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 2}, nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
		logs.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("目標時間があるゴールは進捗を更新する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goal := &model.LearningGoal{ID: 9, UserID: 1, TargetHours: 10, Status: model.GoalStatusActive}
		goals.On("FindByID", mock.Anything, uint(9)).Return(goal, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(300, nil)
		goals.On("Update", mock.Anything, mock.MatchedBy(func(g *model.LearningGoal) bool {
			return g.Progress == 50 && g.Status == model.GoalStatusActive
		})).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", Duration: 60, GoalID: &goalID})

		assert.NoError(t, err)
		goals.AssertExpectations(t)
	})

	t.Run("進捗が 100% に達したゴールは完了へ遷移する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goal := &model.LearningGoal{ID: 9, UserID: 1, TargetHours: 1, Status: model.GoalStatusActive}
		goals.On("FindByID", mock.Anything, uint(9)).Return(goal, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(120, nil)
		goals.On("Update", mock.Anything, mock.MatchedBy(func(g *model.LearningGoal) bool {
			return g.Progress == 100 && g.Status == model.GoalStatusCompleted && g.CompletedAt != nil
		})).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", Duration: 120, GoalID: &goalID})

		assert.NoError(t, err)
		goals.AssertExpectations(t)
	})

	t.Run("目標時間が未設定のゴールは進捗を更新しない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 1}, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		assert.NoError(t, err)
		logs.AssertNotCalled(t, "SumDurationByGoalID", mock.Anything, mock.Anything)
		goals.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("集計に失敗してもログ作成は成功のまま", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).
			Return(&model.LearningGoal{ID: 9, UserID: 1, TargetHours: 10}, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(0, errors.New("db error"))
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		assert.NoError(t, err)
		goals.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("ゴール更新の失敗もログ作成には影響しない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).
			Return(&model.LearningGoal{ID: 9, UserID: 1, TargetHours: 10}, nil)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(60, nil)
		goals.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateLearningLogUseCase(logs, goals)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		assert.NoError(t, err)
	})

	t.Run("ゴール連携が無効ならゴール ID は無視する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateLearningLogUseCase(logs, nil)

		goalID := uint(9)
		err := uc.Execute(context.Background(), &model.LearningLog{UserID: 1, Title: "題", GoalID: &goalID})

		assert.NoError(t, err)
	})
}

// ============================================================
// BatchCreate
// ============================================================

func TestBatchCreateLearningLogsUseCase_Execute(t *testing.T) {
	t.Run("全件にユーザー ID を入れて保存する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(l []model.LearningLog) bool {
			return len(l) == 2 && l[0].UserID == 3 && l[1].UserID == 3 &&
				l[0].Category == model.LogCategoryOther
		})).Return(nil)
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		result, err := uc.Execute(context.Background(), 3, []model.LearningLog{
			{Title: "1件目"}, {Title: "2件目"},
		})

		require.NoError(t, err)
		assert.Len(t, result, 2)
		logs.AssertExpectations(t)
	})

	t.Run("0 件は 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		_, err := uc.Execute(context.Background(), 3, nil)

		require.Error(t, err)
		assert.NotNil(t, domain.GetDomainError(err))
		logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
	})

	t.Run("51 件以上は 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		many := make([]model.LearningLog, 51)
		for i := range many {
			many[i] = model.LearningLog{Title: "題"}
		}
		_, err := uc.Execute(context.Background(), 3, many)

		require.Error(t, err)
		logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
	})

	t.Run("ちょうど 50 件は保存できる", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		many := make([]model.LearningLog, 50)
		for i := range many {
			many[i] = model.LearningLog{Title: "題"}
		}
		_, err := uc.Execute(context.Background(), 3, many)

		assert.NoError(t, err)
	})

	t.Run("1 件でも無効なら保存しない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		_, err := uc.Execute(context.Background(), 3, []model.LearningLog{
			{Title: "1件目"}, {Title: "2件目", Category: "unknown"},
		})

		require.Error(t, err)
		logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
	})

	t.Run("書き込み失敗は内部エラーに包む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewBatchCreateLearningLogsUseCase(logs)

		_, err := uc.Execute(context.Background(), 3, []model.LearningLog{{Title: "題"}})

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "学習ログの一括作成に失敗しました", domainErr.Message)
	})
}

// ============================================================
// 取得 / 一覧 / 件数
// ============================================================

func TestGetLearningLogUseCase_Execute(t *testing.T) {
	t.Run("所有者は取得できる", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewGetLearningLogUseCase(logs)

		log, err := uc.Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(1), log.ID)
	})

	t.Run("他人のログは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewGetLearningLogUseCase(logs)

		_, err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetLearningLogUseCase(logs)

		_, err := uc.Execute(context.Background(), 1, 5)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

func TestListLearningLogsUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("GetByUserID", mock.Anything, uint(5), 20, 0).
		Return([]model.LearningLog{*ownedLearningLog(1, 5)}, int64(1), nil)
	uc := usecase.NewListLearningLogsUseCase(logs)

	result, total, err := uc.Execute(context.Background(), 5, 20, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
}

func TestCountLearningLogsUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("CountByUserID", mock.Anything, uint(5)).Return(int64(7), nil)
	uc := usecase.NewCountLearningLogsUseCase(logs)

	count, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, int64(7), count)
}

// ============================================================
// Update / Delete / お気に入り
// ============================================================

func TestUpdateLearningLogUseCase_Execute(t *testing.T) {
	t.Run("指定したフィールドだけを更新する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		logs.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		log, err := uc.Execute(context.Background(), 1, 5, &model.LearningLog{Title: "  更新後  "})

		require.NoError(t, err)
		assert.Equal(t, "更新後", log.Title)
		assert.Equal(t, "本文", log.Content)
		assert.Equal(t, 60, log.Duration)
	})

	t.Run("学習時間 0 は変更なしとして扱う", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		logs.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		log, err := uc.Execute(context.Background(), 1, 5, &model.LearningLog{Duration: 0})

		require.NoError(t, err)
		assert.Equal(t, 60, log.Duration)
	})

	t.Run("空白のみのフィールドは据え置く", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		logs.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		log, err := uc.Execute(context.Background(), 1, 5, &model.LearningLog{Title: "   ", Content: "  "})

		require.NoError(t, err)
		assert.Equal(t, "既存ログ", log.Title)
		assert.Equal(t, "本文", log.Content)
	})

	t.Run("上限超過は 400 で書き込まない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		_, err := uc.Execute(context.Background(), 1, 5, &model.LearningLog{Title: strings.Repeat("あ", 201)})

		require.Error(t, err)
		assert.NotNil(t, domain.GetDomainError(err))
		logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("学習時間の範囲外は 400 で書き込まない", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		_, err := uc.Execute(context.Background(), 1, 5, &model.LearningLog{Duration: 1441})

		require.Error(t, err)
		logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("他人のログは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewUpdateLearningLogUseCase(logs)

		_, err := uc.Execute(context.Background(), 1, 9, &model.LearningLog{Title: "更新後"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestDeleteLearningLogUseCase_Execute(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		logs.On("Delete", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewDeleteLearningLogUseCase(logs)

		err := uc.Execute(context.Background(), 1, 5)

		assert.NoError(t, err)
		logs.AssertExpectations(t)
	})

	t.Run("他人のログは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewDeleteLearningLogUseCase(logs)

		err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		logs.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestFavoriteLearningLogUseCase_Execute(t *testing.T) {
	t.Run("お気に入りに設定する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		logs.On("Update", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
			return l.IsFavorite
		})).Return(nil)
		uc := usecase.NewFavoriteLearningLogUseCase(logs)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		logs.AssertExpectations(t)
	})

	t.Run("他人のログは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("FindByID", mock.Anything, uint(1)).Return(ownedLearningLog(1, 5), nil)
		uc := usecase.NewFavoriteLearningLogUseCase(logs)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1, 9), domain.ErrForbidden)
		logs.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestUnfavoriteLearningLogUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	favorited := ownedLearningLog(1, 5)
	favorited.IsFavorite = true
	logs.On("FindByID", mock.Anything, uint(1)).Return(favorited, nil)
	logs.On("Update", mock.Anything, mock.MatchedBy(func(l *model.LearningLog) bool {
		return !l.IsFavorite
	})).Return(nil)
	uc := usecase.NewUnfavoriteLearningLogUseCase(logs)

	assert.NoError(t, uc.Execute(context.Background(), 1, 5))
	logs.AssertExpectations(t)
}

// ============================================================
// カテゴリ / ソース / 集計
// ============================================================

func TestListLearningLogsByCategoryUseCase_Execute(t *testing.T) {
	t.Run("有効なカテゴリで絞り込む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetByCategory", mock.Anything, uint(5), "coding").
			Return([]model.LearningLog{*ownedLearningLog(1, 5)}, nil)
		uc := usecase.NewListLearningLogsByCategoryUseCase(logs)

		result, err := uc.Execute(context.Background(), 5, "coding")

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("無効なカテゴリは 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewListLearningLogsByCategoryUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, "unknown")

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "無効なカテゴリです", domainErr.Message)
		logs.AssertNotCalled(t, "GetByCategory", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestListLearningLogsBySourceUseCase_Execute(t *testing.T) {
	t.Run("有効なソースで絞り込む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetBySource", mock.Anything, uint(5), "manual").
			Return([]model.LearningLog{*ownedLearningLog(1, 5)}, nil)
		uc := usecase.NewListLearningLogsBySourceUseCase(logs)

		result, err := uc.Execute(context.Background(), 5, "manual")

		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("無効なソースは 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewListLearningLogsBySourceUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, "unknown")

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "無効なソースです", domainErr.Message)
	})
}

func TestGetWeeklyLearningDurationUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("SumDurationByPeriod", mock.Anything, uint(5), 7).Return(180, nil)
	uc := usecase.NewGetWeeklyLearningDurationUseCase(logs)

	duration, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, 180, duration)
	logs.AssertExpectations(t)
}

func TestListRecentLearningCategoriesUseCase_Execute(t *testing.T) {
	t.Run("上位 5 件を返す", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetRecentCategories", mock.Anything, uint(5), 5).
			Return([]string{"coding", "reading"}, nil)
		uc := usecase.NewListRecentLearningCategoriesUseCase(logs)

		categories, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, []string{"coding", "reading"}, categories)
		logs.AssertExpectations(t)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetRecentCategories", mock.Anything, uint(5), 5).Return([]string(nil), errors.New("db error"))
		uc := usecase.NewListRecentLearningCategoriesUseCase(logs)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestGetLearningLogMonthlySummaryUseCase_Execute(t *testing.T) {
	t.Run("範囲内の months で取得する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetMonthlySummary", mock.Anything, uint(5), 12).
			Return([]model.MonthlySummary{{Month: "2026-01-01", TotalMinutes: 120, LogCount: 2}}, nil)
		uc := usecase.NewGetLearningLogMonthlySummaryUseCase(logs)

		summaries, err := uc.Execute(context.Background(), 5, 12)

		require.NoError(t, err)
		assert.Len(t, summaries, 1)
	})

	t.Run("範囲外の months は 400", func(t *testing.T) {
		for _, months := range []int{0, 25} {
			logs := new(mockLearningLogRepo)
			uc := usecase.NewGetLearningLogMonthlySummaryUseCase(logs)

			_, err := uc.Execute(context.Background(), 5, months)

			require.Error(t, err)
			assert.NotNil(t, domain.GetDomainError(err))
			logs.AssertNotCalled(t, "GetMonthlySummary", mock.Anything, mock.Anything, mock.Anything)
		}
	})
}

func TestGetLearningStreakUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("GetStreakInfo", mock.Anything, uint(5)).
		Return(&model.StreakInfo{CurrentStreak: 3, LongestStreak: 5}, nil)
	uc := usecase.NewGetLearningStreakUseCase(logs)

	info, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, 3, info.CurrentStreak)
}

func TestGetLearningCalendarUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("GetCalendarData", mock.Anything, uint(5)).
		Return([]model.CalendarEntry{{Date: "2026-01-01", Count: 2}}, nil)
	uc := usecase.NewGetLearningCalendarUseCase(logs)

	entries, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestListFavoriteLearningLogsUseCase_Execute(t *testing.T) {
	logs := new(mockLearningLogRepo)
	logs.On("GetFavorites", mock.Anything, uint(5), 20, 0).
		Return([]model.LearningLog{*ownedLearningLog(1, 5)}, int64(1), nil)
	uc := usecase.NewListFavoriteLearningLogsUseCase(logs)

	result, total, err := uc.Execute(context.Background(), 5, 20, 0)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
}

// ============================================================
// ゴール連携
// ============================================================

func TestListGoalLinkedLogsUseCase_Execute(t *testing.T) {
	t.Run("所有ゴールの紐付けログを返す", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 5}, nil)
		logs.On("GetByGoalID", mock.Anything, uint(9), 20, 0).
			Return([]model.LearningLog{*ownedLearningLog(1, 5)}, int64(1), nil)
		uc := usecase.NewListGoalLinkedLogsUseCase(logs, goals)

		result, total, err := uc.Execute(context.Background(), 9, 5, 20, 0)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("他人のゴールは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 2}, nil)
		uc := usecase.NewListGoalLinkedLogsUseCase(logs, goals)

		_, _, err := uc.Execute(context.Background(), 9, 5, 20, 0)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
		logs.AssertNotCalled(t, "GetByGoalID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ゴールが無ければ 404", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(nil, nil)
		uc := usecase.NewListGoalLinkedLogsUseCase(logs, goals)

		_, _, err := uc.Execute(context.Background(), 9, 5, 20, 0)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
	})

	t.Run("ゴール連携が無効なら 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewListGoalLinkedLogsUseCase(logs, nil)

		_, _, err := uc.Execute(context.Background(), 9, 5, 20, 0)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "ゴール連携が有効ではありません", domainErr.Message)
	})
}

func TestGetGoalProgressUseCase_Execute(t *testing.T) {
	t.Run("実績と目標から進捗を算出する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).
			Return(&model.LearningGoal{ID: 9, UserID: 5, TargetHours: 10}, nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(300, nil)
		uc := usecase.NewGetGoalProgressUseCase(logs, goals)

		progress, err := uc.Execute(context.Background(), 9, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(9), progress.GoalID)
		assert.Equal(t, 10, progress.TargetHours)
		assert.Equal(t, 300, progress.ActualMinutes)
		assert.Equal(t, 50, progress.Percentage)
	})

	t.Run("他人のゴールは 403", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).Return(&model.LearningGoal{ID: 9, UserID: 2}, nil)
		uc := usecase.NewGetGoalProgressUseCase(logs, goals)

		_, err := uc.Execute(context.Background(), 9, 5)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
	})

	t.Run("集計の DB 障害はそのまま伝播する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		goals := new(mockLearningGoalLinker)
		goals.On("FindByID", mock.Anything, uint(9)).
			Return(&model.LearningGoal{ID: 9, UserID: 5, TargetHours: 10}, nil)
		logs.On("SumDurationByGoalID", mock.Anything, uint(9)).Return(0, errors.New("db error"))
		uc := usecase.NewGetGoalProgressUseCase(logs, goals)

		_, err := uc.Execute(context.Background(), 9, 5)

		assert.EqualError(t, err, "db error")
	})

	t.Run("ゴール連携が無効なら 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewGetGoalProgressUseCase(logs, nil)

		_, err := uc.Execute(context.Background(), 9, 5)

		assert.NotNil(t, domain.GetDomainError(err))
	})
}

func TestCalculateGoalProgressPercentage(t *testing.T) {
	cases := []struct {
		name         string
		totalMinutes int
		targetHours  int
		want         int
	}{
		{"目標未設定は 0", 600, 0, 0},
		{"目標が負なら 0", 600, -1, 0},
		{"半分", 300, 10, 50},
		{"ちょうど 100%", 600, 10, 100},
		{"100% を超えても 100 で止める", 1200, 10, 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, usecase.CalculateGoalProgressPercentage(c.totalMinutes, c.targetHours))
		})
	}
}

// ============================================================
// エクスポート / インポート
// ============================================================

func TestExportLearningLogsCSVUseCase_Execute(t *testing.T) {
	t.Run("BOM 付きヘッダーと行を出力する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		created := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
		logs.On("GetByPeriod", mock.Anything, uint(5), 30).Return([]model.LearningLog{
			{Title: "Go学習", Content: "メモ", Category: model.LogCategoryCoding, Duration: 60, CreatedAt: created},
		}, nil)
		uc := usecase.NewExportLearningLogsCSVUseCase(logs)

		data, err := uc.Execute(context.Background(), 5, 30)

		require.NoError(t, err)
		out := string(data)
		assert.True(t, strings.HasPrefix(out, "\xef\xbb\xbf"), "BOM で始まる")
		assert.Contains(t, out, "日付,カテゴリ,タイトル,学習時間(分),メモ")
		assert.Contains(t, out, "2026-01-15,coding,Go学習,60,メモ")
	})

	t.Run("期間が負なら 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewExportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, -1)

		require.Error(t, err)
		assert.NotNil(t, domain.GetDomainError(err))
		logs.AssertNotCalled(t, "GetByPeriod", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestExportLearningLogsJSONUseCase_Execute(t *testing.T) {
	t.Run("日付を整形して出力する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		created := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
		logs.On("GetByPeriod", mock.Anything, uint(5), 0).Return([]model.LearningLog{
			{Title: "Go学習", Content: "メモ", Category: model.LogCategoryCoding, Duration: 60, CreatedAt: created},
		}, nil)
		uc := usecase.NewExportLearningLogsJSONUseCase(logs)

		data, err := uc.Execute(context.Background(), 5, 0)

		require.NoError(t, err)
		var entries []map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &entries))
		require.Len(t, entries, 1)
		assert.Equal(t, "2026-01-15", entries[0]["date"])
		assert.Equal(t, "Go学習", entries[0]["title"])
		assert.Equal(t, float64(60), entries[0]["duration"])
	})

	t.Run("0 件でも空配列を出力する", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("GetByPeriod", mock.Anything, uint(5), 30).Return([]model.LearningLog(nil), nil)
		uc := usecase.NewExportLearningLogsJSONUseCase(logs)

		data, err := uc.Execute(context.Background(), 5, 30)

		require.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("期間が負なら 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewExportLearningLogsJSONUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, -1)

		require.Error(t, err)
		logs.AssertNotCalled(t, "GetByPeriod", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestImportLearningLogsCSVUseCase_Execute(t *testing.T) {
	const header = "日付,カテゴリ,タイトル,学習時間(分),メモ\n"

	t.Run("BOM 付き CSV を取り込む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(l []model.LearningLog) bool {
			return len(l) == 1 && l[0].UserID == 5 && l[0].Title == "Go学習" &&
				l[0].Duration == 60 && l[0].Source == model.LogSourceManual &&
				l[0].CreatedAt.Format("2006-01-02") == "2026-01-15"
		})).Return(nil)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		result, err := uc.Execute(context.Background(), 5, []byte("\xef\xbb\xbf"+header+"2026-01-15,coding,Go学習,60,メモ\n"))

		require.NoError(t, err)
		assert.Len(t, result, 1)
		logs.AssertExpectations(t)
	})

	t.Run("メモが空ならタイトルで補う", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(l []model.LearningLog) bool {
			return l[0].Content == "Go学習"
		})).Return(nil)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, []byte(header+"2026-01-15,coding,Go学習,60,\n"))

		require.NoError(t, err)
	})

	t.Run("カテゴリが空ならその他で補う", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.MatchedBy(func(l []model.LearningLog) bool {
			return l[0].Category == model.LogCategoryOther
		})).Return(nil)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, []byte(header+"2026-01-15,,Go学習,60,メモ\n"))

		require.NoError(t, err)
	})

	t.Run("行ごとのエラーは行番号つきで返す", func(t *testing.T) {
		cases := []struct {
			name string
			body string
			msg  string
		}{
			{"日付が不正", "xxxx,coding,Go学習,60,メモ\n", "CSV 2行目の日付形式が不正です"},
			{"無効なカテゴリ", "2026-01-15,unknown,Go学習,60,メモ\n", "CSV 2行目の無効なカテゴリです"},
			{"タイトルが空", "2026-01-15,coding, ,60,メモ\n", "CSV 2行目のタイトルが空です"},
			{"学習時間が数値でない", "2026-01-15,coding,Go学習,abc,メモ\n", "CSV 2行目の学習時間が不正です"},
			{"学習時間が範囲外", "2026-01-15,coding,Go学習,1441,メモ\n", "CSV 2行目: 学習時間は0〜1440分の範囲で指定してください"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				logs := new(mockLearningLogRepo)
				uc := usecase.NewImportLearningLogsCSVUseCase(logs)

				_, err := uc.Execute(context.Background(), 5, []byte(header+c.body))

				domainErr := domain.GetDomainError(err)
				require.NotNil(t, domainErr)
				assert.Contains(t, domainErr.Message, c.msg)
				logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("ヘッダーだけなら 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, []byte(header))

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "インポートするデータがありません", domainErr.Message)
	})

	t.Run("空データは 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, nil)

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "CSVファイルの読み取りに失敗しました", domainErr.Message)
	})

	t.Run("カラム数が足りないヘッダーは 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, []byte("日付,カテゴリ\n"))

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Contains(t, domainErr.Message, "5列必要")
	})

	t.Run("51 件以上は 400", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		body := strings.Builder{}
		body.WriteString(header)
		for i := 0; i < 51; i++ {
			body.WriteString("2026-01-15,coding,Go学習,60,メモ\n")
		}
		_, err := uc.Execute(context.Background(), 5, []byte(body.String()))

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "一度にインポートできるのは50件までです", domainErr.Message)
		logs.AssertNotCalled(t, "CreateBatch", mock.Anything, mock.Anything)
	})

	t.Run("書き込み失敗は内部エラーに包む", func(t *testing.T) {
		logs := new(mockLearningLogRepo)
		logs.On("CreateBatch", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewImportLearningLogsCSVUseCase(logs)

		_, err := uc.Execute(context.Background(), 5, []byte(header+"2026-01-15,coding,Go学習,60,メモ\n"))

		domainErr := domain.GetDomainError(err)
		require.NotNil(t, domainErr)
		assert.Equal(t, "学習ログのインポートに失敗しました", domainErr.Message)
	})
}
