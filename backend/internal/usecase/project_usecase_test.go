package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockProjectRepo は usecase/repository.ProjectRepository のモック。
type mockProjectRepo struct{ mock.Mock }

func (m *mockProjectRepo) Create(ctx context.Context, project *model.Project) error {
	return m.Called(ctx, project).Error(0)
}
func (m *mockProjectRepo) Update(ctx context.Context, project *model.Project) error {
	return m.Called(ctx, project).Error(0)
}
func (m *mockProjectRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockProjectRepo) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Project)
	return p, args.Error(1)
}
func (m *mockProjectRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}
func (m *mockProjectRepo) FindFeaturedByUserID(ctx context.Context, userID uint) ([]model.Project, error) {
	args := m.Called(ctx, userID)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Error(1)
}
func (m *mockProjectRepo) FindArchivedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}
func (m *mockProjectRepo) FindAll(ctx context.Context, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}
func (m *mockProjectRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}
func (m *mockProjectRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockProjectRepo) Archive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockProjectRepo) Unarchive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ownedProject は指定ユーザーが所有するプロジェクトを返すテスト用ヘルパー。
func ownedProject(id, userID uint) *model.Project {
	return &model.Project{
		ID: id, UserID: userID,
		Title: "既存プロジェクト", Description: "説明", TechStack: "Go,React",
		DemoURL: "https://demo.example.com", GithubURL: "https://github.com/example/repo",
		ImageURL: "https://img.example.com/a.png", Role: "Lead Developer",
	}
}

// ============================================================
// Create
// ============================================================

func TestCreateProjectUseCase_Execute(t *testing.T) {
	t.Run("検証を通れば作成する", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Project")).Return(nil)
		uc := usecase.NewCreateProjectUseCase(repo)

		err := uc.Execute(context.Background(), &model.Project{UserID: 1, Title: "新規", Description: "説明"})

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	// 説明は DTO では任意だが、usecase の検証では必須（移行前から変わらない）。
	t.Run("URL は空でも作成できるが説明は必須", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateProjectUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), &model.Project{UserID: 1, Title: "新規", Description: "説明"}))

		err := usecase.NewCreateProjectUseCase(new(mockProjectRepo)).
			Execute(context.Background(), &model.Project{UserID: 1, Title: "新規"})
		require.Error(t, err)
		assert.NotNil(t, domain.GetDomainError(err))
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name    string
			project *model.Project
		}{
			{"タイトルが空", &model.Project{Title: "", Description: "説明"}},
			{"タイトルが空白のみ", &model.Project{Title: "   ", Description: "説明"}},
			{"説明が空", &model.Project{Title: "題"}},
			{"デモ URL が不正", &model.Project{Title: "題", Description: "説明", DemoURL: "not-a-url"}},
			{"GitHub URL が不正", &model.Project{Title: "題", Description: "説明", GithubURL: "not-a-url"}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockProjectRepo)
				uc := usecase.NewCreateProjectUseCase(repo)

				err := uc.Execute(context.Background(), c.project)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateProjectUseCase(repo)

		err := uc.Execute(context.Background(), &model.Project{UserID: 1, Title: "新規", Description: "説明"})

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// 取得 / 一覧 / 検索 / 件数
// ============================================================

func TestGetProjectUseCase_Execute(t *testing.T) {
	t.Run("所有者は取得できる", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewGetProjectUseCase(repo)

		project, err := uc.Execute(context.Background(), 1, 5)

		require.NoError(t, err)
		assert.Equal(t, uint(1), project.ID)
	})

	t.Run("他人のプロジェクトは 403", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewGetProjectUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 9)

		assert.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetProjectUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestListProjectsByUserUseCase_Execute(t *testing.T) {
	repo := new(mockProjectRepo)
	repo.On("FindByUserID", mock.Anything, uint(5), 20, 0).
		Return([]model.Project{*ownedProject(1, 5)}, int64(1), nil)
	uc := usecase.NewListProjectsByUserUseCase(repo)

	projects, total, err := uc.Execute(context.Background(), 5, 20, 0)

	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, int64(1), total)
}

func TestListFeaturedProjectsUseCase_Execute(t *testing.T) {
	repo := new(mockProjectRepo)
	featured := ownedProject(1, 5)
	featured.Featured = true
	repo.On("FindFeaturedByUserID", mock.Anything, uint(5)).Return([]model.Project{*featured}, nil)
	uc := usecase.NewListFeaturedProjectsUseCase(repo)

	projects, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Len(t, projects, 1)
}

func TestListArchivedProjectsUseCase_Execute(t *testing.T) {
	repo := new(mockProjectRepo)
	repo.On("FindArchivedByUserID", mock.Anything, uint(5), 20, 0).
		Return([]model.Project{*ownedProject(1, 5)}, int64(1), nil)
	uc := usecase.NewListArchivedProjectsUseCase(repo)

	projects, total, err := uc.Execute(context.Background(), 5, 20, 0)

	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, int64(1), total)
}

func TestListAllProjectsUseCase_Execute(t *testing.T) {
	repo := new(mockProjectRepo)
	repo.On("FindAll", mock.Anything, 20, 0).
		Return([]model.Project{*ownedProject(1, 5)}, int64(1), nil)
	uc := usecase.NewListAllProjectsUseCase(repo)

	projects, total, err := uc.Execute(context.Background(), 20, 0)

	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, int64(1), total)
}

func TestSearchProjectsUseCase_Execute(t *testing.T) {
	t.Run("前後の空白を落として検索する", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("Search", mock.Anything, "go", 20, 0).
			Return([]model.Project{*ownedProject(1, 5)}, int64(1), nil)
		uc := usecase.NewSearchProjectsUseCase(repo)

		projects, total, err := uc.Execute(context.Background(), "  go  ", 20, 0)

		require.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, int64(1), total)
		repo.AssertExpectations(t)
	})

	t.Run("空のキーワードは 400", func(t *testing.T) {
		for _, q := range []string{"", "   "} {
			repo := new(mockProjectRepo)
			uc := usecase.NewSearchProjectsUseCase(repo)

			_, _, err := uc.Execute(context.Background(), q, 20, 0)

			require.Error(t, err)
			assert.NotNil(t, domain.GetDomainError(err))
			repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		}
	})
}

func TestCountProjectsUseCase_Execute(t *testing.T) {
	repo := new(mockProjectRepo)
	repo.On("CountByUserID", mock.Anything, uint(5)).Return(int64(3), nil)
	uc := usecase.NewCountProjectsUseCase(repo)

	count, err := uc.Execute(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// ============================================================
// Update / Featured
// ============================================================

func TestUpdateProjectUseCase_Execute(t *testing.T) {
	t.Run("指定したフィールドだけを更新する", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateProjectUseCase(repo)

		project, err := uc.Execute(context.Background(), 1, 5, &model.Project{Title: "  新タイトル  "})

		require.NoError(t, err)
		assert.Equal(t, "新タイトル", project.Title)
		assert.Equal(t, "説明", project.Description)
		assert.Equal(t, "Go,React", project.TechStack)
		assert.Equal(t, "Lead Developer", project.Role)
	})

	t.Run("日付とリポジトリ ID はポインタが渡されたときだけ変更する", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateProjectUseCase(repo)

		start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		repoID := uint(42)
		project, err := uc.Execute(context.Background(), 1, 5, &model.Project{
			StartDate: &start, GithubRepoID: &repoID,
		})

		require.NoError(t, err)
		require.NotNil(t, project.StartDate)
		assert.Equal(t, start, *project.StartDate)
		require.NotNil(t, project.GithubRepoID)
		assert.Equal(t, uint(42), *project.GithubRepoID)
		assert.Nil(t, project.EndDate)
	})

	t.Run("空白のみのフィールドは据え置く", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Update", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewUpdateProjectUseCase(repo)

		project, err := uc.Execute(context.Background(), 1, 5, &model.Project{TechStack: "   ", Role: "  "})

		require.NoError(t, err)
		assert.Equal(t, "Go,React", project.TechStack)
		assert.Equal(t, "Lead Developer", project.Role)
	})

	t.Run("検証エラーでは書き込まない", func(t *testing.T) {
		cases := []struct {
			name    string
			updates *model.Project
		}{
			{"タイトルが上限超過", &model.Project{Title: strings.Repeat("あ", 201)}},
			{"デモ URL が不正", &model.Project{DemoURL: "not-a-url"}},
			{"GitHub URL が不正", &model.Project{GithubURL: "not-a-url"}},
			{"技術スタックが上限超過", &model.Project{TechStack: strings.Repeat("a", 501)}},
			{"役割が上限超過", &model.Project{Role: strings.Repeat("あ", 101)}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				repo := new(mockProjectRepo)
				repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
				uc := usecase.NewUpdateProjectUseCase(repo)

				_, err := uc.Execute(context.Background(), 1, 5, c.updates)

				require.Error(t, err)
				assert.NotNil(t, domain.GetDomainError(err))
				repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			})
		}
	})

	t.Run("他人のプロジェクトは 403", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewUpdateProjectUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 9, &model.Project{Title: "題"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestUpdateProjectFeaturedUseCase_Execute(t *testing.T) {
	t.Run("注目指定を切り替える", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
			return p.Featured
		})).Return(nil)
		uc := usecase.NewUpdateProjectFeaturedUseCase(repo)

		project, err := uc.Execute(context.Background(), 1, 5, true)

		require.NoError(t, err)
		assert.True(t, project.Featured)
		repo.AssertExpectations(t)
	})

	t.Run("他人のプロジェクトは 403", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewUpdateProjectFeaturedUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 9, true)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

// ============================================================
// アーカイブ / 削除
// ============================================================

func TestArchiveProjectUseCase_Execute(t *testing.T) {
	t.Run("未アーカイブならアーカイブできる", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Archive", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewArchiveProjectUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("アーカイブ済みなら 400", func(t *testing.T) {
		repo := new(mockProjectRepo)
		archived := ownedProject(1, 5)
		archived.IsArchived = true
		repo.On("FindByID", mock.Anything, uint(1)).Return(archived, nil)
		uc := usecase.NewArchiveProjectUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrBadRequest)
		repo.AssertNotCalled(t, "Archive", mock.Anything, mock.Anything)
	})

	t.Run("他人のプロジェクトは 403", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewArchiveProjectUseCase(repo)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1, 9), domain.ErrForbidden)
		repo.AssertNotCalled(t, "Archive", mock.Anything, mock.Anything)
	})
}

func TestUnarchiveProjectUseCase_Execute(t *testing.T) {
	t.Run("アーカイブ済みなら解除できる", func(t *testing.T) {
		repo := new(mockProjectRepo)
		archived := ownedProject(1, 5)
		archived.IsArchived = true
		repo.On("FindByID", mock.Anything, uint(1)).Return(archived, nil)
		repo.On("Unarchive", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewUnarchiveProjectUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("アーカイブされていなければ 400", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewUnarchiveProjectUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrBadRequest)
		repo.AssertNotCalled(t, "Unarchive", mock.Anything, mock.Anything)
	})
}

func TestDeleteProjectUseCase_Execute(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		uc := usecase.NewDeleteProjectUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("他人のプロジェクトは 403", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(ownedProject(1, 5), nil)
		uc := usecase.NewDeleteProjectUseCase(repo)

		assert.ErrorIs(t, uc.Execute(context.Background(), 1, 9), domain.ErrForbidden)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("不在は 404 を返す", func(t *testing.T) {
		repo := new(mockProjectRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewDeleteProjectUseCase(repo)

		err := uc.Execute(context.Background(), 1, 5)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}
