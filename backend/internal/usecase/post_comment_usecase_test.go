package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPostCommentRepo は usecase/repository.PostCommentRepository のモック。
type mockPostCommentRepo struct{ mock.Mock }

func (m *mockPostCommentRepo) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.Comment)
	return c, args.Error(1)
}

func (m *mockPostCommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	return m.Called(ctx, comment).Error(0)
}

func (m *mockPostCommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	return m.Called(ctx, comment).Error(0)
}

func (m *mockPostCommentRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostCommentRepo) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	args := m.Called(ctx, postID)
	cs, _ := args.Get(0).([]model.Comment)
	return cs, args.Error(1)
}

func (m *mockPostCommentRepo) ListReplies(ctx context.Context, parentID uint) ([]model.Comment, error) {
	args := m.Called(ctx, parentID)
	cs, _ := args.Get(0).([]model.Comment)
	return cs, args.Error(1)
}

// ============================================================
// 作成
// ============================================================

func TestCreatePostCommentUseCase(t *testing.T) {
	t.Run("コメントを作成する", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
			return c.UserID == 1 && c.PostID == 5 && c.Content == "hello" && c.ParentID == nil
		})).Return(nil)

		got, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "  hello  ", nil)
		require.NoError(t, err)
		assert.Equal(t, "hello", got.Content)
		comments.AssertExpectations(t)
	})

	t.Run("空白だけの本文は弾く", func(t *testing.T) {
		comments := new(mockPostCommentRepo)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "   ", nil)
		assert.Error(t, err)
		comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("1000 文字を超える本文は弾く", func(t *testing.T) {
		comments := new(mockPostCommentRepo)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, strings.Repeat("a", 1001), nil)
		assert.Error(t, err)
		comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("返信は親コメントを検証して作成する", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		parentID := uint(10)
		comments.On("FindCommentByID", mock.Anything, parentID).Return(&model.Comment{ID: 10, PostID: 5}, nil)
		comments.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
			return c.ParentID != nil && *c.ParentID == parentID
		})).Return(nil)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "reply", &parentID)
		require.NoError(t, err)
		comments.AssertExpectations(t)
	})

	t.Run("親コメントが無ければ 404 相当のエラー", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		parentID := uint(99)
		comments.On("FindCommentByID", mock.Anything, parentID).Return(nil, pgx.ErrNoRows)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "reply", &parentID)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeNotFound, domainErr.Code)
		comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("親コメントが別投稿なら 400 相当のエラー", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		parentID := uint(10)
		comments.On("FindCommentByID", mock.Anything, parentID).Return(&model.Comment{ID: 10, PostID: 20}, nil)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "reply", &parentID)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("返信への返信は許可しない", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		grandParentID := uint(3)
		parentID := uint(10)
		comments.On("FindCommentByID", mock.Anything, parentID).
			Return(&model.Comment{ID: 10, PostID: 5, ParentID: &grandParentID}, nil)

		_, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "reply", &parentID)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("保存に失敗したらエラーを返す", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		createErr := errors.New("db error")
		comments.On("Create", mock.Anything, mock.Anything).Return(createErr)

		got, err := usecase.NewCreatePostCommentUseCase(comments).Execute(context.Background(), 1, 5, "hello", nil)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, createErr)
	})
}

// ============================================================
// 参照
// ============================================================

func TestListPostCommentUseCases(t *testing.T) {
	comments := new(mockPostCommentRepo)

	comments.On("ListByPostID", mock.Anything, uint(5)).Return([]model.Comment{{ID: 1}}, nil)
	list, err := usecase.NewListPostCommentsUseCase(comments).Execute(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	comments.On("ListReplies", mock.Anything, uint(10)).Return([]model.Comment{{ID: 2}, {ID: 3}}, nil)
	replies, err := usecase.NewListCommentRepliesUseCase(comments).Execute(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, replies, 2)

	comments.AssertExpectations(t)
}

// ============================================================
// 編集
// ============================================================

func TestEditPostCommentUseCase(t *testing.T) {
	t.Run("所有者は本文を更新できる", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1, Content: "old"}, nil)
		comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
			return c.ID == 3 && c.Content == "new"
		})).Return(nil)

		got, err := usecase.NewEditPostCommentUseCase(comments).Execute(context.Background(), 3, 1, "  new  ")
		require.NoError(t, err)
		assert.Equal(t, "new", got.Content)
		comments.AssertExpectations(t)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 999}, nil)

		_, err := usecase.NewEditPostCommentUseCase(comments).Execute(context.Background(), 3, 1, "new")
		assert.ErrorIs(t, err, domain.ErrForbidden)
		comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("取得エラーはそのまま返す", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(nil, pgx.ErrNoRows)

		_, err := usecase.NewEditPostCommentUseCase(comments).Execute(context.Background(), 3, 1, "new")
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("本文の検証は所有権チェックの後に行う", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1}, nil)

		_, err := usecase.NewEditPostCommentUseCase(comments).Execute(context.Background(), 3, 1, "   ")
		assert.Error(t, err)
		comments.AssertExpectations(t)
		comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

// ============================================================
// 削除・非表示
// ============================================================

func TestDeletePostCommentUseCase(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1}, nil)
		comments.On("Delete", mock.Anything, uint(3)).Return(nil)

		require.NoError(t, usecase.NewDeletePostCommentUseCase(comments).Execute(context.Background(), 3, 1))
		comments.AssertExpectations(t)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 999}, nil)

		err := usecase.NewDeletePostCommentUseCase(comments).Execute(context.Background(), 3, 1)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		comments.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

func TestHideAndUnhidePostCommentUseCase(t *testing.T) {
	t.Run("非表示にする", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1}, nil)
		comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool { return c.IsHidden })).Return(nil)

		require.NoError(t, usecase.NewHidePostCommentUseCase(comments).Execute(context.Background(), 3, 1))
		comments.AssertExpectations(t)
	})

	t.Run("非表示を解除する", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1, IsHidden: true}, nil)
		comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool { return !c.IsHidden })).Return(nil)

		require.NoError(t, usecase.NewUnhidePostCommentUseCase(comments).Execute(context.Background(), 3, 1))
		comments.AssertExpectations(t)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		comments := new(mockPostCommentRepo)
		comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 999}, nil)

		assert.ErrorIs(t, usecase.NewHidePostCommentUseCase(comments).Execute(context.Background(), 3, 1), domain.ErrForbidden)
		assert.ErrorIs(t, usecase.NewUnhidePostCommentUseCase(comments).Execute(context.Background(), 3, 1), domain.ErrForbidden)
		comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}
