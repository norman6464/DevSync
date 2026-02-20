package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestAnswerService はAnswerServiceのテスト用インスタンスを生成するヘルパー。
func newTestAnswerService() (*AnswerService, *MockAnswerRepository, *MockQuestionRepository) {
	answerRepo := new(MockAnswerRepository)
	questionRepo := new(MockQuestionRepository)
	svc := NewAnswerService(answerRepo, questionRepo)
	return svc, answerRepo, questionRepo
}

// ============================================================
// 回答作成テスト
// ============================================================

func TestAnswerCreate_Success(t *testing.T) {
	svc, answerRepo, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	answer := &model.Answer{QuestionID: 10, UserID: 2, Body: "Answer body"}
	answerRepo.On("Create", answer).Return(nil)

	err := svc.Create(answer)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

func TestAnswerCreate_QuestionNotFound(t *testing.T) {
	svc, _, questionRepo := newTestAnswerService()

	questionRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	answer := &model.Answer{QuestionID: 999, UserID: 2, Body: "Answer body"}
	err := svc.Create(answer)
	assert.ErrorIs(t, err, ErrNotFound)
	questionRepo.AssertExpectations(t)
}

func TestAnswerCreate_WhitespaceOnlyBody(t *testing.T) {
	svc, _, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	// 空白のみの回答 → エラーになるべき
	answer := &model.Answer{QuestionID: 10, UserID: 2, Body: "   "}
	err := svc.Create(answer)
	assert.Error(t, err)
}

func TestAnswerCreate_EmptyBody(t *testing.T) {
	svc, _, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	// 空の回答 → エラーになるべき
	answer := &model.Answer{QuestionID: 10, UserID: 2, Body: ""}
	err := svc.Create(answer)
	assert.Error(t, err)
}

// ============================================================
// 回答更新テスト
// ============================================================

func TestAnswerUpdate_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{Body: "Old Body", UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)
	answerRepo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "New Body")
	assert.NoError(t, err)
	assert.Equal(t, "New Body", result.Body)
	answerRepo.AssertExpectations(t)
}

func TestAnswerUpdate_Forbidden(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 999, "New Body")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	answerRepo.AssertExpectations(t)
}

func TestAnswerUpdate_NotFound(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.Update(999, 1, "New Body")
	assert.Error(t, err)
	assert.Nil(t, result)
	answerRepo.AssertExpectations(t)
}

func TestAnswerUpdate_RepoError(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{Body: "Old Body", UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)
	answerRepo.On("Update", existing).Return(errors.New("db error"))

	result, err := svc.Update(1, 1, "New Body")
	assert.Error(t, err)
	assert.Nil(t, result)
	answerRepo.AssertExpectations(t)
}

func TestAnswerUpdate_EmptyBody(t *testing.T) {
	svc, _, _ := newTestAnswerService()

	result, err := svc.Update(1, 1, "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "回答内容は必須です")
}

func TestAnswerUpdate_WhitespaceBody(t *testing.T) {
	svc, _, _ := newTestAnswerService()

	result, err := svc.Update(1, 1, "   ")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "回答内容は必須です")
}

// ============================================================
// 回答削除テスト
// ============================================================

func TestAnswerDelete_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)
	answerRepo.On("Delete", existing).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerDelete_Forbidden(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	answerRepo.AssertExpectations(t)
}

func TestAnswerDelete_NotFound(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Delete(999, 1)
	assert.Error(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerDelete_RepoError(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	existing := &model.Answer{UserID: 1}
	existing.ID = 1

	answerRepo.On("FindByID", uint(1)).Return(existing, nil)
	answerRepo.On("Delete", existing).Return(errors.New("db error"))

	err := svc.Delete(1, 1)
	assert.Error(t, err)
	answerRepo.AssertExpectations(t)
}

// ============================================================
// ベストアンサー設定テスト
// ============================================================

func TestSetBestAnswer_Success(t *testing.T) {
	svc, answerRepo, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	answer := &model.Answer{QuestionID: 10}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)
	answerRepo.On("SetBestAnswer", uint(10), uint(5)).Return(nil)

	err := svc.SetBestAnswer(10, 5, 1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

func TestSetBestAnswer_NotQuestionOwner(t *testing.T) {
	svc, _, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	// userID=999 は質問の所有者ではない
	err := svc.SetBestAnswer(10, 5, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	questionRepo.AssertExpectations(t)
}

func TestSetBestAnswer_AnswerBelongsToDifferentQuestion(t *testing.T) {
	svc, answerRepo, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	// 回答は別の質問に属している
	answer := &model.Answer{QuestionID: 20}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)

	err := svc.SetBestAnswer(10, 5, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	answerRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

func TestSetBestAnswer_QuestionNotFound(t *testing.T) {
	svc, _, questionRepo := newTestAnswerService()

	questionRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.SetBestAnswer(999, 5, 1)
	assert.ErrorIs(t, err, ErrNotFound)
	questionRepo.AssertExpectations(t)
}

func TestSetBestAnswer_AnswerNotFound(t *testing.T) {
	svc, answerRepo, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	answerRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.SetBestAnswer(10, 999, 1)
	assert.ErrorIs(t, err, ErrNotFound)
	answerRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

func TestSetBestAnswer_RepoError(t *testing.T) {
	svc, answerRepo, questionRepo := newTestAnswerService()

	question := &model.Question{UserID: 1}
	question.ID = 10
	questionRepo.On("FindByID", uint(10)).Return(question, nil)

	answer := &model.Answer{QuestionID: 10}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)
	answerRepo.On("SetBestAnswer", uint(10), uint(5)).Return(errors.New("db error"))

	err := svc.SetBestAnswer(10, 5, 1)
	assert.Error(t, err)
	answerRepo.AssertExpectations(t)
	questionRepo.AssertExpectations(t)
}

// ============================================================
// 回答投票テスト
// ============================================================

func TestAnswerVote_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	// 回答者はuserID=99（投票者のuserID=1とは異なる）
	answer := &model.Answer{UserID: 99}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)
	answerRepo.On("Vote", uint(1), uint(5), 1).Return(nil)

	err := svc.Vote(1, 5, 1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerVote_SelfVote_Forbidden(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	// 自分の回答に投票しようとする（userID=1の回答にuserID=1が投票）
	answer := &model.Answer{UserID: 1}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)

	err := svc.Vote(1, 5, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	answerRepo.AssertNotCalled(t, "Vote")
}

func TestAnswerVote_AnswerNotFound(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByID", uint(99)).Return(nil, ErrNotFound)

	err := svc.Vote(1, 99, 1)
	assert.ErrorIs(t, err, ErrNotFound)
	answerRepo.AssertNotCalled(t, "Vote")
}

func TestAnswerRemoveVote_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answer := &model.Answer{UserID: 99}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)
	answerRepo.On("RemoveVote", uint(1), uint(5)).Return(nil)

	err := svc.RemoveVote(1, 5)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerRemoveVote_RepoError(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answer := &model.Answer{UserID: 99}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)
	answerRepo.On("RemoveVote", uint(1), uint(5)).Return(errors.New("db error"))

	err := svc.RemoveVote(1, 5)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	answerRepo.AssertExpectations(t)
}

func TestAnswerRemoveVote_SelfVote_Forbidden(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answer := &model.Answer{UserID: 1}
	answer.ID = 5
	answerRepo.On("FindByID", uint(5)).Return(answer, nil)

	err := svc.RemoveVote(1, 5)
	assert.ErrorIs(t, err, ErrForbidden)
	answerRepo.AssertNotCalled(t, "RemoveVote")
}

func TestAnswerRemoveVote_AnswerNotFound(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.RemoveVote(1, 99)
	assert.ErrorIs(t, err, ErrNotFound)
	answerRepo.AssertNotCalled(t, "RemoveVote")
}

// ============================================================
// 質問IDによる回答取得テスト
// ============================================================

func TestAnswerGetByQuestionID_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answers := []model.Answer{
		{Body: "Answer 1", QuestionID: 10, UserID: 1},
		{Body: "Answer 2", QuestionID: 10, UserID: 2},
	}
	answerRepo.On("FindByQuestionID", uint(10)).Return(answers, nil)

	result, err := svc.GetByQuestionID(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	answerRepo.AssertExpectations(t)
}

func TestAnswerGetByQuestionID_Empty(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByQuestionID", uint(99)).Return([]model.Answer{}, nil)

	result, err := svc.GetByQuestionID(99)
	assert.NoError(t, err)
	assert.Empty(t, result)
	answerRepo.AssertExpectations(t)
}

// ============================================================
// ユーザー投票取得テスト
// ============================================================

func TestAnswerGetUserVotes_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	votes := map[uint]int{1: 1, 2: -1, 3: 1}
	answerRepo.On("GetUserVotes", uint(1), []uint{1, 2, 3}).Return(votes, nil)

	result, err := svc.GetUserVotes(1, []uint{1, 2, 3})
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[1])
	assert.Equal(t, -1, result[2])
	answerRepo.AssertExpectations(t)
}

func TestAnswerGetUserVotes_Empty(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("GetUserVotes", uint(1), []uint{}).Return(map[uint]int{}, nil)

	result, err := svc.GetUserVotes(1, []uint{})
	assert.NoError(t, err)
	assert.Empty(t, result)
	answerRepo.AssertExpectations(t)
}

// ============================================================
// Vote バリデーションテスト
// ============================================================

func TestAnswerVote_InvalidValue(t *testing.T) {
	svc, _, _ := newTestAnswerService()

	// 0は無効な投票値
	err := svc.Vote(1, 1, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "投票値は1または-1")
}

func TestAnswerVote_OutOfRangeValue(t *testing.T) {
	svc, _, _ := newTestAnswerService()

	// 99は無効な投票値
	err := svc.Vote(1, 1, 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "投票値は1または-1")
}

func TestAnswerVote_ValidUpvote(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	// 回答者はuserID=99（投票者のuserID=1とは異なる）
	answer := &model.Answer{UserID: 99}
	answer.ID = 1
	answerRepo.On("FindByID", uint(1)).Return(answer, nil)
	answerRepo.On("Vote", uint(1), uint(1), 1).Return(nil)

	err := svc.Vote(1, 1, 1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerVote_ValidDownvote(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	// 回答者はuserID=99（投票者のuserID=1とは異なる）
	answer := &model.Answer{UserID: 99}
	answer.ID = 1
	answerRepo.On("FindByID", uint(1)).Return(answer, nil)
	answerRepo.On("Vote", uint(1), uint(1), -1).Return(nil)

	err := svc.Vote(1, 1, -1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

// ============================================================
// 投票範囲による回答取得テスト
// ============================================================

func TestAnswerGetByVoteRange_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	expected := []model.Answer{
		{Body: "高評価回答", QuestionID: 10, VoteCount: 5},
		{Body: "中評価回答", QuestionID: 10, VoteCount: 3},
	}
	answerRepo.On("FindByVoteRange", uint(10), 3, 10).Return(expected, nil)

	result, err := svc.GetByVoteRange(10, 3, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	answerRepo.AssertExpectations(t)
}

func TestAnswerGetByVoteRange_InvalidRange(t *testing.T) {
	svc, _, _ := newTestAnswerService()

	// minVote > maxVote
	result, err := svc.GetByVoteRange(10, 10, 3)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "投票範囲が無効です")
}

func TestAnswerGetByVoteRange_NegativeRange(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	expected := []model.Answer{
		{Body: "低評価回答", QuestionID: 10, VoteCount: -2},
	}
	answerRepo.On("FindByVoteRange", uint(10), -5, 0).Return(expected, nil)

	result, err := svc.GetByVoteRange(10, -5, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	answerRepo.AssertExpectations(t)
}

func TestAnswerGetByVoteRange_RepoError(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("FindByVoteRange", uint(10), 0, 5).Return([]model.Answer(nil), errors.New("db error"))

	result, err := svc.GetByVoteRange(10, 0, 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	answerRepo.AssertExpectations(t)
}

func TestAnswerUpdate_TrimsPaddedBody(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()
	existing := &model.Answer{UserID: 1, Body: "Original"}
	existing.ID = 1
	answerRepo.On("FindByID", uint(1)).Return(existing, nil)
	answerRepo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "  Updated Body  ")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Body", result.Body)
}
