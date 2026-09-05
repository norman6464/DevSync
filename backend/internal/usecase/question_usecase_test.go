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

// mockQuestionRepo は usecase/repository.QuestionRepository のモック。
type mockQuestionRepo struct{ mock.Mock }

func (m *mockQuestionRepo) Create(ctx context.Context, q *model.Question) error {
	return m.Called(ctx, q).Error(0)
}

func (m *mockQuestionRepo) FindByID(ctx context.Context, id uint) (*model.Question, error) {
	args := m.Called(ctx, id)
	q, _ := args.Get(0).(*model.Question)
	return q, args.Error(1)
}

func (m *mockQuestionRepo) Update(ctx context.Context, q *model.Question) error {
	return m.Called(ctx, q).Error(0)
}

func (m *mockQuestionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockQuestionRepo) FindAll(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset, tag, sort)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) FindSolved(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) FindUnanswered(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

func (m *mockQuestionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockQuestionRepo) Vote(ctx context.Context, userID, questionID uint, value int) error {
	return m.Called(ctx, userID, questionID, value).Error(0)
}

func (m *mockQuestionRepo) RemoveVote(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}

func (m *mockQuestionRepo) GetUserVote(ctx context.Context, userID, questionID uint) (int, error) {
	args := m.Called(ctx, userID, questionID)
	return args.Int(0), args.Error(1)
}

func (m *mockQuestionRepo) Bookmark(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}

func (m *mockQuestionRepo) Unbookmark(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}

func (m *mockQuestionRepo) HasBookmarked(ctx context.Context, userID, questionID uint) (bool, error) {
	args := m.Called(ctx, userID, questionID)
	return args.Bool(0), args.Error(1)
}

func (m *mockQuestionRepo) FindBookmarkedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	q, _ := args.Get(0).([]model.Question)
	return q, args.Get(1).(int64), args.Error(2)
}

// assertQuestionCode は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertQuestionCode(t *testing.T, err error, want domain.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.Code)
}

func TestCreateQuestionUseCase_Execute(t *testing.T) {
	t.Run("検証を通れば作成する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Question")).Return(nil)
		uc := usecase.NewCreateQuestionUseCase(repo)

		err := uc.Execute(context.Background(), &model.Question{Title: "Go の質問", Body: "本文"})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("タイトルが空なら 400 で作成しない", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		uc := usecase.NewCreateQuestionUseCase(repo)

		err := uc.Execute(context.Background(), &model.Question{Title: "  ", Body: "本文"})

		assertQuestionCode(t, err, domain.ErrCodeValidation)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("タイトルが 501 文字なら 400", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		uc := usecase.NewCreateQuestionUseCase(repo)

		err := uc.Execute(context.Background(), &model.Question{
			Title: strings.Repeat("あ", 501), Body: "本文",
		})

		assertQuestionCode(t, err, domain.ErrCodeValidation)
	})

	t.Run("タグが 301 文字なら 400", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		uc := usecase.NewCreateQuestionUseCase(repo)

		err := uc.Execute(context.Background(), &model.Question{
			Title: "題", Body: "本文", Tags: strings.Repeat("a", 301),
		})

		assertQuestionCode(t, err, domain.ErrCodeValidation)
	})
}

func TestGetQuestionUseCase_Execute(t *testing.T) {
	t.Run("見つかれば返す", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		expected := &model.Question{ID: 5, Title: "Q"}
		repo.On("FindByID", mock.Anything, uint(5)).Return(expected, nil)
		uc := usecase.NewGetQuestionUseCase(repo)

		got, err := uc.Execute(context.Background(), 5)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewGetQuestionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetQuestionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.ErrorContains(t, err, "db error")
	})
}

func TestUpdateQuestionUseCase_Execute(t *testing.T) {
	newQuestion := func() *model.Question {
		return &model.Question{ID: 5, UserID: 1, Title: "旧題", Body: "旧本文", Tags: `["go"]`}
	}

	t.Run("指定した項目だけ更新する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(newQuestion(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Question")).Return(nil)
		uc := usecase.NewUpdateQuestionUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, "  新題  ", "", "")

		require.NoError(t, err)
		assert.Equal(t, "新題", got.Title, "前後の空白は除去される")
		assert.Equal(t, "旧本文", got.Body, "空文字列は変更なし")
		assert.Equal(t, `["go"]`, got.Tags)
	})

	t.Run("空白のみの入力は変更なしとして扱う", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(newQuestion(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Question")).Return(nil)
		uc := usecase.NewUpdateQuestionUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1, "   ", "   ", "   ")

		require.NoError(t, err)
		assert.Equal(t, "旧題", got.Title)
		assert.Equal(t, "旧本文", got.Body)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 99}, nil)
		uc := usecase.NewUpdateQuestionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, "新題", "", "")

		assertQuestionCode(t, err, domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewUpdateQuestionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, "新題", "", "")

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("タイトルが 501 文字なら 400", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(newQuestion(), nil)
		uc := usecase.NewUpdateQuestionUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1, strings.Repeat("あ", 501), "", "")

		assertQuestionCode(t, err, domain.ErrCodeValidation)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestDeleteQuestionUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 1}, nil)
		repo.On("Delete", mock.Anything, uint(5)).Return(nil)
		uc := usecase.NewDeleteQuestionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 5, 1))
		repo.AssertExpectations(t)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 99}, nil)
		uc := usecase.NewDeleteQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 5, 1), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

func TestVoteQuestionUseCase_Execute(t *testing.T) {
	t.Run("他人の質問には投票できる", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 99}, nil)
		repo.On("Vote", mock.Anything, uint(1), uint(5), 1).Return(nil)
		uc := usecase.NewVoteQuestionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5, 1))
		repo.AssertExpectations(t)
	})

	t.Run("投票値が 0 なら 400 で質問も引かない", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		uc := usecase.NewVoteQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5, 0), domain.ErrCodeValidation)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("自分の質問には投票できない（403）", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 1}, nil)
		uc := usecase.NewVoteQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5, 1), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Vote", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("不在の質問への投票は 404", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewVoteQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5, 1), domain.ErrCodeNotFound)
	})

	t.Run("取得の DB 障害も 404 に潰れる", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))
		uc := usecase.NewVoteQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5, 1), domain.ErrCodeNotFound)
	})
}

func TestRemoveQuestionVoteUseCase_Execute(t *testing.T) {
	t.Run("他人の質問への投票は取り消せる", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 99}, nil)
		repo.On("RemoveVote", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewRemoveQuestionVoteUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("自分の質問なら 403", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 1}, nil)
		uc := usecase.NewRemoveQuestionVoteUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5), domain.ErrCodeForbidden)
	})
}

func TestBookmarkQuestionUseCase_Execute(t *testing.T) {
	t.Run("未ブックマークなら追加する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 99}, nil)
		repo.On("HasBookmarked", mock.Anything, uint(1), uint(5)).Return(false, nil)
		repo.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewBookmarkQuestionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("ブックマーク済みなら 409", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5}, nil)
		repo.On("HasBookmarked", mock.Anything, uint(1), uint(5)).Return(true, nil)
		uc := usecase.NewBookmarkQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5), domain.ErrCodeConflict)
		repo.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("不在なら 404", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)
		uc := usecase.NewBookmarkQuestionUseCase(repo)

		assertQuestionCode(t, uc.Execute(context.Background(), 1, 5), domain.ErrCodeNotFound)
	})

	t.Run("自分の質問でもブックマークできる", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{ID: 5, UserID: 1}, nil)
		repo.On("HasBookmarked", mock.Anything, uint(1), uint(5)).Return(false, nil)
		repo.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewBookmarkQuestionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
	})
}

func TestUnbookmarkQuestionUseCase_Execute(t *testing.T) {
	t.Run("存在確認をせずに解除する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("Unbookmark", mock.Anything, uint(1), uint(5)).Return(nil)
		uc := usecase.NewUnbookmarkQuestionUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})
}

func TestQuestionPassThroughUseCases(t *testing.T) {
	ctx := context.Background()

	t.Run("一覧はタグとソートをそのまま渡す", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindAll", mock.Anything, 20, 0, "go", "votes").
			Return([]model.Question{{Title: "Q"}}, int64(1), nil)
		uc := usecase.NewListQuestionsUseCase(repo)

		got, total, err := uc.Execute(ctx, 20, 0, "go", "votes")

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("検索は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("Search", mock.Anything, "go", 20, 0).Return([]model.Question{}, int64(0), nil)
		uc := usecase.NewSearchQuestionsUseCase(repo)

		_, total, err := uc.Execute(ctx, "go", 20, 0)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("ユーザー別一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindByUserID", mock.Anything, uint(7), 20, 0).
			Return([]model.Question{{Title: "Q"}}, int64(1), nil)
		uc := usecase.NewListQuestionsByUserUseCase(repo)

		got, _, err := uc.Execute(ctx, 7, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("解決済み一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindSolved", mock.Anything, 20, 0).Return([]model.Question{}, int64(0), nil)
		uc := usecase.NewListSolvedQuestionsUseCase(repo)

		_, _, err := uc.Execute(ctx, 20, 0)

		assert.NoError(t, err)
	})

	t.Run("未回答一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindUnanswered", mock.Anything, 20, 0).Return([]model.Question{}, int64(0), nil)
		uc := usecase.NewListUnansweredQuestionsUseCase(repo)

		_, _, err := uc.Execute(ctx, 20, 0)

		assert.NoError(t, err)
	})

	t.Run("ブックマーク一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("FindBookmarkedByUserID", mock.Anything, uint(1), 20, 0).
			Return([]model.Question{{Title: "Q"}}, int64(1), nil)
		uc := usecase.NewListBookmarkedQuestionsUseCase(repo)

		got, _, err := uc.Execute(ctx, 1, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("投票値は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("GetUserVote", mock.Anything, uint(1), uint(5)).Return(-1, nil)
		uc := usecase.NewGetQuestionUserVoteUseCase(repo)

		got, err := uc.Execute(ctx, 1, 5)

		assert.NoError(t, err)
		assert.Equal(t, -1, got)
	})

	t.Run("質問数は repo に委譲する", func(t *testing.T) {
		repo := new(mockQuestionRepo)
		repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
		uc := usecase.NewCountQuestionsUseCase(repo)

		got, err := uc.Execute(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), got)
	})
}
