package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockWidgetSettingsRepo は usecase/repository.WidgetSettingsRepository のモック。
type mockWidgetSettingsRepo struct{ mock.Mock }

func (m *mockWidgetSettingsRepo) FindByUserID(ctx context.Context, userID uint) (*model.WidgetSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.WidgetSettings)
	return s, args.Error(1)
}

func (m *mockWidgetSettingsRepo) Upsert(ctx context.Context, settings *model.WidgetSettings) error {
	return m.Called(ctx, settings).Error(0)
}

func TestGetWidgetSettingsUseCase_Execute(t *testing.T) {
	t.Run("登録済みの設定をそのまま返す", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		expected := &model.WidgetSettings{UserID: 3, Settings: `[{"key":"level","visible":true,"order":0}]`}
		repo.On("FindByUserID", mock.Anything, uint(3)).Return(expected, nil)
		uc := usecase.NewGetWidgetSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), 3)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("未登録ならデフォルト設定を返す", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(3)).
			Return((*model.WidgetSettings)(nil), errors.New("record not found"))
		uc := usecase.NewGetWidgetSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), 3)

		assert.NoError(t, err)
		assert.Equal(t, uint(3), got.UserID)
		// デフォルトは 14 項目のウィジェット配置
		assert.Contains(t, got.Settings, `"key":"userProfile"`)
		assert.Contains(t, got.Settings, `"key":"quickStats"`)
		repo.AssertExpectations(t)
	})
}

func TestUpdateWidgetSettingsUseCase_Execute(t *testing.T) {
	valid := `[{"key":"userProfile","visible":true,"order":0}]`

	t.Run("JSON 配列なら Upsert する", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		repo.On("Upsert", mock.Anything, mock.MatchedBy(func(s *model.WidgetSettings) bool {
			return s.UserID == 3 && s.Settings == valid
		})).Return(nil)
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		err := uc.Execute(context.Background(), 3, valid)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("空文字は 400（Upsert しない）", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		err := uc.Execute(context.Background(), 3, "")

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Upsert")
	})

	t.Run("10000 文字超は 400（Upsert しない）", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		long := "[\"" + strings.Repeat("a", 10001) + "\"]"
		err := uc.Execute(context.Background(), 3, long)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Upsert")
	})

	t.Run("不正な JSON は 400（Upsert しない）", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		err := uc.Execute(context.Background(), 3, "not json")

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Upsert")
	})

	t.Run("JSON だが配列でない場合は 400（Upsert しない）", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		err := uc.Execute(context.Background(), 3, `{"key":"userProfile"}`)

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Upsert")
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockWidgetSettingsRepo)
		repo.On("Upsert", mock.Anything, mock.AnythingOfType("*model.WidgetSettings")).
			Return(errors.New("db error"))
		uc := usecase.NewUpdateWidgetSettingsUseCase(repo)

		err := uc.Execute(context.Background(), 3, valid)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
