package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockNoteRepo は usecase/repository.NoteRepository のモック。
type mockNoteRepo struct{ mock.Mock }

func (m *mockNoteRepo) Create(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}
func (m *mockNoteRepo) Update(ctx context.Context, note *model.Note) error {
	return m.Called(ctx, note).Error(0)
}
func (m *mockNoteRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) FindByFolderID(ctx context.Context, folderID, userID uint) ([]model.Note, error) {
	args := m.Called(ctx, folderID, userID)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) Search(ctx context.Context, userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	args := m.Called(ctx, userID, query, limit, offset)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Get(1).(int64), args.Error(2)
}
func (m *mockNoteRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockNoteRepo) ToggleFavorite(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindFavorites(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) CountFavoritesByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockNoteRepo) Archive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) Unarchive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockNoteRepo) FindArchived(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(ctx, userID, page, limit)
	n, _ := args.Get(0).([]model.Note)
	return n, args.Error(1)
}
func (m *mockNoteRepo) CountArchivedByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// assertNoteStatus は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertNoteStatus(t *testing.T, err error, want int) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.HTTPStatus())
}

// ownedNoteFixture は所有者 1 のノートを返す。
func ownedNoteFixture() *model.Note {
	return &model.Note{ID: 1, UserID: 1, Title: "旧題", Content: "旧本文", Tags: "go"}
}

func TestCreateNoteUseCase_Execute(t *testing.T) {
	t.Run("検証を通れば作成する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)
		uc := usecase.NewCreateNoteUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), &model.Note{Title: "題", Content: "本文"}))
		repo.AssertExpectations(t)
	})

	t.Run("タイトルが空なら 400 で作成しない", func(t *testing.T) {
		repo := new(mockNoteRepo)
		uc := usecase.NewCreateNoteUseCase(repo)

		assertNoteStatus(t, uc.Execute(context.Background(), &model.Note{Content: "本文"}), http.StatusBadRequest)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestGetNoteUseCase_Execute(t *testing.T) {
	t.Run("所有者なら取得できる", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFixture(), nil)
		uc := usecase.NewGetNoteUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), got.ID)
	})

	t.Run("他人のノートは 403", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)
		uc := usecase.NewGetNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assertNoteStatus(t, err, http.StatusForbidden)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.ErrorContains(t, err, "db error")
	})
}

func TestUpdateNoteUseCase_Execute(t *testing.T) {
	t.Run("空文字列は変更なしとして扱う", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFixture(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)
		uc := usecase.NewUpdateNoteUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, "新題", "", "", nil)
		require.NoError(t, err)
		assert.Equal(t, "新題", got.Title)
		assert.Equal(t, "旧本文", got.Content)
		assert.Equal(t, "go", got.Tags)
	})

	t.Run("空白のみのタイトルは専用メッセージで 400", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFixture(), nil)
		uc := usecase.NewUpdateNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, "   ", "", "", nil)
		assertNoteStatus(t, err, http.StatusBadRequest)
		assert.Equal(t, "タイトルは空白のみにできません", domain.GetDomainError(err).Message)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("空白のみの本文も専用メッセージで 400", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFixture(), nil)
		uc := usecase.NewUpdateNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, "", "   ", "", nil)
		assertNoteStatus(t, err, http.StatusBadRequest)
		assert.Equal(t, "本文は空白のみにできません", domain.GetDomainError(err).Message)
	})

	t.Run("フォルダの付け替えができる", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedNoteFixture(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)
		uc := usecase.NewUpdateNoteUseCase(repo)

		folderID := uint(7)
		got, err := uc.Execute(context.Background(), 1, 1, "", "", "", &folderID)
		require.NoError(t, err)
		require.NotNil(t, got.FolderID)
		assert.Equal(t, uint(7), *got.FolderID)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)
		uc := usecase.NewUpdateNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, "新題", "", "", nil)
		assertNoteStatus(t, err, http.StatusForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestNoteOwnershipGuardedUseCases(t *testing.T) {
	others := &model.Note{ID: 1, UserID: 99}

	cases := []struct {
		name string
		call func(repo *mockNoteRepo) error
		skip string
	}{
		{"削除", func(repo *mockNoteRepo) error {
			return usecase.NewDeleteNoteUseCase(repo).Execute(context.Background(), 1, 1)
		}, "Delete"},
		{"お気に入り切替", func(repo *mockNoteRepo) error {
			return usecase.NewToggleNoteFavoriteUseCase(repo).Execute(context.Background(), 1, 1)
		}, "ToggleFavorite"},
		{"アーカイブ", func(repo *mockNoteRepo) error {
			return usecase.NewArchiveNoteUseCase(repo).Execute(context.Background(), 1, 1)
		}, "Archive"},
		{"アーカイブ解除", func(repo *mockNoteRepo) error {
			return usecase.NewUnarchiveNoteUseCase(repo).Execute(context.Background(), 1, 1)
		}, "Unarchive"},
	}

	for _, tc := range cases {
		t.Run(tc.name+"は他人のノートなら 403", func(t *testing.T) {
			repo := new(mockNoteRepo)
			repo.On("FindByID", mock.Anything, uint(1)).Return(others, nil)

			assertNoteStatus(t, tc.call(repo), http.StatusForbidden)
			repo.AssertNotCalled(t, tc.skip, mock.Anything, mock.Anything)
		})
	}
}

func TestDuplicateNoteUseCase_Execute(t *testing.T) {
	t.Run("タイトルに (コピー) を付けて状態をリセットする", func(t *testing.T) {
		repo := new(mockNoteRepo)
		folderID := uint(3)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{
			ID: 1, UserID: 1, Title: "元ノート", Content: "本文", Tags: "go",
			FolderID: &folderID, IsFavorite: true, IsArchived: true,
		}, nil)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Note")).Return(nil)
		uc := usecase.NewDuplicateNoteUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, "元ノート (コピー)", got.Title)
		assert.Equal(t, "本文", got.Content)
		assert.False(t, got.IsFavorite, "お気に入りはリセットされる")
		assert.False(t, got.IsArchived, "アーカイブはリセットされる")
		require.NotNil(t, got.FolderID)
		assert.Equal(t, uint(3), *got.FolderID, "フォルダは引き継ぐ")
	})

	t.Run("他人のノートは複製できない", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)
		uc := usecase.NewDuplicateNoteUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assertNoteStatus(t, err, http.StatusForbidden)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestExportNoteMarkdownUseCase_Execute(t *testing.T) {
	created := time.Date(2026, 1, 2, 3, 4, 0, 0, time.UTC)
	updated := time.Date(2026, 2, 3, 4, 5, 0, 0, time.UTC)

	t.Run("Markdown の書式を保つ", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{
			ID: 1, UserID: 1, Title: "ノート", Content: "本文です", Tags: "go,test",
			CreatedAt: created, UpdatedAt: updated,
		}, nil)
		uc := usecase.NewExportNoteMarkdownUseCase(repo)

		data, title, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, "ノート", title)
		assert.Equal(t, "# ノート\n\n**Tags:** go,test\n**Created:** 2026-01-02 03:04\n**Updated:** 2026-02-03 04:05\n\n---\n\n本文です\n", string(data))
	})

	t.Run("タグが空なら Tags 行を出さない", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{
			ID: 1, UserID: 1, Title: "ノート", Content: "本文", CreatedAt: created, UpdatedAt: updated,
		}, nil)
		uc := usecase.NewExportNoteMarkdownUseCase(repo)

		data, _, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.NotContains(t, string(data), "**Tags:**")
	})

	t.Run("他人のノートは書き出せない", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)
		uc := usecase.NewExportNoteMarkdownUseCase(repo)

		_, _, err := uc.Execute(context.Background(), 1, 1)
		assertNoteStatus(t, err, http.StatusForbidden)
	})
}

func TestListNoteTagsUseCase_Execute(t *testing.T) {
	t.Run("カンマ区切りのタグを重複なく抽出する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), 1, 1000).Return([]model.Note{
			{Tags: "go, test"},
			{Tags: "go,web"},
			{Tags: ""},
			{Tags: "  ,  "},
		}, nil)
		uc := usecase.NewListNoteTagsUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, []string{"go", "test", "web"}, got)
	})

	t.Run("取得に失敗したらそのまま返す", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), 1, 1000).
			Return([]model.Note(nil), errors.New("db error"))
		uc := usecase.NewListNoteTagsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1)
		require.Error(t, err)
	})
}

func TestExtractUniqueNoteTags(t *testing.T) {
	t.Run("タグが無ければ空を返す", func(t *testing.T) {
		assert.Empty(t, usecase.ExtractUniqueNoteTags([]model.Note{{Tags: ""}}))
	})

	t.Run("前後の空白を除いて重複を排除する", func(t *testing.T) {
		got := usecase.ExtractUniqueNoteTags([]model.Note{{Tags: " a , b "}, {Tags: "b,c"}})
		assert.Equal(t, []string{"a", "b", "c"}, got)
	})
}

func TestNotePassThroughUseCases(t *testing.T) {
	ctx := context.Background()

	t.Run("検索はページからオフセットを計算する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("Search", mock.Anything, uint(1), "go", 20, 40).
			Return([]model.Note{{ID: 1}}, int64(1), nil)
		got, total, err := usecase.NewSearchNotesUseCase(repo).Execute(ctx, 1, "go", 3, 20)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), 1, 20).Return([]model.Note{{ID: 1}}, nil)
		got, err := usecase.NewListNotesUseCase(repo).Execute(ctx, 1, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("フォルダ別一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindByFolderID", mock.Anything, uint(5), uint(1)).Return([]model.Note{{ID: 1}}, nil)
		got, err := usecase.NewListNotesByFolderUseCase(repo).Execute(ctx, 5, 1)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("お気に入り一覧と件数は repo に委譲する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindFavorites", mock.Anything, uint(1), 1, 20).Return([]model.Note{{ID: 1}}, nil)
		repo.On("CountFavoritesByUserID", mock.Anything, uint(1)).Return(int64(1), nil)
		notes, err := usecase.NewListFavoriteNotesUseCase(repo).Execute(ctx, 1, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		count, err := usecase.NewCountFavoriteNotesUseCase(repo).Execute(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("アーカイブ一覧と件数は repo に委譲する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("FindArchived", mock.Anything, uint(1), 1, 20).Return([]model.Note{{ID: 1}}, nil)
		repo.On("CountArchivedByUserID", mock.Anything, uint(1)).Return(int64(1), nil)
		notes, err := usecase.NewListArchivedNotesUseCase(repo).Execute(ctx, 1, 1, 20)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		count, err := usecase.NewCountArchivedNotesUseCase(repo).Execute(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("件数は repo に委譲する", func(t *testing.T) {
		repo := new(mockNoteRepo)
		repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(9), nil)
		got, err := usecase.NewCountNotesUseCase(repo).Execute(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(9), got)
	})
}
