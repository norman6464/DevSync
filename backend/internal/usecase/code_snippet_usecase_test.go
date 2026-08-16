package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockCodeSnippetRepo は usecase/repository.CodeSnippetRepository のモック。
type mockCodeSnippetRepo struct{ mock.Mock }

func (m *mockCodeSnippetRepo) Create(ctx context.Context, s *model.CodeSnippet) error {
	return m.Called(ctx, s).Error(0)
}

func (m *mockCodeSnippetRepo) FindByID(ctx context.Context, id uint) (*model.CodeSnippet, error) {
	args := m.Called(ctx, id)
	s, _ := args.Get(0).(*model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) FindByPostID(ctx context.Context, postID uint) ([]model.CodeSnippet, error) {
	args := m.Called(ctx, postID)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) FindByUserIDAndLanguage(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error) {
	args := m.Called(ctx, userID, language)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Error(1)
}

func (m *mockCodeSnippetRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Get(1).(int64), args.Error(2)
}

func (m *mockCodeSnippetRepo) Update(ctx context.Context, s *model.CodeSnippet) error {
	return m.Called(ctx, s).Error(0)
}

func (m *mockCodeSnippetRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCodeSnippetRepo) CreateComment(ctx context.Context, c *model.SnippetComment) error {
	return m.Called(ctx, c).Error(0)
}

func (m *mockCodeSnippetRepo) GetComments(ctx context.Context, snippetID uint) ([]model.SnippetComment, error) {
	args := m.Called(ctx, snippetID)
	c, _ := args.Get(0).([]model.SnippetComment)
	return c, args.Error(1)
}

func (m *mockCodeSnippetRepo) FindCommentByID(ctx context.Context, id uint) (*model.SnippetComment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.SnippetComment)
	return c, args.Error(1)
}

func (m *mockCodeSnippetRepo) DeleteComment(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCodeSnippetRepo) IncrementForkCount(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockCodeSnippetRepo) Favorite(ctx context.Context, userID, snippetID uint) error {
	return m.Called(ctx, userID, snippetID).Error(0)
}

func (m *mockCodeSnippetRepo) Unfavorite(ctx context.Context, userID, snippetID uint) error {
	return m.Called(ctx, userID, snippetID).Error(0)
}

func (m *mockCodeSnippetRepo) HasFavorited(ctx context.Context, userID, snippetID uint) (bool, error) {
	args := m.Called(ctx, userID, snippetID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCodeSnippetRepo) FindFavoritedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	s, _ := args.Get(0).([]model.CodeSnippet)
	return s, args.Get(1).(int64), args.Error(2)
}

func (m *mockCodeSnippetRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ownedSnippetOf は所有者が userID=1 のスニペットを返す。
func ownedSnippetOf(id uint) *model.CodeSnippet {
	s := &model.CodeSnippet{Language: "go", Code: "package main", FileName: "main.go"}
	s.ID = id
	s.UserID = 1
	return s
}

func TestCreateCodeSnippetUseCase_Execute(t *testing.T) {
	t.Run("前後空白を除いて作成し、再取得した値を返す", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		snippets.On("Create", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
			return s.Language == "go" && s.Code == "package main" && s.FileName == "main.go"
		})).Return(nil)
		stored := ownedSnippetOf(1)
		snippets.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(stored, nil)
		uc := usecase.NewCreateCodeSnippetUseCase(snippets, posts)

		got, err := uc.Execute(context.Background(), &model.CodeSnippet{
			PostID: 5, Language: "  go  ", Code: "  package main  ", FileName: "  main.go  ",
		})

		assert.NoError(t, err)
		assert.Equal(t, stored, got)
		snippets.AssertExpectations(t)
	})

	// 再取得に失敗しても作成した値を返す（移行前の挙動）。
	t.Run("再取得に失敗しても作成した値を返す", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		snippets.On("Create", mock.Anything, mock.Anything).Return(nil)
		snippets.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
			Return(nil, errors.New("db error"))
		uc := usecase.NewCreateCodeSnippetUseCase(snippets, posts)

		got, err := uc.Execute(context.Background(), &model.CodeSnippet{
			PostID: 5, Language: "go", Code: "x",
		})

		assert.NoError(t, err)
		assert.Equal(t, "go", got.Language)
	})

	t.Run("投稿が不在なら NotFound（作成しない）", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewCreateCodeSnippetUseCase(snippets, posts)

		_, err := uc.Execute(context.Background(), &model.CodeSnippet{PostID: 5, Language: "go", Code: "x"})

		assertDomainCode(t, err, domain.ErrCodeNotFound)
		snippets.AssertNotCalled(t, "Create")
	})

	t.Run("入力が不正なら投稿も読まない", func(t *testing.T) {
		cases := map[string]*model.CodeSnippet{
			"言語が空":          {PostID: 5, Language: "", Code: "x"},
			"コードが空":         {PostID: 5, Language: "go", Code: ""},
			"言語が 101 文字":    {PostID: 5, Language: strings.Repeat("a", 101), Code: "x"},
			"ファイル名が 201 文字": {PostID: 5, Language: "go", Code: "x", FileName: strings.Repeat("a", 201)},
		}
		for name, in := range cases {
			t.Run(name, func(t *testing.T) {
				snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
				uc := usecase.NewCreateCodeSnippetUseCase(snippets, posts)

				_, err := uc.Execute(context.Background(), in)

				assert.Error(t, err)
				posts.AssertNotCalled(t, "FindByID")
				snippets.AssertNotCalled(t, "Create")
			})
		}
	})

	t.Run("作成の DB 障害は Internal に包む", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		snippets.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateCodeSnippetUseCase(snippets, posts)

		_, err := uc.Execute(context.Background(), &model.CodeSnippet{PostID: 5, Language: "go", Code: "x"})

		assertDomainCode(t, err, domain.ErrCodeInternal)
	})
}

func TestSearchCodeSnippetsUseCase_Execute(t *testing.T) {
	t.Run("前後空白を除いて検索する", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("Search", mock.Anything, "go", 20, 0).
			Return([]model.CodeSnippet{*ownedSnippetOf(1)}, int64(1), nil)
		uc := usecase.NewSearchCodeSnippetsUseCase(snippets)

		got, total, err := uc.Execute(context.Background(), "  go  ", 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		snippets.AssertExpectations(t)
	})

	t.Run("空のキーワードは BadRequest（検索しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		uc := usecase.NewSearchCodeSnippetsUseCase(snippets)

		_, _, err := uc.Execute(context.Background(), "   ", 20, 0)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		snippets.AssertNotCalled(t, "Search")
	})
}

func TestListCodeSnippetsByLanguageUseCase_Execute(t *testing.T) {
	t.Run("言語が空なら BadRequest（DB を触らない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		uc := usecase.NewListCodeSnippetsByLanguageUseCase(snippets)

		_, err := uc.Execute(context.Background(), 1, "")

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		snippets.AssertNotCalled(t, "FindByUserIDAndLanguage")
	})

	t.Run("言語で絞り込む", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByUserIDAndLanguage", mock.Anything, uint(1), "go").
			Return([]model.CodeSnippet{*ownedSnippetOf(1)}, nil)
		uc := usecase.NewListCodeSnippetsByLanguageUseCase(snippets)

		got, err := uc.Execute(context.Background(), 1, "go")

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		snippets.AssertExpectations(t)
	})
}

func TestUpdateCodeSnippetUseCase_Execute(t *testing.T) {
	t.Run("空の項目は据え置く部分更新", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("Update", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
			return s.Language == "rust" && s.Code == "package main" && s.FileName == "main.go"
		})).Return(nil)
		uc := usecase.NewUpdateCodeSnippetUseCase(snippets)

		got, err := uc.Execute(context.Background(), usecase.UpdateCodeSnippetInput{
			ID: 1, UserID: 1, Language: "  rust  ",
		})

		assert.NoError(t, err)
		assert.Equal(t, "rust", got.Language)
		snippets.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		other := ownedSnippetOf(1)
		other.UserID = 999
		snippets.On("FindByID", mock.Anything, uint(1)).Return(other, nil)
		uc := usecase.NewUpdateCodeSnippetUseCase(snippets)

		_, err := uc.Execute(context.Background(), usecase.UpdateCodeSnippetInput{ID: 1, UserID: 1, Language: "rust"})

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		snippets.AssertNotCalled(t, "Update")
	})

	t.Run("保存の DB 障害は Internal に包む", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateCodeSnippetUseCase(snippets)

		_, err := uc.Execute(context.Background(), usecase.UpdateCodeSnippetInput{ID: 1, UserID: 1, Language: "rust"})

		assertDomainCode(t, err, domain.ErrCodeInternal)
	})
}

func TestDeleteCodeSnippetUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteCodeSnippetUseCase(snippets)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1))
		snippets.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（削除しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		other := ownedSnippetOf(1)
		other.UserID = 999
		snippets.On("FindByID", mock.Anything, uint(1)).Return(other, nil)
		uc := usecase.NewDeleteCodeSnippetUseCase(snippets)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
		snippets.AssertNotCalled(t, "Delete")
	})
}

func TestCreateSnippetCommentUseCase_Execute(t *testing.T) {
	t.Run("スニペットがあればコメントを作成する", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("CreateComment", mock.Anything, mock.MatchedBy(func(c *model.SnippetComment) bool {
			return c.Content == "コメント"
		})).Return(nil)
		uc := usecase.NewCreateSnippetCommentUseCase(snippets)

		err := uc.Execute(context.Background(), &model.SnippetComment{SnippetID: 1, Content: "  コメント  "})

		assert.NoError(t, err)
		snippets.AssertExpectations(t)
	})

	t.Run("スニペットが不在なら NotFound（作成しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewCreateSnippetCommentUseCase(snippets)

		err := uc.Execute(context.Background(), &model.SnippetComment{SnippetID: 1, Content: "x"})

		assertDomainCode(t, err, domain.ErrCodeNotFound)
		snippets.AssertNotCalled(t, "CreateComment")
	})

	t.Run("内容が空なら作成しない", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		uc := usecase.NewCreateSnippetCommentUseCase(snippets)

		err := uc.Execute(context.Background(), &model.SnippetComment{SnippetID: 1, Content: ""})

		assert.Error(t, err)
		snippets.AssertNotCalled(t, "FindByID")
		snippets.AssertNotCalled(t, "CreateComment")
	})
}

func TestForkCodeSnippetUseCase_Execute(t *testing.T) {
	t.Run("自分の投稿へフォークし、元のカウンターを加算する", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		posts.On("FindByID", mock.Anything, uint(9)).Return(ownedPostOf(9), nil)
		snippets.On("Create", mock.Anything, mock.MatchedBy(func(s *model.CodeSnippet) bool {
			return s.PostID == 9 && s.UserID == 1 && s.ForkedFromID != nil && *s.ForkedFromID == 1
		})).Return(nil)
		snippets.On("IncrementForkCount", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewForkCodeSnippetUseCase(snippets, posts)

		got, err := uc.Execute(context.Background(), 1, 1, 9)

		assert.NoError(t, err)
		assert.Equal(t, uint(9), got.PostID)
		snippets.AssertExpectations(t)
	})

	// カウンター加算の失敗はフォーク自体を失敗させない（移行前の挙動）。
	t.Run("カウンター加算に失敗してもフォークは成功する", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		posts.On("FindByID", mock.Anything, uint(9)).Return(ownedPostOf(9), nil)
		snippets.On("Create", mock.Anything, mock.Anything).Return(nil)
		snippets.On("IncrementForkCount", mock.Anything, uint(1)).Return(errors.New("db error"))
		uc := usecase.NewForkCodeSnippetUseCase(snippets, posts)

		_, err := uc.Execute(context.Background(), 1, 1, 9)

		assert.NoError(t, err)
		snippets.AssertExpectations(t)
	})

	// 403 は専用メッセージを返す（共通 helper の汎用文言ではない）。
	t.Run("他ユーザーの投稿へは Forbidden で専用メッセージ", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		other := ownedPostOf(9)
		other.UserID = 999
		posts.On("FindByID", mock.Anything, uint(9)).Return(other, nil)
		uc := usecase.NewForkCodeSnippetUseCase(snippets, posts)

		_, err := uc.Execute(context.Background(), 1, 1, 9)

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		var de *domain.DomainError
		if assert.ErrorAs(t, err, &de) {
			assert.Equal(t, "自分の投稿にのみフォークできます。投稿の編集権限がありません", de.Message)
		}
		snippets.AssertNotCalled(t, "Create")
	})

	t.Run("フォーク元が不在なら NotFound", func(t *testing.T) {
		snippets, posts := new(mockCodeSnippetRepo), new(mockPostReader)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewForkCodeSnippetUseCase(snippets, posts)

		_, err := uc.Execute(context.Background(), 1, 1, 9)

		assertDomainCode(t, err, domain.ErrCodeNotFound)
		posts.AssertNotCalled(t, "FindByID")
	})
}

func TestFavoriteCodeSnippetUseCase_Execute(t *testing.T) {
	t.Run("未登録なら追加する", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("HasFavorited", mock.Anything, uint(2), uint(1)).Return(false, nil)
		snippets.On("Favorite", mock.Anything, uint(2), uint(1)).Return(nil)
		uc := usecase.NewFavoriteCodeSnippetUseCase(snippets)

		assert.NoError(t, uc.Execute(context.Background(), 2, 1))
		snippets.AssertExpectations(t)
	})

	t.Run("登録済みなら Conflict（追加しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(ownedSnippetOf(1), nil)
		snippets.On("HasFavorited", mock.Anything, uint(2), uint(1)).Return(true, nil)
		uc := usecase.NewFavoriteCodeSnippetUseCase(snippets)

		assertDomainCode(t, uc.Execute(context.Background(), 2, 1), domain.ErrCodeConflict)
		snippets.AssertNotCalled(t, "Favorite")
	})

	t.Run("スニペットが不在なら NotFound", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewFavoriteCodeSnippetUseCase(snippets)

		assertDomainCode(t, uc.Execute(context.Background(), 2, 1), domain.ErrCodeNotFound)
		snippets.AssertNotCalled(t, "HasFavorited")
	})
}

func TestCountCodeSnippetsUseCase_Execute(t *testing.T) {
	snippets := new(mockCodeSnippetRepo)
	snippets.On("CountByUserID", mock.Anything, uint(1)).Return(int64(7), nil)
	uc := usecase.NewCountCodeSnippetsUseCase(snippets)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(7), got)
	snippets.AssertExpectations(t)
}

func TestDeleteSnippetCommentUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("所有者はコメントを削除できる", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindCommentByID", mock.Anything, uint(3)).
			Return(&model.SnippetComment{ID: 3, UserID: 1}, nil)
		snippets.On("DeleteComment", mock.Anything, uint(3)).Return(nil)
		uc := usecase.NewDeleteSnippetCommentUseCase(snippets)

		require.NoError(t, uc.Execute(ctx, 3, 1))
		snippets.AssertExpectations(t)
	})

	t.Run("他ユーザーのコメントは Forbidden（削除しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindCommentByID", mock.Anything, uint(3)).
			Return(&model.SnippetComment{ID: 3, UserID: 999}, nil)
		uc := usecase.NewDeleteSnippetCommentUseCase(snippets)

		err := uc.Execute(ctx, 3, 1)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		snippets.AssertNotCalled(t, "DeleteComment", mock.Anything, mock.Anything)
	})

	t.Run("不在は 404 を返し削除しない", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		snippets.On("FindCommentByID", mock.Anything, uint(3)).Return((*model.SnippetComment)(nil), nil)
		uc := usecase.NewDeleteSnippetCommentUseCase(snippets)

		err := uc.Execute(ctx, 3, 1)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		snippets.AssertNotCalled(t, "DeleteComment", mock.Anything, mock.Anything)
	})

	t.Run("取得の DB 障害は伝播する（削除しない）", func(t *testing.T) {
		snippets := new(mockCodeSnippetRepo)
		dbErr := errors.New("db down")
		snippets.On("FindCommentByID", mock.Anything, uint(3)).Return((*model.SnippetComment)(nil), dbErr)
		uc := usecase.NewDeleteSnippetCommentUseCase(snippets)

		assert.ErrorIs(t, uc.Execute(ctx, 3, 1), dbErr)
		snippets.AssertNotCalled(t, "DeleteComment", mock.Anything, mock.Anything)
	})
}
