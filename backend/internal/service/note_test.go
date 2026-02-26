package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoteRepository はNoteRepositoryのモック実装。
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Create(note *model.Note) error {
	return m.Called(note).Error(0)
}

func (m *MockNoteRepository) FindByID(id uint) (*model.Note, error) {
	args := m.Called(id)
	if note := args.Get(0); note != nil {
		return note.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteRepository) FindByUserID(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteRepository) FindByFolderID(folderID, userID uint) ([]model.Note, error) {
	args := m.Called(folderID, userID)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(note *model.Note) error {
	return m.Called(note).Error(0)
}

func (m *MockNoteRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteRepository) Search(userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	args := m.Called(userID, query, limit, offset)
	return args.Get(0).([]model.Note), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteRepository) ToggleFavorite(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteRepository) Archive(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteRepository) Unarchive(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockNoteRepository) FindArchived(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteRepository) CountArchivedByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNoteRepository) FindFavorites(userID uint, page, limit int) ([]model.Note, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Note), args.Error(1)
}

func (m *MockNoteRepository) CountFavoritesByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// newTestNoteService はテスト用のNoteServiceを生成する。
func newTestNoteService() (*NoteService, *MockNoteRepository) {
	repo := new(MockNoteRepository)
	svc := NewNoteService(repo)
	return svc, repo
}

// ============================================================
// Create テスト
// ============================================================

func TestNoteService_Create(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{
		UserID:  1,
		Title:   "テストノート",
		Content: "これはテストです",
		Tags:    "Go,TDD",
	}

	repo.On("Create", note).Return(nil)

	err := svc.Create(note)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Create_ValidationError(t *testing.T) {
	svc, _ := newTestNoteService()

	note := &model.Note{
		UserID:  1,
		Title:   "",
		Content: "内容",
	}

	err := svc.Create(note)
	assert.Error(t, err)
}

func TestNoteService_Create_RepoError(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{
		UserID:  1,
		Title:   "テストノート",
		Content: "内容",
	}

	repo.On("Create", note).Return(errors.New("db error"))

	err := svc.Create(note)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestNoteService_GetByID(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "テストノート",
		Content: "内容",
	}

	repo.On("FindByID", uint(1)).Return(note, nil)

	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, note, result)
	repo.AssertExpectations(t)
}

func TestNoteService_GetByID_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteService_GetByID_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "他人のノート",
		Content: "内容",
	}

	repo.On("FindByID", uint(1)).Return(note, nil)

	result, err := svc.GetByID(1, 999) // 別ユーザーがアクセス
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestNoteService_GetByUserID(t *testing.T) {
	svc, repo := newTestNoteService()

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "ノート1"},
		{ID: 2, UserID: 1, Title: "ノート2"},
	}

	repo.On("FindByUserID", uint(1), 1, 20).Return(notes, nil)

	result, err := svc.GetByUserID(1, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestNoteService_Update(t *testing.T) {
	svc, repo := newTestNoteService()

	existing := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "元タイトル",
		Content: "元内容",
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.MatchedBy(func(n *model.Note) bool {
		return n.Title == "更新後タイトル" && n.Content == "更新後内容"
	})).Return(nil)

	result, err := svc.Update(1, 1, "更新後タイトル", "更新後内容", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "更新後タイトル", result.Title)
	assert.Equal(t, "更新後内容", result.Content)
	repo.AssertExpectations(t)
}

func TestNoteService_Update_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	existing := &model.Note{ID: 1, UserID: 1, Title: "元タイトル"}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 999, "更新", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteService_Update_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.Update(99, 1, "更新", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteService_Delete(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Delete_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)

	err := svc.Delete(1, 999)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Delete_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Delete(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// Search テスト
// ============================================================

func TestNoteService_Search(t *testing.T) {
	svc, repo := newTestNoteService()

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "Go学習", Content: "Goを学習中"},
	}

	repo.On("Search", uint(1), "Go", 20, 0).Return(notes, int64(1), nil)

	result, total, err := svc.Search(1, "Go", 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(result))
	repo.AssertExpectations(t)
}

// ============================================================
// GetByFolderID テスト
// ============================================================

func TestNoteService_GetByFolderID(t *testing.T) {
	svc, repo := newTestNoteService()

	notes := []model.Note{
		{ID: 1, UserID: 1, Title: "フォルダ内ノート1"},
		{ID: 2, UserID: 1, Title: "フォルダ内ノート2"},
	}

	repo.On("FindByFolderID", uint(5), uint(1)).Return(notes, nil)

	result, err := svc.GetByFolderID(5, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

// ============================================================
// CountByUserID テスト
// ============================================================

func TestNoteService_CountByUserID(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("CountByUserID", uint(1)).Return(int64(10), nil)

	count, err := svc.CountByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)
	repo.AssertExpectations(t)
}

// ============================================================
// ToggleFavorite テスト
// ============================================================

func TestNoteService_ToggleFavorite(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)
	repo.On("ToggleFavorite", uint(1)).Return(nil)

	err := svc.ToggleFavorite(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_ToggleFavorite_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)

	err := svc.ToggleFavorite(1, 999)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_ToggleFavorite_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.ToggleFavorite(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// Archive / Unarchive テスト
// ============================================================

func TestNoteService_Archive(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)
	repo.On("Archive", uint(1)).Return(nil)

	err := svc.Archive(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Archive_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)

	err := svc.Archive(1, 999)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Archive_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Archive(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Unarchive(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)
	repo.On("Unarchive", uint(1)).Return(nil)

	err := svc.Unarchive(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Unarchive_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 1, Title: "ノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)

	err := svc.Unarchive(1, 999)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNoteService_Unarchive_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Unarchive(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetArchived テスト
// ============================================================

func TestNoteService_GetArchived(t *testing.T) {
	svc, repo := newTestNoteService()

	notes := []model.Note{
		{ID: 3, UserID: 1, Title: "アーカイブノート", IsArchived: true},
	}

	repo.On("FindArchived", uint(1), 1, 20).Return(notes, nil)

	result, err := svc.GetArchived(1, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

// ============================================================
// CountArchivedByUserID テスト
// ============================================================

func TestNoteService_CountArchivedByUserID(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("CountArchivedByUserID", uint(1)).Return(int64(3), nil)

	count, err := svc.CountArchivedByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	repo.AssertExpectations(t)
}

// ============================================================
// Duplicate テスト
// ============================================================

func TestNoteService_Duplicate(t *testing.T) {
	svc, repo := newTestNoteService()

	original := &model.Note{
		ID:       1,
		UserID:   1,
		Title:    "元ノート",
		Content:  "元の内容",
		Tags:     "Go,TDD",
		FolderID: nil,
	}

	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("Create", mock.MatchedBy(func(n *model.Note) bool {
		return n.Title == "元ノート (コピー)" &&
			n.Content == "元の内容" &&
			n.Tags == "Go,TDD" &&
			n.UserID == 1 &&
			!n.IsFavorite &&
			!n.IsArchived
	})).Return(nil)

	result, err := svc.Duplicate(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "元ノート (コピー)", result.Title)
	assert.Equal(t, "元の内容", result.Content)
	assert.False(t, result.IsFavorite)
	assert.False(t, result.IsArchived)
	repo.AssertExpectations(t)
}

func TestNoteService_Duplicate_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	original := &model.Note{ID: 1, UserID: 1, Title: "元ノート"}
	repo.On("FindByID", uint(1)).Return(original, nil)

	result, err := svc.Duplicate(1, 999)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteService_Duplicate_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.Duplicate(99, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestNoteService_Update_RepoError(t *testing.T) {
	svc, repo := newTestNoteService()
	existing := &model.Note{ID: 1, UserID: 1, Title: "元タイトル", Content: "内容"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything).Return(errors.New("db error"))
	result, err := svc.Update(1, 1, "新タイトル", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteService_Duplicate_RepoError(t *testing.T) {
	svc, repo := newTestNoteService()
	original := &model.Note{ID: 1, UserID: 1, Title: "元ノート", Content: "内容", Tags: "Go"}
	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("Create", mock.Anything).Return(errors.New("db error"))
	result, err := svc.Duplicate(1, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteService_Update_WithAllFields(t *testing.T) {
	svc, repo := newTestNoteService()
	folderID := uint(5)
	existing := &model.Note{ID: 1, UserID: 1, Title: "元", Content: "元内容", Tags: "Go"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything).Return(nil)
	result, err := svc.Update(1, 1, "新タイトル", "新内容", "React,Go", &folderID)
	assert.NoError(t, err)
	assert.Equal(t, "新タイトル", result.Title)
	assert.Equal(t, "新内容", result.Content)
	assert.Equal(t, "React,Go", result.Tags)
	assert.Equal(t, &folderID, result.FolderID)
}

func TestNoteService_Update_ValidationError(t *testing.T) {
	svc, repo := newTestNoteService()
	existing := &model.Note{ID: 1, UserID: 1, Title: "元"}

	repo.On("FindByID", uint(1)).Return(existing, nil)

	// 201文字のタイトル → ValidateUpdateNote でバリデーションエラー
	longTitle := strings.Repeat("あ", 201)
	result, err := svc.Update(1, 1, longTitle, "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNoteService_Duplicate_ValidationError(t *testing.T) {
	svc, repo := newTestNoteService()

	// 198文字のタイトル → " (コピー)" (12バイト) 付加後 210バイト > 200 → バリデーションエラー
	longTitle := strings.Repeat("a", 198)
	existing := &model.Note{ID: 1, UserID: 1, Title: longTitle}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Duplicate(1, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================================
// GetFavorites テスト
// ============================================================

func TestNoteService_GetFavorites_Success(t *testing.T) {
	svc, repo := newTestNoteService()

	favorites := []model.Note{
		{Title: "お気に入りノート1", UserID: 1, IsFavorite: true},
		{Title: "お気に入りノート2", UserID: 1, IsFavorite: true},
	}
	repo.On("FindFavorites", uint(1), 1, 10).Return(favorites, nil)

	result, err := svc.GetFavorites(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.True(t, result[0].IsFavorite)
	repo.AssertExpectations(t)
}

func TestNoteService_GetFavorites_Empty(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindFavorites", uint(1), 1, 10).Return([]model.Note{}, nil)

	result, err := svc.GetFavorites(1, 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestNoteService_GetFavorites_RepoError(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindFavorites", uint(1), 1, 10).Return([]model.Note{}, errors.New("db error"))

	result, err := svc.GetFavorites(1, 1, 10)
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Update 空白バイパス テスト
// ============================================================

func TestNoteService_Update_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestNoteService()

	existing := &model.Note{ID: 1, UserID: 1, Title: "元のタイトル", Content: "元の内容"}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 1, "   \t\n  ", "", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "タイトルは空白のみにできません")
}

func TestNoteService_Update_WhitespaceContent(t *testing.T) {
	svc, repo := newTestNoteService()

	existing := &model.Note{ID: 1, UserID: 1, Title: "タイトル", Content: "元の内容"}
	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.Update(1, 1, "", "   \t\n  ", "", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "本文は空白のみにできません")
}

// ============================================================
// CountFavoritesByUserID テスト
// ============================================================

func TestNoteService_CountFavoritesByUserID_Success(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("CountFavoritesByUserID", uint(1)).Return(int64(5), nil)

	count, err := svc.CountFavoritesByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	repo.AssertExpectations(t)
}

func TestNoteService_CountFavoritesByUserID_Zero(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("CountFavoritesByUserID", uint(1)).Return(int64(0), nil)

	count, err := svc.CountFavoritesByUserID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
	repo.AssertExpectations(t)
}

func TestNoteService_CountFavoritesByUserID_RepoError(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("CountFavoritesByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	count, err := svc.CountFavoritesByUserID(1)
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	repo.AssertExpectations(t)
}

// ============================================================
// ExportMarkdown テスト
// ============================================================

func TestNoteService_ExportMarkdown_Success(t *testing.T) {
	svc, repo := newTestNoteService()

	now := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	note := &model.Note{
		ID:        1,
		UserID:    1,
		Title:     "Go学習メモ",
		Content:   "## インターフェース\n\nGoのインターフェースは暗黙的に実装される。",
		Tags:      "Go,学習",
		CreatedAt: now,
		UpdatedAt: now,
	}
	repo.On("FindByID", uint(1)).Return(note, nil)

	data, title, err := svc.ExportMarkdown(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Go学習メモ", title)
	content := string(data)
	assert.Contains(t, content, "# Go学習メモ")
	assert.Contains(t, content, "**Tags:** Go,学習")
	assert.Contains(t, content, "**Created:** 2025-06-15 10:30")
	assert.Contains(t, content, "---")
	assert.Contains(t, content, "## インターフェース")
	repo.AssertExpectations(t)
}

func TestNoteService_ExportMarkdown_NoTags(t *testing.T) {
	svc, repo := newTestNoteService()

	now := time.Now()
	note := &model.Note{
		ID: 1, UserID: 1, Title: "タグなしノート", Content: "内容", Tags: "",
		CreatedAt: now, UpdatedAt: now,
	}
	repo.On("FindByID", uint(1)).Return(note, nil)

	data, _, err := svc.ExportMarkdown(1, 1)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "**Tags:**")
}

func TestNoteService_ExportMarkdown_NotFound(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, _, err := svc.ExportMarkdown(999, 1)
	assert.Error(t, err)
}

func TestNoteService_ExportMarkdown_Forbidden(t *testing.T) {
	svc, repo := newTestNoteService()

	note := &model.Note{ID: 1, UserID: 99, Title: "他人のノート"}
	repo.On("FindByID", uint(1)).Return(note, nil)

	_, _, err := svc.ExportMarkdown(1, 1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)
}
