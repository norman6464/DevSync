package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestPostService はPostServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostService() (*PostService, *MockPostRepository, *MockNotificationRepository) {
	postRepo := new(MockPostRepository)
	notifRepo := new(MockNotificationRepository)
	notifService := NewNotificationService(notifRepo)
	svc := NewPostService(postRepo, notifService)
	return svc, postRepo, notifRepo
}

// ============================================================
// 投稿作成テスト
// ============================================================

func TestPostCreate_Success(t *testing.T) {
	svc, postRepo, notifRepo := newTestPostService()

	post := &model.Post{Title: "Test Post", Content: "Content", UserID: 1}

	postRepo.On("Create", post).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 10
	}).Return(nil)

	createdPost := &model.Post{Title: "Test Post", Content: "Content", UserID: 1}
	createdPost.ID = 10
	postRepo.On("FindByID", uint(10)).Return(createdPost, nil)

	// NotifyFollowers はgoroutineで呼ばれるため、モックを設定
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil).Maybe()
	notifRepo.On("CreateBatch", mock.Anything).Return(nil).Maybe()

	result, err := svc.Create(post)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	postRepo.AssertCalled(t, "Create", post)
}

// ============================================================
// 投稿更新テスト
// ============================================================

func TestPostUpdate_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "New Title", "New Content", "")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Content", result.Content)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Title", Content: "Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 999, "New Title", "", "")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_PartialUpdate(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	// タイトルのみ更新（contentは空文字）
	result, err := svc.Update(1, 1, "New Title", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "Old Content", result.Content) // 変更されない
	postRepo.AssertExpectations(t)
}

// ============================================================
// 投稿削除テスト
// ============================================================

func TestPostDelete_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostDelete_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertExpectations(t)
}

func TestPostDelete_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Delete(999, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}
