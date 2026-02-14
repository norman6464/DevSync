package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NoteService は学習ノートのビジネスロジックを提供する。
type NoteService struct {
	repo repository.NoteRepositoryInterface
}

// NewNoteService は新しいNoteServiceインスタンスを生成する。
func NewNoteService(repo repository.NoteRepositoryInterface) *NoteService {
	return &NoteService{repo: repo}
}

// Create は新しいノートを作成する。
func (s *NoteService) Create(note *model.Note) error {
	// バリデーション
	v := validator.NewNoteValidator()
	if err := v.ValidateCreateNote(note.Title, note.Content, note.Tags); err != nil {
		return err
	}

	return s.repo.Create(note)
}

// GetByID は指定IDのノートを取得する。
func (s *NoteService) GetByID(id uint) (*model.Note, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーのノート一覧を取得する。
func (s *NoteService) GetByUserID(userID uint, page, limit int) ([]model.Note, error) {
	return s.repo.FindByUserID(userID, page, limit)
}

// GetByFolderID は指定フォルダ内のノート一覧を取得する。
func (s *NoteService) GetByFolderID(folderID uint) ([]model.Note, error) {
	return s.repo.FindByFolderID(folderID)
}

// Update はノートを更新する。
func (s *NoteService) Update(note *model.Note) error {
	// バリデーション（部分更新対応）
	v := validator.NewNoteValidator()
	if err := v.ValidateUpdateNote(note.Title, note.Content, note.Tags); err != nil {
		return err
	}

	return s.repo.Update(note)
}

// Delete はノートを削除する。
func (s *NoteService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// Search はキーワードでノートを検索する。
func (s *NoteService) Search(userID uint, query string, page, limit int) ([]model.Note, int64, error) {
	offset := (page - 1) * limit
	return s.repo.Search(userID, query, limit, offset)
}

// CountByUserID は指定ユーザーのノート総数を取得する。
func (s *NoteService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}

// ToggleFavorite はノートのお気に入り状態を切り替える。
func (s *NoteService) ToggleFavorite(id uint) error {
	return s.repo.ToggleFavorite(id)
}
