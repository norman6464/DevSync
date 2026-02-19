package service

import (
	"errors"
	"testing"

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

	postRepo.On("Bookmark", uint(1), uint(10)).Return(nil)

	err := svc.Bookmark(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostBookmark_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Bookmark", uint(1), uint(10)).Return(errors.New("already bookmarked"))

	err := svc.Bookmark(1, 10)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnbookmark_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Unbookmark", uint(1), uint(10)).Return(nil)

	err := svc.Unbookmark(1, 10)
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostUnbookmark_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("Unbookmark", uint(1), uint(10)).Return(errors.New("not bookmarked"))

	err := svc.Unbookmark(1, 10)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
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

	postRepo.On("AddReaction", uint(1), uint(10), "🎉").Return(errors.New("duplicate"))

	err := svc.AddReaction(1, 10, "🎉")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostRemoveReaction_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("RemoveReaction", uint(1), uint(10), "👍").Return(nil)

	err := svc.RemoveReaction(1, 10, "👍")
	assert.NoError(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostRemoveReaction_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("RemoveReaction", uint(1), uint(10), "👍").Return(errors.New("not found"))

	err := svc.RemoveReaction(1, 10, "👍")
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
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
