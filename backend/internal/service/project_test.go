package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestProjectService はProjectServiceのテスト用インスタンスを生成するヘルパー。
func newTestProjectService() (*ProjectService, *MockProjectRepository) {
	repo := new(MockProjectRepository)
	svc := NewProjectService(repo)
	return svc, repo
}

// ============================================================
// プロジェクト更新テスト
// ============================================================

func TestProjectUpdate_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{Title: "Old", Description: "Old Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Project{Title: "New", TechStack: "Go, React"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New", result.Title)
	assert.Equal(t, "Go, React", result.TechStack)
	assert.Equal(t, "Old Desc", result.Description)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.Project{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.Project{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_ValidationError(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, Title: "Old Title"}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.Project{DemoURL: "not-a-valid-url"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Update")
}

// ============================================================
// 注目プロジェクト設定テスト
// ============================================================

func TestProjectUpdateFeatured_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, Featured: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.UpdateFeatured(1, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.Featured)
	repo.AssertExpectations(t)
}

func TestProjectUpdateFeatured_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.UpdateFeatured(1, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// プロジェクト削除テスト
// ============================================================

func TestProjectDelete_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProjectDelete_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestProjectDelete_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Delete(999, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// プロジェクト作成テスト
// ============================================================

func TestProjectCreate_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	project := &model.Project{
		Title:       "My Project",
		Description: "A great project",
		UserID:      1,
	}

	repo.On("Create", project).Return(nil)

	err := svc.Create(project)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProjectCreate_ValidationError(t *testing.T) {
	svc, _ := newTestProjectService()

	project := &model.Project{
		Title:       "",
		Description: "A great project",
		UserID:      1,
	}

	err := svc.Create(project)
	assert.Error(t, err)
}

func TestProjectCreate_RepoError(t *testing.T) {
	svc, repo := newTestProjectService()

	project := &model.Project{
		Title:       "My Project",
		Description: "A great project",
		UserID:      1,
	}

	repo.On("Create", project).Return(errors.New("db error"))

	err := svc.Create(project)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// プロジェクト取得テスト
// ============================================================

func TestProjectGetByID_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	expected := &model.Project{Title: "Test", UserID: 1}
	expected.ID = 1

	repo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Title)
	repo.AssertExpectations(t)
}

func TestProjectGetByID_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectGetByID_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	project := &model.Project{Title: "Test", UserID: 1}
	project.ID = 1
	repo.On("FindByID", uint(1)).Return(project, nil)

	result, err := svc.GetByID(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectGetByUserID_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	expected := []model.Project{
		{Title: "Project 1", UserID: 1},
		{Title: "Project 2", UserID: 1},
	}

	repo.On("FindByUserID", uint(1), 20, 0).Return(expected, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestProjectGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByUserID", uint(1), 20, 0).Return([]model.Project{}, int64(0), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

func TestProjectGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByUserID", uint(1), 10, 10).Return([]model.Project{}, int64(15), nil)

	result, total, err := svc.GetByUserID(1, 10, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(15), total)
	repo.AssertExpectations(t)
}

func TestProjectGetFeaturedByUserID_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	expected := []model.Project{
		{Title: "Featured 1", UserID: 1, Featured: true},
	}

	repo.On("FindFeaturedByUserID", uint(1)).Return(expected, nil)

	result, err := svc.GetFeaturedByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].Featured)
	repo.AssertExpectations(t)
}

func TestProjectGetFeaturedByUserID_Empty(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindFeaturedByUserID", uint(1)).Return([]model.Project{}, nil)

	result, err := svc.GetFeaturedByUserID(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestProjectGetAll_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	expected := []model.Project{
		{Title: "Project 1"},
		{Title: "Project 2"},
	}

	repo.On("FindAll", 10, 0).Return(expected, int64(2), nil)

	result, total, err := svc.GetAll(10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

// ============================================================
// 更新・削除エラーケーステスト
// ============================================================

func TestProjectUpdate_UpdateError(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{Title: "Old", Description: "Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))

	updates := &model.Project{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectUpdateFeatured_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.UpdateFeatured(999, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_AllFields(t *testing.T) {
	svc, repo := newTestProjectService()

	now := time.Now()
	repoID := uint(42)
	existing := &model.Project{Title: "Old", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Project{
		Title:        "New Title",
		Description:  "New Description",
		TechStack:    "Go, React, PostgreSQL",
		DemoURL:      "https://demo.example.com",
		GithubURL:    "https://github.com/example/repo",
		ImageURL:     "https://example.com/image.jpg",
		Role:         "Backend Developer",
		StartDate:    &now,
		EndDate:      &now,
		GithubRepoID: &repoID,
	}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Description", result.Description)
	assert.Equal(t, "Go, React, PostgreSQL", result.TechStack)
	assert.Equal(t, "https://demo.example.com", result.DemoURL)
	assert.Equal(t, "https://github.com/example/repo", result.GithubURL)
	assert.Equal(t, "https://example.com/image.jpg", result.ImageURL)
	assert.Equal(t, "Backend Developer", result.Role)
	repo.AssertExpectations(t)
}

func TestProjectUpdateFeatured_UpdateError(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, Featured: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))

	result, err := svc.UpdateFeatured(1, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 空白バイパステスト
// ============================================================

func TestProjectUpdate_WhitespaceTechStack(t *testing.T) {
	svc, repo := newTestProjectService()
	existing := &model.Project{Title: "Project", Description: "Desc", TechStack: "Go, React", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Project{TechStack: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Go, React", result.TechStack)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_WhitespaceImageURL(t *testing.T) {
	svc, repo := newTestProjectService()
	existing := &model.Project{Title: "Project", Description: "Desc", ImageURL: "https://example.com/img.png", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Project{ImageURL: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/img.png", result.ImageURL)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_WhitespaceRole(t *testing.T) {
	svc, repo := newTestProjectService()
	existing := &model.Project{Title: "Project", Description: "Desc", Role: "Backend Engineer", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Project{Role: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Backend Engineer", result.Role)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_TrimsPaddedTitle(t *testing.T) {
	svc, repo := newTestProjectService()
	existing := &model.Project{Title: "Original", Description: "Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Project{Title: "  New Title  "})
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
}

// ============================================================
// アーカイブテスト
// ============================================================

func TestProjectArchive_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, IsArchived: false}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Archive", uint(1)).Return(nil)

	err := svc.Archive(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProjectArchive_AlreadyArchived(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, IsArchived: true}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Archive(1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestProjectArchive_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 2, IsArchived: false}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Archive(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestProjectArchive_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Archive(999, 1)
	assert.Error(t, err)
}

func TestProjectUnarchive_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, IsArchived: true}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Unarchive", uint(1)).Return(nil)

	err := svc.Unarchive(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProjectUnarchive_NotArchived(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, IsArchived: false}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Unarchive(1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestProjectUnarchive_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 2, IsArchived: true}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Unarchive(1, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestProjectGetArchivedByUserID_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	expected := []model.Project{
		{Title: "Archived 1", UserID: 1, IsArchived: true},
	}
	repo.On("FindArchivedByUserID", uint(1), 20, 0).Return(expected, int64(1), nil)

	result, total, err := svc.GetArchivedByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}
