package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// assertBadRequestError はエラーがdomain.ErrCodeBadRequestのDomainErrorであることを検証する。
func assertBadRequestError(t *testing.T, err error) {
	t.Helper()
	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr), "エラーはDomainError型であるべき")
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

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

// ============================================================
// バリデーションテスト
// ============================================================

func TestPostViewRecordView_InvalidUserID(t *testing.T) {
	svc, repo := newTestPostViewService()

	err := svc.RecordView(0, 10)
	assertBadRequestError(t, err)
	repo.AssertNotCalled(t, "HasViewed")
}

func TestPostViewRecordView_InvalidPostID(t *testing.T) {
	svc, repo := newTestPostViewService()

	err := svc.RecordView(1, 0)
	assertBadRequestError(t, err)
	repo.AssertNotCalled(t, "HasViewed")
}

func TestPostViewGetViewCount_InvalidPostID(t *testing.T) {
	svc, repo := newTestPostViewService()

	count, err := svc.GetViewCount(0)
	assertBadRequestError(t, err)
	assert.Equal(t, int64(0), count)
	repo.AssertNotCalled(t, "GetViewCount")
}

func TestPostViewHasViewed_InvalidUserID(t *testing.T) {
	svc, repo := newTestPostViewService()

	viewed, err := svc.HasViewed(0, 10)
	assertBadRequestError(t, err)
	assert.False(t, viewed)
	repo.AssertNotCalled(t, "HasViewed")
}

func TestPostViewHasViewed_InvalidPostID(t *testing.T) {
	svc, repo := newTestPostViewService()

	viewed, err := svc.HasViewed(1, 0)
	assertBadRequestError(t, err)
	assert.False(t, viewed)
	repo.AssertNotCalled(t, "HasViewed")
}

func TestPostViewGetMostViewed_InvalidLimit_Zero(t *testing.T) {
	svc, repo := newTestPostViewService()

	result, err := svc.GetMostViewed(0)
	assertBadRequestError(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "GetMostViewed")
}

func TestPostViewGetMostViewed_InvalidLimit_Negative(t *testing.T) {
	svc, repo := newTestPostViewService()

	result, err := svc.GetMostViewed(-5)
	assertBadRequestError(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "GetMostViewed")
}

func TestPostViewGetMostViewed_InvalidLimit_TooLarge(t *testing.T) {
	svc, repo := newTestPostViewService()

	result, err := svc.GetMostViewed(200)
	assertBadRequestError(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "GetMostViewed")
}

func TestPostViewHasViewed_RepoError(t *testing.T) {
	svc, repo := newTestPostViewService()

	repo.On("HasViewed", uint(1), uint(10)).Return(false, errors.New("db error"))

	viewed, err := svc.HasViewed(1, 10)
	assert.Error(t, err)
	assert.False(t, viewed)
	repo.AssertExpectations(t)
}
