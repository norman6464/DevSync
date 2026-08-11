package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockPostTagRepo は usecase/repository.PostTagRepository のモック（ctx 付き）。
type mockPostTagRepo struct{ mock.Mock }

func (m *mockPostTagRepo) SetTags(ctx context.Context, postID uint, tags []string) error {
	return m.Called(ctx, postID, tags).Error(0)
}

func (m *mockPostTagRepo) GetByPostID(ctx context.Context, postID uint) ([]string, error) {
	args := m.Called(ctx, postID)
	t, _ := args.Get(0).([]string)
	return t, args.Error(1)
}

func (m *mockPostTagRepo) FindPostsByTag(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, tag, limit, offset)
	p, _ := args.Get(0).([]model.Post)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostTagRepo) GetPopularTags(ctx context.Context, limit int) ([]model.TagCount, error) {
	args := m.Called(ctx, limit)
	t, _ := args.Get(0).([]model.TagCount)
	return t, args.Error(1)
}

// PostReader のモックは post_pin_test.go の mockPostReader を再利用する。

// setupPostTagHandler は本物の usecase と port モックで PostTagHandler を組む。
func setupPostTagHandler() (*PostTagHandler, *mockPostTagRepo, *mockPostReader) {
	tags := new(mockPostTagRepo)
	posts := new(mockPostReader)
	h := NewPostTagHandler(
		usecase.NewSetPostTagsUseCase(tags, posts),
		usecase.NewGetPostTagsUseCase(tags),
		usecase.NewFindPostsByTagUseCase(tags),
		usecase.NewGetPopularTagsUseCase(tags),
	)
	return h, tags, posts
}

// ownedPost は所有者が userID=1 の投稿を返す。
func ownedPost(id uint) *model.Post {
	p := &model.Post{}
	p.ID = id
	p.UserID = 1
	return p
}

func TestPostTag_SetTags_Success(t *testing.T) {
	h, tags, posts := setupPostTagHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPost(5), nil)
	// 正規化されて渡る（小文字化・重複除外）
	tags.On("SetTags", mock.Anything, uint(5), []string{"go", "web"}).Return(nil)

	r := newRouter(1)
	r.PUT("/tags/posts/:postId", h.SetTags)
	w := doRequest(r, http.MethodPut, "/tags/posts/5", map[string]interface{}{
		"tags": []string{"Go", " web ", "GO"},
	})

	assertStatus(t, w, http.StatusOK)
	tags.AssertExpectations(t)
}

// 所有者以外のタグ設定は 403 を返し、保存しない。
func TestPostTag_SetTags_Forbidden(t *testing.T) {
	h, tags, posts := setupPostTagHandler()
	other := &model.Post{}
	other.ID = 5
	other.UserID = 999
	posts.On("FindByID", mock.Anything, uint(5)).Return(other, nil)

	r := newRouter(1)
	r.PUT("/tags/posts/:postId", h.SetTags)
	w := doRequest(r, http.MethodPut, "/tags/posts/5", map[string]interface{}{"tags": []string{"go"}})

	assertStatus(t, w, http.StatusForbidden)
	tags.AssertNotCalled(t, "SetTags")
}

// 正規化後に 11 個以上あれば 400 を返し、保存しない。
func TestPostTag_SetTags_TooMany(t *testing.T) {
	h, tags, posts := setupPostTagHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPost(5), nil)

	many := make([]string, 0, 11)
	for _, s := range []string{"a1", "b2", "c3", "d4", "e5", "f6", "g7", "h8", "i9", "j10", "k11"} {
		many = append(many, s)
	}

	r := newRouter(1)
	r.PUT("/tags/posts/:postId", h.SetTags)
	w := doRequest(r, http.MethodPut, "/tags/posts/5", map[string]interface{}{"tags": many})

	assertStatus(t, w, http.StatusBadRequest)
	tags.AssertNotCalled(t, "SetTags")
}

// 51 文字以上のタグは 400 を返し、保存しない。
func TestPostTag_SetTags_TagTooLong(t *testing.T) {
	h, tags, posts := setupPostTagHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPost(5), nil)

	r := newRouter(1)
	r.PUT("/tags/posts/:postId", h.SetTags)
	w := doRequest(r, http.MethodPut, "/tags/posts/5", map[string]interface{}{
		"tags": []string{strings.Repeat("a", 51)},
	})

	assertStatus(t, w, http.StatusBadRequest)
	tags.AssertNotCalled(t, "SetTags")
}

// 投稿が存在しない場合は 500（移行前の挙動を維持している）。
func TestPostTag_SetTags_PostNotFound(t *testing.T) {
	h, tags, posts := setupPostTagHandler()
	posts.On("FindByID", mock.Anything, uint(5)).Return(nil, errors.New("record not found"))

	r := newRouter(1)
	r.PUT("/tags/posts/:postId", h.SetTags)
	w := doRequest(r, http.MethodPut, "/tags/posts/5", map[string]interface{}{"tags": []string{"go"}})

	assertStatus(t, w, http.StatusInternalServerError)
	tags.AssertNotCalled(t, "SetTags")
}

func TestPostTag_GetByPostID_Success(t *testing.T) {
	h, tags, _ := setupPostTagHandler()
	tags.On("GetByPostID", mock.Anything, uint(5)).Return([]string{"go", "web"}, nil)

	r := newRouter(1)
	r.GET("/tags/posts/:postId", h.GetByPostID)
	w := doRequest(r, http.MethodGet, "/tags/posts/5", nil)

	assertStatus(t, w, http.StatusOK)
	tags.AssertExpectations(t)
}

func TestPostTag_GetByPostID_RepoError(t *testing.T) {
	h, tags, _ := setupPostTagHandler()
	tags.On("GetByPostID", mock.Anything, uint(5)).Return([]string(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/tags/posts/:postId", h.GetByPostID)
	w := doRequest(r, http.MethodGet, "/tags/posts/5", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestPostTag_FindPostsByTag_Success(t *testing.T) {
	h, tags, _ := setupPostTagHandler()
	tags.On("FindPostsByTag", mock.Anything, "go", 20, 0).
		Return([]model.Post{{Title: "t"}}, int64(1), nil)

	r := newRouter(1)
	r.GET("/tags/search", h.FindPostsByTag)
	w := doRequest(r, http.MethodGet, "/tags/search?tag=go", nil)

	assertStatus(t, w, http.StatusOK)
	tags.AssertExpectations(t)
}

func TestPostTag_FindPostsByTag_MissingParam(t *testing.T) {
	h, _, _ := setupPostTagHandler()

	r := newRouter(1)
	r.GET("/tags/search", h.FindPostsByTag)
	w := doRequest(r, http.MethodGet, "/tags/search", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTag_FindPostsByTag_TooLong(t *testing.T) {
	h, _, _ := setupPostTagHandler()

	r := newRouter(1)
	r.GET("/tags/search", h.FindPostsByTag)
	w := doRequest(r, http.MethodGet, "/tags/search?tag="+strings.Repeat("a", 101), nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTag_GetPopularTags_Success(t *testing.T) {
	h, tags, _ := setupPostTagHandler()
	tags.On("GetPopularTags", mock.Anything, 20).
		Return([]model.TagCount{{Tag: "go", Count: 3}}, nil)

	r := newRouter(1)
	r.GET("/tags/popular", h.GetPopularTags)
	w := doRequest(r, http.MethodGet, "/tags/popular", nil)

	assertStatus(t, w, http.StatusOK)
	tags.AssertExpectations(t)
}

func TestPostTag_GetPopularTags_RepoError(t *testing.T) {
	h, tags, _ := setupPostTagHandler()
	tags.On("GetPopularTags", mock.Anything, 20).
		Return([]model.TagCount(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/tags/popular", h.GetPopularTags)
	w := doRequest(r, http.MethodGet, "/tags/popular", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}
