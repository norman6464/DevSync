package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestQuestionService はQuestionServiceのテスト用インスタンスを生成するヘルパー。
func newTestQuestionService() (*QuestionService, *MockQuestionRepository) {
	repo := new(MockQuestionRepository)
	svc := NewQuestionService(repo)
	return svc, repo
}

// ============================================================
// 質問更新テスト
// ============================================================

func TestQuestionUpdate_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	existing := &model.Question{Title: "Old", Body: "Old Body", Tags: "go", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, "New Title", "New Body", "go,rust")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Body", result.Body)
	assert.Equal(t, "go,rust", result.Tags)
	repo.AssertExpectations(t)
}

func TestQuestionUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestQuestionService()

	existing := &model.Question{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 999, "Title", "Body", "")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestQuestionUpdate_NotFound(t *testing.T) {
	svc, repo := newTestQuestionService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.Update(999, 1, "Title", "Body", "")
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestQuestionUpdate_PartialUpdate(t *testing.T) {
	svc, repo := newTestQuestionService()

	existing := &model.Question{Title: "Old Title", Body: "Old Body", Tags: "go", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// タイトルのみ更新
	result, err := svc.Update(1, 1, "New Title", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "Old Body", result.Body)
	assert.Equal(t, "go", result.Tags)
	repo.AssertExpectations(t)
}

// ============================================================
// 質問削除テスト
// ============================================================

func TestQuestionDelete_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	existing := &model.Question{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestQuestionDelete_Forbidden(t *testing.T) {
	svc, repo := newTestQuestionService()

	existing := &model.Question{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// 質問投票テスト
// ============================================================

func TestQuestionVote_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	repo.On("Vote", uint(1), uint(10), 1).Return(nil)

	err := svc.Vote(1, 10, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestQuestionRemoveVote_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	repo.On("RemoveVote", uint(1), uint(10)).Return(nil)

	err := svc.RemoveVote(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
