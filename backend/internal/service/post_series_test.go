package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestPostSeriesService はPostSeriesServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostSeriesService() (*PostSeriesService, *MockPostSeriesRepository) {
	repo := new(MockPostSeriesRepository)
	svc := NewPostSeriesService(repo)
	return svc, repo
}

// ============================================================
// シリーズ作成テスト
// ============================================================

func TestPostSeriesCreate_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	series := &model.PostSeries{
		Title:       "Go入門シリーズ",
		Description: "Go言語を基礎から学ぶ連載",
		UserID:      1,
	}

	repo.On("Create", series).Return(nil)

	err := svc.Create(series)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesCreate_ValidationError(t *testing.T) {
	svc, _ := newTestPostSeriesService()

	series := &model.PostSeries{
		Title:  "",
		UserID: 1,
	}

	err := svc.Create(series)
	assert.Error(t, err)
}

func TestPostSeriesCreate_RepoError(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	series := &model.PostSeries{
		Title:  "Go入門シリーズ",
		UserID: 1,
	}

	repo.On("Create", series).Return(errors.New("db error"))

	err := svc.Create(series)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// シリーズ取得テスト
// ============================================================

func TestPostSeriesGetByID_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	expected := &model.PostSeries{Title: "Go入門", UserID: 1}
	expected.ID = 1

	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Go入門", result.Title)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByID_NotFound(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	expected := []model.PostSeries{
		{Title: "Go入門", UserID: 1},
		{Title: "React入門", UserID: 1},
	}

	repo.On("FindByUserID", uint(1), 0, 10).Return(expected, nil)

	result, err := svc.GetByUserID(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByUserID", uint(1), 0, 10).Return([]model.PostSeries{}, nil)

	result, err := svc.GetByUserID(1, 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	expected := []model.PostSeries{
		{Title: "シリーズ11", UserID: 1},
	}

	repo.On("FindByUserID", uint(1), 10, 10).Return(expected, nil)

	result, err := svc.GetByUserID(1, 2, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestPostSeriesCountByUser_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("CountByUser", uint(1)).Return(int64(5), nil)

	count, err := svc.CountByUser(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	repo.AssertExpectations(t)
}

func TestPostSeriesCountByUser_Zero(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("CountByUser", uint(1)).Return(int64(0), nil)

	count, err := svc.CountByUser(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
	repo.AssertExpectations(t)
}

func TestPostSeriesCountByUser_RepoError(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("CountByUser", uint(1)).Return(int64(0), errors.New("db error"))

	count, err := svc.CountByUser(1)
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

// ============================================================
// シリーズ更新テスト
// ============================================================

func TestPostSeriesUpdate_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{Title: "Old Title", Description: "Old Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.PostSeries{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "Old Desc", result.Description)
	repo.AssertExpectations(t)
}

func TestPostSeriesUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.PostSeries{Title: "New Title"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestPostSeriesUpdate_NotFound(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.PostSeries{Title: "New Title"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestPostSeriesUpdate_RepoError(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{Title: "Old", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))

	updates := &model.PostSeries{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// シリーズ削除テスト
// ============================================================

func TestPostSeriesDelete_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesDelete_Forbidden(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestPostSeriesDelete_NotFound(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Delete(999, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// シリーズへの投稿追加テスト
// ============================================================

func TestPostSeriesAddPost_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("HasPost", uint(1), uint(10)).Return(false, nil)
	repo.On("AddPost", &model.PostSeriesItem{SeriesID: 1, PostID: 10, OrderIndex: 0}).Return(nil)

	err := svc.AddPost(1, 10, 0, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesAddPost_Forbidden(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.AddPost(1, 10, 0, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestPostSeriesAddPost_NotFound(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.AddPost(999, 10, 0, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesAddPost_Duplicate(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("HasPost", uint(1), uint(10)).Return(true, nil)

	err := svc.AddPost(1, 10, 0, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "すでに追加済み")
	repo.AssertNotCalled(t, "AddPost")
	repo.AssertExpectations(t)
}

func TestPostSeriesAddPost_HasPostError(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("HasPost", uint(1), uint(10)).Return(false, errors.New("db error"))

	err := svc.AddPost(1, 10, 0, 1)
	assert.Error(t, err)
	repo.AssertNotCalled(t, "AddPost")
	repo.AssertExpectations(t)
}

// ============================================================
// シリーズからの投稿削除テスト
// ============================================================

func TestPostSeriesRemovePost_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("RemovePost", uint(1), uint(10)).Return(nil)

	err := svc.RemovePost(1, 10, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesRemovePost_Forbidden(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.RemovePost(1, 10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// シリーズの投稿一覧取得テスト
// ============================================================

func TestPostSeriesGetPosts_Success(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	expected := []model.PostSeriesItem{
		{SeriesID: 1, PostID: 1, OrderIndex: 0},
		{SeriesID: 1, PostID: 2, OrderIndex: 1},
	}

	repo.On("GetPostsBySeriesID", uint(1)).Return(expected, nil)

	result, err := svc.GetPosts(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 0, result[0].OrderIndex)
	assert.Equal(t, 1, result[1].OrderIndex)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetPosts_Empty(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("GetPostsBySeriesID", uint(1)).Return([]model.PostSeriesItem{}, nil)

	result, err := svc.GetPosts(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestPostSeriesRemovePost_NotFound(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.RemovePost(99, 5, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostSeriesCreate_WhitespaceTitle(t *testing.T) {
	svc, _ := newTestPostSeriesService()

	series := &model.PostSeries{
		Title:  "   ",
		UserID: 1,
	}

	err := svc.Create(series)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトルは必須です")
}

func TestPostSeriesUpdate_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{Title: "Old Title", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.PostSeries{Title: "  \t  "}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "空白のみ")
}

func TestPostSeriesUpdate_WhitespaceDescription(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{Title: "Title", Description: "Old Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.PostSeries{Description: "  \t  "}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "空白のみ")
}

func TestPostSeriesUpdate_Description(t *testing.T) {
	svc, repo := newTestPostSeriesService()

	existing := &model.PostSeries{Title: "Title", Description: "Old Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.PostSeries{Description: "New Desc"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "Title", result.Title)
	assert.Equal(t, "New Desc", result.Description)
	repo.AssertExpectations(t)
}
