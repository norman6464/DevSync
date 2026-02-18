package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestPostViewService はPostViewServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostViewService() (*PostViewService, *MockPostViewRepository) {
	repo := new(MockPostViewRepository)
	svc := NewPostViewService(repo)
	return svc, repo
}

// ============================================================
// 閲覧記録テスト
// ============================================================

func TestPostViewRecordView_Success(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(false, nil)
	repo.On("RecordView", mock.MatchedBy(func(v *model.PostView) bool {
		return v.UserID == 1 && v.PostID == 10
	})).Return(nil)

	err := svc.RecordView(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostViewRecordView_AlreadyViewed(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(true, nil)

	err := svc.RecordView(1, 10)
	assert.NoError(t, err)
	// 既に閲覧済みの場合はRecordViewが呼ばれない
	repo.AssertNotCalled(t, "RecordView")
}

func TestPostViewRecordView_HasViewedError(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(false, errors.New("db error"))

	err := svc.RecordView(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostViewRecordView_RepoError(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(false, nil)
	repo.On("RecordView", mock.Anything).Return(errors.New("db error"))

	err := svc.RecordView(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// 閲覧数取得テスト
// ============================================================

func TestPostViewGetViewCount_Success(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("GetViewCount", uint(10)).Return(int64(42), nil)

	count, err := svc.GetViewCount(10)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)
	repo.AssertExpectations(t)
}

func TestPostViewGetViewCount_RepoError(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("GetViewCount", uint(10)).Return(int64(0), errors.New("db error"))

	count, err := svc.GetViewCount(10)
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	repo.AssertExpectations(t)
}

// ============================================================
// 閲覧済み確認テスト
// ============================================================

func TestPostViewHasViewed_True(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(true, nil)

	viewed, err := svc.HasViewed(1, 10)
	assert.NoError(t, err)
	assert.True(t, viewed)
	repo.AssertExpectations(t)
}

func TestPostViewHasViewed_False(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(false, nil)

	viewed, err := svc.HasViewed(1, 10)
	assert.NoError(t, err)
	assert.False(t, viewed)
	repo.AssertExpectations(t)
}

// ============================================================
// 人気投稿取得テスト
// ============================================================

func TestPostViewGetMostViewed_Success(t *testing.T) {
	svc, repo := newTestPostViewService()

	viewCounts := []model.ViewCount{
		{PostID: 1, Count: 100},
		{PostID: 2, Count: 50},
	}
	repo.On("GetMostViewed", 10).Return(viewCounts, nil)

	result, err := svc.GetMostViewed(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 100, result[0].Count)
	repo.AssertExpectations(t)
}

func TestPostViewGetMostViewed_Empty(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("GetMostViewed", 10).Return([]model.ViewCount{}, nil)

	result, err := svc.GetMostViewed(10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestPostViewGetMostViewed_RepoError(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("GetMostViewed", 10).Return([]model.ViewCount(nil), errors.New("db error"))

	result, err := svc.GetMostViewed(10)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
