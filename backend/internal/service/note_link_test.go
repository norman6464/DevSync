package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNoteLinkRepository はNoteLinkRepositoryのモック実装。
type MockNoteLinkRepository struct {
	mock.Mock
}

func (m *MockNoteLinkRepository) Create(link *model.NoteLink) error {
	return m.Called(link).Error(0)
}

func (m *MockNoteLinkRepository) FindBySourceNoteID(sourceNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(sourceNoteID)
	return args.Get(0).([]model.NoteLink), args.Error(1)
}

func (m *MockNoteLinkRepository) FindByTargetNoteID(targetNoteID uint) ([]model.NoteLink, error) {
	args := m.Called(targetNoteID)
	return args.Get(0).([]model.NoteLink), args.Error(1)
}

func (m *MockNoteLinkRepository) Delete(sourceNoteID, targetNoteID uint) error {
	return m.Called(sourceNoteID, targetNoteID).Error(0)
}

func (m *MockNoteLinkRepository) Exists(sourceNoteID, targetNoteID uint) (bool, error) {
	args := m.Called(sourceNoteID, targetNoteID)
	return args.Bool(0), args.Error(1)
}

// newTestNoteLinkService はテスト用のNoteLinkServiceを生成する。
func newTestNoteLinkService() (*NoteLinkService, *MockNoteLinkRepository, *MockNoteRepository) {
	linkRepo := new(MockNoteLinkRepository)
	noteRepo := new(MockNoteRepository)
	svc := NewNoteLinkService(linkRepo, noteRepo)
	return svc, linkRepo, noteRepo
}

// ============================================================
// CreateLink テスト
// ============================================================

func TestNoteLinkService_CreateLink(t *testing.T) {
	svc, linkRepo, noteRepo := newTestNoteLinkService()

	targetNote := &model.Note{ID: 2, UserID: 1, Title: "リンク先ノート"}

	noteRepo.On("FindByID", uint(2)).Return(targetNote, nil)
	linkRepo.On("Exists", uint(1), uint(2)).Return(false, nil)
	linkRepo.On("Create", mock.MatchedBy(func(l *model.NoteLink) bool {
		return l.SourceNoteID == 1 && l.TargetNoteID == 2
	})).Return(nil)

	err := svc.CreateLink(1, 2)
	assert.NoError(t, err)
	linkRepo.AssertExpectations(t)
	noteRepo.AssertExpectations(t)
}

func TestNoteLinkService_CreateLink_SelfLink(t *testing.T) {
	svc, _, _ := newTestNoteLinkService()

	err := svc.CreateLink(1, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "同じノートへのリンクは作成できません")
}

func TestNoteLinkService_CreateLink_AlreadyExists(t *testing.T) {
	svc, linkRepo, noteRepo := newTestNoteLinkService()

	targetNote := &model.Note{ID: 2, UserID: 1, Title: "リンク先ノート"}

	noteRepo.On("FindByID", uint(2)).Return(targetNote, nil)
	linkRepo.On("Exists", uint(1), uint(2)).Return(true, nil)

	err := svc.CreateLink(1, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "このリンクは既に存在します")
	linkRepo.AssertExpectations(t)
	noteRepo.AssertExpectations(t)
}

// ============================================================
// GetLinks テスト
// ============================================================

func TestNoteLinkService_GetLinks(t *testing.T) {
	svc, linkRepo, _ := newTestNoteLinkService()

	links := []model.NoteLink{
		{ID: 1, SourceNoteID: 1, TargetNoteID: 2},
		{ID: 2, SourceNoteID: 1, TargetNoteID: 3},
	}

	linkRepo.On("FindBySourceNoteID", uint(1)).Return(links, nil)

	result, err := svc.GetLinks(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	linkRepo.AssertExpectations(t)
}

// ============================================================
// GetBacklinks テスト
// ============================================================

func TestNoteLinkService_GetBacklinks(t *testing.T) {
	svc, linkRepo, _ := newTestNoteLinkService()

	backlinks := []model.NoteLink{
		{ID: 1, SourceNoteID: 2, TargetNoteID: 1},
		{ID: 2, SourceNoteID: 3, TargetNoteID: 1},
	}

	linkRepo.On("FindByTargetNoteID", uint(1)).Return(backlinks, nil)

	result, err := svc.GetBacklinks(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	linkRepo.AssertExpectations(t)
}

// ============================================================
// DeleteLink テスト
// ============================================================

func TestNoteLinkService_DeleteLink(t *testing.T) {
	svc, linkRepo, _ := newTestNoteLinkService()

	linkRepo.On("Delete", uint(1), uint(2)).Return(nil)

	err := svc.DeleteLink(1, 2)
	assert.NoError(t, err)
	linkRepo.AssertExpectations(t)
}
