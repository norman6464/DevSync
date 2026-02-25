package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// ---------- SchedulePublish ----------

func TestSchedulePublish_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Update", mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt != nil
	})).Return(nil)

	scheduledAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": scheduledAt,
	})

	assertStatus(t, w, http.StatusOK)
}

func TestSchedulePublish_InvalidDate(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": "invalid-date",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestSchedulePublish_PastDate(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	postRepo.On("FindByID", uint(1)).Return(post, nil)

	pastTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": pastTime,
	})

	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CancelSchedule ----------

func TestCancelSchedule_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/cancel-schedule", h.CancelSchedule)

	scheduled := time.Now().Add(24 * time.Hour)
	post := &model.Post{ID: 1, UserID: 1, IsDraft: true, ScheduledAt: &scheduled}
	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Update", mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt == nil
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/1/cancel-schedule", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetScheduled ----------

func TestGetScheduled_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/scheduled", h.GetScheduled)

	scheduled := time.Now().Add(24 * time.Hour)
	postRepo.On("FindScheduledByUserID", uint(1)).Return([]model.Post{
		{ID: 1, Title: "予定投稿", ScheduledAt: &scheduled},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/scheduled", nil)
	assertStatus(t, w, http.StatusOK)
}
