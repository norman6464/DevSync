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

// mockNoteFolderRepo は usecase/repository.NoteFolderRepository のモック。
type mockNoteFolderRepo struct{ mock.Mock }

func (m *mockNoteFolderRepo) Create(ctx context.Context, folder *model.NoteFolder) error {
	return m.Called(ctx, folder).Error(0)
}

func (m *mockNoteFolderRepo) FindByID(ctx context.Context, id uint) (*model.NoteFolder, error) {
	args := m.Called(ctx, id)
	f, _ := args.Get(0).(*model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Get(1).(int64), args.Error(2)
}

func (m *mockNoteFolderRepo) FindByParentID(ctx context.Context, parentID uint) ([]model.NoteFolder, error) {
	args := m.Called(ctx, parentID)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) FindRootsByUserID(ctx context.Context, userID uint) ([]model.NoteFolder, error) {
	args := m.Called(ctx, userID)
	f, _ := args.Get(0).([]model.NoteFolder)
	return f, args.Error(1)
}

func (m *mockNoteFolderRepo) Update(ctx context.Context, folder *model.NoteFolder) error {
	return m.Called(ctx, folder).Error(0)
}

func (m *mockNoteFolderRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockNoteFolderRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// ownedNoteFolder は所有者が userID=1 のフォルダを返す。
func ownedNoteFolder() *model.NoteFolder {
	return &model.NoteFolder{ID: 1, UserID: 1, Name: "元の名前"}
}

func TestCreateNoteFolderUseCase_Execute(t *testing.T) {
	t.Run("名前が妥当なら作成する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		parent := uint(7)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
			return f.UserID == 1 && f.Name == "新規" && f.ParentID != nil && *f.ParentID == 7
		})).Return(nil)
		uc := usecase.NewCreateNoteFolderUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.CreateNoteFolderInput{
			UserID: 1, Name: "新規", ParentID: &parent,
		})

		assert.NoError(t, err)
		assert.Equal(t, "新規", got.Name)
		repo.AssertExpectations(t)
	})

	t.Run("名前が不正なら作成しない", func(t *testing.T) {
		cases := []struct {
			name  string
			input string
		}{
			{"空文字", ""},
			{"空白のみ", "   "},
			{"101 文字", strings.Repeat("a", 101)},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				repo := new(mockNoteFolderRepo)
				uc := usecase.NewCreateNoteFolderUseCase(repo)

				_, err := uc.Execute(context.Background(), usecase.CreateNoteFolderInput{UserID: 1, Name: tc.input})

				assert.Error(t, err)
				repo.AssertNotCalled(t, "Create")
			})
		}
	})

	t.Run("100 文字ちょうどは許可する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.CreateNoteFolderInput{
			UserID: 1, Name: strings.Repeat("a", 100),
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		uc := usecase.NewCreateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.CreateNoteFolderInput{Name: "新規"})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Create")
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.CreateNoteFolderInput{UserID: 1, Name: "新規"})

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetNoteFolderUseCase_Execute(t *testing.T) {
	t.Run("所有権を検証せずに返す", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.NoteFolder{ID: 1, UserID: 999, Name: "他人のフォルダ"}, nil)
		uc := usecase.NewGetNoteFolderUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, uint(999), got.UserID)
		repo.AssertExpectations(t)
	})

	// 不在は DomainError にしない（handler で 500 のままにするため）。
	t.Run("不在なら DomainError ではないエラーを返す", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		var de *domain.DomainError
		assert.False(t, errors.As(err, &de), "500 を維持するため DomainError にしない")
		repo.AssertExpectations(t)
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestListNoteFoldersUseCase_Execute(t *testing.T) {
	t.Run("一覧と総件数を返す", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).
			Return([]model.NoteFolder{{ID: 1, UserID: 1}}, int64(3), nil)
		uc := usecase.NewListNoteFoldersUseCase(repo)

		got, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(3), total)
		repo.AssertExpectations(t)
	})

	t.Run("ユーザーID が未指定なら BadRequest（DB を触らない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		uc := usecase.NewListNoteFoldersUseCase(repo)

		_, _, err := uc.Execute(context.Background(), 0, 20, 0)

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "FindByUserID")
	})
}

func TestListChildNoteFoldersUseCase_Execute(t *testing.T) {
	repo := new(mockNoteFolderRepo)
	repo.On("FindByParentID", mock.Anything, uint(5)).
		Return([]model.NoteFolder{{ID: 6, UserID: 1}}, nil)
	uc := usecase.NewListChildNoteFoldersUseCase(repo)

	got, err := uc.Execute(context.Background(), 5)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestListRootNoteFoldersUseCase_Execute(t *testing.T) {
	repo := new(mockNoteFolderRepo)
	repo.On("FindRootsByUserID", mock.Anything, uint(1)).
		Return([]model.NoteFolder{{ID: 1, UserID: 1}}, nil)
	uc := usecase.NewListRootNoteFoldersUseCase(repo)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestCountNoteFoldersUseCase_Execute(t *testing.T) {
	repo := new(mockNoteFolderRepo)
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(4), nil)
	uc := usecase.NewCountNoteFoldersUseCase(repo)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(4), got)
	repo.AssertExpectations(t)
}

func TestUpdateNoteFolderUseCase_Execute(t *testing.T) {
	t.Run("名前を前後の空白を除いて更新する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
			return f.Name == "新しい名前"
		})).Return(nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, Name: "  新しい名前  ",
		})

		assert.NoError(t, err)
		assert.Equal(t, "新しい名前", got.Name)
		repo.AssertExpectations(t)
	})

	// 空文字は「変更しない」を意味する部分更新。
	t.Run("名前が空文字なら据え置く", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
			return f.Name == "元の名前"
		})).Return(nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		got, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{ID: 1, UserID: 1})

		assert.NoError(t, err)
		assert.Equal(t, "元の名前", got.Name)
		repo.AssertExpectations(t)
	})

	t.Run("空白のみの名前は BadRequest（保存しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, Name: "   ",
		})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("101 文字の名前は保存しない", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, Name: strings.Repeat("a", 101),
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.NoteFolder{ID: 1, UserID: 999}, nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, Name: "変更",
		})

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("自分自身を親にすると BadRequest（子孫探索もしない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)
		self := uint(1)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, ParentID: &self,
		})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "FindByParentID")
		repo.AssertNotCalled(t, "Update")
	})

	// 孫を親にしようとする場合も再帰的に検出する。
	t.Run("子孫を親にすると BadRequest（保存しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("FindByParentID", mock.Anything, uint(1)).
			Return([]model.NoteFolder{{ID: 2, UserID: 1}}, nil)
		repo.On("FindByParentID", mock.Anything, uint(2)).
			Return([]model.NoteFolder{{ID: 3, UserID: 1}}, nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)
		grandChild := uint(3)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, ParentID: &grandChild,
		})

		assertDomainCode(t, err, domain.ErrCodeBadRequest)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("子孫でない親は設定できる", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("FindByParentID", mock.Anything, uint(1)).Return([]model.NoteFolder{}, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(f *model.NoteFolder) bool {
			return f.ParentID != nil && *f.ParentID == 9
		})).Return(nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)
		parent := uint(9)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, ParentID: &parent,
		})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("子孫探索の DB 障害を伝播する（保存しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("FindByParentID", mock.Anything, uint(1)).
			Return([]model.NoteFolder(nil), errors.New("db error"))
		uc := usecase.NewUpdateNoteFolderUseCase(repo)
		parent := uint(9)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{
			ID: 1, UserID: 1, ParentID: &parent,
		})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "Update")
	})

	t.Run("不在なら DomainError ではないエラーを返す（保存しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewUpdateNoteFolderUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.UpdateNoteFolderInput{ID: 1, UserID: 1})

		assert.Error(t, err)
		var de *domain.DomainError
		assert.False(t, errors.As(err, &de))
		repo.AssertNotCalled(t, "Update")
	})
}

func TestDeleteNoteFolderUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFolder(), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteNoteFolderUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（削除しない）", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).
			Return(&model.NoteFolder{ID: 1, UserID: 999}, nil)
		uc := usecase.NewDeleteNoteFolderUseCase(repo)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 1), domain.ErrCodeForbidden)
		repo.AssertNotCalled(t, "Delete")
	})

	t.Run("不在なら削除しない", func(t *testing.T) {
		repo := new(mockNoteFolderRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewDeleteNoteFolderUseCase(repo)

		assert.Error(t, uc.Execute(context.Background(), 1, 1))
		repo.AssertNotCalled(t, "Delete")
	})
}
