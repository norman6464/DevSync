package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newPostPinTestService() (*PostPinService, *MockPostPinRepository, *MockPostRepository) {
	pinRepo := new(MockPostPinRepository)
	postRepo := new(MockPostRepository)
	svc := NewPostPinService(pinRepo, postRepo)
	return svc, pinRepo, postRepo
}

// --- Pin ---

func TestPostPinService_Pin_Success(t *testing.T) {
	svc, pinRepo, postRepo := newPostPinTestService()
	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	pinRepo.On("CountByUserID", uint(1)).Return(int64(1), nil)
	pinRepo.On("Pin", mock.MatchedBy(func(p *model.PostPin) bool {
		return p.UserID == 1 && p.PostID == 10
	})).Return(nil)

	err := svc.Pin(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
	pinRepo.AssertExpectations(t)
}

func TestPostPinService_Pin_Forbidden(t *testing.T) {
	svc, _, postRepo := newPostPinTestService()
	post := &model.Post{UserID: 2}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.Pin(1, 10)
	assert.Error(t, err)
}

func TestPostPinService_Pin_PostNotFound(t *testing.T) {
	svc, _, postRepo := newPostPinTestService()
	postRepo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	err := svc.Pin(1, 99)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestPostPinService_Pin_MaxReached(t *testing.T) {
	svc, pinRepo, postRepo := newPostPinTestService()
	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	pinRepo.On("CountByUserID", uint(1)).Return(int64(3), nil)

	err := svc.Pin(1, 10)
	assert.Error(t, err)
}

func TestPostPinService_Pin_RepoError(t *testing.T) {
	svc, pinRepo, postRepo := newPostPinTestService()
	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	pinRepo.On("CountByUserID", uint(1)).Return(int64(0), nil)
	pinRepo.On("Pin", mock.Anything).Return(errors.New("db error"))

	err := svc.Pin(1, 10)
	assert.Error(t, err)
}

// --- Unpin ---

func TestPostPinService_Unpin_Success(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("IsPinned", uint(1), uint(10)).Return(true, nil)
	pinRepo.On("Unpin", uint(1), uint(10)).Return(nil)

	err := svc.Unpin(1, 10)
	assert.NoError(t, err)
	pinRepo.AssertExpectations(t)
}

func TestPostPinService_Unpin_NotPinned(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("IsPinned", uint(1), uint(10)).Return(false, nil)

	err := svc.Unpin(1, 10)
	assert.Error(t, err)
}

func TestPostPinService_Unpin_RepoError(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("IsPinned", uint(1), uint(10)).Return(true, nil)
	pinRepo.On("Unpin", uint(1), uint(10)).Return(errors.New("db error"))

	err := svc.Unpin(1, 10)
	assert.Error(t, err)
}

// --- GetByUserID ---

func TestPostPinService_GetByUserID_Success(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	expected := []model.PostPin{
		{UserID: 1, PostID: 10, PinOrder: 0},
		{UserID: 1, PostID: 20, PinOrder: 1},
	}
	pinRepo.On("GetByUserID", uint(1)).Return(expected, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	pinRepo.AssertExpectations(t)
}

func TestPostPinService_GetByUserID_Empty(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("GetByUserID", uint(1)).Return([]model.PostPin{}, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// --- Reorder ---

func TestPostPinService_Reorder_Success(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("UpdateOrder", uint(1), []uint{20, 10}).Return(nil)

	err := svc.Reorder(1, []uint{20, 10})
	assert.NoError(t, err)
	pinRepo.AssertExpectations(t)
}

func TestPostPinService_Reorder_TooMany(t *testing.T) {
	svc, _, _ := newPostPinTestService()

	err := svc.Reorder(1, []uint{1, 2, 3, 4})
	assert.Error(t, err)
}

func TestPostPinService_Reorder_RepoError(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("UpdateOrder", uint(1), []uint{10}).Return(errors.New("db error"))

	err := svc.Reorder(1, []uint{10})
	assert.Error(t, err)
}

// --- IsPinned ---

func TestPostPinService_IsPinned_True(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("IsPinned", uint(1), uint(10)).Return(true, nil)

	result, err := svc.IsPinned(1, 10)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestPostPinService_IsPinned_False(t *testing.T) {
	svc, pinRepo, _ := newPostPinTestService()
	pinRepo.On("IsPinned", uint(1), uint(10)).Return(false, nil)

	result, err := svc.IsPinned(1, 10)
	assert.NoError(t, err)
	assert.False(t, result)
}
