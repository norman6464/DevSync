package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockReminderSettingsRepo は usecase/repository.ReminderSettingsRepository のモック。
type mockReminderSettingsRepo struct{ mock.Mock }

func (m *mockReminderSettingsRepo) GetOrCreateDefault(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ReminderSettings)
	return s, args.Error(1)
}

func (m *mockReminderSettingsRepo) FindByUserID(ctx context.Context, userID uint) (*model.ReminderSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ReminderSettings)
	return s, args.Error(1)
}

func (m *mockReminderSettingsRepo) Save(ctx context.Context, settings *model.ReminderSettings) error {
	return m.Called(ctx, settings).Error(0)
}

// existingReminderSettings は更新テストの起点となる既存設定を返す。
func existingReminderSettings() *model.ReminderSettings {
	return &model.ReminderSettings{
		ID:               10,
		UserID:           1,
		Enabled:          true,
		Frequency:        model.ReminderFrequencyDaily,
		NotificationTime: "09:00",
		InactiveDays:     3,
		EnableWeb:        true,
		EnableEmail:      false,
	}
}

func TestGetReminderSettingsUseCase_Execute(t *testing.T) {
	t.Run("未登録ならデフォルト作成に委譲する", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		want := existingReminderSettings()
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(want, nil)
		uc := usecase.NewGetReminderSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		uc := usecase.NewGetReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), 0)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "GetOrCreateDefault")
	})
}

func TestUpdateReminderSettingsUseCase_Execute(t *testing.T) {
	t.Run("指定された項目を反映して保存する", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(existingReminderSettings(), nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.ReminderSettings) bool {
			return s.ID == 10 &&
				s.Frequency == model.ReminderFrequencyWeekly &&
				s.NotificationTime == "21:30" &&
				s.InactiveDays == 7 &&
				!s.Enabled && !s.EnableWeb && s.EnableEmail
		})).Return(nil)
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{
			UserID:           1,
			Enabled:          false,
			Frequency:        model.ReminderFrequencyWeekly,
			NotificationTime: "21:30",
			InactiveDays:     7,
			EnableWeb:        false,
			EnableEmail:      true,
		})

		assert.NoError(t, err)
		assert.Equal(t, model.ReminderFrequencyWeekly, got.Frequency)
		assert.Equal(t, "21:30", got.NotificationTime)
		assert.Equal(t, 7, got.InactiveDays)
		repo.AssertExpectations(t)
	})

	// 文字列と数値は「空文字 / 0 なら据え置き」。真偽値は常に上書きされる。
	t.Run("空文字と 0 は既存値を保つが真偽値は常に上書きする", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(existingReminderSettings(), nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.ReminderSettings) bool {
			return s.Frequency == model.ReminderFrequencyDaily &&
				s.NotificationTime == "09:00" &&
				s.InactiveDays == 3 &&
				!s.Enabled && !s.EnableWeb && !s.EnableEmail
		})).Return(nil)
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{UserID: 1})

		assert.NoError(t, err)
		assert.Equal(t, model.ReminderFrequencyDaily, got.Frequency)
		assert.Equal(t, "09:00", got.NotificationTime)
		assert.Equal(t, 3, got.InactiveDays)
		assert.False(t, got.Enabled)
		repo.AssertExpectations(t)
	})

	t.Run("バリデーション違反は BadRequest（保存しない）", func(t *testing.T) {
		cases := []struct {
			name    string
			in      usecase.UpdateReminderSettingsInput
			wantMsg string
		}{
			{"未知の頻度", usecase.UpdateReminderSettingsInput{UserID: 1, Frequency: "hourly"}, "頻度はdailyまたはweeklyのみ有効です"},
			{"時刻の形式違反", usecase.UpdateReminderSettingsInput{UserID: 1, NotificationTime: "9:00"}, "通知時間はHH:MM形式で指定してください"},
			{"時刻の範囲外", usecase.UpdateReminderSettingsInput{UserID: 1, NotificationTime: "24:00"}, "通知時間はHH:MM形式で指定してください"},
			{"非活動日数の上限超過", usecase.UpdateReminderSettingsInput{UserID: 1, InactiveDays: 31}, "非活動日数は1〜30の範囲で指定してください"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				repo := new(mockReminderSettingsRepo)
				repo.On("FindByUserID", mock.Anything, uint(1)).Return(existingReminderSettings(), nil)
				uc := usecase.NewUpdateReminderSettingsUseCase(repo)

				_, err := uc.Execute(context.Background(), tc.in)

				assertDomainCode(t, err, domain.ErrCodeBadRequest)
				var de *domain.DomainError
				if assert.ErrorAs(t, err, &de) {
					assert.Equal(t, tc.wantMsg, de.Message)
				}
				repo.AssertNotCalled(t, "Save")
			})
		}
	})

	// 境界値: 上限ちょうどは許可される。
	t.Run("非活動日数の上限ちょうどは許可する", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(existingReminderSettings(), nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.ReminderSettings) bool {
			return s.InactiveDays == 30
		})).Return(nil)
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{UserID: 1, InactiveDays: 30})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("設定が未登録なら DomainError ではないエラーを返す（保存しない）", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{UserID: 1})

		assert.Error(t, err)
		var de *domain.DomainError
		assert.False(t, errors.As(err, &de), "500 を維持するため DomainError にしない")
		repo.AssertNotCalled(t, "Save")
	})

	t.Run("取得時の DB 障害を伝播する（保存しない）", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{UserID: 1})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Save")
	})

	t.Run("保存時の DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		repo.On("FindByUserID", mock.Anything, uint(1)).Return(existingReminderSettings(), nil)
		repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{UserID: 1})

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockReminderSettingsRepo)
		uc := usecase.NewUpdateReminderSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateReminderSettingsInput{})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "FindByUserID")
		repo.AssertNotCalled(t, "Save")
	})
}
