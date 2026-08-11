package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockBookReviewRepo は usecase/repository.BookReviewRepository のモック（ctx 付き）。
type mockBookReviewRepo struct{ mock.Mock }

func (m *mockBookReviewRepo) Create(ctx context.Context, r *model.BookReview) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockBookReviewRepo) FindByID(ctx context.Context, id uint) (*model.BookReview, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.BookReview)
	return r, args.Error(1)
}

func (m *mockBookReviewRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) FindAll(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) FindByRating(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	args := m.Called(ctx, userID, minRating, maxRating)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Error(1)
}

func (m *mockBookReviewRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) Update(ctx context.Context, r *model.BookReview) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockBookReviewRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockBookReviewRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupBookReviewHandler は本物の usecase と port モックで BookReviewHandler を組む。
func setupBookReviewHandler() (*BookReviewHandler, *mockBookReviewRepo) {
	repo := new(mockBookReviewRepo)
	h := NewBookReviewHandler(
		usecase.NewCreateBookReviewUseCase(repo),
		usecase.NewGetBookReviewUseCase(repo),
		usecase.NewListBookReviewsByUserUseCase(repo),
		usecase.NewListAllBookReviewsUseCase(repo),
		usecase.NewListBookReviewsByRatingUseCase(repo),
		usecase.NewSearchBookReviewsUseCase(repo),
		usecase.NewUpdateBookReviewUseCase(repo),
		usecase.NewUpdateBookReviewStatusUseCase(repo),
		usecase.NewArchiveBookReviewUseCase(repo),
		usecase.NewUpdateBookReviewProgressUseCase(repo),
		usecase.NewDeleteBookReviewUseCase(repo),
		usecase.NewCountBookReviewsUseCase(repo),
	)
	return h, repo
}

// ownedReview は所有者が userID=1 の書籍レビューを返す。
func ownedReview(id uint) *model.BookReview {
	r := &model.BookReview{
		UserID: 1, Title: "元のタイトル", Author: "元の著者", Rating: 3,
		TotalPages: 100, Status: model.ReviewStatusNotStarted,
	}
	r.ID = id
	return r
}

// --- Create ---

func TestBookReview_Create_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		return r.Title == "書名" && r.Author == "著者" && r.Rating == 4
	})).Return(nil)

	r := newRouter(1)
	r.POST("/book-reviews", h.Create)
	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"title": "  書名  ", "author": "  著者  ", "rating": 4,
	})

	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

// 評価が範囲外なら 400 を返し、作成しない。
func TestBookReview_Create_InvalidRating(t *testing.T) {
	for _, rating := range []int{0, 6} {
		h, repo := setupBookReviewHandler()

		r := newRouter(1)
		r.POST("/book-reviews", h.Create)
		w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
			"title": "書名", "rating": rating,
		})

		assertStatus(t, w, http.StatusBadRequest)
		repo.AssertNotCalled(t, "Create")
	}
}

// 総ページ数が上限超過なら 400 を返し、作成しない。
func TestBookReview_Create_TotalPagesTooLarge(t *testing.T) {
	h, repo := setupBookReviewHandler()

	r := newRouter(1)
	r.POST("/book-reviews", h.Create)
	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"title": "書名", "rating": 3, "total_pages": 100000,
	})

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create")
}

// --- Get / List ---

func TestBookReview_GetByID_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)

	r := newRouter(1)
	r.GET("/book-reviews/:id", h.GetByID)
	w := doRequest(r, http.MethodGet, "/book-reviews/1", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 存在しない ID は 500（移行前の挙動を維持している）。
func TestBookReview_GetByID_NotFound(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/book-reviews/:id", h.GetByID)
	w := doRequest(r, http.MethodGet, "/book-reviews/1", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookReview_GetAll_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindAll", mock.Anything, 20, 0).
		Return([]model.BookReview{*ownedReview(1)}, int64(1), nil)

	r := newRouter(1)
	r.GET("/book-reviews", h.GetAll)
	w := doRequest(r, http.MethodGet, "/book-reviews", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 評価範囲が不正なら 400 を返し、DB を触らない。
func TestBookReview_GetByRating_InvalidRange(t *testing.T) {
	h, repo := setupBookReviewHandler()

	r := newRouter(1)
	r.GET("/book-reviews/rating", h.GetByRating)
	// パラメータ名は min_rating / max_rating。誤った名前だとクエリ解析側で 400 になり、
	// 評価範囲の検証そのものを通らないため、正しい名前で最小 > 最大を渡す。
	w := doRequest(r, http.MethodGet, "/book-reviews/rating?min_rating=4&max_rating=2", nil)

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByRating")
}

func TestBookReview_Search_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("Search", mock.Anything, "go", 20, 0).
		Return([]model.BookReview{*ownedReview(1)}, int64(1), nil)

	r := newRouter(1)
	r.GET("/book-reviews/search", h.Search)
	w := doRequest(r, http.MethodGet, "/book-reviews/search?q=go", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// --- Update ---

func TestBookReview_Update_PartialKeepsOthers(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		// タイトルだけ変わり、著者と評価は据え置き
		return r.Title == "新タイトル" && r.Author == "元の著者" && r.Rating == 3
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)
	w := doRequest(r, http.MethodPut, "/book-reviews/1", map[string]interface{}{"title": "新タイトル"})

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 所有者以外の更新は 403 を返し、保存しない。
func TestBookReview_Update_Forbidden(t *testing.T) {
	h, repo := setupBookReviewHandler()
	other := ownedReview(1)
	other.UserID = 999
	repo.On("FindByID", mock.Anything, uint(1)).Return(other, nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)
	w := doRequest(r, http.MethodPut, "/book-reviews/1", map[string]interface{}{"title": "乗っ取り"})

	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Update")
}

// --- Status / Archive ---

// 無効なステータスは 400 を返し、DB を触らない。
func TestBookReview_UpdateStatus_Invalid(t *testing.T) {
	h, repo := setupBookReviewHandler()

	r := newRouter(1)
	r.PUT("/book-reviews/:id/status", h.UpdateStatus)
	w := doRequest(r, http.MethodPut, "/book-reviews/1/status", map[string]interface{}{"status": "unknown"})

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID")
	repo.AssertNotCalled(t, "Update")
}

func TestBookReview_Archive_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		return r.IsArchived
	})).Return(nil)

	r := newRouter(1)
	r.POST("/book-reviews/:id/archive", h.Archive)
	w := doRequest(r, http.MethodPost, "/book-reviews/1/archive", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestBookReview_Unarchive_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	archived := ownedReview(1)
	archived.IsArchived = true
	repo.On("FindByID", mock.Anything, uint(1)).Return(archived, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		return !r.IsArchived
	})).Return(nil)

	r := newRouter(1)
	r.POST("/book-reviews/:id/unarchive", h.Unarchive)
	w := doRequest(r, http.MethodPost, "/book-reviews/1/unarchive", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// --- Progress ---

// 総ページ数に到達すると読了へ自動遷移する。
func TestBookReview_UpdateProgress_CompletesAtTotal(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		return r.CurrentPage == 100 && r.Status == model.ReviewStatusCompleted
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id/progress", h.UpdateProgress)
	w := doRequest(r, http.MethodPut, "/book-reviews/1/progress", map[string]interface{}{"current_page": 100})

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 未読から 1 ページ以上で読中へ自動遷移する。
func TestBookReview_UpdateProgress_StartsReading(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
		return r.CurrentPage == 10 && r.Status == model.ReviewStatusReading
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id/progress", h.UpdateProgress)
	w := doRequest(r, http.MethodPut, "/book-reviews/1/progress", map[string]interface{}{"current_page": 10})

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 総ページ数を超える指定は 400 を返し、保存しない。
func TestBookReview_UpdateProgress_ExceedsTotal(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id/progress", h.UpdateProgress)
	w := doRequest(r, http.MethodPut, "/book-reviews/1/progress", map[string]interface{}{"current_page": 101})

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Update")
}

// 総ページ数が未設定なら 400 を返し、保存しない。
func TestBookReview_UpdateProgress_NoTotalPages(t *testing.T) {
	h, repo := setupBookReviewHandler()
	noTotal := ownedReview(1)
	noTotal.TotalPages = 0
	repo.On("FindByID", mock.Anything, uint(1)).Return(noTotal, nil)

	r := newRouter(1)
	r.PUT("/book-reviews/:id/progress", h.UpdateProgress)
	w := doRequest(r, http.MethodPut, "/book-reviews/1/progress", map[string]interface{}{"current_page": 10})

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Update")
}

// --- Delete / Count ---

func TestBookReview_Delete_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReview(1), nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/book-reviews/:id", h.Delete)
	w := doRequest(r, http.MethodDelete, "/book-reviews/1", nil)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestBookReview_Delete_Forbidden(t *testing.T) {
	h, repo := setupBookReviewHandler()
	other := ownedReview(1)
	other.UserID = 999
	repo.On("FindByID", mock.Anything, uint(1)).Return(other, nil)

	r := newRouter(1)
	r.DELETE("/book-reviews/:id", h.Delete)
	w := doRequest(r, http.MethodDelete, "/book-reviews/1", nil)

	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete")
}

func TestBookReview_GetMyCount_Success(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(5), nil)

	r := newRouter(1)
	r.GET("/book-reviews/my/count", h.GetMyCount)
	w := doRequest(r, http.MethodGet, "/book-reviews/my/count", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(5), body["count"])
}

func TestBookReview_GetMyCount_RepoError(t *testing.T) {
	h, repo := setupBookReviewHandler()
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/book-reviews/my/count", h.GetMyCount)
	w := doRequest(r, http.MethodGet, "/book-reviews/my/count", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

// タイトルが長すぎる場合は 400。
func TestBookReview_Create_TitleTooLong(t *testing.T) {
	h, repo := setupBookReviewHandler()

	r := newRouter(1)
	r.POST("/book-reviews", h.Create)
	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"title": strings.Repeat("a", 201), "rating": 3,
	})

	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create")
}
