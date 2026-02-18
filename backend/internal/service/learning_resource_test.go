package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestLearningResourceService はLearningResourceServiceのテスト用インスタンスを生成するヘルパー。
func newTestLearningResourceService() (*LearningResourceService, *MockLearningResourceRepository) {
	repo := new(MockLearningResourceRepository)
	svc := NewLearningResourceService(repo)
	return svc, repo
}

// ============================================================
// GetByID（可視性チェック）
// ============================================================

func TestLearningResourceGetByID_PublicResource(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Public Resource", UserID: 1, IsPublic: true}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 他人がアクセスしても公開なのでOK
	result, err := svc.GetByID(1, 999)
	assert.NoError(t, err)
	assert.Equal(t, "Public Resource", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_PrivateResourceOwner(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Private Resource", UserID: 1, IsPublic: false}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 所有者がアクセス → OK
	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Private Resource", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_PrivateResourceForbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Private Resource", UserID: 1, IsPublic: false}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 他人がアクセス → Forbidden
	result, err := svc.GetByID(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID（自分 vs 他人の可視性）
// ============================================================

func TestLearningResourceGetByUserID_Self(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{
		{Title: "Public", IsPublic: true},
		{Title: "Private", IsPublic: false},
	}
	// 自分 → includePrivate=true
	repo.On("FindByUserID", uint(1), true).Return(resources, nil)

	result, err := svc.GetByUserID(1, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByUserID_Other(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{
		{Title: "Public Only", IsPublic: true},
	}
	// 他人 → includePrivate=false
	repo.On("FindByUserID", uint(1), false).Return(resources, nil)

	result, err := svc.GetByUserID(1, 999)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習リソース更新テスト
// ============================================================

func TestLearningResourceUpdate_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{Title: "Old", Description: "Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningResource{Title: "New Title", Category: "video"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, model.ResourceCategory("video"), result.Category)
	assert.Equal(t, "Desc", result.Description) // 変更なし
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningResource{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdate_AllFields(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{Title: "Old", Description: "Old Desc", URL: "https://old.com", UserID: 1, Category: "article", Difficulty: "beginner"}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningResource{
		Title:       "New Title",
		Description: "New Desc",
		URL:         "https://new.com",
		Category:    "video",
		Difficulty:  "intermediate",
		Tags:        "go,web",
		ImageURL:    "https://img.com/pic.png",
	}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, "New Desc", result.Description)
	assert.Equal(t, "https://new.com", result.URL)
	assert.Equal(t, model.ResourceCategory("video"), result.Category)
	assert.Equal(t, model.ResourceDifficulty("intermediate"), result.Difficulty)
	assert.Equal(t, "go,web", result.Tags)
	assert.Equal(t, "https://img.com/pic.png", result.ImageURL)
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdate_RepoError(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{Title: "Old", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))

	updates := &model.LearningResource{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdate_NotFound(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.LearningResource{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 公開設定変更テスト
// ============================================================

func TestLearningResourceUpdateVisibility_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1, IsPublic: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.UpdateVisibility(1, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdateVisibility_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.UpdateVisibility(1, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習リソース削除テスト
// ============================================================

func TestLearningResourceDelete_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceDelete_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// Create（バリデーション＋リポジトリ呼び出し）
// ============================================================

func TestLearningResourceCreate_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{
		Title:       "Go入門",
		Description: "Go言語の基礎",
		URL:         "https://example.com/go",
		Category:    "article",
		Difficulty:  "beginner",
		UserID:      1,
	}

	repo.On("Create", resource).Return(nil)

	err := svc.Create(resource)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceCreate_ValidationError(t *testing.T) {
	svc, _ := newTestLearningResourceService()

	// タイトル空 → バリデーションエラー
	resource := &model.LearningResource{
		Title:       "",
		Description: "説明",
		URL:         "https://example.com",
		Category:    "article",
		Difficulty:  "beginner",
	}

	err := svc.Create(resource)
	assert.Error(t, err)
}

func TestLearningResourceCreate_RepoError(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{
		Title:       "Go入門",
		Description: "Go言語の基礎",
		URL:         "https://example.com/go",
		Category:    "article",
		Difficulty:  "beginner",
	}

	repo.On("Create", resource).Return(errors.New("db error"))

	err := svc.Create(resource)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// HasLiked
// ============================================================

func TestLearningResourceHasLiked_True(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasLiked", uint(1), uint(10)).Return(true, nil)

	liked, err := svc.HasLiked(1, 10)
	assert.NoError(t, err)
	assert.True(t, liked)
	repo.AssertExpectations(t)
}

func TestLearningResourceHasLiked_False(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasLiked", uint(1), uint(10)).Return(false, nil)

	liked, err := svc.HasLiked(1, 10)
	assert.NoError(t, err)
	assert.False(t, liked)
	repo.AssertExpectations(t)
}

func TestLearningResourceHasLiked_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasLiked", uint(1), uint(10)).Return(false, errors.New("db error"))

	liked, err := svc.HasLiked(1, 10)
	assert.Error(t, err)
	assert.False(t, liked)
	repo.AssertExpectations(t)
}

// ============================================================
// HasSaved
// ============================================================

func TestLearningResourceHasSaved_True(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasSaved", uint(1), uint(10)).Return(true, nil)

	saved, err := svc.HasSaved(1, 10)
	assert.NoError(t, err)
	assert.True(t, saved)
	repo.AssertExpectations(t)
}

func TestLearningResourceHasSaved_False(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasSaved", uint(1), uint(10)).Return(false, nil)

	saved, err := svc.HasSaved(1, 10)
	assert.NoError(t, err)
	assert.False(t, saved)
	repo.AssertExpectations(t)
}

func TestLearningResourceHasSaved_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("HasSaved", uint(1), uint(10)).Return(false, errors.New("db error"))

	saved, err := svc.HasSaved(1, 10)
	assert.Error(t, err)
	assert.False(t, saved)
	repo.AssertExpectations(t)
}

// ============================================================
// GetPublic（ページネーション・フィルタ）
// ============================================================

func TestLearningResourceGetPublic_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{
		{Title: "Resource1"},
		{Title: "Resource2"},
	}
	repo.On("FindPublic", 10, 0, "article", "beginner").Return(resources, int64(2), nil)

	result, total, err := svc.GetPublic(10, 0, "article", "beginner")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetPublic_Empty(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("FindPublic", 10, 0, "", "").Return([]model.LearningResource{}, int64(0), nil)

	result, total, err := svc.GetPublic(10, 0, "", "")
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

// ============================================================
// Search
// ============================================================

func TestLearningResourceSearch_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{{Title: "Go Tutorial"}}
	repo.On("Search", "Go", 10, 0).Return(resources, int64(1), nil)

	result, total, err := svc.Search("Go", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestLearningResourceSearch_Empty(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("Search", "nonexistent", 10, 0).Return([]model.LearningResource{}, int64(0), nil)

	result, total, err := svc.Search("nonexistent", 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

// ============================================================
// Like / Unlike
// ============================================================

func TestLearningResourceLike_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Like", uint(1), uint(10)).Return(nil)

	err := svc.Like(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceLike_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Like", uint(1), uint(10)).Return(errors.New("already liked"))

	err := svc.Like(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceUnlike_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Unlike", uint(1), uint(10)).Return(nil)

	err := svc.Unlike(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceUnlike_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Unlike", uint(1), uint(10)).Return(errors.New("not liked"))

	err := svc.Unlike(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// Save / Unsave
// ============================================================

func TestLearningResourceSave_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Save", uint(1), uint(10)).Return(nil)

	err := svc.Save(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceSave_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Save", uint(1), uint(10)).Return(errors.New("already saved"))

	err := svc.Save(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceUnsave_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Unsave", uint(1), uint(10)).Return(nil)

	err := svc.Unsave(1, 10)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceUnsave_Error(t *testing.T) {
	svc, repo := newTestLearningResourceService()
	repo.On("Unsave", uint(1), uint(10)).Return(errors.New("not saved"))

	err := svc.Unsave(1, 10)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetSavedByUserID
// ============================================================

func TestLearningResourceGetSavedByUserID_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{{Title: "Saved Resource"}}
	repo.On("FindSavedByUserID", uint(1), 10, 0).Return(resources, int64(1), nil)

	result, total, err := svc.GetSavedByUserID(1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetSavedByUserID_Empty(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("FindSavedByUserID", uint(1), 10, 0).Return([]model.LearningResource{}, int64(0), nil)

	result, total, err := svc.GetSavedByUserID(1, 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}
