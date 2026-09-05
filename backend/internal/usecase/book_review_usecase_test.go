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
)

// mockBookReviewRepo は usecase/repository.BookReviewRepository のモック。
type mockBookReviewRepo struct{ mock.Mock }

func (m *mockBookReviewRepo) Create(ctx context.Context, r *model.BookReview) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockBookReviewRepo) FindByID(ctx context.Context, id uint) (*model.BookReview, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.BookReview)
	return r, args.Error(1)
}

func (m *mockBookReviewRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) FindAll(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) FindByRating(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	args := m.Called(ctx, userID, minRating, maxRating)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Error(1)
}

func (m *mockBookReviewRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	r, _ := args.Get(0).([]model.BookReview)
	return r, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookReviewRepo) Update(ctx context.Context, r *model.BookReview) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockBookReviewRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockBookReviewRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ownedReviewOf は所有者が userID=1 の書籍レビューを返す。
func ownedReviewOf(id uint) *model.BookReview {
	r := &model.BookReview{
		UserID: 1, Title: "元のタイトル", Author: "元の著者", Rating: 3,
		TotalPages: 100, Status: model.ReviewStatusNotStarted,
	}
	r.ID = id
	return r
}

func TestCreateBookReviewUseCase_Execute(t *testing.T) {
	t.Run("画像URLが2000文字を超えると検証エラーになる", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		uc := usecase.NewCreateBookReviewUseCase(repo)

		err := uc.Execute(context.Background(), &model.BookReview{
			Title: "書名", Rating: 4,
			ImageURL: "https://example.com/" + strings.Repeat("a", 2000),
		})

		assert.Error(t, err)
		assert.True(t, domain.IsDomainError(err), "DomainError（400系）として返す: %v", err)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("2000文字ちょうどの画像URLは作成できる", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		imageURL := "https://example.com/" + strings.Repeat("a", 2000-len("https://example.com/"))
		repo.On("Create", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.ImageURL == imageURL
		})).Return(nil)
		uc := usecase.NewCreateBookReviewUseCase(repo)

		err := uc.Execute(context.Background(), &model.BookReview{
			Title: "書名", Rating: 4, ImageURL: imageURL,
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("前後空白を除いて作成する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.Title == "書名" && r.Author == "著者" && r.ISBN == "123" && r.Review == "感想"
		})).Return(nil)
		uc := usecase.NewCreateBookReviewUseCase(repo)

		err := uc.Execute(context.Background(), &model.BookReview{
			Title: "  書名  ", Author: "  著者  ", ISBN: "  123  ", Review: "  感想  ", Rating: 4,
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("入力が不正なら作成しない", func(t *testing.T) {
		cases := map[string]*model.BookReview{
			"タイトルが空":       {Title: "", Rating: 3},
			"タイトルが 201 文字": {Title: strings.Repeat("a", 201), Rating: 3},
			"評価 0":         {Title: "t", Rating: 0},
			"評価 6":         {Title: "t", Rating: 6},
			"総ページ数が上限超過":   {Title: "t", Rating: 3, TotalPages: 100000},
			"総ページ数が負":      {Title: "t", Rating: 3, TotalPages: -1},
		}
		for name, in := range cases {
			t.Run(name, func(t *testing.T) {
				repo := new(mockBookReviewRepo)
				uc := usecase.NewCreateBookReviewUseCase(repo)

				assert.Error(t, uc.Execute(context.Background(), in))
				repo.AssertNotCalled(t, "Create")
			})
		}
	})

	// 作成時は著者・ISBN・本文が空でもよい（更新時は 1 文字以上を要求する点と異なる）。
	t.Run("著者や本文は空でもよい", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateBookReviewUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), &model.BookReview{Title: "t", Rating: 3}))
		repo.AssertExpectations(t)
	})
}

func TestUpdateBookReviewUseCase_Execute(t *testing.T) {
	t.Run("空の項目と 0 の数値は据え置く部分更新", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.Title == "新題" && r.Author == "元の著者" && r.Rating == 3 && r.TotalPages == 100
		})).Return(nil)
		uc := usecase.NewUpdateBookReviewUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, &model.BookReview{Title: "  新題  "})

		assert.NoError(t, err)
		assert.Equal(t, "新題", got.Title)
		repo.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		other := ownedReviewOf(1)
		other.UserID = 999
		repo.On("FindByID", mock.Anything, uint(1)).Return(other, nil)
		uc := usecase.NewUpdateBookReviewUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.BookReview{Title: "x"})

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("評価が範囲外なら保存しない", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		uc := usecase.NewUpdateBookReviewUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.BookReview{Rating: 9})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestListBookReviewsByRatingUseCase_Execute(t *testing.T) {
	t.Run("評価範囲が不正なら BadRequest（DB を触らない）", func(t *testing.T) {
		cases := map[string][2]int{
			"最小が 0":   {0, 5},
			"最大が 6":   {1, 6},
			"最小 > 最大": {4, 2},
		}
		for name, rng := range cases {
			t.Run(name, func(t *testing.T) {
				repo := new(mockBookReviewRepo)
				uc := usecase.NewListBookReviewsByRatingUseCase(repo)

				_, err := uc.Execute(context.Background(), 1, rng[0], rng[1])

				assertDomainCode(t, err, domain.ErrCodeBadRequest)
				repo.AssertNotCalled(t, "FindByRating")
			})
		}
	})

	t.Run("正しい範囲なら委譲する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByRating", mock.Anything, uint(1), 4, 5).
			Return([]model.BookReview{*ownedReviewOf(1)}, nil)
		uc := usecase.NewListBookReviewsByRatingUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 4, 5)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		repo.AssertExpectations(t)
	})
}

func TestUpdateBookReviewProgressUseCase_Execute(t *testing.T) {
	t.Run("総ページ到達で読了へ遷移する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.CurrentPage == 100 && r.Status == model.ReviewStatusCompleted
		})).Return(nil)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, 100)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("未読から 1 以上で読中へ遷移する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.CurrentPage == 1 && r.Status == model.ReviewStatusReading
		})).Return(nil)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	// 進捗 0 はステータスを変えない。
	t.Run("進捗 0 ではステータスが変わらない", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.CurrentPage == 0 && r.Status == model.ReviewStatusNotStarted
		})).Return(nil)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, 0)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	// 検証の順序は「負数 → 所有権 → 総ページ数未設定 → 超過」。
	t.Run("負のページ数は所有権も見ずに BadRequest", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, -1)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "FindByID")
	})

	t.Run("総ページ数が未設定なら BadRequest（保存しない）", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		noTotal := ownedReviewOf(1)
		noTotal.TotalPages = 0
		repo.On("FindByID", mock.Anything, uint(1)).Return(noTotal, nil)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, 10)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("総ページ数を超えると BadRequest（保存しない）", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		uc := usecase.NewUpdateBookReviewProgressUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, 101)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestUpdateBookReviewStatusUseCase_Execute(t *testing.T) {
	t.Run("無効なステータスは所有権も見ずに BadRequest", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		uc := usecase.NewUpdateBookReviewStatusUseCase(repo)

		err := uc.Execute(context.Background(), 1, 1, model.ReviewStatus("unknown"))

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "FindByID")
	})

	t.Run("有効なステータスなら更新する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
			return r.Status == model.ReviewStatusReading
		})).Return(nil)
		uc := usecase.NewUpdateBookReviewStatusUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1, model.ReviewStatusReading))
		repo.AssertExpectations(t)
	})
}

func TestArchiveBookReviewUseCase_Execute(t *testing.T) {
	t.Run("アーカイブと解除を切り替える", func(t *testing.T) {
		for _, archived := range []bool{true, false} {
			repo := new(mockBookReviewRepo)
			repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
			repo.On("Update", mock.Anything, mock.MatchedBy(func(r *model.BookReview) bool {
				return r.IsArchived == archived
			})).Return(nil)
			uc := usecase.NewArchiveBookReviewUseCase(repo)

			assert.NoError(t, uc.Execute(context.Background(), 1, 1, archived))
			repo.AssertExpectations(t)
		}
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		other := ownedReviewOf(1)
		other.UserID = 999
		repo.On("FindByID", mock.Anything, uint(1)).Return(other, nil)
		uc := usecase.NewArchiveBookReviewUseCase(repo)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 1, true), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeleteBookReviewUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedReviewOf(1), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteBookReviewUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertExpectations(t)
	})

	t.Run("不在なら削除しない", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewDeleteBookReviewUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertNotCalled(t, "Delete")
	})
}

func TestSearchBookReviewsUseCase_Execute(t *testing.T) {
	t.Run("空のキーワードは BadRequest（検索しない）", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		uc := usecase.NewSearchBookReviewsUseCase(repo)

		_, _, err := uc.Execute(context.Background(), "   ", 20, 0)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Search")
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockBookReviewRepo)
		repo.On("Search", mock.Anything, "go", 20, 0).
			Return([]model.BookReview(nil), int64(0), errors.New("db error"))
		uc := usecase.NewSearchBookReviewsUseCase(repo)

		_, _, err := uc.Execute(context.Background(), "go", 20, 0)

		assert.Error(t, err)
	})
}
