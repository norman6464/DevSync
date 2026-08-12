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

// リアクションは DIP へ移行済み。テストは handler/usecase 側（port モック）に移した。

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
// CountDraftsByUserID テスト
// ============================================================

func TestPostCountDraftsByUserID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountDraftsByUserID", uint(1)).Return(int64(3), nil)

	count, err := svc.CountDraftsByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	postRepo.AssertExpectations(t)
}

func TestPostCountDraftsByUserID_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountDraftsByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	_, err := svc.CountDraftsByUserID(1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

// ============================================================
// CountScheduledByUserID テスト
// ============================================================

func TestPostCountScheduledByUserID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountScheduledByUserID", uint(1)).Return(int64(2), nil)

	count, err := svc.CountScheduledByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	postRepo.AssertExpectations(t)
}

func TestPostCountScheduledByUserID_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountScheduledByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	_, err := svc.CountScheduledByUserID(1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

// ============================================================
// CountBookmarkedByUserID テスト
// ============================================================

func TestPostCountBookmarkedByUserID_Success(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountBookmarkedByUserID", uint(1)).Return(int64(7), nil)

	count, err := svc.CountBookmarkedByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), count)
	postRepo.AssertExpectations(t)
}

func TestPostCountBookmarkedByUserID_Error(t *testing.T) {
	svc, postRepo, _ := newTestPostService()

	postRepo.On("CountBookmarkedByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	_, err := svc.CountBookmarkedByUserID(1)
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}
