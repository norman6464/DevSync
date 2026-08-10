package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCodeSnippetStatsRepo は usecase/repository.CodeSnippetStatsRepository のモック。
type mockCodeSnippetStatsRepo struct{ mock.Mock }

func (m *mockCodeSnippetStatsRepo) GetCodeSnippetStats(ctx context.Context, userID uint) (*model.CodeSnippetStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.CodeSnippetStats)
	return s, args.Error(1)
}

func TestGetCodeSnippetStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockCodeSnippetStatsRepo)
		expected := &model.CodeSnippetStats{TotalSnippets: 4, TotalComments: 9, LanguageCount: 3}
		stats.On("GetCodeSnippetStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetCodeSnippetStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockCodeSnippetStatsRepo)
		uc := usecase.NewGetCodeSnippetStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetCodeSnippetStats")
	})
}
