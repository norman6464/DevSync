package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNoteStatsRepo は usecase/repository.NoteStatsRepository のモック。
type mockNoteStatsRepo struct{ mock.Mock }

func (m *mockNoteStatsRepo) GetNoteStats(ctx context.Context, userID uint) (*model.NoteStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NoteStats)
	return s, args.Error(1)
}

func TestGetNoteStatsUseCase_Execute(t *testing.T) {
	t.Run("集計統計を返す", func(t *testing.T) {
		stats := new(mockNoteStatsRepo)
		expected := &model.NoteStats{TotalNotes: 5, ArchivedNotes: 1, FavoriteNotes: 2, TotalFolders: 3, NotesThisWeek: 4}
		stats.On("GetNoteStats", mock.Anything, uint(10)).Return(expected, nil)
		uc := usecase.NewGetNoteStatsUseCase(stats)

		got, err := uc.Execute(context.Background(), 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		stats.AssertExpectations(t)
	})

	t.Run("userID が 0 は 400（repo を呼ばない）", func(t *testing.T) {
		stats := new(mockNoteStatsRepo)
		uc := usecase.NewGetNoteStatsUseCase(stats)

		_, err := uc.Execute(context.Background(), 0)

		assert.Error(t, err)
		stats.AssertNotCalled(t, "GetNoteStats")
	})
}
