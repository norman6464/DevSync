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

// ============================================================
// 回答投票テスト
// ============================================================

func TestAnswerVote_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("Vote", uint(1), uint(5), 1).Return(nil)

	err := svc.Vote(1, 5, 1)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}

func TestAnswerRemoveVote_Success(t *testing.T) {
	svc, answerRepo, _ := newTestAnswerService()

	answerRepo.On("RemoveVote", uint(1), uint(5)).Return(nil)

	err := svc.RemoveVote(1, 5)
	assert.NoError(t, err)
	answerRepo.AssertExpectations(t)
}
