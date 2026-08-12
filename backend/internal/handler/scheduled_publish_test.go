package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- SchedulePublish ----------

func TestSchedulePublish_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	ports.Posts.On("FindByID", mock.Anything, uint(1)).Return(post, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt != nil
	})).Return(nil)

	scheduledAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": scheduledAt,
	})

	assertStatus(t, w, http.StatusOK)
}

func TestSchedulePublish_InvalidDate(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": "invalid-date",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestSchedulePublish_PastDate(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	ports.Posts.On("FindByID", mock.Anything, uint(1)).Return(post, nil)

	pastTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{
		"scheduled_at": pastTime,
	})

	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CancelSchedule ----------

func TestCancelSchedule_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/cancel-schedule", h.CancelSchedule)

	scheduled := time.Now().Add(24 * time.Hour)
	post := &model.Post{ID: 1, UserID: 1, IsDraft: true, ScheduledAt: &scheduled}
	ports.Posts.On("FindByID", mock.Anything, uint(1)).Return(post, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt == nil
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/1/cancel-schedule", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetScheduled ----------

func TestGetScheduled_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/scheduled", h.GetScheduled)

	scheduled := time.Now().Add(24 * time.Hour)
	ports.Posts.On("FindScheduledByUserID", mock.Anything, uint(1)).Return([]model.Post{
		{ID: 1, Title: "予定投稿", ScheduledAt: &scheduled},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/scheduled", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestSchedulePublish_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	ports.Posts.On("FindByID", mock.Anything, uint(1)).
		Return(&model.Post{ID: 1, UserID: 999, IsDraft: true}, nil)

	scheduledAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{"scheduled_at": scheduledAt})

	assertStatus(t, w, http.StatusForbidden)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 公開済みの投稿はスケジュールできない。
func TestSchedulePublish_NotDraft(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/schedule", h.SchedulePublish)

	ports.Posts.On("FindByID", mock.Anything, uint(1)).
		Return(&model.Post{ID: 1, UserID: 1, IsDraft: false}, nil)

	scheduledAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	w := doRequest(r, http.MethodPut, "/posts/1/schedule", map[string]string{"scheduled_at": scheduledAt})

	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// スケジュールされていない投稿の解除は 400 になる。
func TestCancelSchedule_NotScheduled(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/cancel-schedule", h.CancelSchedule)

	ports.Posts.On("FindByID", mock.Anything, uint(1)).
		Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)

	w := doRequest(r, http.MethodPut, "/posts/1/cancel-schedule", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 予約投稿が無ければ null ではなく空配列を返す。
func TestGetScheduled_Empty(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/scheduled", h.GetScheduled)

	ports.Posts.On("FindScheduledByUserID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/posts/scheduled", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}
