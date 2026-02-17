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

// ============================================================
// 投稿取得テスト
// ============================================================

func TestPostGetByID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	expected := &model.Post{Title: "Test", UserID: 1}
	expected.ID = 1
	postRepo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Title)
}

func TestPostGetByID_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostGetAll_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	posts := []model.Post{{Title: "Post 1"}, {Title: "Post 2"}}
	postRepo.On("FindAll", 1, 20).Return(posts, nil)

	result, err := svc.GetAll(1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPostGetByUserID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	posts := []model.Post{{Title: "My Post", UserID: 1}}
	postRepo.On("FindByUserID", uint(1)).Return(posts, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPostGetDrafts_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	drafts := []model.Post{{Title: "Draft", IsDraft: true, UserID: 1}}
	postRepo.On("FindDraftsByUserID", uint(1)).Return(drafts, nil)

	result, err := svc.GetDrafts(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].IsDraft)
}

func TestPostTimeline_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	posts := []model.Post{{Title: "Timeline Post"}}
	postRepo.On("Timeline", uint(1), 1, 10).Return(posts, nil)

	result, err := svc.Timeline(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// ============================================================
// いいねテスト
// ============================================================

func TestPostLike_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Like", uint(1), uint(10)).Return(nil)

	err := svc.Like(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnlike_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Unlike", uint(1), uint(10)).Return(nil)

	err := svc.Unlike(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostHasLiked(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("HasLiked", uint(1), uint(10)).Return(true)

	result := svc.HasLiked(1, 10)
	assert.True(t, result)
}

// ============================================================
// コメントテスト
// ============================================================

func TestPostCreateComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{Content: "Nice post!", UserID: 1, PostID: 10}
	postRepo.On("CreateComment", comment).Return(nil)

	err := svc.CreateComment(comment)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostCreateComment_ValidationError(t *testing.T) {
	svc, _, _ := newTestPostService()

	// 空コメントはバリデーションエラー
	comment := &model.Comment{Content: "", UserID: 1, PostID: 10}

	err := svc.CreateComment(comment)
	assert.Error(t, err)
}

func TestPostGetComments_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comments := []model.Comment{{Content: "Comment 1"}, {Content: "Comment 2"}}
	postRepo.On("GetComments", uint(10)).Return(comments, nil)

	result, err := svc.GetComments(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPostDeleteComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("DeleteComment", uint(5), uint(1)).Return(nil)

	err := svc.DeleteComment(5, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

// ============================================================
// 投稿公開テスト
// ============================================================

func TestPostPublish_Success(t *testing.T) {
	svc, postRepo, notifRepo := newTestPostService()

	draft := &model.Post{Title: "Draft Post", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)
	postRepo.On("Update", draft).Return(nil)
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{2, 3}, nil).Maybe()
	notifRepo.On("CreateBatch", mock.Anything).Return(nil).Maybe()

	result, err := svc.Publish(1, 1)
	assert.NoError(t, err)
	assert.False(t, result.IsDraft)
}

func TestPostPublish_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "Draft", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)

	result, err := svc.Publish(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestPostPublish_AlreadyPublished(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "Published", UserID: 1, IsDraft: false}
	published.ID = 1

	postRepo.On("FindByID", uint(1)).Return(published, nil)

	result, err := svc.Publish(1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
}

// ============================================================
// 投稿作成 追加テスト
// ============================================================

func TestPostCreate_Draft(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{Title: "Draft Post", Content: "Content", UserID: 1, IsDraft: true}

	postRepo.On("Create", post).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 20
	}).Return(nil)

	createdPost := &model.Post{Title: "Draft Post", Content: "Content", UserID: 1, IsDraft: true}
	createdPost.ID = 20
	postRepo.On("FindByID", uint(20)).Return(createdPost, nil)

	result, err := svc.Create(post)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// 下書きの場合はNotifyFollowersは呼ばれない
}

func TestPostCreate_ValidationError(t *testing.T) {
	svc, _, _ := newTestPostService()

	// タイトルが空の場合バリデーションエラー
	post := &model.Post{Title: "", Content: "Content", UserID: 1}

	result, err := svc.Create(post)
	assert.Error(t, err)
	assert.Nil(t, result)
}
