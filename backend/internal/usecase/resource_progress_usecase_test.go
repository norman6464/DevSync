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

// mockResourceProgressRepo は usecase/repository.ResourceProgressRepository のモック。
// （LearningResourceReader のモックは resource_review のテストで定義済みのものを再利用する）
type mockResourceProgressRepo struct{ mock.Mock }

func (m *mockResourceProgressRepo) Upsert(ctx context.Context, progress *model.ResourceProgress) error {
	return m.Called(ctx, progress).Error(0)
}
func (m *mockResourceProgressRepo) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error) {
	args := m.Called(ctx, userID, resourceID)
	p, _ := args.Get(0).(*model.ResourceProgress)
	return p, args.Error(1)
}
func (m *mockResourceProgressRepo) FindByUserID(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	list, _ := args.Get(0).([]model.ResourceProgress)
	return list, args.Get(1).(int64), args.Error(2)
}

func TestUpsertResourceProgressUseCase_Execute(t *testing.T) {
	t.Run("in_progress で作成し StartedAt が設定される", func(t *testing.T) {
		progress := new(mockResourceProgressRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		progress.On("Upsert", mock.Anything, mock.MatchedBy(func(p *model.ResourceProgress) bool {
			return p.Status == model.ResourceProgressInProgress && p.StartedAt != nil && p.CompletedAt == nil && p.CompletionPercent == 50
		})).Return(nil)
		progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
			Return(&model.ResourceProgress{UserID: 1, ResourceID: 10, CompletionPercent: 50}, nil)
		uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

		result, err := uc.Execute(context.Background(), 1, 10, "in_progress", 50, "学習中")

		assert.NoError(t, err)
		assert.Equal(t, 50, result.CompletionPercent)
		progress.AssertExpectations(t)
	})

	t.Run("completed で StartedAt と CompletedAt が設定される", func(t *testing.T) {
		progress := new(mockResourceProgressRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		progress.On("Upsert", mock.Anything, mock.MatchedBy(func(p *model.ResourceProgress) bool {
			return p.Status == model.ResourceProgressCompleted && p.StartedAt != nil && p.CompletedAt != nil
		})).Return(nil)
		progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
			Return(&model.ResourceProgress{UserID: 1, ResourceID: 10}, nil)
		uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

		_, err := uc.Execute(context.Background(), 1, 10, "completed", 100, "")

		assert.NoError(t, err)
		progress.AssertExpectations(t)
	})

	t.Run("存在しないリソースは 404", func(t *testing.T) {
		progress := new(mockResourceProgressRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(99)).Return((*model.LearningResource)(nil), errors.New("not found"))
		uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

		_, err := uc.Execute(context.Background(), 1, 99, "in_progress", 50, "")

		var de *domain.DomainError
		assert.ErrorAs(t, err, &de)
		assert.Equal(t, domain.ErrCodeNotFound, de.Code)
		progress.AssertNotCalled(t, "Upsert")
	})

	t.Run("無効なステータスは 400", func(t *testing.T) {
		progress := new(mockResourceProgressRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

		_, err := uc.Execute(context.Background(), 1, 10, "invalid", 50, "")

		var de *domain.DomainError
		assert.ErrorAs(t, err, &de)
		assert.Equal(t, domain.ErrCodeValidation, de.Code)
		progress.AssertNotCalled(t, "Upsert")
	})

	t.Run("進捗率が範囲外は 400", func(t *testing.T) {
		for _, pct := range []int{-1, 101} {
			progress := new(mockResourceProgressRepo)
			resources := new(mockLearningResourceReader)
			resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
			uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

			_, err := uc.Execute(context.Background(), 1, 10, "in_progress", pct, "")

			assert.Error(t, err)
			progress.AssertNotCalled(t, "Upsert")
		}
	})

	t.Run("メモが長すぎると 400", func(t *testing.T) {
		progress := new(mockResourceProgressRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		uc := usecase.NewUpsertResourceProgressUseCase(progress, resources)

		_, err := uc.Execute(context.Background(), 1, 10, "in_progress", 50, string(make([]rune, 1001)))

		assert.Error(t, err)
		progress.AssertNotCalled(t, "Upsert")
	})
}

func TestGetResourceProgressUseCase_Execute(t *testing.T) {
	progress := new(mockResourceProgressRepo)
	expected := &model.ResourceProgress{UserID: 1, ResourceID: 10, CompletionPercent: 75}
	progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).Return(expected, nil)
	uc := usecase.NewGetResourceProgressUseCase(progress)

	got, err := uc.Execute(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	progress.AssertExpectations(t)
}

func TestListResourceProgressUseCase_Execute(t *testing.T) {
	progress := new(mockResourceProgressRepo)
	expected := []model.ResourceProgress{{ID: 1}, {ID: 2}}
	progress.On("FindByUserID", mock.Anything, uint(1), "completed", 20, 0).Return(expected, int64(2), nil)
	uc := usecase.NewListResourceProgressUseCase(progress)

	list, total, err := uc.Execute(context.Background(), 1, "completed", 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, expected, list)
	assert.Equal(t, int64(2), total)
	progress.AssertExpectations(t)
}
