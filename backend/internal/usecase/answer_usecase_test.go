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

// mockAnswerRepo は usecase/repository.AnswerRepository のモック。
type mockAnswerRepo struct{ mock.Mock }

func (m *mockAnswerRepo) Create(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}

func (m *mockAnswerRepo) FindByID(ctx context.Context, id uint) (*model.Answer, error) {
	args := m.Called(ctx, id)
	a, _ := args.Get(0).(*model.Answer)
	return a, args.Error(1)
}

func (m *mockAnswerRepo) Update(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}

func (m *mockAnswerRepo) Delete(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}

func (m *mockAnswerRepo) FindByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error) {
	args := m.Called(ctx, questionID)
	a, _ := args.Get(0).([]model.Answer)
	return a, args.Error(1)
}

func (m *mockAnswerRepo) FindByVoteRange(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	args := m.Called(ctx, questionID, minVote, maxVote)
	a, _ := args.Get(0).([]model.Answer)
	return a, args.Error(1)
}

func (m *mockAnswerRepo) SetBestAnswer(ctx context.Context, questionID, answerID uint) error {
	return m.Called(ctx, questionID, answerID).Error(0)
}

func (m *mockAnswerRepo) Vote(ctx context.Context, userID, answerID uint, value int) error {
	return m.Called(ctx, userID, answerID, value).Error(0)
}

func (m *mockAnswerRepo) RemoveVote(ctx context.Context, userID, answerID uint) error {
	return m.Called(ctx, userID, answerID).Error(0)
}

// mockQuestionReader は usecase/repository.QuestionReader のモック。
type mockQuestionReader struct{ mock.Mock }

func (m *mockQuestionReader) FindByID(ctx context.Context, id uint) (*model.Question, error) {
	args := m.Called(ctx, id)
	q, _ := args.Get(0).(*model.Question)
	return q, args.Error(1)
}

// assertAnswerCode は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertAnswerCode(t *testing.T, err error, want domain.ErrorCode) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.Code)
}

func TestCreateAnswerUseCase_Execute(t *testing.T) {
	t.Run("質問が存在すれば作成し、本文の前後空白を除く", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(&model.Question{ID: 3}, nil)
		answers.On("Create", mock.Anything, mock.AnythingOfType("*model.Answer")).Return(nil)
		uc := usecase.NewCreateAnswerUseCase(answers, questions)

		answer := &model.Answer{QuestionID: 3, UserID: 1, Body: "  回答です  "}
		err := uc.Execute(context.Background(), answer)

		assert.NoError(t, err)
		assert.Equal(t, "回答です", answer.Body)
		answers.AssertExpectations(t)
	})

	t.Run("本文が空なら 400 で質問も引かない", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		uc := usecase.NewCreateAnswerUseCase(answers, questions)

		err := uc.Execute(context.Background(), &model.Answer{QuestionID: 3, Body: "   "})

		assertAnswerCode(t, err, domain.ErrCodeValidation)
		questions.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
		answers.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("本文が 10001 文字なら 400", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		uc := usecase.NewCreateAnswerUseCase(answers, questions)

		err := uc.Execute(context.Background(), &model.Answer{
			QuestionID: 3, Body: strings.Repeat("あ", 10001),
		})

		assertAnswerCode(t, err, domain.ErrCodeValidation)
	})

	t.Run("質問が不在なら 404 で作成しない", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(nil, nil)
		uc := usecase.NewCreateAnswerUseCase(answers, questions)

		err := uc.Execute(context.Background(), &model.Answer{QuestionID: 3, Body: "回答"})

		assertAnswerCode(t, err, domain.ErrCodeNotFound)
		answers.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("質問取得の DB 障害も 404 に潰れる", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(nil, errors.New("db error"))
		uc := usecase.NewCreateAnswerUseCase(answers, questions)

		err := uc.Execute(context.Background(), &model.Answer{QuestionID: 3, Body: "回答"})

		assertAnswerCode(t, err, domain.ErrCodeNotFound)
	})
}

func TestUpdateAnswerUseCase_Execute(t *testing.T) {
	t.Run("所有者なら本文を更新する", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).
			Return(&model.Answer{ID: 10, UserID: 1, Body: "旧本文"}, nil)
		answers.On("Update", mock.Anything, mock.AnythingOfType("*model.Answer")).Return(nil)
		uc := usecase.NewUpdateAnswerUseCase(answers)

		got, err := uc.Execute(context.Background(), 10, 1, "  新本文  ")

		require.NoError(t, err)
		assert.Equal(t, "新本文", got.Body)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).
			Return(&model.Answer{ID: 10, UserID: 99}, nil)
		uc := usecase.NewUpdateAnswerUseCase(answers)

		_, err := uc.Execute(context.Background(), 10, 1, "新本文")

		assertAnswerCode(t, err, domain.ErrCodeForbidden)
		answers.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewUpdateAnswerUseCase(answers)

		_, err := uc.Execute(context.Background(), 10, 1, "新本文")

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("空本文は 400 で保存しない", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).
			Return(&model.Answer{ID: 10, UserID: 1, Body: "旧本文"}, nil)
		uc := usecase.NewUpdateAnswerUseCase(answers)

		_, err := uc.Execute(context.Background(), 10, 1, "   ")

		assertAnswerCode(t, err, domain.ErrCodeValidation)
		answers.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestDeleteAnswerUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		existing := &model.Answer{ID: 10, UserID: 1, QuestionID: 3}
		answers.On("FindByID", mock.Anything, uint(10)).Return(existing, nil)
		answers.On("Delete", mock.Anything, existing).Return(nil)
		uc := usecase.NewDeleteAnswerUseCase(answers)

		assert.NoError(t, uc.Execute(context.Background(), 10, 1))
		answers.AssertExpectations(t)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, UserID: 99}, nil)
		uc := usecase.NewDeleteAnswerUseCase(answers)

		assertAnswerCode(t, uc.Execute(context.Background(), 10, 1), domain.ErrCodeForbidden)
		answers.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

func TestSetBestAnswerUseCase_Execute(t *testing.T) {
	t.Run("質問の投稿者なら設定できる", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(&model.Question{ID: 3, UserID: 1}, nil)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, QuestionID: 3}, nil)
		answers.On("SetBestAnswer", mock.Anything, uint(3), uint(10)).Return(nil)
		uc := usecase.NewSetBestAnswerUseCase(answers, questions)

		assert.NoError(t, uc.Execute(context.Background(), 3, 10, 1))
		answers.AssertExpectations(t)
	})

	t.Run("質問が不在なら 404", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(nil, nil)
		uc := usecase.NewSetBestAnswerUseCase(answers, questions)

		assertAnswerCode(t, uc.Execute(context.Background(), 3, 10, 1), domain.ErrCodeNotFound)
		answers.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("質問の投稿者でなければ 403", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(&model.Question{ID: 3, UserID: 99}, nil)
		uc := usecase.NewSetBestAnswerUseCase(answers, questions)

		assertAnswerCode(t, uc.Execute(context.Background(), 3, 10, 1), domain.ErrCodeForbidden)
	})

	t.Run("回答が不在なら 404", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(&model.Question{ID: 3, UserID: 1}, nil)
		answers.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewSetBestAnswerUseCase(answers, questions)

		assertAnswerCode(t, uc.Execute(context.Background(), 3, 10, 1), domain.ErrCodeNotFound)
	})

	t.Run("別の質問の回答なら 400", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		questions := new(mockQuestionReader)
		questions.On("FindByID", mock.Anything, uint(3)).Return(&model.Question{ID: 3, UserID: 1}, nil)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, QuestionID: 77}, nil)
		uc := usecase.NewSetBestAnswerUseCase(answers, questions)

		assertAnswerCode(t, uc.Execute(context.Background(), 3, 10, 1), domain.ErrCodeBadRequest)
		answers.AssertNotCalled(t, "SetBestAnswer", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestVoteAnswerUseCase_Execute(t *testing.T) {
	t.Run("他人の回答には投票できる", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, UserID: 99}, nil)
		answers.On("Vote", mock.Anything, uint(1), uint(10), -1).Return(nil)
		uc := usecase.NewVoteAnswerUseCase(answers)

		assert.NoError(t, uc.Execute(context.Background(), 1, 10, -1))
		answers.AssertExpectations(t)
	})

	t.Run("投票値が 2 なら 400 で回答も引かない", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		uc := usecase.NewVoteAnswerUseCase(answers)

		assertAnswerCode(t, uc.Execute(context.Background(), 1, 10, 2), domain.ErrCodeValidation)
		answers.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("自分の回答には投票できない（403）", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, UserID: 1}, nil)
		uc := usecase.NewVoteAnswerUseCase(answers)

		assertAnswerCode(t, uc.Execute(context.Background(), 1, 10, 1), domain.ErrCodeForbidden)
		answers.AssertNotCalled(t, "Vote", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("不在の回答への投票は 404", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewVoteAnswerUseCase(answers)

		assertAnswerCode(t, uc.Execute(context.Background(), 1, 10, 1), domain.ErrCodeNotFound)
	})
}

func TestRemoveAnswerVoteUseCase_Execute(t *testing.T) {
	t.Run("他人の回答への投票は取り消せる", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, UserID: 99}, nil)
		answers.On("RemoveVote", mock.Anything, uint(1), uint(10)).Return(nil)
		uc := usecase.NewRemoveAnswerVoteUseCase(answers)

		assert.NoError(t, uc.Execute(context.Background(), 1, 10))
		answers.AssertExpectations(t)
	})

	t.Run("自分の回答なら 403", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByID", mock.Anything, uint(10)).Return(&model.Answer{ID: 10, UserID: 1}, nil)
		uc := usecase.NewRemoveAnswerVoteUseCase(answers)

		assertAnswerCode(t, uc.Execute(context.Background(), 1, 10), domain.ErrCodeForbidden)
	})
}

func TestListAnswersByVoteRangeUseCase_Execute(t *testing.T) {
	t.Run("下限が上限を上回るなら 400 で repo を引かない", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		uc := usecase.NewListAnswersByVoteRangeUseCase(answers)

		_, err := uc.Execute(context.Background(), 3, 10, 5)

		assertAnswerCode(t, err, domain.ErrCodeBadRequest)
		assert.Equal(t, "投票範囲が無効です", domain.GetDomainError(err).Message)
		answers.AssertNotCalled(t, "FindByVoteRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("同値の範囲は許容する", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByVoteRange", mock.Anything, uint(3), 5, 5).
			Return([]model.Answer{{ID: 1}}, nil)
		uc := usecase.NewListAnswersByVoteRangeUseCase(answers)

		got, err := uc.Execute(context.Background(), 3, 5, 5)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})
}

func TestListAnswersUseCase_Execute(t *testing.T) {
	t.Run("repo に委譲する", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByQuestionID", mock.Anything, uint(3)).
			Return([]model.Answer{{ID: 1}, {ID: 2}}, nil)
		uc := usecase.NewListAnswersUseCase(answers)

		got, err := uc.Execute(context.Background(), 3)

		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		answers := new(mockAnswerRepo)
		answers.On("FindByQuestionID", mock.Anything, uint(3)).
			Return([]model.Answer(nil), errors.New("db error"))
		uc := usecase.NewListAnswersUseCase(answers)

		_, err := uc.Execute(context.Background(), 3)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}
