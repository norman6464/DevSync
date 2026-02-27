package service

import (
	"fmt"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
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

// GetByID は指定IDのノートを所有権検証付きで取得する。
func (s *NoteService) GetByID(id, userID uint) (*model.Note, error) {
	return s.findAndCheckOwnership(id, userID)
}

// GetByUserID は指定ユーザーのノート一覧を取得する。
func (s *NoteService) GetByUserID(userID uint, page, limit int) ([]model.Note, error) {
	return s.repo.FindByUserID(userID, page, limit)
}

// GetByFolderID は指定フォルダ内のノート一覧を所有権検証付きで取得する。
func (s *NoteService) GetByFolderID(folderID, userID uint) ([]model.Note, error) {
	return s.repo.FindByFolderID(folderID, userID)
}

// findAndCheckOwnership は指定IDのノートを取得し、所有権を検証する共通ヘルパー。
// 取得失敗または所有権不一致の場合はエラーを返す。
func (s *NoteService) findAndCheckOwnership(id, userID uint) (*model.Note, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(n *model.Note) uint { return n.UserID })
}

// Update は所有権を検証した後、ノートを更新する。
func (s *NoteService) Update(id, userID uint, title, content, tags string, folderID *uint) (*model.Note, error) {
	note, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		t := strings.TrimSpace(title)
		if t == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "タイトルは空白のみにできません", nil)
		}
		note.Title = t
	}
	if content != "" {
		c := strings.TrimSpace(content)
		if c == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "本文は空白のみにできません", nil)
		}
		note.Content = c
	}
	if tags != "" {
		note.Tags = strings.TrimSpace(tags)
	}
	if folderID != nil {
		note.FolderID = folderID
	}

	v := validator.NewNoteValidator()
	if err := v.ValidateUpdateNote(note.Title, note.Content, note.Tags); err != nil {
		return nil, err
	}

	if err := s.repo.Update(note); err != nil {
		return nil, err
	}
	return note, nil
}

// Delete は所有権を検証した後、ノートを削除する。
func (s *NoteService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
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

// ToggleFavorite は所有権を検証した後、ノートのお気に入り状態を切り替える。
func (s *NoteService) ToggleFavorite(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.ToggleFavorite(id)
}

// Archive は所有権を検証した後、ノートをアーカイブする。
func (s *NoteService) Archive(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Archive(id)
}

// Unarchive は所有権を検証した後、ノートのアーカイブを解除する。
func (s *NoteService) Unarchive(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Unarchive(id)
}

// GetFavorites は指定ユーザーのお気に入りノート一覧をページネーション付きで取得する。
func (s *NoteService) GetFavorites(userID uint, page, limit int) ([]model.Note, error) {
	return s.repo.FindFavorites(userID, page, limit)
}

// CountFavoritesByUserID は指定ユーザーのお気に入りノート総数を取得する。
func (s *NoteService) CountFavoritesByUserID(userID uint) (int64, error) {
	return s.repo.CountFavoritesByUserID(userID)
}

// GetArchived は指定ユーザーのアーカイブ済みノート一覧を取得する。
func (s *NoteService) GetArchived(userID uint, page, limit int) ([]model.Note, error) {
	return s.repo.FindArchived(userID, page, limit)
}

// CountArchivedByUserID は指定ユーザーのアーカイブ済みノート総数を取得する。
func (s *NoteService) CountArchivedByUserID(userID uint) (int64, error) {
	return s.repo.CountArchivedByUserID(userID)
}

// GetTags は指定ユーザーのノートで使用されているタグを重複なしで取得する。
func (s *NoteService) GetTags(userID uint) ([]string, error) {
	notes, err := s.repo.FindByUserID(userID, 1, 1000)
	if err != nil {
		return nil, err
	}
	return ExtractUniqueTags(notes), nil
}

// ExtractUniqueTags はノート一覧からユニークなタグを抽出する純粋関数。
func ExtractUniqueTags(notes []model.Note) []string {
	seen := make(map[string]bool)
	var tags []string
	for _, note := range notes {
		if note.Tags == "" {
			continue
		}
		for _, tag := range strings.Split(note.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" && !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

// ExportMarkdown はノートをMarkdown形式でエクスポートする。
// メタデータ（タグ、作成日、更新日）をヘッダーに含む。
func (s *NoteService) ExportMarkdown(id, userID uint) ([]byte, string, error) {
	note, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, "", err
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("# %s\n\n", note.Title))

	if note.Tags != "" {
		b.WriteString(fmt.Sprintf("**Tags:** %s\n", note.Tags))
	}
	b.WriteString(fmt.Sprintf("**Created:** %s\n", note.CreatedAt.Format("2006-01-02 15:04")))
	b.WriteString(fmt.Sprintf("**Updated:** %s\n", note.UpdatedAt.Format("2006-01-02 15:04")))
	b.WriteString("\n---\n\n")
	b.WriteString(note.Content)
	b.WriteString("\n")

	return []byte(b.String()), note.Title, nil
}

// Duplicate は既存のノートを複製する。
// タイトルに「(コピー)」を付与し、アーカイブ・お気に入り状態はリセットされる。
func (s *NoteService) Duplicate(id uint, userID uint) (*model.Note, error) {
	original, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	// 複製ノートを作成
	duplicate := &model.Note{
		UserID:     userID,
		Title:      original.Title + " (コピー)",
		Content:    original.Content,
		Tags:       original.Tags,
		FolderID:   original.FolderID,
		IsFavorite: false,
		IsArchived: false,
	}

	// バリデーション
	v := validator.NewNoteValidator()
	if err := v.ValidateCreateNote(duplicate.Title, duplicate.Content, duplicate.Tags); err != nil {
		return nil, err
	}

	// 保存
	if err := s.repo.Create(duplicate); err != nil {
		return nil, err
	}

	return duplicate, nil
}
