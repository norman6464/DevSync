package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// makeComment はテスト用のCommentオブジェクトを生成する。
func makeComment(id, userID uint) *model.Comment {
	c := &model.Comment{UserID: userID, PostID: 1, Content: "test"}
	c.ID = id
	return c
}

func newTestCommentLikeService() (*CommentLikeService, *MockCommentLikeRepository, *MockPostRepository) {
	likeRepo := new(MockCommentLikeRepository)
	postRepo := new(MockPostRepository)
	svc := NewCommentLikeService(likeRepo, postRepo)
	return svc, likeRepo, postRepo
}

// --- Like ---

func TestCommentLikeService_Like_Success(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, nil)
	likeRepo.On("Like", uint(10), uint(1)).Return(nil)

	err := svc.Like(10, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_Like_CommentNotFound(t *testing.T) {
	svc, _, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Like(10, 99)
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertExpectations(t)
}

func TestCommentLikeService_Like_AlreadyLiked(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(true, nil)

	err := svc.Like(10, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_Like_SelfLike_Forbidden(t *testing.T) {
	svc, _, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 10), nil)

	err := svc.Like(10, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertExpectations(t)
}

func TestCommentLikeService_Like_HasLikedError(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, errors.New("db error"))

	err := svc.Like(10, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

// --- Unlike ---

func TestCommentLikeService_Unlike_SelfUnlike_Forbidden(t *testing.T) {
	svc, _, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 10), nil)

	err := svc.Unlike(10, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertExpectations(t)
}

func TestCommentLikeService_Unlike_Success(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(true, nil)
	likeRepo.On("Unlike", uint(10), uint(1)).Return(nil)

	err := svc.Unlike(10, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_Unlike_CommentNotFound(t *testing.T) {
	svc, _, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Unlike(10, 99)
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertExpectations(t)
}

func TestCommentLikeService_Unlike_NotLiked(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, nil)

	err := svc.Unlike(10, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_Unlike_HasLikedError(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, errors.New("db error"))

	err := svc.Unlike(10, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

// --- GetStatus ---

func TestCommentLikeService_GetStatus_Success(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(true, nil)
	likeRepo.On("CountByCommentID", uint(1)).Return(int64(5), nil)

	liked, count, err := svc.GetStatus(10, 1)
	assert.NoError(t, err)
	assert.True(t, liked)
	assert.Equal(t, int64(5), count)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_GetStatus_CommentNotFound(t *testing.T) {
	svc, _, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	liked, count, err := svc.GetStatus(10, 99)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.False(t, liked)
	assert.Equal(t, int64(0), count)
	postRepo.AssertExpectations(t)
}

func TestCommentLikeService_GetStatus_HasLikedError(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, errors.New("db error"))

	liked, count, err := svc.GetStatus(10, 1)
	assert.Error(t, err)
	assert.False(t, liked)
	assert.Equal(t, int64(0), count)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}

func TestCommentLikeService_GetStatus_CountError(t *testing.T) {
	svc, likeRepo, postRepo := newTestCommentLikeService()
	postRepo.On("FindCommentByID", uint(1)).Return(makeComment(1, 2), nil)
	likeRepo.On("HasLiked", uint(10), uint(1)).Return(false, nil)
	likeRepo.On("CountByCommentID", uint(1)).Return(int64(0), errors.New("db error"))

	liked, count, err := svc.GetStatus(10, 1)
	assert.Error(t, err)
	assert.False(t, liked)
	assert.Equal(t, int64(0), count)
	postRepo.AssertExpectations(t)
	likeRepo.AssertExpectations(t)
}
