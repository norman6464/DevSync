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

// CreateLink は新しいリンクを作成する。
func (s *NoteLinkService) CreateLink(sourceNoteID, targetNoteID uint) error {
	// 同じノートへのリンクは作成できない
	if sourceNoteID == targetNoteID {
		return domain.NewError(domain.ErrCodeValidation, "同じノートへのリンクは作成できません", nil)
	}

	// ターゲットノートが存在するかチェック
	_, err := s.noteRepo.FindByID(targetNoteID)
	if err != nil {
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

// DeleteLink はリンクを削除する。
func (s *NoteLinkService) DeleteLink(sourceNoteID, targetNoteID uint) error {
	return s.repo.Delete(sourceNoteID, targetNoteID)
}
