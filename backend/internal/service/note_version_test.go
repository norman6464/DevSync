package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockNoteVersionRepository
// ============================================================

type MockNoteVersionRepository struct {
	mock.Mock
}

var _ repository.NoteVersionRepositoryInterface = (*MockNoteVersionRepository)(nil)

func (m *MockNoteVersionRepository) Create(version *model.NoteVersion) error {
	return m.Called(version).Error(0)
}

func (m *MockNoteVersionRepository) FindByNoteID(noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	args := m.Called(noteID, limit, offset)
	return args.Get(0).([]model.NoteVersion), args.Get(1).(int64), args.Error(2)
}

func (m *MockNoteVersionRepository) FindByID(id uint) (*model.NoteVersion, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*model.NoteVersion), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteVersionRepository) GetLatestVersionNumber(noteID uint) (int, error) {
	args := m.Called(noteID)
	return args.Int(0), args.Error(1)
}

// MockNoteRepository2 はNoteVersionテスト用のNoteRepositoryモック。
// note_test.goのMockNoteRepositoryと同パッケージ内で重複を避ける。
type MockNoteRepository2 struct {
	mock.Mock
}

func (m *MockNoteRepository2) FindByID(id uint) (*model.Note, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteRepository2) Update(note *model.Note) error {
	return m.Called(note).Error(0)
}

// 以下はインターフェース充足のためのスタブ
func (m *MockNoteRepository2) Create(note *model.Note) error          { return nil }
func (m *MockNoteRepository2) FindByUserID(uint, int, int) ([]model.Note, error) {
	return nil, nil
}
func (m *MockNoteRepository2) FindByFolderID(uint, uint) ([]model.Note, error) { return nil, nil }
func (m *MockNoteRepository2) Delete(uint) error                               { return nil }
func (m *MockNoteRepository2) Search(uint, string, int, int) ([]model.Note, int64, error) {
	return nil, 0, nil
}
func (m *MockNoteRepository2) CountByUserID(uint) (int64, error)             { return 0, nil }
func (m *MockNoteRepository2) ToggleFavorite(uint) error                     { return nil }
func (m *MockNoteRepository2) FindFavorites(uint, int, int) ([]model.Note, error) {
	return nil, nil
}
func (m *MockNoteRepository2) CountFavoritesByUserID(uint) (int64, error)    { return 0, nil }
func (m *MockNoteRepository2) Archive(uint) error                           { return nil }
func (m *MockNoteRepository2) Unarchive(uint) error                         { return nil }
func (m *MockNoteRepository2) FindArchived(uint, int, int) ([]model.Note, error) {
	return nil, nil
}
func (m *MockNoteRepository2) CountArchivedByUserID(uint) (int64, error) { return 0, nil }

func newTestNoteVersionService() (*NoteVersionService, *MockNoteRepository2, *MockNoteVersionRepository) {
	noteRepo := new(MockNoteRepository2)
	versionRepo := new(MockNoteVersionRepository)
	svc := NewNoteVersionService(noteRepo, versionRepo)
	return svc, noteRepo, versionRepo
}

// ============================================================
// SaveVersion テスト
// ============================================================

func TestNoteVersionSaveVersion_Success(t *testing.T) {
	svc, noteRepo, versionRepo := newTestNoteVersionService()

	note := &model.Note{ID: 1, UserID: 10, Title: "Title", Content: "Content", Tags: "go"}
	noteRepo.On("FindByID", uint(1)).Return(note, nil)
	versionRepo.On("GetLatestVersionNumber", uint(1)).Return(2, nil)
	versionRepo.On("Create", mock.MatchedBy(func(v *model.NoteVersion) bool {
		return v.NoteID == 1 && v.VersionNumber == 3 && v.Title == "Title" && v.Content == "Content" && v.Tags == "go"
	})).Return(nil)

	err := svc.SaveVersion(1, 10)
	assert.NoError(t, err)
	versionRepo.AssertExpectations(t)
}

func TestNoteVersionSaveVersion_NoteNotFound(t *testing.T) {
	svc, noteRepo, _ := newTestNoteVersionService()
	noteRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	err := svc.SaveVersion(1, 10)
	assert.Error(t, err)
}

func TestNoteVersionSaveVersion_Forbidden(t *testing.T) {
	svc, noteRepo, _ := newTestNoteVersionService()
	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)

	err := svc.SaveVersion(1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FORBIDDEN")
}

// ============================================================
// GetVersions テスト
// ============================================================

func TestNoteVersionGetVersions_Success(t *testing.T) {
	svc, noteRepo, versionRepo := newTestNoteVersionService()

	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 10}, nil)
	versions := []model.NoteVersion{
		{ID: 2, NoteID: 1, VersionNumber: 2, Title: "V2", CreatedAt: time.Now()},
		{ID: 1, NoteID: 1, VersionNumber: 1, Title: "V1", CreatedAt: time.Now()},
	}
	versionRepo.On("FindByNoteID", uint(1), 20, 0).Return(versions, int64(2), nil)

	result, total, err := svc.GetVersions(1, 10, 20, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestNoteVersionGetVersions_Forbidden(t *testing.T) {
	svc, noteRepo, _ := newTestNoteVersionService()
	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)

	_, _, err := svc.GetVersions(1, 10, 20, 0)
	assert.Error(t, err)
}

// ============================================================
// GetVersion テスト
// ============================================================

func TestNoteVersionGetVersion_Success(t *testing.T) {
	svc, noteRepo, versionRepo := newTestNoteVersionService()

	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 10}, nil)
	version := &model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 3, Title: "V3"}
	versionRepo.On("FindByID", uint(5)).Return(version, nil)

	result, err := svc.GetVersion(1, 5, 10)
	assert.NoError(t, err)
	assert.Equal(t, uint(5), result.ID)
}

func TestNoteVersionGetVersion_WrongNote(t *testing.T) {
	svc, noteRepo, versionRepo := newTestNoteVersionService()

	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 10}, nil)
	// バージョンは別のノートに属する
	versionRepo.On("FindByID", uint(5)).Return(&model.NoteVersion{ID: 5, NoteID: 999}, nil)

	_, err := svc.GetVersion(1, 5, 10)
	assert.Error(t, err)
}

// ============================================================
// RestoreVersion テスト
// ============================================================

func TestNoteVersionRestore_Success(t *testing.T) {
	svc, noteRepo, versionRepo := newTestNoteVersionService()

	note := &model.Note{ID: 1, UserID: 10, Title: "Current", Content: "Current Content"}
	noteRepo.On("FindByID", uint(1)).Return(note, nil)

	version := &model.NoteVersion{ID: 5, NoteID: 1, VersionNumber: 2, Title: "Old Title", Content: "Old Content", Tags: "tag1"}
	versionRepo.On("FindByID", uint(5)).Return(version, nil)

	// 復元前に現在の状態をバージョン保存
	versionRepo.On("GetLatestVersionNumber", uint(1)).Return(3, nil)
	versionRepo.On("Create", mock.MatchedBy(func(v *model.NoteVersion) bool {
		return v.NoteID == 1 && v.VersionNumber == 4 && v.Title == "Current"
	})).Return(nil)

	noteRepo.On("Update", mock.MatchedBy(func(n *model.Note) bool {
		return n.Title == "Old Title" && n.Content == "Old Content" && n.Tags == "tag1"
	})).Return(nil)

	result, err := svc.RestoreVersion(1, 5, 10)
	assert.NoError(t, err)
	assert.Equal(t, "Old Title", result.Title)
	assert.Equal(t, "Old Content", result.Content)
}

func TestNoteVersionRestore_Forbidden(t *testing.T) {
	svc, noteRepo, _ := newTestNoteVersionService()
	noteRepo.On("FindByID", uint(1)).Return(&model.Note{ID: 1, UserID: 99}, nil)

	_, err := svc.RestoreVersion(1, 5, 10)
	assert.Error(t, err)
}
