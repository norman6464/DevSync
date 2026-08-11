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

// mockNotificationSettingsRepo は usecase/repository.NotificationSettingsRepository のモック。
type mockNotificationSettingsRepo struct{ mock.Mock }

func (m *mockNotificationSettingsRepo) GetOrCreateDefault(ctx context.Context, userID uint) (*model.NotificationSettings, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NotificationSettings)
	return s, args.Error(1)
}

func (m *mockNotificationSettingsRepo) Save(ctx context.Context, settings *model.NotificationSettings) error {
	return m.Called(ctx, settings).Error(0)
}

// allEnabledSettings は 8 項目すべて有効な設定を返す。
func allEnabledSettings() *model.NotificationSettings {
	return &model.NotificationSettings{
		ID: 1, UserID: 1,
		EnableLikes: true, EnableComments: true, EnableFollows: true, EnableMessages: true,
		EnableMentions: true, EnableWebPush: true, EnableEmail: true, EnableSound: true,
	}
}

func TestGetNotificationSettingsUseCase_Execute(t *testing.T) {
	t.Run("未登録ならデフォルト作成に委譲する", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		want := allEnabledSettings()
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(want, nil)
		uc := usecase.NewGetNotificationSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
		repo.AssertExpectations(t)
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		uc := usecase.NewGetNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), 0)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "GetOrCreateDefault")
	})
}

func TestUpdateNotificationSettingsUseCase_Execute(t *testing.T) {
	// 8 項目すべてが入力の値で上書きされる（部分更新ではない）。
	t.Run("すべての項目を上書きする", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(allEnabledSettings(), nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.NotificationSettings) bool {
			return !s.EnableLikes && s.EnableComments && !s.EnableFollows && !s.EnableMessages &&
				!s.EnableMentions && !s.EnableWebPush && !s.EnableEmail && !s.EnableSound
		})).Return(nil)
		uc := usecase.NewUpdateNotificationSettingsUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.UpdateNotificationSettingsInput{
			UserID: 1, EnableComments: true,
		})

		assert.NoError(t, err)
		assert.False(t, got.EnableLikes)
		assert.True(t, got.EnableComments)
		repo.AssertExpectations(t)
	})

	t.Run("すべて有効にもできる", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		off := &model.NotificationSettings{ID: 1, UserID: 1}
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(off, nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(s *model.NotificationSettings) bool {
			return s.EnableLikes && s.EnableComments && s.EnableFollows && s.EnableMessages &&
				s.EnableMentions && s.EnableWebPush && s.EnableEmail && s.EnableSound
		})).Return(nil)
		uc := usecase.NewUpdateNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNotificationSettingsInput{
			UserID: 1, EnableLikes: true, EnableComments: true, EnableFollows: true, EnableMessages: true,
			EnableMentions: true, EnableWebPush: true, EnableEmail: true, EnableSound: true,
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("取得時の DB 障害を伝播する（保存しない）", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewUpdateNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNotificationSettingsInput{UserID: 1})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Save")
	})

	t.Run("保存時の DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		repo.On("GetOrCreateDefault", mock.Anything, uint(1)).Return(allEnabledSettings(), nil)
		repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNotificationSettingsInput{UserID: 1})

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockNotificationSettingsRepo)
		uc := usecase.NewUpdateNotificationSettingsUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNotificationSettingsInput{})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "GetOrCreateDefault")
		repo.AssertNotCalled(t, "Save")
	})
}
