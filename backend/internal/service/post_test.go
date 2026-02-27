package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
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
	postRepo.On("FindByUserID", uint(1), 20, 0).Return(posts, int64(1), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
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

func TestPostGetDrafts_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindDraftsByUserID", uint(1)).Return([]model.Post(nil), errors.New("db error"))

	result, err := svc.GetDrafts(1)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostTimeline_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	posts := []model.Post{{Title: "Timeline Post"}}
	postRepo.On("Timeline", uint(1), 1, 10).Return(posts, nil)

	result, err := svc.Timeline(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPostTimeline_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Timeline", uint(1), 1, 10).Return([]model.Post(nil), errors.New("db error"))

	result, err := svc.Timeline(1, 1, 10)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

// ============================================================
// いいねテスト
// ============================================================

func TestPostLike_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 2, Title: "Other's post"}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Like", uint(1), uint(10)).Return(nil)

	err := svc.Like(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostLike_SelfLike_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 1, Title: "My post"}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.Like(1, 10)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "Like")
}

func TestPostLike_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Like(1, 999)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestPostUnlike_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 2, Title: "Other's post"}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Unlike", uint(1), uint(10)).Return(nil)

	err := svc.Unlike(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnlike_SelfUnlike_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 1, Title: "My post"}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.Unlike(1, 10)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "Unlike")
}

func TestPostUnlike_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Unlike(1, 999)
	assert.ErrorIs(t, err, ErrNotFound)
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

func TestPostGetComments_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetComments", uint(10)).Return([]model.Comment(nil), errors.New("db error"))

	result, err := svc.GetComments(10)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostDeleteComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)
	postRepo.On("DeleteComment", uint(5)).Return(nil)

	err := svc.DeleteComment(5, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostDeleteComment_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	err := svc.DeleteComment(5, 999)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
	postRepo.AssertExpectations(t)
}

func TestPostDeleteComment_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.DeleteComment(99, 1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

// ============================================================
// コメント編集テスト
// ============================================================

func TestPostEditComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "old content"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)
	postRepo.On("UpdateComment", mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 5 && c.Content == "new content"
	})).Return(nil)

	result, err := svc.EditComment(5, 1, "new content")
	assert.NoError(t, err)
	assert.Equal(t, "new content", result.Content)
	postRepo.AssertExpectations(t)
}

func TestPostEditComment_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "old"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	result, err := svc.EditComment(5, 999, "new content")
	assert.Nil(t, result)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
}

func TestPostEditComment_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.EditComment(99, 1, "new content")
	assert.Nil(t, result)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostEditComment_EmptyContent(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "old"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	result, err := svc.EditComment(5, 1, "")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestPostEditComment_TooLongContent(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "old"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	longContent := strings.Repeat("a", 5001)
	result, err := svc.EditComment(5, 1, longContent)
	assert.Nil(t, result)
	assert.Error(t, err)
}

// ============================================================
// スレッドコメント（返信）テスト
// ============================================================

func TestPostCreateReply_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	parentID := uint(5)
	parentComment := &model.Comment{PostID: 10, ParentID: nil}
	parentComment.ID = 5
	postRepo.On("FindCommentByID", uint(5)).Return(parentComment, nil)

	reply := &model.Comment{Content: "Great reply!", UserID: 2, PostID: 10, ParentID: &parentID}
	postRepo.On("CreateComment", reply).Return(nil)

	err := svc.CreateComment(reply)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostCreateReply_ParentNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	parentID := uint(999)
	reply := &model.Comment{Content: "Reply", UserID: 2, PostID: 10, ParentID: &parentID}
	postRepo.On("FindCommentByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.CreateComment(reply)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "親コメントが見つかりません")
}

func TestPostCreateReply_ParentOnDifferentPost(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	parentID := uint(5)
	parentComment := &model.Comment{PostID: 20, ParentID: nil}
	parentComment.ID = 5
	postRepo.On("FindCommentByID", uint(5)).Return(parentComment, nil)

	reply := &model.Comment{Content: "Reply", UserID: 2, PostID: 10, ParentID: &parentID}

	err := svc.CreateComment(reply)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "別の投稿に属しています")
}

func TestPostCreateReply_NestedReplyNotAllowed(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	grandParentID := uint(3)
	parentID := uint(5)
	parentComment := &model.Comment{PostID: 10, ParentID: &grandParentID}
	parentComment.ID = 5
	postRepo.On("FindCommentByID", uint(5)).Return(parentComment, nil)

	reply := &model.Comment{Content: "Nested reply", UserID: 2, PostID: 10, ParentID: &parentID}

	err := svc.CreateComment(reply)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "返信への返信はできません")
}

func TestPostGetReplies_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	parentID := uint(5)
	replies := []model.Comment{
		{ID: 10, Content: "Reply 1", ParentID: &parentID},
		{ID: 11, Content: "Reply 2", ParentID: &parentID},
	}
	postRepo.On("GetReplies", uint(5)).Return(replies, nil)

	result, err := svc.GetReplies(5)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	postRepo.AssertExpectations(t)
}

func TestPostGetReplies_Empty(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReplies", uint(5)).Return([]model.Comment{}, nil)

	result, err := svc.GetReplies(5)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestPostGetReplies_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReplies", uint(5)).Return([]model.Comment(nil), errors.New("db error"))

	result, err := svc.GetReplies(5)
	assert.Error(t, err)
	assert.Nil(t, result)
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

func TestPostPublish_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.Publish(99, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostPublish_UpdateError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "Draft", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)
	postRepo.On("Update", draft).Return(errors.New("db error"))

	result, err := svc.Publish(1, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

// ============================================================
// 投稿非公開化（Unpublish）テスト
// ============================================================

func TestPostUnpublish_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "Published Post", UserID: 1, IsDraft: false}
	published.ID = 1

	postRepo.On("FindByID", uint(1)).Return(published, nil)
	postRepo.On("Update", published).Return(nil)

	result, err := svc.Unpublish(1, 1)
	assert.NoError(t, err)
	assert.True(t, result.IsDraft)
}

func TestPostUnpublish_AlreadyDraft(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "Draft", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)

	result, err := svc.Unpublish(1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
}

func TestPostUnpublish_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "Published", UserID: 1, IsDraft: false}
	published.ID = 1

	postRepo.On("FindByID", uint(1)).Return(published, nil)

	result, err := svc.Unpublish(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestPostUnpublish_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.Unpublish(99, 1)
	assert.Error(t, err)
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

// ============================================================
// ブックマークテスト
// ============================================================

func TestPostBookmark_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("Bookmark", uint(1), uint(10)).Return(nil)

	err := svc.Bookmark(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostBookmark_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("Bookmark", uint(1), uint(10)).Return(errors.New("already bookmarked"))

	err := svc.Bookmark(1, 10)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostBookmark_SelfBookmark_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 1}, nil)

	err := svc.Bookmark(1, 10)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "Bookmark")
}

func TestPostBookmark_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Bookmark(1, 99)
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertNotCalled(t, "Bookmark")
}

func TestPostUnbookmark_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("Unbookmark", uint(1), uint(10)).Return(nil)

	err := svc.Unbookmark(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnbookmark_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("Unbookmark", uint(1), uint(10)).Return(errors.New("not bookmarked"))

	err := svc.Unbookmark(1, 10)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnbookmark_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(nil, errors.New("not found"))

	err := svc.Unbookmark(1, 10)
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertNotCalled(t, "Unbookmark")
}

func TestPostUnbookmark_SelfUnbookmark_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 1}, nil)

	err := svc.Unbookmark(1, 10)
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "Unbookmark")
}

func TestPostHasBookmarked_True(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("HasBookmarked", uint(1), uint(10)).Return(true)

	result := svc.HasBookmarked(1, 10)
	assert.True(t, result)
}

func TestPostHasBookmarked_False(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("HasBookmarked", uint(1), uint(10)).Return(false)

	result := svc.HasBookmarked(1, 10)
	assert.False(t, result)
}

func TestPostGetBookmarks_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	posts := []model.Post{{Title: "Bookmarked Post 1"}, {Title: "Bookmarked Post 2"}}
	postRepo.On("FindBookmarkedByUserID", uint(1), 1, 20).Return(posts, int64(2), nil)

	result, total, err := svc.GetBookmarks(1, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestPostGetBookmarks_Empty(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindBookmarkedByUserID", uint(1), 1, 20).Return([]model.Post{}, int64(0), nil)

	result, total, err := svc.GetBookmarks(1, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	assert.Equal(t, int64(0), total)
}

// ============================================================
// リアクションテスト
// ============================================================

func TestPostAddReaction_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("AddReaction", uint(1), uint(10), "👍").Return(nil)

	err := svc.AddReaction(1, 10, "👍")
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostAddReaction_InvalidEmoji(t *testing.T) {
	svc, _, _ := newTestPostService()

	err := svc.AddReaction(1, 10, "invalid")
	assert.Error(t, err)
}

func TestPostAddReaction_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 2}, nil)
	postRepo.On("AddReaction", uint(1), uint(10), "🎉").Return(errors.New("duplicate"))

	err := svc.AddReaction(1, 10, "🎉")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostAddReaction_SelfReaction_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(10)).Return(&model.Post{UserID: 1}, nil)

	err := svc.AddReaction(1, 10, "👍")
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "AddReaction")
}

func TestPostAddReaction_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.AddReaction(1, 99, "👍")
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertNotCalled(t, "AddReaction")
}

func TestPostRemoveReaction_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 99}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("RemoveReaction", uint(1), uint(10), "👍").Return(nil)

	err := svc.RemoveReaction(1, 10, "👍")
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostRemoveReaction_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 99}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("RemoveReaction", uint(1), uint(10), "👍").Return(errors.New("not found"))

	err := svc.RemoveReaction(1, 10, "👍")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostRemoveReaction_SelfReaction_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.RemoveReaction(1, 10, "👍")
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertNotCalled(t, "RemoveReaction")
}

func TestPostRemoveReaction_PostNotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.RemoveReaction(1, 99, "👍")
	assert.ErrorIs(t, err, ErrNotFound)
	postRepo.AssertNotCalled(t, "RemoveReaction")
}

func TestPostRemoveReaction_InvalidEmoji_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	err := svc.RemoveReaction(1, 10, "malicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "許可されていない絵文字です")
	postRepo.AssertNotCalled(t, "FindByID")
	postRepo.AssertNotCalled(t, "RemoveReaction")
}

func TestPostGetReactionsByPostID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	counts := []model.ReactionCount{
		{Emoji: "👍", Count: 5},
		{Emoji: "❤️", Count: 3},
	}
	postRepo.On("GetReactionsByPostID", uint(10)).Return(counts, nil)

	result, err := svc.GetReactionsByPostID(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "👍", result[0].Emoji)
	assert.Equal(t, 5, result[0].Count)
	postRepo.AssertExpectations(t)
}

func TestPostGetReactionsByPostID_Empty(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReactionsByPostID", uint(10)).Return([]model.ReactionCount{}, nil)

	result, err := svc.GetReactionsByPostID(10)
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	postRepo.AssertExpectations(t)
}

func TestPostGetUserReactions_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetUserReactions", uint(1), uint(10)).Return([]string{"👍", "🎉"}, nil)

	result, err := svc.GetUserReactions(1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "👍")
	assert.Contains(t, result, "🎉")
	postRepo.AssertExpectations(t)
}

// ============================================================
// 読了時間推定テスト
// ============================================================

func TestEstimateReadTime_ShortContent(t *testing.T) {
	// 短いコンテンツ → 最低1分
	result := EstimateReadTime("Hello World")
	assert.Equal(t, 1, result)
}

func TestEstimateReadTime_JapaneseContent(t *testing.T) {
	// 500文字の日本語 → 1分
	content := ""
	for i := 0; i < 500; i++ {
		content += "あ"
	}
	result := EstimateReadTime(content)
	assert.Equal(t, 1, result)
}

func TestEstimateReadTime_LongJapaneseContent(t *testing.T) {
	// 1500文字の日本語 → 3分
	content := ""
	for i := 0; i < 1500; i++ {
		content += "あ"
	}
	result := EstimateReadTime(content)
	assert.Equal(t, 3, result)
}

func TestEstimateReadTime_EmptyContent(t *testing.T) {
	result := EstimateReadTime("")
	assert.Equal(t, 1, result)
}

func TestEstimateReadTime_MixedContent(t *testing.T) {
	// 1000文字（混合） → 2分
	content := ""
	for i := 0; i < 1000; i++ {
		content += "a"
	}
	result := EstimateReadTime(content)
	assert.Equal(t, 2, result)
}

func TestPostCreate_FindByIDErrorAfterCreate(t *testing.T) {
	svc, postRepo, notifRepo := newTestPostService()

	post := &model.Post{Title: "Test", Content: "Content", UserID: 1}

	postRepo.On("Create", post).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 10
	}).Return(nil)
	// Createは成功するがFindByIDが失敗 → return post, nil
	postRepo.On("FindByID", uint(10)).Return(nil, errors.New("not found"))

	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil).Maybe()
	notifRepo.On("CreateBatch", mock.Anything).Return(nil).Maybe()

	result, err := svc.Create(post)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, post, result) // FindByIDが失敗した場合は元のpostを返す
	postRepo.AssertExpectations(t)
}

func TestPostCreate_SetsEstimatedReadTime(t *testing.T) {
	svc, postRepo, notifRepo := newTestPostService()

	post := &model.Post{Title: "Test", Content: "テスト投稿の本文です。", UserID: 1}

	postRepo.On("Create", mock.MatchedBy(func(p *model.Post) bool {
		return p.EstimatedReadTime >= 1
	})).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 10
	}).Return(nil)

	createdPost := &model.Post{Title: "Test", Content: "テスト投稿の本文です。", UserID: 1, EstimatedReadTime: 1}
	createdPost.ID = 10
	postRepo.On("FindByID", uint(10)).Return(createdPost, nil)

	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil).Maybe()
	notifRepo.On("CreateBatch", mock.Anything).Return(nil).Maybe()

	result, err := svc.Create(post)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.EstimatedReadTime)
}

// ============================================================
// Update 追加テスト
// ============================================================

func TestPostUpdate_WithImageURLs(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	// imageUrlsをJSON配列形式で指定してアップデート
	result, err := svc.Update(1, 1, "New Title", "New Content", `["https://example.com/img.jpg"]`)
	assert.NoError(t, err)
	assert.Equal(t, `["https://example.com/img.jpg"]`, result.ImageURLs)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(errors.New("db error"))

	_, err := svc.Update(1, 1, "New Title", "New Content", "")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	_, err := svc.Update(99, 1, "Title", "Content", "")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostCreate_CreateRepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{Title: "Test Post", Content: "Content", UserID: 1}

	postRepo.On("Create", post).Return(errors.New("db error"))

	result, err := svc.Create(post)
	assert.Error(t, err)
	assert.Nil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_ValidationError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)

	// タイトルが長すぎる場合 → バリデーションエラー
	longTitle := strings.Repeat("a", 300)
	_, err := svc.Update(1, 1, longTitle, "", "")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_WhitespaceOnlyTitle(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	// 空白のみのタイトル → 変更なし（エラーにならない）
	result, err := svc.Update(1, 1, "   ", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "Old Title", result.Title) // 元のタイトルが維持される
	postRepo.AssertExpectations(t)
}

func TestPostUpdate_WhitespaceOnlyContent(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "Old Title", Content: "Old Content", UserID: 1}
	existing.ID = 1

	postRepo.On("FindByID", uint(1)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	// 空白のみの本文 → 変更なし（エラーにならない）
	result, err := svc.Update(1, 1, "", "   ", "")
	assert.NoError(t, err)
	assert.Equal(t, "Old Content", result.Content) // 元の本文が維持される
	postRepo.AssertExpectations(t)
}

// ============================================================
// 投稿非公開化（Unpublish）追加テスト
// ============================================================

func TestPostUnpublish_UpdateError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "Published", UserID: 1, IsDraft: false}
	published.ID = 1

	postRepo.On("FindByID", uint(1)).Return(published, nil)
	postRepo.On("Update", published).Return(errors.New("db error"))

	result, err := svc.Unpublish(1, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "db error")
	postRepo.AssertExpectations(t)
}

// ============================================================
// 投稿取得系 エラーパステスト
// ============================================================

func TestPostCountAll_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountAll").Return(int64(42), nil)

	count, err := svc.CountAll()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)
	postRepo.AssertExpectations(t)
}

func TestPostCountAll_Zero(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountAll").Return(int64(0), nil)

	count, err := svc.CountAll()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
	postRepo.AssertExpectations(t)
}

func TestPostCountAll_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountAll").Return(int64(0), errors.New("db error"))

	count, err := svc.CountAll()
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	postRepo.AssertExpectations(t)
}

func TestPostGetAll_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindAll", 1, 10).Return([]model.Post{}, errors.New("db error"))

	result, err := svc.GetAll(1, 10)
	assert.Error(t, err)
	assert.Empty(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostGetByUserID_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByUserID", uint(1), 20, 0).Return([]model.Post{}, int64(0), errors.New("db error"))

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	postRepo.AssertExpectations(t)
}

func TestPostDeleteComment_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{UserID: 1, PostID: 1, Content: "test"}
	comment.ID = 10

	postRepo.On("FindCommentByID", uint(10)).Return(comment, nil)
	postRepo.On("DeleteComment", uint(10)).Return(errors.New("db error"))

	err := svc.DeleteComment(10, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	postRepo.AssertExpectations(t)
}

// ============================================================
// コメント非表示テスト
// ============================================================

func TestPostHideComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "test"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)
	postRepo.On("UpdateComment", mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 5 && c.IsHidden
	})).Return(nil)

	err := svc.HideComment(5, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostHideComment_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "test"}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	err := svc.HideComment(5, 999)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
}

func TestPostHideComment_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindCommentByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.HideComment(99, 1)
	assert.Error(t, err)
}

func TestPostUnhideComment_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "test", IsHidden: true}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)
	postRepo.On("UpdateComment", mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 5 && !c.IsHidden
	})).Return(nil)

	err := svc.UnhideComment(5, 1)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnhideComment_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{PostID: 10, Content: "test", IsHidden: true}
	comment.ID = 5
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(5)).Return(comment, nil)

	err := svc.UnhideComment(5, 999)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeForbidden, domainErr.Code)
}

func TestPostCreateComment_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	comment := &model.Comment{UserID: 1, PostID: 1, Content: "コメント内容"}

	postRepo.On("CreateComment", comment).Return(errors.New("db error"))

	err := svc.CreateComment(comment)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	postRepo.AssertExpectations(t)
}

// ============================================================
// SchedulePublish テスト
// ============================================================

func TestPostSchedulePublish_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "下書き", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)
	postRepo.On("Update", draft).Return(nil)

	futureTime := time.Now().Add(24 * time.Hour)
	result, err := svc.SchedulePublish(1, 1, futureTime)
	assert.NoError(t, err)
	assert.NotNil(t, result.ScheduledAt)
	postRepo.AssertExpectations(t)
}

func TestPostSchedulePublish_NotDraft(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "公開済み", UserID: 1, IsDraft: false}
	published.ID = 1

	postRepo.On("FindByID", uint(1)).Return(published, nil)

	futureTime := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, futureTime)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestPostSchedulePublish_PastTime(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "下書き", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)

	pastTime := time.Now().Add(-1 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, pastTime)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestPostSchedulePublish_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "他人の下書き", UserID: 2, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)

	futureTime := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, futureTime)
	assert.Error(t, err)
}

func TestPostSchedulePublish_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	futureTime := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(99, 1, futureTime)
	assert.Error(t, err)
}

func TestPostSchedulePublish_UpdateError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	draft := &model.Post{Title: "下書き", UserID: 1, IsDraft: true}
	draft.ID = 1

	postRepo.On("FindByID", uint(1)).Return(draft, nil)
	postRepo.On("Update", draft).Return(errors.New("db error"))

	futureTime := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, futureTime)
	assert.Error(t, err)
}

// ============================================================
// CancelSchedule テスト
// ============================================================

func TestPostCancelSchedule_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	scheduledTime := time.Now().Add(24 * time.Hour)
	post := &model.Post{Title: "予約投稿", UserID: 1, IsDraft: true, ScheduledAt: &scheduledTime}
	post.ID = 1

	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Update", post).Return(nil)

	result, err := svc.CancelSchedule(1, 1)
	assert.NoError(t, err)
	assert.Nil(t, result.ScheduledAt)
	postRepo.AssertExpectations(t)
}

func TestPostCancelSchedule_NotScheduled(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	post := &model.Post{Title: "通常の下書き", UserID: 1, IsDraft: true, ScheduledAt: nil}
	post.ID = 1

	postRepo.On("FindByID", uint(1)).Return(post, nil)

	_, err := svc.CancelSchedule(1, 1)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestPostCancelSchedule_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	scheduledTime := time.Now().Add(24 * time.Hour)
	post := &model.Post{Title: "他人の予約投稿", UserID: 2, ScheduledAt: &scheduledTime}
	post.ID = 1

	postRepo.On("FindByID", uint(1)).Return(post, nil)

	_, err := svc.CancelSchedule(1, 1)
	assert.Error(t, err)
}

func TestPostCancelSchedule_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	_, err := svc.CancelSchedule(99, 1)
	assert.Error(t, err)
}

func TestPostCancelSchedule_UpdateError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	scheduledTime := time.Now().Add(24 * time.Hour)
	post := &model.Post{Title: "予約投稿", UserID: 1, IsDraft: true, ScheduledAt: &scheduledTime}
	post.ID = 1

	postRepo.On("FindByID", uint(1)).Return(post, nil)
	postRepo.On("Update", post).Return(errors.New("db error"))

	_, err := svc.CancelSchedule(1, 1)
	assert.Error(t, err)
}

// ============================================================
// GetScheduled テスト
// ============================================================

func TestPostGetScheduled_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	scheduled := []model.Post{
		{Title: "予約1"},
		{Title: "予約2"},
	}
	postRepo.On("FindScheduledByUserID", uint(1)).Return(scheduled, nil)

	result, err := svc.GetScheduled(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPostGetScheduled_Empty(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindScheduledByUserID", uint(1)).Return([]model.Post{}, nil)

	result, err := svc.GetScheduled(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestPostGetScheduled_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindScheduledByUserID", uint(1)).Return([]model.Post{}, errors.New("db error"))

	_, err := svc.GetScheduled(1)
	assert.Error(t, err)
}

// ============================================================
// 下書き自動保存（AutoSaveDraft）テスト
// ============================================================

func TestPostAutoSaveDraft_CreateNew(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Create", mock.MatchedBy(func(p *model.Post) bool {
		return p.UserID == 1 && p.Title == "タイトル" && p.Content == "本文" && p.IsDraft
	})).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 100
	}).Return(nil)

	result, err := svc.AutoSaveDraft(1, 0, "タイトル", "本文", "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(100), result.ID)
	assert.True(t, result.IsDraft)
	postRepo.AssertExpectations(t)
}

func TestPostAutoSaveDraft_UpdateExisting(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "旧タイトル", Content: "旧本文", UserID: 1, IsDraft: true}
	existing.ID = 5

	postRepo.On("FindByID", uint(5)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	result, err := svc.AutoSaveDraft(1, 5, "新タイトル", "新本文", "")
	assert.NoError(t, err)
	assert.Equal(t, "新タイトル", result.Title)
	assert.Equal(t, "新本文", result.Content)
	postRepo.AssertExpectations(t)
}

func TestPostAutoSaveDraft_AllowEmptyTitle(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Create", mock.MatchedBy(func(p *model.Post) bool {
		return p.Title == "" && p.Content == "書き始め" && p.IsDraft
	})).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 101
	}).Return(nil)

	result, err := svc.AutoSaveDraft(1, 0, "", "書き始め", "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostAutoSaveDraft_AllowEmptyContent(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Create", mock.MatchedBy(func(p *model.Post) bool {
		return p.Title == "タイトルのみ" && p.Content == "" && p.IsDraft
	})).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 102
	}).Return(nil)

	result, err := svc.AutoSaveDraft(1, 0, "タイトルのみ", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	postRepo.AssertExpectations(t)
}

func TestPostAutoSaveDraft_Forbidden(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{Title: "他人の下書き", UserID: 2, IsDraft: true}
	existing.ID = 5

	postRepo.On("FindByID", uint(5)).Return(existing, nil)

	result, err := svc.AutoSaveDraft(1, 5, "ハック", "", "")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestPostAutoSaveDraft_NotDraft(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	published := &model.Post{Title: "公開済み", UserID: 1, IsDraft: false}
	published.ID = 5

	postRepo.On("FindByID", uint(5)).Return(published, nil)

	result, err := svc.AutoSaveDraft(1, 5, "更新", "", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostAutoSaveDraft_TitleTooLong(t *testing.T) {
	svc, _, _ := newTestPostService()

	longTitle := strings.Repeat("あ", 201)
	result, err := svc.AutoSaveDraft(1, 0, longTitle, "", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostAutoSaveDraft_NotFound(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.AutoSaveDraft(1, 99, "タイトル", "本文", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostAutoSaveDraft_CreateRepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Create", mock.Anything).Return(errors.New("db error"))

	result, err := svc.AutoSaveDraft(1, 0, "タイトル", "本文", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// リアクション一括取得テスト
// ============================================================

func TestPostGetReactionsBatch_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postIDs := []uint{1, 2, 3}

	postRepo.On("GetReactionsBatch", postIDs).Return(map[uint][]model.ReactionCount{
		1: {{Emoji: "👍", Count: 5}, {Emoji: "🎉", Count: 2}},
		2: {{Emoji: "❤️", Count: 1}},
	}, nil)
	postRepo.On("GetUserReactionsBatch", uint(10), postIDs).Return(map[uint][]string{
		1: {"👍"},
	}, nil)

	reactions, userReactions, err := svc.GetReactionsBatch(10, postIDs)
	assert.NoError(t, err)
	assert.Len(t, reactions[1], 2)
	assert.Len(t, reactions[2], 1)
	assert.Equal(t, []model.ReactionCount{}, reactions[3])
	assert.Equal(t, []string{"👍"}, userReactions[1])
	assert.Equal(t, []string{}, userReactions[3])
	postRepo.AssertExpectations(t)
}

func TestPostGetReactionsBatch_EmptyPostIDs(t *testing.T) {
	svc, _, _ := newTestPostService()

	reactions, userReactions, err := svc.GetReactionsBatch(1, []uint{})
	assert.NoError(t, err)
	assert.Empty(t, reactions)
	assert.Empty(t, userReactions)
}

func TestPostGetReactionsBatch_TooManyPostIDs(t *testing.T) {
	svc, _, _ := newTestPostService()

	postIDs := make([]uint, 51)
	for i := range postIDs {
		postIDs[i] = uint(i + 1)
	}

	_, _, err := svc.GetReactionsBatch(1, postIDs)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestPostGetReactionsBatch_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postIDs := []uint{1, 2}
	postRepo.On("GetReactionsBatch", postIDs).Return(
		map[uint][]model.ReactionCount(nil), errors.New("db error"),
	)

	_, _, err := svc.GetReactionsBatch(1, postIDs)
	assert.Error(t, err)
}

func TestPostGetReactionsBatch_UserReactionsRepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postIDs := []uint{1}
	postRepo.On("GetReactionsBatch", postIDs).Return(map[uint][]model.ReactionCount{}, nil)
	postRepo.On("GetUserReactionsBatch", uint(1), postIDs).Return(
		map[uint][]string(nil), errors.New("db error"),
	)

	_, _, err := svc.GetReactionsBatch(1, postIDs)
	assert.Error(t, err)
}

func TestPostAutoSaveDraft_ContentTooLong(t *testing.T) {
	svc, _, _ := newTestPostService()

	longContent := strings.Repeat("あ", 50001)
	_, err := svc.AutoSaveDraft(1, 0, "title", longContent, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "50000文字")
}

func TestPostAutoSaveDraft_UpdateRepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{UserID: 1, Title: "下書き", IsDraft: true}
	existing.ID = 5
	postRepo.On("FindByID", uint(5)).Return(existing, nil)
	postRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	_, err := svc.AutoSaveDraft(1, 5, "更新タイトル", "更新内容", "")
	assert.Error(t, err)
}

func TestPostAutoSaveDraft_UpdateWithImageURLs(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	existing := &model.Post{UserID: 1, Title: "下書き", IsDraft: true}
	existing.ID = 5
	postRepo.On("FindByID", uint(5)).Return(existing, nil)
	postRepo.On("Update", mock.MatchedBy(func(p *model.Post) bool {
		return p.ImageURLs == "https://example.com/img.png"
	})).Return(nil)

	result, err := svc.AutoSaveDraft(1, 5, "title", "content", "https://example.com/img.png")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/img.png", result.ImageURLs)
}

// ============================================================
// NormalizeReactionMaps テスト
// ============================================================

func TestNormalizeReactionMaps_NilToEmpty(t *testing.T) {
	reactions := map[uint][]model.ReactionCount{
		1: {{Emoji: "👍", Count: 3}},
	}
	userReactions := map[uint][]string{
		1: {"👍"},
	}
	postIDs := []uint{1, 2, 3}

	NormalizeReactionMaps(reactions, userReactions, postIDs)

	assert.Equal(t, []model.ReactionCount{{Emoji: "👍", Count: 3}}, reactions[1])
	assert.Equal(t, []model.ReactionCount{}, reactions[2])
	assert.Equal(t, []model.ReactionCount{}, reactions[3])
	assert.Equal(t, []string{"👍"}, userReactions[1])
	assert.Equal(t, []string{}, userReactions[2])
	assert.Equal(t, []string{}, userReactions[3])
}

func TestNormalizeReactionMaps_EmptyPostIDs(t *testing.T) {
	reactions := map[uint][]model.ReactionCount{}
	userReactions := map[uint][]string{}
	NormalizeReactionMaps(reactions, userReactions, []uint{})
	assert.Empty(t, reactions)
	assert.Empty(t, userReactions)
}

// ============================================================
// GetReactionsWithUser テスト
// ============================================================

func TestGetReactionsWithUser_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReactionsByPostID", uint(1)).Return([]model.ReactionCount{
		{Emoji: "👍", Count: 5},
	}, nil)
	postRepo.On("GetUserReactions", uint(10), uint(1)).Return([]string{"👍"}, nil)

	reactions, userReactions, err := svc.GetReactionsWithUser(10, 1)
	assert.NoError(t, err)
	assert.Len(t, reactions, 1)
	assert.Equal(t, []string{"👍"}, userReactions)
	postRepo.AssertExpectations(t)
}

func TestGetReactionsWithUser_NilNormalized(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReactionsByPostID", uint(1)).Return([]model.ReactionCount(nil), nil)
	postRepo.On("GetUserReactions", uint(10), uint(1)).Return([]string(nil), nil)

	reactions, userReactions, err := svc.GetReactionsWithUser(10, 1)
	assert.NoError(t, err)
	assert.Equal(t, []model.ReactionCount{}, reactions)
	assert.Equal(t, []string{}, userReactions)
}

func TestGetReactionsWithUser_RepoError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReactionsByPostID", uint(1)).Return([]model.ReactionCount(nil), assert.AnError)

	_, _, err := svc.GetReactionsWithUser(10, 1)
	assert.Error(t, err)
}

func TestGetReactionsWithUser_UserReactionsError(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("GetReactionsByPostID", uint(1)).Return([]model.ReactionCount{}, nil)
	postRepo.On("GetUserReactions", uint(10), uint(1)).Return([]string(nil), assert.AnError)

	_, _, err := svc.GetReactionsWithUser(10, 1)
	assert.Error(t, err)
}
