package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockMentionRepo は usecase/repository.MentionRepository のモック。
type mockMentionRepo struct{ mock.Mock }

func (m *mockMentionRepo) Create(ctx context.Context, mention *model.Mention) error {
	return m.Called(ctx, mention).Error(0)
}

func (m *mockMentionRepo) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error) {
	args := m.Called(ctx, userID, page, limit)
	ms, _ := args.Get(0).([]model.Mention)
	return ms, args.Error(1)
}

func (m *mockMentionRepo) FindByPostID(ctx context.Context, postID uint) ([]model.Mention, error) {
	args := m.Called(ctx, postID)
	ms, _ := args.Get(0).([]model.Mention)
	return ms, args.Error(1)
}

func (m *mockMentionRepo) DeleteByPostID(ctx context.Context, postID uint) error {
	return m.Called(ctx, postID).Error(0)
}

func (m *mockMentionRepo) DeleteByCommentID(ctx context.Context, commentID uint) error {
	return m.Called(ctx, commentID).Error(0)
}

// mockNotificationCreatorPort は usecase/repository.NotificationCreator のモック。
type mockNotificationCreatorPort struct{ mock.Mock }

func (m *mockNotificationCreatorPort) Create(ctx context.Context, notification *model.Notification) error {
	return m.Called(ctx, notification).Error(0)
}

// mockUsernameLookup は usecase/repository.UsernameLookup のモック。
type mockUsernameLookup struct{ mock.Mock }

func (m *mockUsernameLookup) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

// ============================================================
// @username の抽出
// ============================================================

func TestParseMentions(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{"単一のメンション", "こんにちは @alice さん", []string{"alice"}},
		{"複数のメンション", "@alice と @bob へ", []string{"alice", "bob"}},
		{"重複は 1 回だけ", "@alice @alice", []string{"alice"}},
		{"大文字は小文字に揃える", "@Alice", []string{"alice"}},
		{"行頭のメンション", "@alice おはよう", []string{"alice"}},
		{"アンダースコアとハイフン", "@my_user-name", []string{"my_user-name"}},
		{"メールアドレスは拾わない", "連絡は foo@example.com まで", nil},
		{"メンションが無い", "ただのテキスト", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, usecase.ParseMentions(tt.text))
		})
	}
}

// ============================================================
// メンション処理
// ============================================================

func TestProcessMentionsUseCase(t *testing.T) {
	postID := uint(10)

	t.Run("メンションを記録して通知する", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		users.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
		mentions.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Mention) bool {
			return m.UserID == 2 && m.ActorID == 1 && m.PostID != nil && *m.PostID == postID
		})).Return(nil)
		notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
			return n.UserID == 2 && n.ActorID == 1 && n.Type == model.NotificationTypeMention &&
				n.PostID != nil && *n.PostID == postID
		})).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 1, "やあ @alice", &postID, nil))
		mentions.AssertExpectations(t)
		users.AssertExpectations(t)
		notifications.AssertExpectations(t)
	})

	t.Run("存在しないユーザーはスキップする", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		users.On("FindByUsername", mock.Anything, "ghost").Return(nil, nil)

		require.NoError(t, uc.Execute(context.Background(), 1, "@ghost", &postID, nil))
		mentions.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("ユーザー取得に失敗してもスキップして続行する", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		users.On("FindByUsername", mock.Anything, "broken").Return(nil, errors.New("db error"))
		users.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
		mentions.On("Create", mock.Anything, mock.Anything).Return(nil)
		notifications.On("Create", mock.Anything, mock.Anything).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 1, "@broken と @alice", &postID, nil))
		mentions.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("自分自身へのメンションはスキップする", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		users.On("FindByUsername", mock.Anything, "me").Return(&model.User{ID: 1}, nil)

		require.NoError(t, uc.Execute(context.Background(), 1, "@me メモ", &postID, nil))
		mentions.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("メンションが無ければ何もしない", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		require.NoError(t, uc.Execute(context.Background(), 1, "メンション無し", &postID, nil))
		users.AssertNotCalled(t, "FindByUsername", mock.Anything, mock.Anything)
		mentions.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("コメントへのメンションは通知に投稿 ID を入れない", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		commentID := uint(20)
		users.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
		mentions.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Mention) bool {
			return m.CommentID != nil && *m.CommentID == commentID && m.PostID == nil
		})).Return(nil)
		notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
			return n.PostID == nil
		})).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 1, "@alice", nil, &commentID))
		mentions.AssertExpectations(t)
		notifications.AssertExpectations(t)
	})

	t.Run("メンションの保存に失敗したらエラーを返す", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		createErr := errors.New("db error")
		users.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
		mentions.On("Create", mock.Anything, mock.Anything).Return(createErr)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1, "@alice", &postID, nil), createErr)
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("通知の失敗は無視する", func(t *testing.T) {
		mentions := new(mockMentionRepo)
		users := new(mockUsernameLookup)
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewProcessMentionsUseCase(mentions, users, notifications)

		users.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
		mentions.On("Create", mock.Anything, mock.Anything).Return(nil)
		notifications.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		require.NoError(t, uc.Execute(context.Background(), 1, "@alice", &postID, nil))
	})
}

// ============================================================
// 参照・削除
// ============================================================

func TestListMentionUseCases(t *testing.T) {
	mentions := new(mockMentionRepo)

	mentions.On("FindByUserID", mock.Anything, uint(1), 1, 20).Return([]model.Mention{{ID: 1}}, nil)
	mine, err := usecase.NewListUserMentionsUseCase(mentions).Execute(context.Background(), 1, 1, 20)
	require.NoError(t, err)
	assert.Len(t, mine, 1)

	mentions.On("FindByPostID", mock.Anything, uint(7)).Return([]model.Mention{{ID: 2}}, nil)
	byPost, err := usecase.NewListPostMentionsUseCase(mentions).Execute(context.Background(), 7)
	require.NoError(t, err)
	assert.Len(t, byPost, 1)

	mentions.AssertExpectations(t)
}

func TestDeleteMentionUseCases(t *testing.T) {
	mentions := new(mockMentionRepo)

	mentions.On("DeleteByPostID", mock.Anything, uint(7)).Return(nil)
	require.NoError(t, usecase.NewDeletePostMentionsUseCase(mentions).Execute(context.Background(), 7))

	mentions.On("DeleteByCommentID", mock.Anything, uint(20)).Return(nil)
	require.NoError(t, usecase.NewDeleteCommentMentionsUseCase(mentions).Execute(context.Background(), 20))

	mentions.AssertExpectations(t)
}
