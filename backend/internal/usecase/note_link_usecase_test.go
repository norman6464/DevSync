package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNoteLinkRepo は usecase/repository.NoteLinkRepository のモック。
type mockNoteLinkRepo struct{ mock.Mock }

func (m *mockNoteLinkRepo) Create(ctx context.Context, link *model.NoteLink) error {
	return m.Called(ctx, link).Error(0)
}

func (m *mockNoteLinkRepo) FindBySourceNoteID(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(ctx, sourceNoteID)
	l, _ := args.Get(0).([]model.NoteLink)
	return l, args.Error(1)
}

func (m *mockNoteLinkRepo) FindByTargetNoteID(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(ctx, targetNoteID)
	l, _ := args.Get(0).([]model.NoteLink)
	return l, args.Error(1)
}

func (m *mockNoteLinkRepo) Delete(ctx context.Context, sourceNoteID, targetNoteID uint) error {
	return m.Called(ctx, sourceNoteID, targetNoteID).Error(0)
}

func (m *mockNoteLinkRepo) Exists(ctx context.Context, sourceNoteID, targetNoteID uint) (bool, error) {
	args := m.Called(ctx, sourceNoteID, targetNoteID)
	return args.Bool(0), args.Error(1)
}

func (m *mockNoteLinkRepo) CountBySourceNoteID(ctx context.Context, noteID uint) (int64, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNoteLinkRepo) CountByTargetNoteID(ctx context.Context, noteID uint) (int64, error) {
	args := m.Called(ctx, noteID)
	return args.Get(0).(int64), args.Error(1)
}

// mockNoteReader は usecase/repository.NoteReader のモック。
type mockNoteReader struct{ mock.Mock }

func (m *mockNoteReader) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*model.Note)
	return n, args.Error(1)
}

// ownedNoteOf は所有者が userID=1 のノートを返す。
func ownedNoteOf(id uint) *model.Note { return &model.Note{ID: id, UserID: 1} }

func TestCreateNoteLinkUseCase_Execute(t *testing.T) {
	t.Run("検査を通ればリンクを作成する", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		notes.On("FindByID", mock.Anything, uint(2)).Return(ownedNoteOf(2), nil)
		links.On("Exists", mock.Anything, uint(1), uint(2)).Return(false, nil)
		links.On("Create", mock.Anything, mock.MatchedBy(func(l *model.NoteLink) bool {
			return l.SourceNoteID == 1 && l.TargetNoteID == 2
		})).Return(nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assert.NoError(t, uc.Execute(context.Background(), 1, 2, 1))
		links.AssertExpectations(t)
		notes.AssertExpectations(t)
	})

	// 自己リンクは最初に弾くため、ノートも読まない。
	t.Run("同じノートへのリンクは Validation エラー", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 1, 1), domain.ErrCodeValidation)
		notes.AssertNotCalled(t, "FindByID")
		links.AssertNotCalled(t, "Create")
	})

	t.Run("ソースノートの所有者でなければ Forbidden", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 2, 1), domain.ErrCodeForbidden)
		links.AssertNotCalled(t, "Create")
	})

	t.Run("ソースノートが不在なら NotFound", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 2, 1), domain.ErrCodeNotFound)
		links.AssertNotCalled(t, "Create")
	})

	t.Run("リンク先が不在なら NotFound", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		notes.On("FindByID", mock.Anything, uint(2)).Return(nil, nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 2, 1), domain.ErrCodeNotFound)
		links.AssertNotCalled(t, "Create")
	})

	// リンク先は他ユーザーのノートでもよい（所有権は確認しない）。
	t.Run("リンク先が他ユーザーのノートでも作成できる", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		notes.On("FindByID", mock.Anything, uint(2)).Return(&model.Note{ID: 2, UserID: 999}, nil)
		links.On("Exists", mock.Anything, uint(1), uint(2)).Return(false, nil)
		links.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assert.NoError(t, uc.Execute(context.Background(), 1, 2, 1))
		links.AssertExpectations(t)
	})

	t.Run("既に存在すれば Validation エラー（作成しない）", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		notes.On("FindByID", mock.Anything, uint(2)).Return(ownedNoteOf(2), nil)
		links.On("Exists", mock.Anything, uint(1), uint(2)).Return(true, nil)
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 2, 1), domain.ErrCodeValidation)
		links.AssertNotCalled(t, "Create")
	})

	t.Run("重複確認の DB 障害を伝播する（作成しない）", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		notes.On("FindByID", mock.Anything, uint(2)).Return(ownedNoteOf(2), nil)
		links.On("Exists", mock.Anything, uint(1), uint(2)).Return(false, errors.New("db error"))
		uc := usecase.NewCreateNoteLinkUseCase(links, notes)

		assert.Error(t, uc.Execute(context.Background(), 1, 2, 1))
		links.AssertNotCalled(t, "Create")
	})
}

func TestListNoteLinksUseCase_Execute(t *testing.T) {
	links := new(mockNoteLinkRepo)
	links.On("FindBySourceNoteID", mock.Anything, uint(1)).
		Return([]model.NoteLink{{ID: 1, SourceNoteID: 1, TargetNoteID: 2}}, nil)
	uc := usecase.NewListNoteLinksUseCase(links)

	got, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	links.AssertExpectations(t)
}

func TestListNoteBacklinksUseCase_Execute(t *testing.T) {
	links := new(mockNoteLinkRepo)
	links.On("FindByTargetNoteID", mock.Anything, uint(2)).
		Return([]model.NoteLink{{ID: 1, SourceNoteID: 1, TargetNoteID: 2}}, nil)
	uc := usecase.NewListNoteBacklinksUseCase(links)

	got, err := uc.Execute(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	links.AssertExpectations(t)
}

func TestDeleteNoteLinkUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		links.On("Delete", mock.Anything, uint(1), uint(2)).Return(nil)
		uc := usecase.NewDeleteNoteLinkUseCase(links, notes)

		assert.NoError(t, uc.Execute(context.Background(), 1, 2, 1))
		links.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（削除しない）", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)
		uc := usecase.NewDeleteNoteLinkUseCase(links, notes)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 2, 1), domain.ErrCodeForbidden)
		links.AssertNotCalled(t, "Delete")
	})
}

func TestGetNoteLinkStatsUseCase_Execute(t *testing.T) {
	t.Run("フォワードとバックリンクの件数を返す", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		links.On("CountBySourceNoteID", mock.Anything, uint(1)).Return(int64(3), nil)
		links.On("CountByTargetNoteID", mock.Anything, uint(1)).Return(int64(5), nil)
		uc := usecase.NewGetNoteLinkStatsUseCase(links, notes)

		got, err := uc.Execute(context.Background(), 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), got.NoteID)
		assert.Equal(t, int64(3), got.ForwardLinkCount)
		assert.Equal(t, int64(5), got.BacklinkCount)
		links.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（集計しない）", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)
		uc := usecase.NewGetNoteLinkStatsUseCase(links, notes)

		_, err := uc.Execute(context.Background(), 1, 1)

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		links.AssertNotCalled(t, "CountBySourceNoteID")
	})

	// 作成・削除は 404 だが、統計だけは DomainError にならず 500 になる（移行前の挙動）。
	t.Run("ノートが不在なら 404 を返す", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetNoteLinkStatsUseCase(links, notes)

		_, err := uc.Execute(context.Background(), 1, 1)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		links.AssertNotCalled(t, "CountBySourceNoteID")
	})

	t.Run("バックリンク集計の DB 障害を伝播する", func(t *testing.T) {
		links, notes := new(mockNoteLinkRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		links.On("CountBySourceNoteID", mock.Anything, uint(1)).Return(int64(3), nil)
		links.On("CountByTargetNoteID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))
		uc := usecase.NewGetNoteLinkStatsUseCase(links, notes)

		_, err := uc.Execute(context.Background(), 1, 1)

		assert.Error(t, err)
		links.AssertExpectations(t)
	})
}
