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

// ============================================================
// 質問作成テスト
// ============================================================

func TestQuestionCreate_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	question := &model.Question{
		Title:  "Goのエラーハンドリングについて",
		Body:   "Goでエラーハンドリングのベストプラクティスを教えてください。",
		Tags:   "go,error-handling",
		UserID: 1,
	}

	repo.On("Create", question).Return(nil)

	err := svc.Create(question)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestQuestionCreate_ValidationError(t *testing.T) {
	svc, _ := newTestQuestionService()

	// タイトルが空
	question := &model.Question{
		Title:  "",
		Body:   "本文",
		Tags:   "go",
		UserID: 1,
	}

	err := svc.Create(question)
	assert.Error(t, err)
}

// ============================================================
// 質問取得テスト
// ============================================================

func TestQuestionGetByID_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	expected := &model.Question{Title: "テスト質問", Body: "テスト本文", UserID: 1}
	expected.ID = 1

	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "テスト質問", result.Title)
	repo.AssertExpectations(t)
}

func TestQuestionGetByID_NotFound(t *testing.T) {
	svc, repo := newTestQuestionService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestQuestionGetAll_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	questions := []model.Question{
		{Title: "質問1", Tags: "go"},
		{Title: "質問2", Tags: "go,rust"},
	}

	repo.On("FindAll", 10, 0, "go", "newest").Return(questions, int64(2), nil)

	result, total, err := svc.GetAll(10, 0, "go", "newest")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestQuestionSearch_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	questions := []model.Question{
		{Title: "Goのエラーハンドリング"},
	}

	repo.On("Search", "エラー", 10, 0).Return(questions, int64(1), nil)

	result, total, err := svc.Search("エラー", 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, "Goのエラーハンドリング", result[0].Title)
	repo.AssertExpectations(t)
}

func TestQuestionGetByUserID_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	questions := []model.Question{
		{Title: "ユーザー1の質問", UserID: 1},
	}

	repo.On("FindByUserID", uint(1)).Return(questions, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, uint(1), result[0].UserID)
	repo.AssertExpectations(t)
}

func TestQuestionGetUserVote_Success(t *testing.T) {
	svc, repo := newTestQuestionService()

	repo.On("GetUserVote", uint(1), uint(10)).Return(1, nil)

	vote, err := svc.GetUserVote(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, vote)
	repo.AssertExpectations(t)
}

// ============================================================
// 投票バリデーションエラーテスト
// ============================================================

func TestQuestionVote_InvalidValue(t *testing.T) {
	svc, _ := newTestQuestionService()

	// 0は無効な投票値
	err := svc.Vote(1, 10, 0)
	assert.Error(t, err)
}

func TestQuestionVote_InvalidValueTwo(t *testing.T) {
	svc, _ := newTestQuestionService()

	// 2は無効な投票値
	err := svc.Vote(1, 10, 2)
	assert.Error(t, err)
}
