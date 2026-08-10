package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockNoteStatsRepo は usecase/repository.NoteStatsRepository のモック（ctx 付き）。
type mockNoteStatsRepo struct{ mock.Mock }

func (m *mockNoteStatsRepo) GetNoteStats(ctx context.Context, userID uint) (*model.NoteStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NoteStats)
	return s, args.Error(1)
}

func setupNoteStatsHandler() (*NoteStatsHandler, *mockNoteStatsRepo) {
	stats := new(mockNoteStatsRepo)
	h := NewNoteStatsHandler(usecase.NewGetNoteStatsUseCase(stats))
	return h, stats
}

func TestNoteStats_GetStats_Success(t *testing.T) {
	h, stats := setupNoteStatsHandler()
	stats.On("GetNoteStats", mock.Anything, uint(5)).Return(
		&model.NoteStats{TotalNotes: 5, ArchivedNotes: 1, FavoriteNotes: 2, TotalFolders: 3, NotesThisWeek: 4}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/notes", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_notes"].(float64) != 5 {
		t.Errorf("expected total_notes 5, got %v", data["total_notes"])
	}
	stats.AssertExpectations(t)
}

func TestNoteStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupNoteStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/notes", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNoteStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupNoteStatsHandler()
	stats.On("GetNoteStats", mock.Anything, uint(5)).Return((*model.NoteStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/notes", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/notes", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
