package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestMentionService はMentionServiceのテスト用インスタンスを生成するヘルパー。
func newTestMentionService() (*MentionService, *MockMentionRepository, *MockUserRepository, *MockNotificationService) {
	mentionRepo := new(MockMentionRepository)
	userRepo := new(MockUserRepository)
	notifSvc := new(MockNotificationService)
	svc := NewMentionService(mentionRepo, userRepo, notifSvc)
	return svc, mentionRepo, userRepo, notifSvc
}

// ============================================================
// ParseMentions テスト（テキストから @username を抽出）
// ============================================================

func TestParseMentions_SingleMention(t *testing.T) {
	usernames := ParseMentions("こんにちは @alice さん")
	assert.Equal(t, []string{"alice"}, usernames)
}

func TestParseMentions_MultipleMentions(t *testing.T) {
	usernames := ParseMentions("@alice と @bob に聞いてみて")
	assert.Equal(t, []string{"alice", "bob"}, usernames)
}

func TestParseMentions_NoDuplicates(t *testing.T) {
	usernames := ParseMentions("@alice と @alice にメンション")
	assert.Equal(t, []string{"alice"}, usernames)
}

func TestParseMentions_NoMentions(t *testing.T) {
	usernames := ParseMentions("メンションなしのテキスト")
	assert.Empty(t, usernames)
}

func TestParseMentions_EmailNotMatched(t *testing.T) {
	usernames := ParseMentions("test@example.com は無視")
	assert.Empty(t, usernames)
}

func TestParseMentions_AtStartOfText(t *testing.T) {
	usernames := ParseMentions("@alice こんにちは")
	assert.Equal(t, []string{"alice"}, usernames)
}

func TestParseMentions_WithHyphenAndUnderscore(t *testing.T) {
	usernames := ParseMentions("@user-name と @user_name2 に")
	assert.Equal(t, []string{"user-name", "user_name2"}, usernames)
}

// ============================================================
// ProcessMentions テスト（メンション処理 + 通知）
// ============================================================

func TestProcessMentions_PostMention_Success(t *testing.T) {
	svc, mentionRepo, userRepo, notifSvc := newTestMentionService()
	postID := uint(10)

	userRepo.On("FindByUsername", "alice").Return(&model.User{
		ID: 2, Username: "alice",
	}, nil)

	mentionRepo.On("Create", mock.MatchedBy(func(m *model.Mention) bool {
		return m.UserID == 2 && m.ActorID == 1 && *m.PostID == 10 && m.CommentID == nil
	})).Return(nil)

	notifSvc.On("CreateNotification", mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 2 && n.ActorID == 1 && n.Type == model.NotificationTypeMention
	})).Return(nil)

	err := svc.ProcessMentions(1, "こんにちは @alice さん", &postID, nil)
	assert.NoError(t, err)
	mentionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	notifSvc.AssertExpectations(t)
}

func TestProcessMentions_CommentMention_Success(t *testing.T) {
	svc, mentionRepo, userRepo, notifSvc := newTestMentionService()
	commentID := uint(20)

	userRepo.On("FindByUsername", "bob").Return(&model.User{
		ID: 3, Username: "bob",
	}, nil)

	mentionRepo.On("Create", mock.MatchedBy(func(m *model.Mention) bool {
		return m.UserID == 3 && m.ActorID == 1 && m.PostID == nil && *m.CommentID == 20
	})).Return(nil)

	notifSvc.On("CreateNotification", mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 3 && n.ActorID == 1 && n.Type == model.NotificationTypeMention
	})).Return(nil)

	err := svc.ProcessMentions(1, "@bob よろしく", nil, &commentID)
	assert.NoError(t, err)
	mentionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	notifSvc.AssertExpectations(t)
}

func TestProcessMentions_SkipSelfMention(t *testing.T) {
	svc, mentionRepo, userRepo, _ := newTestMentionService()
	postID := uint(10)

	// actorID=1 が自分自身をメンション
	userRepo.On("FindByUsername", "myself").Return(&model.User{
		ID: 1, Username: "myself",
	}, nil)

	err := svc.ProcessMentions(1, "@myself テスト", &postID, nil)
	assert.NoError(t, err)
	// 自分自身へのメンションはCreateされない
	mentionRepo.AssertNotCalled(t, "Create")
}

func TestProcessMentions_UserNotFound(t *testing.T) {
	svc, mentionRepo, userRepo, _ := newTestMentionService()
	postID := uint(10)

	userRepo.On("FindByUsername", "unknown").Return(nil, errors.New("not found"))

	err := svc.ProcessMentions(1, "@unknown テスト", &postID, nil)
	// 存在しないユーザーは無視してエラーにしない
	assert.NoError(t, err)
	mentionRepo.AssertNotCalled(t, "Create")
}

func TestProcessMentions_MultipleMentions(t *testing.T) {
	svc, mentionRepo, userRepo, notifSvc := newTestMentionService()
	postID := uint(10)

	userRepo.On("FindByUsername", "alice").Return(&model.User{
		ID: 2, Username: "alice",
	}, nil)
	userRepo.On("FindByUsername", "bob").Return(&model.User{
		ID: 3, Username: "bob",
	}, nil)

	mentionRepo.On("Create", mock.MatchedBy(func(m *model.Mention) bool {
		return m.UserID == 2 && m.ActorID == 1
	})).Return(nil)
	mentionRepo.On("Create", mock.MatchedBy(func(m *model.Mention) bool {
		return m.UserID == 3 && m.ActorID == 1
	})).Return(nil)

	notifSvc.On("CreateNotification", mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 2
	})).Return(nil)
	notifSvc.On("CreateNotification", mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 3
	})).Return(nil)

	err := svc.ProcessMentions(1, "@alice と @bob へ", &postID, nil)
	assert.NoError(t, err)
	mentionRepo.AssertNumberOfCalls(t, "Create", 2)
	notifSvc.AssertNumberOfCalls(t, "CreateNotification", 2)
}

func TestProcessMentions_NoMentions(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()
	postID := uint(10)

	err := svc.ProcessMentions(1, "メンションなし", &postID, nil)
	assert.NoError(t, err)
	mentionRepo.AssertNotCalled(t, "Create")
}

func TestProcessMentions_CreateError(t *testing.T) {
	svc, mentionRepo, userRepo, _ := newTestMentionService()
	postID := uint(10)

	userRepo.On("FindByUsername", "alice").Return(&model.User{
		ID: 2, Username: "alice",
	}, nil)

	mentionRepo.On("Create", mock.Anything).Return(errors.New("db error"))

	err := svc.ProcessMentions(1, "@alice テスト", &postID, nil)
	assert.Error(t, err)
}

// ============================================================
// GetMentionsByUserID テスト
// ============================================================

func TestGetMentionsByUserID_Success(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	postID := uint(10)
	mentions := []model.Mention{
		{ID: 1, UserID: 2, ActorID: 1, PostID: &postID},
	}
	mentionRepo.On("FindByUserID", uint(2), 1, 20).Return(mentions, nil)

	result, err := svc.GetMentionsByUserID(2, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mentionRepo.AssertExpectations(t)
}

func TestGetMentionsByUserID_Empty(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("FindByUserID", uint(2), 1, 20).Return([]model.Mention{}, nil)

	result, err := svc.GetMentionsByUserID(2, 1, 20)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetMentionsByUserID_RepoError(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("FindByUserID", uint(2), 1, 20).Return([]model.Mention(nil), errors.New("db error"))

	result, err := svc.GetMentionsByUserID(2, 1, 20)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// GetMentionsByPostID テスト
// ============================================================

func TestGetMentionsByPostID_Success(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	postID := uint(10)
	mentions := []model.Mention{
		{ID: 1, UserID: 2, ActorID: 1, PostID: &postID},
	}
	mentionRepo.On("FindByPostID", uint(10)).Return(mentions, nil)

	result, err := svc.GetMentionsByPostID(10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mentionRepo.AssertExpectations(t)
}

func TestGetMentionsByPostID_RepoError(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("FindByPostID", uint(10)).Return([]model.Mention(nil), errors.New("db error"))

	result, err := svc.GetMentionsByPostID(10)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// DeleteMentionsByPostID テスト
// ============================================================

func TestDeleteMentionsByPostID_Success(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("DeleteByPostID", uint(10)).Return(nil)

	err := svc.DeleteMentionsByPostID(10)
	assert.NoError(t, err)
	mentionRepo.AssertExpectations(t)
}

func TestDeleteMentionsByPostID_RepoError(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("DeleteByPostID", uint(10)).Return(errors.New("db error"))

	err := svc.DeleteMentionsByPostID(10)
	assert.Error(t, err)
}

// ============================================================
// DeleteMentionsByCommentID テスト
// ============================================================

func TestDeleteMentionsByCommentID_Success(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("DeleteByCommentID", uint(20)).Return(nil)

	err := svc.DeleteMentionsByCommentID(20)
	assert.NoError(t, err)
	mentionRepo.AssertExpectations(t)
}

func TestDeleteMentionsByCommentID_RepoError(t *testing.T) {
	svc, mentionRepo, _, _ := newTestMentionService()

	mentionRepo.On("DeleteByCommentID", uint(20)).Return(errors.New("db error"))

	err := svc.DeleteMentionsByCommentID(20)
	assert.Error(t, err)
}
