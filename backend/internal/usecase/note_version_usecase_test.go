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

// mockNoteVersionRepo は usecase/repository.NoteVersionRepository のモック。
type mockNoteVersionRepo struct{ mock.Mock }

func (m *mockNoteVersionRepo) Create(ctx context.Context, version *model.NoteVersion) error {
	return m.Called(ctx, version).Error(0)
}

func (m *mockNoteVersionRepo) FindByNoteID(ctx context.Context, noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	args := m.Called(ctx, noteID, limit, offset)
	v, _ := args.Get(0).([]model.NoteVersion)
	return v, args.Get(1).(int64), args.Error(2)
}

func (m *mockNoteVersionRepo) FindByID(ctx context.Context, id uint) (*model.NoteVersion, error) {
	args := m.Called(ctx, id)
	v, _ := args.Get(0).(*model.NoteVersion)
	return v, args.Error(1)
}

func (m *mockNoteVersionRepo) GetLatestVersionNumber(ctx context.Context, noteID uint) (int, error) {
	args := m.Called(ctx, noteID)
	return args.Int(0), args.Error(1)
}

// mockNoteUpdater は usecase/repository.NoteUpdater のモック。
type mockNoteUpdater struct{ mock.Mock }

func (m *mockNoteUpdater) Update(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}

func TestListNoteVersionsUseCase_Execute(t *testing.T) {
	t.Run("所有者には履歴と総件数を返す", func(t *testing.T) {
		versions, notes := new(mockNoteVersionRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		versions.On("FindByNoteID", mock.Anything, uint(1), 20, 0).
			Return([]model.NoteVersion{{ID: 5, NoteID: 1}}, int64(3), nil)
		uc := usecase.NewListNoteVersionsUseCase(versions, notes)

		got, total, err := uc.Execute(context.Background(), 1, 1, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(3), total)
		versions.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（履歴を読まない）", func(t *testing.T) {
		versions, notes := new(mockNoteVersionRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)
		uc := usecase.NewListNoteVersionsUseCase(versions, notes)

		_, _, err := uc.Execute(context.Background(), 1, 1, 20, 0)

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		versions.AssertNotCalled(t, "FindByNoteID")
	})
}

func TestGetNoteVersionUseCase_Execute(t *testing.T) {
	t.Run("対象ノートのバージョンを返す", func(t *testing.T) {
		versions, notes := new(mockNoteVersionRepo), new(mockNoteReader)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		versions.On("FindByID", mock.Anything, uint(5)).
			Return(&model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 2}, nil)
		uc := usecase.NewGetNoteVersionUseCase(versions, notes)

		got, err := uc.Execute(context.Background(), 1, 5, 1)

		assert.NoError(t, err)
		assert.Equal(t, 2, got.VersionNumber)
		versions.AssertExpectations(t)
	})

	// 不在と「他ノートのバージョン」を同じ 404 に揃えている（移行前の挙動）。
	t.Run("不在も他ノートのバージョンも NotFound", func(t *testing.T) {
		cases := []struct {
			name    string
			version *model.NoteVersion
		}{
			{"不在", nil},
			{"他ノートのバージョン", &model.NoteVersion{ID: 5, NoteID: 999}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				versions, notes := new(mockNoteVersionRepo), new(mockNoteReader)
				notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
				versions.On("FindByID", mock.Anything, uint(5)).Return(tc.version, nil)
				uc := usecase.NewGetNoteVersionUseCase(versions, notes)

				_, err := uc.Execute(context.Background(), 1, 5, 1)

				assertDomainCode(t, err, domain.ErrCodeNotFound)
				var de *domain.DomainError
				if assert.ErrorAs(t, err, &de) {
					assert.Equal(t, "バージョンが見つかりません", de.Message)
				}
			})
		}
	})
}

func TestRestoreNoteVersionUseCase_Execute(t *testing.T) {
	t.Run("復元前の内容を新しいバージョンとして残してから書き戻す", func(t *testing.T) {
		versions, notes, writer := new(mockNoteVersionRepo), new(mockNoteReader), new(mockNoteUpdater)
		current := &model.Note{ID: 1, UserID: 1, Title: "いまの題", Content: "いまの本文", Tags: "now"}
		notes.On("FindByID", mock.Anything, uint(1)).Return(current, nil)
		versions.On("FindByID", mock.Anything, uint(5)).
			Return(&model.NoteVersion{ID: 5, NoteID: 1, Title: "むかしの題", Content: "むかしの本文", Tags: "old"}, nil)
		versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(3, nil)
		versions.On("Create", mock.Anything, mock.MatchedBy(func(v *model.NoteVersion) bool {
			// 採番は最新 + 1、内容は復元前のもの
			return v.VersionNumber == 4 && v.Title == "いまの題" && v.Content == "いまの本文" && v.Tags == "now"
		})).Return(nil)
		writer.On("Update", mock.Anything, mock.MatchedBy(func(n *model.Note) bool {
			return n.Title == "むかしの題" && n.Content == "むかしの本文" && n.Tags == "old"
		})).Return(nil)
		uc := usecase.NewRestoreNoteVersionUseCase(versions, notes, writer)

		got, err := uc.Execute(context.Background(), 1, 5, 1)

		assert.NoError(t, err)
		assert.Equal(t, "むかしの題", got.Title)
		versions.AssertExpectations(t)
		writer.AssertExpectations(t)
	})

	// バージョンが 1 つも無い場合は 1 番から採番する。
	t.Run("バージョンが無ければ 1 番から採番する", func(t *testing.T) {
		versions, notes, writer := new(mockNoteVersionRepo), new(mockNoteReader), new(mockNoteUpdater)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		versions.On("FindByID", mock.Anything, uint(5)).Return(&model.NoteVersion{ID: 5, NoteID: 1}, nil)
		versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(0, nil)
		versions.On("Create", mock.Anything, mock.MatchedBy(func(v *model.NoteVersion) bool {
			return v.VersionNumber == 1
		})).Return(nil)
		writer.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewRestoreNoteVersionUseCase(versions, notes, writer)

		_, err := uc.Execute(context.Background(), 1, 5, 1)

		assert.NoError(t, err)
		versions.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（バージョンも読まない）", func(t *testing.T) {
		versions, notes, writer := new(mockNoteVersionRepo), new(mockNoteReader), new(mockNoteUpdater)
		notes.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 999}, nil)
		uc := usecase.NewRestoreNoteVersionUseCase(versions, notes, writer)

		_, err := uc.Execute(context.Background(), 1, 5, 1)

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		versions.AssertNotCalled(t, "FindByID")
		writer.AssertNotCalled(t, "Update")
	})

	t.Run("採番の DB 障害を伝播する（保存も書き戻しもしない）", func(t *testing.T) {
		versions, notes, writer := new(mockNoteVersionRepo), new(mockNoteReader), new(mockNoteUpdater)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		versions.On("FindByID", mock.Anything, uint(5)).Return(&model.NoteVersion{ID: 5, NoteID: 1}, nil)
		versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(0, errors.New("db error"))
		uc := usecase.NewRestoreNoteVersionUseCase(versions, notes, writer)

		_, err := uc.Execute(context.Background(), 1, 5, 1)

		assert.Error(t, err)
		versions.AssertNotCalled(t, "Create")
		writer.AssertNotCalled(t, "Update")
	})

	t.Run("書き戻しの DB 障害を伝播する", func(t *testing.T) {
		versions, notes, writer := new(mockNoteVersionRepo), new(mockNoteReader), new(mockNoteUpdater)
		notes.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteOf(1), nil)
		versions.On("FindByID", mock.Anything, uint(5)).Return(&model.NoteVersion{ID: 5, NoteID: 1}, nil)
		versions.On("GetLatestVersionNumber", mock.Anything, uint(1)).Return(1, nil)
		versions.On("Create", mock.Anything, mock.Anything).Return(nil)
		writer.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewRestoreNoteVersionUseCase(versions, notes, writer)

		_, err := uc.Execute(context.Background(), 1, 5, 1)

		assert.Error(t, err)
		writer.AssertExpectations(t)
	})
}
