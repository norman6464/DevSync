package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NoteVersionService はノートバージョン履歴のビジネスロジックを提供する。
type NoteVersionService struct {
	noteRepo    repository.NoteRepositoryInterface
	versionRepo repository.NoteVersionRepositoryInterface
}

// NewNoteVersionService は新しいNoteVersionServiceインスタンスを生成する。
func NewNoteVersionService(noteRepo repository.NoteRepositoryInterface, versionRepo repository.NoteVersionRepositoryInterface) *NoteVersionService {
	return &NoteVersionService{noteRepo: noteRepo, versionRepo: versionRepo}
}

// SaveVersion はノートの現在状態をバージョンとして保存する。
func (s *NoteVersionService) SaveVersion(noteID, userID uint) error {
	note, err := s.findAndCheckOwnership(noteID, userID)
	if err != nil {
		return err
	}

	latestNum, err := s.versionRepo.GetLatestVersionNumber(noteID)
	if err != nil {
		return err
	}

	version := &model.NoteVersion{
		NoteID:        noteID,
		VersionNumber: latestNum + 1,
		Title:         note.Title,
		Content:       note.Content,
		Tags:          note.Tags,
	}
	return s.versionRepo.Create(version)
}

// GetVersions はノートのバージョン履歴一覧を取得する。
func (s *NoteVersionService) GetVersions(noteID, userID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	if _, err := s.findAndCheckOwnership(noteID, userID); err != nil {
		return nil, 0, err
	}
	return s.versionRepo.FindByNoteID(noteID, limit, offset)
}

// GetVersion は指定バージョンの詳細を取得する。
func (s *NoteVersionService) GetVersion(noteID, versionID, userID uint) (*model.NoteVersion, error) {
	if _, err := s.findAndCheckOwnership(noteID, userID); err != nil {
		return nil, err
	}

	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", err)
	}

	if version.NoteID != noteID {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", nil)
	}

	return version, nil
}

// RestoreVersion は指定バージョンの内容でノートを復元する。
// 復元前に現在の状態をバージョンとして保存する。
func (s *NoteVersionService) RestoreVersion(noteID, versionID, userID uint) (*model.Note, error) {
	note, err := s.findAndCheckOwnership(noteID, userID)
	if err != nil {
		return nil, err
	}

	version, err := s.versionRepo.FindByID(versionID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", err)
	}

	if version.NoteID != noteID {
		return nil, domain.NewError(domain.ErrCodeNotFound, "バージョンが見つかりません", nil)
	}

	// 復元前に現在の状態をバージョン保存
	latestNum, err := s.versionRepo.GetLatestVersionNumber(noteID)
	if err != nil {
		return nil, err
	}
	currentVersion := &model.NoteVersion{
		NoteID:        noteID,
		VersionNumber: latestNum + 1,
		Title:         note.Title,
		Content:       note.Content,
		Tags:          note.Tags,
	}
	if err := s.versionRepo.Create(currentVersion); err != nil {
		return nil, err
	}

	// ノートを復元
	note.Title = version.Title
	note.Content = version.Content
	note.Tags = version.Tags
	if err := s.noteRepo.Update(note); err != nil {
		return nil, err
	}

	return note, nil
}

// findAndCheckOwnership はノートの所有権を検証する共通ヘルパー。
func (s *NoteVersionService) findAndCheckOwnership(noteID, userID uint) (*model.Note, error) {
	return checkOwnership(s.noteRepo.FindByID, noteID, userID, func(n *model.Note) uint { return n.UserID })
}
