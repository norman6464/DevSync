package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NoteLinkRepositoryInterface はNoteLinkRepositoryのインターフェース。
type NoteLinkRepositoryInterface interface {
	Create(link *model.NoteLink) error
	FindBySourceNoteID(sourceNoteID uint) ([]model.NoteLink, error)
	FindByTargetNoteID(targetNoteID uint) ([]model.NoteLink, error)
	Delete(sourceNoteID, targetNoteID uint) error
	Exists(sourceNoteID, targetNoteID uint) (bool, error)
}

// NoteLinkService はノート間リンクのビジネスロジック。
type NoteLinkService struct {
	repo     NoteLinkRepositoryInterface
	noteRepo repository.NoteRepositoryInterface
}

// NewNoteLinkService は新しいNoteLinkServiceインスタンスを生成する。
func NewNoteLinkService(repo NoteLinkRepositoryInterface, noteRepo repository.NoteRepositoryInterface) *NoteLinkService {
	return &NoteLinkService{
		repo:     repo,
		noteRepo: noteRepo,
	}
}

// findAndCheckSourceNoteOwnership はソースノートの所有権を検証する。
func (s *NoteLinkService) findAndCheckSourceNoteOwnership(sourceNoteID, userID uint) error {
	sourceNote, err := s.noteRepo.FindByID(sourceNoteID)
	if err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "ソースノートが見つかりません", err)
	}
	if sourceNote.UserID != userID {
		return domain.NewError(domain.ErrCodeForbidden, "この操作を行う権限がありません", nil)
	}
	return nil
}

// CreateLink は新しいリンクを作成する。
// ソースノートの所有権を検証した後、リンクを作成する。
func (s *NoteLinkService) CreateLink(sourceNoteID, targetNoteID, userID uint) error {
	// 同じノートへのリンクは作成できない
	if sourceNoteID == targetNoteID {
		return domain.NewError(domain.ErrCodeValidation, "同じノートへのリンクは作成できません", nil)
	}

	// ソースノートの所有権を確認
	if err := s.findAndCheckSourceNoteOwnership(sourceNoteID, userID); err != nil {
		return err
	}

	// ターゲットノートが存在するかチェック
	if _, err := s.noteRepo.FindByID(targetNoteID); err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "リンク先のノートが見つかりません", err)
	}

	// 既にリンクが存在するかチェック
	exists, err := s.repo.Exists(sourceNoteID, targetNoteID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeValidation, "このリンクは既に存在します", nil)
	}

	link := &model.NoteLink{
		SourceNoteID: sourceNoteID,
		TargetNoteID: targetNoteID,
	}

	return s.repo.Create(link)
}

// GetLinks は指定ノートからのリンク一覧を取得する。
func (s *NoteLinkService) GetLinks(sourceNoteID uint) ([]model.NoteLink, error) {
	return s.repo.FindBySourceNoteID(sourceNoteID)
}

// GetBacklinks は指定ノートへのリンク一覧（バックリンク）を取得する。
func (s *NoteLinkService) GetBacklinks(targetNoteID uint) ([]model.NoteLink, error) {
	return s.repo.FindByTargetNoteID(targetNoteID)
}

// DeleteLink はソースノートの所有権を検証した後、リンクを削除する。
func (s *NoteLinkService) DeleteLink(sourceNoteID, targetNoteID, userID uint) error {
	if err := s.findAndCheckSourceNoteOwnership(sourceNoteID, userID); err != nil {
		return err
	}
	return s.repo.Delete(sourceNoteID, targetNoteID)
}
