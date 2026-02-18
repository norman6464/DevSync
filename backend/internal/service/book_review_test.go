package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestBookReviewService はBookReviewServiceのテスト用インスタンスを生成するヘルパー。
func newTestBookReviewService() (*BookReviewService, *MockBookReviewRepository) {
	repo := new(MockBookReviewRepository)
	svc := NewBookReviewService(repo)
	return svc, repo
}

// ============================================================
// 書籍レビュー作成テスト
// ============================================================

func TestBookReviewCreate_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	review := &model.BookReview{
		UserID: 1,
		Title:  "Go言語による並行処理",
		Author: "Katherine Cox-Buday",
		Rating: 5,
		Review: "並行処理の基礎から応用まで丁寧に解説",
	}

	repo.On("Create", review).Return(nil)

	err := svc.Create(review)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookReviewCreate_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	review := &model.BookReview{UserID: 1, Title: "テスト本"}

	repo.On("Create", review).Return(errors.New("db error"))

	err := svc.Create(review)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// 書籍レビュー取得テスト
// ============================================================

func TestBookReviewGetByID_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	expected := &model.BookReview{Title: "Go入門", Author: "著者A", Rating: 4, UserID: 1}
	expected.ID = 1

	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestBookReviewGetByID_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestBookReviewGetByUserID_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	reviews := []model.BookReview{
		{Title: "本A", UserID: 1, Rating: 3},
		{Title: "本B", UserID: 1, Rating: 5},
	}

	repo.On("FindByUserID", uint(1)).Return(reviews, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestBookReviewGetAll_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	reviews := []model.BookReview{
		{Title: "本A", Rating: 4},
		{Title: "本B", Rating: 3},
	}

	repo.On("FindAll", 10, 0).Return(reviews, int64(2), nil)

	result, total, err := svc.GetAll(10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

// ============================================================
// 書籍レビュー更新テスト
// ============================================================

func TestBookReviewUpdate_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{Title: "Old", Author: "Author A", Rating: 3, UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.BookReview{Title: "New Title", Rating: 5}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, 5, result.Rating)
	assert.Equal(t, "Author A", result.Author) // 変更なし
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.BookReview{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Rating: 4}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))

	updates := &model.BookReview{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 書籍レビュー削除テスト
// ============================================================

func TestBookReviewDelete_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookReviewDelete_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestBookReviewDelete_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}
