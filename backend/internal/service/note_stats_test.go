package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestNoteStatsService はテスト用のNoteStatsServiceを生成する。
func newTestNoteStatsService() (*NoteStatsService, *MockNoteStatsRepository) {
	repo := new(MockNoteStatsRepository)
	svc := NewNoteStatsService(repo)
	return svc, repo
}

// ============================================================
// GetNoteStats テスト
// ============================================================

func TestNoteStatsService_GetNoteStats_Success(t *testing.T) {
	svc, repo := newTestNoteStatsService()

	expected := &model.NoteStats{
		TotalNotes:    20,
		ArchivedNotes: 3,
		FavoriteNotes: 5,
		TotalFolders:  4,
		NotesThisWeek: 2,
	}
	repo.On("GetNoteStats", uint(1)).Return(expected, nil)

	result, err := svc.GetNoteStats(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(20), result.TotalNotes)
	assert.Equal(t, int64(3), result.ArchivedNotes)
	assert.Equal(t, int64(5), result.FavoriteNotes)
	assert.Equal(t, int64(4), result.TotalFolders)
	assert.Equal(t, int64(2), result.NotesThisWeek)
	repo.AssertExpectations(t)
}

func TestNoteStatsService_GetNoteStats_NewUser(t *testing.T) {
	svc, repo := newTestNoteStatsService()

	empty := &model.NoteStats{}
	repo.On("GetNoteStats", uint(2)).Return(empty, nil)

	result, err := svc.GetNoteStats(2)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalNotes)
	assert.Equal(t, int64(0), result.TotalFolders)
	repo.AssertExpectations(t)
}

func TestNoteStatsService_GetNoteStats_InvalidUserID(t *testing.T) {
	svc, _ := newTestNoteStatsService()

	result, err := svc.GetNoteStats(0)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestNoteStatsService_GetNoteStats_RepoError(t *testing.T) {
	svc, repo := newTestNoteStatsService()

	repo.On("GetNoteStats", uint(1)).Return((*model.NoteStats)(nil), errors.New("db error"))

	result, err := svc.GetNoteStats(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
