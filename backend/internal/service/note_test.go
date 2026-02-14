package service

import (
	"testing"

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

func (m *MockNoteRepository) FindByFolderID(folderID uint) ([]model.Note, error) {
	args := m.Called(folderID)
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

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, note, result)
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

	note := &model.Note{
		ID:      1,
		UserID:  1,
		Title:   "更新後タイトル",
		Content: "更新後内容",
	}

	repo.On("Update", note).Return(nil)

	err := svc.Update(note)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestNoteService_Delete(t *testing.T) {
	svc, repo := newTestNoteService()

	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
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
