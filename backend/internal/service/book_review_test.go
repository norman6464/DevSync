package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestBookReviewCreate_WhitespaceTitle(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{UserID: 1, Title: "   \t  ", Rating: 4}
	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトルを入力してください")
}

func TestBookReviewCreate_EmptyTitle(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{UserID: 1, Title: "", Rating: 4}
	err := svc.Create(review)
	assert.Error(t, err)
}

func TestBookReviewCreate_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	review := &model.BookReview{UserID: 1, Title: "テスト本", Rating: 3}

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

	repo.On("FindByUserID", uint(1), 20, 0).Return(reviews, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestBookReviewGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("FindByUserID", uint(1), 10, 10).Return([]model.BookReview{}, int64(15), nil)

	result, total, err := svc.GetByUserID(1, 10, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(15), total)
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

func TestBookReviewUpdate_AllFields(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Author: "Old Author"}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 全フィールドを更新
	updates := &model.BookReview{
		Title:    "New Title",
		Author:   "New Author",
		ISBN:     "978-1234567890",
		Rating:   5,
		Review:   "Great book!",
		ImageURL: "https://example.com/image.jpg",
	}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Author", result.Author)
	assert.Equal(t, "978-1234567890", result.ISBN)
	assert.Equal(t, 5, result.Rating)
	assert.Equal(t, "Great book!", result.Review)
	assert.Equal(t, "https://example.com/image.jpg", result.ImageURL)
	repo.AssertExpectations(t)
}

// ============================================================
// Rating バリデーションテスト
// ============================================================

func TestBookReviewCreate_InvalidRatingTooLow(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{
		Title:  "テスト本",
		Author: "テスト著者",
		UserID: 1,
		Rating: 0,
	}

	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "評価は1〜5")
}

func TestBookReviewCreate_InvalidRatingTooHigh(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{
		Title:  "テスト本",
		Author: "テスト著者",
		UserID: 1,
		Rating: 6,
	}

	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "評価は1〜5")
}

// ============================================================
// 評価範囲による取得テスト
// ============================================================

func TestBookReviewGetByRating_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	expected := []model.BookReview{
		{Title: "良書A", UserID: 1, Rating: 4},
		{Title: "良書B", UserID: 1, Rating: 5},
	}
	repo.On("FindByRating", uint(1), 4, 5).Return(expected, nil)

	result, err := svc.GetByRating(1, 4, 5)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestBookReviewGetByRating_InvalidRange(t *testing.T) {
	svc, _ := newTestBookReviewService()

	// minRating > maxRating
	result, err := svc.GetByRating(1, 5, 3)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "評価範囲が無効です")
}

func TestBookReviewGetByRating_OutOfRange(t *testing.T) {
	svc, _ := newTestBookReviewService()

	// minRating < 1
	result, err := svc.GetByRating(1, 0, 5)
	assert.Error(t, err)
	assert.Nil(t, result)

	// maxRating > 5
	result2, err2 := svc.GetByRating(1, 1, 6)
	assert.Error(t, err2)
	assert.Nil(t, result2)
}

func TestBookReviewGetByRating_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("FindByRating", uint(1), 1, 3).Return([]model.BookReview(nil), errors.New("db error"))

	result, err := svc.GetByRating(1, 1, 3)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_InvalidRating(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{Title: "Old", Author: "Author", UserID: 1, Rating: 3}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Rating: -1}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "評価は1〜5")
}

// ============================================================
// 空白バイパス脆弱性テスト
// ============================================================

func TestBookReviewUpdate_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestBookReviewService()
	existing := &model.BookReview{Title: "Original Title", Author: "Author", Rating: 4, Review: "Good", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.BookReview{Title: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Title", result.Title)
}

func TestBookReviewUpdate_WhitespaceAuthor(t *testing.T) {
	svc, repo := newTestBookReviewService()
	existing := &model.BookReview{Title: "Title", Author: "Original Author", Rating: 4, Review: "Good", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.BookReview{Author: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Author", result.Author)
}

func TestBookReviewUpdate_WhitespaceReview(t *testing.T) {
	svc, repo := newTestBookReviewService()
	existing := &model.BookReview{Title: "Title", Author: "Author", Rating: 4, Review: "Original Review", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.BookReview{Review: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Review", result.Review)
}

func TestBookReviewUpdate_TrimsPaddedTitle(t *testing.T) {
	svc, repo := newTestBookReviewService()
	existing := &model.BookReview{Title: "Original", Author: "Author", Rating: 4, Review: "Review", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.BookReview{Title: "  New Title  "})
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
}

// ============================================================
// 書籍レビュー アーカイブ / アーカイブ解除テスト
// ============================================================

func TestBookReviewArchive_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	err := svc.ArchiveReview(1, 1)
	assert.NoError(t, err)
	assert.True(t, existing.IsArchived)
	repo.AssertExpectations(t)
}

func TestBookReviewArchive_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.ArchiveReview(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestBookReviewArchive_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.ArchiveReview(99, 1)
	assert.Error(t, err)
}

func TestBookReviewUnarchive_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", IsArchived: true}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	err := svc.UnarchiveReview(1, 1)
	assert.NoError(t, err)
	assert.False(t, existing.IsArchived)
	repo.AssertExpectations(t)
}

func TestBookReviewUnarchive_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", IsArchived: true}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.UnarchiveReview(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestBookReviewUpdateStatus_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(r *model.BookReview) bool {
		return r.ID == 1 && r.Status == model.ReviewStatusReading
	})).Return(nil)

	err := svc.UpdateStatus(1, 1, model.ReviewStatusReading)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdateStatus_InvalidStatus(t *testing.T) {
	svc, _ := newTestBookReviewService()

	err := svc.UpdateStatus(1, 1, model.ReviewStatus("invalid"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なステータス")
}

func TestBookReviewUpdateStatus_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 99, Title: "他人のレビュー"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.UpdateStatus(1, 1, model.ReviewStatusCompleted)
	assert.ErrorIs(t, err, ErrForbidden)
}

// ============================================================
// 書籍レビュー検索テスト
// ============================================================

func TestBookReviewSearch_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	expected := []model.BookReview{
		{Title: "Go言語入門", Author: "テスト著者", Rating: 4, UserID: 1},
	}
	repo.On("Search", "Go", 20, 0).Return(expected, int64(1), nil)

	result, total, err := svc.Search("Go", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestBookReviewSearch_EmptyQuery(t *testing.T) {
	svc, _ := newTestBookReviewService()

	result, total, err := svc.Search("", 20, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "検索キーワードは必須です")
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

func TestBookReviewSearch_WhitespaceQuery(t *testing.T) {
	svc, _ := newTestBookReviewService()

	result, total, err := svc.Search("   ", 20, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "検索キーワードは必須です")
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

func TestBookReviewSearch_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("Search", "Go", 20, 0).Return([]model.BookReview(nil), int64(0), errors.New("db error"))

	result, total, err := svc.Search("Go", 20, 0)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

// ============================================================
// 読書進捗更新テスト
// ============================================================

func TestBookReviewUpdateProgress_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 300, CurrentPage: 50, Status: model.ReviewStatusReading}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(r *model.BookReview) bool {
		return r.ID == 1 && r.CurrentPage == 150 && r.Status == model.ReviewStatusReading
	})).Return(nil)

	result, err := svc.UpdateProgress(1, 1, 150)
	assert.NoError(t, err)
	assert.Equal(t, 150, result.CurrentPage)
	assert.Equal(t, model.ReviewStatusReading, result.Status)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdateProgress_AutoComplete(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 300, CurrentPage: 200, Status: model.ReviewStatusReading}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(r *model.BookReview) bool {
		return r.ID == 1 && r.CurrentPage == 300 && r.Status == model.ReviewStatusCompleted
	})).Return(nil)

	result, err := svc.UpdateProgress(1, 1, 300)
	assert.NoError(t, err)
	assert.Equal(t, 300, result.CurrentPage)
	assert.Equal(t, model.ReviewStatusCompleted, result.Status)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdateProgress_ExceedsTotalPages(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 300, CurrentPage: 200}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.UpdateProgress(1, 1, 350)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "総ページ数を超えることはできません")
}

func TestBookReviewUpdateProgress_NegativePage(t *testing.T) {
	svc, _ := newTestBookReviewService()

	_, err := svc.UpdateProgress(1, 1, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ページ数は0以上")
}

func TestBookReviewUpdateProgress_NoTotalPages(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 0, CurrentPage: 0, Status: model.ReviewStatusNotStarted}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.UpdateProgress(1, 1, 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "総ページ数が設定されていません")
}

func TestBookReviewUpdateProgress_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 99, Title: "他人の本"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.UpdateProgress(1, 1, 50)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestBookReviewUpdateProgress_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	_, err := svc.UpdateProgress(99, 1, 50)
	assert.Error(t, err)
}

func TestBookReviewUpdateProgress_AutoStartReading(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 300, CurrentPage: 0, Status: model.ReviewStatusNotStarted}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(r *model.BookReview) bool {
		return r.ID == 1 && r.CurrentPage == 50 && r.Status == model.ReviewStatusReading
	})).Return(nil)

	result, err := svc.UpdateProgress(1, 1, 50)
	assert.NoError(t, err)
	assert.Equal(t, 50, result.CurrentPage)
	assert.Equal(t, model.ReviewStatusReading, result.Status)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdateProgress_RepoError(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Go本", TotalPages: 300, CurrentPage: 0, Status: model.ReviewStatusReading}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*model.BookReview")).Return(errors.New("db error"))

	_, err := svc.UpdateProgress(1, 1, 100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

// ============================================================
// タイトル・レビュー本文の文字数バリデーションテスト
// ============================================================

func TestBookReviewCreate_TitleTooLong(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{Title: strings.Repeat("あ", 201), UserID: 1, Rating: 4}
	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトルは200文字以下")
}

func TestBookReviewCreate_ReviewTooLong(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{Title: "テスト本", UserID: 1, Rating: 4, Review: strings.Repeat("あ", 10001)}
	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "レビュー本文は10000文字以下")
}

func TestBookReviewUpdate_TitleTooLong(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Rating: 3}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Title: strings.Repeat("あ", 201)}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "タイトルは200文字以下")
}

func TestBookReviewUpdate_ReviewTooLong(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Rating: 3}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Review: strings.Repeat("あ", 10001)}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "レビュー本文は10000文字以下")
}

func TestBookReviewCreate_AuthorTooLong(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{Title: "テスト本", UserID: 1, Rating: 4, Author: strings.Repeat("あ", 201)}
	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "著者名は200文字以下")
}

func TestBookReviewCreate_ISBNTooLong(t *testing.T) {
	svc, _ := newTestBookReviewService()

	review := &model.BookReview{Title: "テスト本", UserID: 1, Rating: 4, ISBN: strings.Repeat("0", 21)}
	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ISBNは20文字以下")
}

func TestBookReviewUpdate_AuthorTooLong(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Rating: 3}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Author: strings.Repeat("あ", 201)}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "著者名は200文字以下")
}

func TestBookReviewUpdate_ISBNTooLong(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1, Title: "Old", Rating: 3}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{ISBN: strings.Repeat("0", 21)}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ISBNは20文字以下")
}
