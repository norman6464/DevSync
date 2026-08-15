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

// mockResourceReviewRepo は usecase/repository.ResourceReviewRepository のモック。
type mockResourceReviewRepo struct{ mock.Mock }

func (m *mockResourceReviewRepo) Create(ctx context.Context, review *model.ResourceReview) error {
	return m.Called(ctx, review).Error(0)
}
func (m *mockResourceReviewRepo) FindByID(ctx context.Context, id uint) (*model.ResourceReview, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.ResourceReview)
	return r, args.Error(1)
}
func (m *mockResourceReviewRepo) FindByResourceID(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	args := m.Called(ctx, resourceID, limit, offset)
	reviews, _ := args.Get(0).([]model.ResourceReview)
	return reviews, args.Get(1).(int64), args.Error(2)
}
func (m *mockResourceReviewRepo) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceReview, error) {
	args := m.Called(ctx, userID, resourceID)
	r, _ := args.Get(0).(*model.ResourceReview)
	return r, args.Error(1)
}
func (m *mockResourceReviewRepo) Update(ctx context.Context, review *model.ResourceReview) error {
	return m.Called(ctx, review).Error(0)
}
func (m *mockResourceReviewRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockLearningResourceReader は usecase/repository.LearningResourceReader のモック。
type mockLearningResourceReader struct{ mock.Mock }

func (m *mockLearningResourceReader) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.LearningResource)
	return r, args.Error(1)
}

func TestCreateResourceReviewUseCase_Execute(t *testing.T) {
	t.Run("存在するリソースへ重複なくレビューできる", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		reviews.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).Return((*model.ResourceReview)(nil), nil)
		reviews.On("Create", mock.Anything, mock.AnythingOfType("*model.ResourceReview")).Return(nil)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 4, Comment: "とても良い"})

		assert.NoError(t, err)
		reviews.AssertExpectations(t)
	})

	t.Run("存在しないリソースは ErrNotFound", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(99)).Return((*model.LearningResource)(nil), nil)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 99, Rating: 4})

		assert.ErrorIs(t, err, domain.ErrNotFound)
		reviews.AssertNotCalled(t, "Create")
	})

	// 存在確認の失敗は「リソースが無い」こととは別で、404 に変換すると障害が隠れる。
	t.Run("存在確認の DB 障害は 404 にせず伝播する", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		dbErr := errors.New("db down")
		resources.On("FindByID", mock.Anything, uint(10)).Return((*model.LearningResource)(nil), dbErr)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 4})

		assert.ErrorIs(t, err, dbErr)
		assert.NotErrorIs(t, err, domain.ErrNotFound, "障害を 404 に変換しない")
		reviews.AssertNotCalled(t, "Create")
	})

	// 重複チェックが失敗したまま作成すると、障害時に二重レビューを許してしまう。
	t.Run("重複チェックの DB 障害では作成しない", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		dbErr := errors.New("db down")
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		reviews.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
			Return((*model.ResourceReview)(nil), dbErr)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 4})

		assert.ErrorIs(t, err, dbErr)
		reviews.AssertNotCalled(t, "Create")
	})

	t.Run("評価が範囲外なら 400", func(t *testing.T) {
		for _, rating := range []int{0, 6} {
			reviews := new(mockResourceReviewRepo)
			resources := new(mockLearningResourceReader)
			resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
			uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

			err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: rating})

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "1〜5")
			reviews.AssertNotCalled(t, "Create")
		}
	})

	t.Run("二重レビューは ErrConflict", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		reviews.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
			Return(&model.ResourceReview{UserID: 1, ResourceID: 10}, nil)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 4})

		assert.ErrorIs(t, err, domain.ErrConflict)
		reviews.AssertNotCalled(t, "Create")
	})

	t.Run("コメントが長すぎると 400", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		resources := new(mockLearningResourceReader)
		resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
		uc := usecase.NewCreateResourceReviewUseCase(reviews, resources)

		err := uc.Execute(context.Background(), &model.ResourceReview{
			UserID: 1, ResourceID: 10, Rating: 4, Comment: string(make([]rune, 5001)),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "コメント")
		reviews.AssertNotCalled(t, "Create")
	})
}

func TestListResourceReviewsUseCase_Execute(t *testing.T) {
	reviews := new(mockResourceReviewRepo)
	expected := []model.ResourceReview{{Rating: 5}, {Rating: 3}}
	reviews.On("FindByResourceID", mock.Anything, uint(10), 20, 0).Return(expected, int64(2), nil)
	uc := usecase.NewListResourceReviewsUseCase(reviews)

	got, total, err := uc.Execute(context.Background(), 10, 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, int64(2), total)
	reviews.AssertExpectations(t)
}

func TestUpdateResourceReviewUseCase_Execute(t *testing.T) {
	t.Run("所有者は rating と comment を更新できる", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		existing := &model.ResourceReview{UserID: 1, Rating: 3, Comment: "普通"}
		reviews.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
		reviews.On("Update", mock.Anything, existing).Return(nil)
		uc := usecase.NewUpdateResourceReviewUseCase(reviews)

		result, err := uc.Execute(context.Background(), 1, 1, 5, "最高でした")

		assert.NoError(t, err)
		assert.Equal(t, 5, result.Rating)
		assert.Equal(t, "最高でした", result.Comment)
		reviews.AssertExpectations(t)
	})

	t.Run("他人のレビューは Forbidden", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(1)).Return(&model.ResourceReview{UserID: 1}, nil)
		uc := usecase.NewUpdateResourceReviewUseCase(reviews)

		_, err := uc.Execute(context.Background(), 1, 999, 5, "")

		assert.ErrorIs(t, err, domain.ErrForbidden)
		reviews.AssertNotCalled(t, "Update")
	})

	t.Run("存在しないレビューは finder のエラーを返す", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(99)).Return((*model.ResourceReview)(nil), domain.ErrNotFound)
		uc := usecase.NewUpdateResourceReviewUseCase(reviews)

		_, err := uc.Execute(context.Background(), 99, 1, 5, "")

		assert.ErrorIs(t, err, domain.ErrNotFound)
		reviews.AssertNotCalled(t, "Update")
	})

	t.Run("評価が範囲外なら 400", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(1)).Return(&model.ResourceReview{UserID: 1, Rating: 3}, nil)
		uc := usecase.NewUpdateResourceReviewUseCase(reviews)

		_, err := uc.Execute(context.Background(), 1, 1, 6, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "1〜5")
		reviews.AssertNotCalled(t, "Update")
	})

	t.Run("コメントが長すぎると 400", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(1)).Return(&model.ResourceReview{UserID: 1, Rating: 3}, nil)
		uc := usecase.NewUpdateResourceReviewUseCase(reviews)

		_, err := uc.Execute(context.Background(), 1, 1, 0, string(make([]rune, 5001)))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "コメント")
		reviews.AssertNotCalled(t, "Update")
	})
}

func TestDeleteResourceReviewUseCase_Execute(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(1)).Return(&model.ResourceReview{UserID: 1}, nil)
		reviews.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteResourceReviewUseCase(reviews)

		err := uc.Execute(context.Background(), 1, 1)

		assert.NoError(t, err)
		reviews.AssertExpectations(t)
	})

	t.Run("他人のレビューは Forbidden", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(1)).Return(&model.ResourceReview{UserID: 1}, nil)
		uc := usecase.NewDeleteResourceReviewUseCase(reviews)

		err := uc.Execute(context.Background(), 1, 999)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		reviews.AssertNotCalled(t, "Delete")
	})

	t.Run("存在しないレビューは finder のエラーを返す", func(t *testing.T) {
		reviews := new(mockResourceReviewRepo)
		reviews.On("FindByID", mock.Anything, uint(99)).Return((*model.ResourceReview)(nil), domain.ErrNotFound)
		uc := usecase.NewDeleteResourceReviewUseCase(reviews)

		err := uc.Execute(context.Background(), 99, 1)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		reviews.AssertNotCalled(t, "Delete")
	})
}
