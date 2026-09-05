package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ownerOfNote は所有権判定に使う所有者 ID の取り出し。
func ownerOfNote(n *model.Note) uint { return n.UserID }

// findOwnedNote は所有者のノートを取得する。不在なら DomainError ではないエラー、
// 所有者でなければ 403 を返す（移行前と同じステータスになる）。
func findOwnedNote(ctx context.Context, notes repository.NoteRepository, id, userID uint) (*model.Note, error) {
	return ensureOwner(ctx, notes.FindByID, id, userID, ownerOfNote)
}

// CreateNoteUseCase はノートを作成する。
type CreateNoteUseCase struct {
	notes repository.NoteRepository
}

// NewCreateNoteUseCase は CreateNoteUseCase を生成する。
func NewCreateNoteUseCase(notes repository.NoteRepository) *CreateNoteUseCase {
	return &CreateNoteUseCase{notes: notes}
}

// Execute はタイトル・本文・タグを検証したうえでノートを作成する。
func (uc *CreateNoteUseCase) Execute(ctx context.Context, note *model.Note) error {
	v := validator.NewNoteValidator()
	if err := v.ValidateCreateNote(note.Title, note.Content, note.Tags); err != nil {
		return err
	}
	return uc.notes.Create(ctx, note)
}

// GetNoteUseCase はノートを 1 件取得する。
type GetNoteUseCase struct {
	notes repository.NoteRepository
}

// NewGetNoteUseCase は GetNoteUseCase を生成する。
func NewGetNoteUseCase(notes repository.NoteRepository) *GetNoteUseCase {
	return &GetNoteUseCase{notes: notes}
}

// Execute はノートを返す。所有者のみ取得できる。
func (uc *GetNoteUseCase) Execute(ctx context.Context, id, userID uint) (*model.Note, error) {
	return findOwnedNote(ctx, uc.notes, id, userID)
}

// ListNotesUseCase はユーザーのノート一覧を取得する。
type ListNotesUseCase struct {
	notes repository.NoteRepository
}

// NewListNotesUseCase は ListNotesUseCase を生成する。
func NewListNotesUseCase(notes repository.NoteRepository) *ListNotesUseCase {
	return &ListNotesUseCase{notes: notes}
}

// Execute はアーカイブ済みを除いたノートを更新日の新しい順で返す。
func (uc *ListNotesUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return uc.notes.FindByUserID(ctx, userID, page, limit)
}

// ListNotesByFolderUseCase は指定フォルダのノート一覧を取得する。
type ListNotesByFolderUseCase struct {
	notes repository.NoteRepository
}

// NewListNotesByFolderUseCase は ListNotesByFolderUseCase を生成する。
func NewListNotesByFolderUseCase(notes repository.NoteRepository) *ListNotesByFolderUseCase {
	return &ListNotesByFolderUseCase{notes: notes}
}

// Execute は指定フォルダ内の自分のノートを返す。
func (uc *ListNotesByFolderUseCase) Execute(ctx context.Context, folderID, userID uint) ([]model.Note, error) {
	return uc.notes.FindByFolderID(ctx, folderID, userID)
}

// UpdateNoteUseCase はノートを更新する。
type UpdateNoteUseCase struct {
	notes repository.NoteRepository
}

// NewUpdateNoteUseCase は UpdateNoteUseCase を生成する。
func NewUpdateNoteUseCase(notes repository.NoteRepository) *UpdateNoteUseCase {
	return &UpdateNoteUseCase{notes: notes}
}

// Execute はノートを部分更新する。所有者のみ。
// 空文字列は「変更なし」、空白のみは 400 として扱う。
func (uc *UpdateNoteUseCase) Execute(ctx context.Context, id, userID uint, title, content, tags string, folderID *uint) (*model.Note, error) {
	note, err := findOwnedNote(ctx, uc.notes, id, userID)
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

	if err := uc.notes.Update(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNoteUseCase はノートを削除する。
type DeleteNoteUseCase struct {
	notes repository.NoteRepository
}

// NewDeleteNoteUseCase は DeleteNoteUseCase を生成する。
func NewDeleteNoteUseCase(notes repository.NoteRepository) *DeleteNoteUseCase {
	return &DeleteNoteUseCase{notes: notes}
}

// Execute はノートを削除する。所有者のみ。
func (uc *DeleteNoteUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedNote(ctx, uc.notes, id, userID); err != nil {
		return err
	}
	return uc.notes.Delete(ctx, id)
}

// SearchNotesUseCase はノートをキーワード検索する。
type SearchNotesUseCase struct {
	notes repository.NoteRepository
}

// NewSearchNotesUseCase は SearchNotesUseCase を生成する。
func NewSearchNotesUseCase(notes repository.NoteRepository) *SearchNotesUseCase {
	return &SearchNotesUseCase{notes: notes}
}

// Execute はアーカイブ済みを除き、タイトルまたは本文への部分一致で検索する。
func (uc *SearchNotesUseCase) Execute(ctx context.Context, userID uint, query string, page, limit int) ([]model.Note, int64, error) {
	offset := (page - 1) * limit
	return uc.notes.Search(ctx, userID, query, limit, offset)
}

// CountNotesUseCase はノート数を取得する。
type CountNotesUseCase struct {
	notes repository.NoteRepository
}

// NewCountNotesUseCase は CountNotesUseCase を生成する。
func NewCountNotesUseCase(notes repository.NoteRepository) *CountNotesUseCase {
	return &CountNotesUseCase{notes: notes}
}

// Execute はアーカイブ済みを除いたノート総数を返す。
func (uc *CountNotesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.notes.CountByUserID(ctx, userID)
}

// ToggleNoteFavoriteUseCase はノートのお気に入り状態を切り替える。
type ToggleNoteFavoriteUseCase struct {
	notes repository.NoteRepository
}

// NewToggleNoteFavoriteUseCase は ToggleNoteFavoriteUseCase を生成する。
func NewToggleNoteFavoriteUseCase(notes repository.NoteRepository) *ToggleNoteFavoriteUseCase {
	return &ToggleNoteFavoriteUseCase{notes: notes}
}

// Execute はお気に入り状態を反転させる。所有者のみ。
func (uc *ToggleNoteFavoriteUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedNote(ctx, uc.notes, id, userID); err != nil {
		return err
	}
	return uc.notes.ToggleFavorite(ctx, id)
}

// ListFavoriteNotesUseCase はお気に入りノート一覧を取得する。
type ListFavoriteNotesUseCase struct {
	notes repository.NoteRepository
}

// NewListFavoriteNotesUseCase は ListFavoriteNotesUseCase を生成する。
func NewListFavoriteNotesUseCase(notes repository.NoteRepository) *ListFavoriteNotesUseCase {
	return &ListFavoriteNotesUseCase{notes: notes}
}

// Execute はお気に入りのノートを更新日の新しい順で返す。
func (uc *ListFavoriteNotesUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return uc.notes.FindFavorites(ctx, userID, page, limit)
}

// CountFavoriteNotesUseCase はお気に入りノート数を取得する。
type CountFavoriteNotesUseCase struct {
	notes repository.NoteRepository
}

// NewCountFavoriteNotesUseCase は CountFavoriteNotesUseCase を生成する。
func NewCountFavoriteNotesUseCase(notes repository.NoteRepository) *CountFavoriteNotesUseCase {
	return &CountFavoriteNotesUseCase{notes: notes}
}

// Execute はお気に入りのノート総数を返す。
func (uc *CountFavoriteNotesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.notes.CountFavoritesByUserID(ctx, userID)
}

// ArchiveNoteUseCase はノートをアーカイブする。
type ArchiveNoteUseCase struct {
	notes repository.NoteRepository
}

// NewArchiveNoteUseCase は ArchiveNoteUseCase を生成する。
func NewArchiveNoteUseCase(notes repository.NoteRepository) *ArchiveNoteUseCase {
	return &ArchiveNoteUseCase{notes: notes}
}

// Execute はノートをアーカイブする。所有者のみ。
func (uc *ArchiveNoteUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedNote(ctx, uc.notes, id, userID); err != nil {
		return err
	}
	return uc.notes.Archive(ctx, id)
}

// UnarchiveNoteUseCase はノートのアーカイブを解除する。
type UnarchiveNoteUseCase struct {
	notes repository.NoteRepository
}

// NewUnarchiveNoteUseCase は UnarchiveNoteUseCase を生成する。
func NewUnarchiveNoteUseCase(notes repository.NoteRepository) *UnarchiveNoteUseCase {
	return &UnarchiveNoteUseCase{notes: notes}
}

// Execute はアーカイブを解除する。所有者のみ。
func (uc *UnarchiveNoteUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedNote(ctx, uc.notes, id, userID); err != nil {
		return err
	}
	return uc.notes.Unarchive(ctx, id)
}

// ListArchivedNotesUseCase はアーカイブ済みノート一覧を取得する。
type ListArchivedNotesUseCase struct {
	notes repository.NoteRepository
}

// NewListArchivedNotesUseCase は ListArchivedNotesUseCase を生成する。
func NewListArchivedNotesUseCase(notes repository.NoteRepository) *ListArchivedNotesUseCase {
	return &ListArchivedNotesUseCase{notes: notes}
}

// Execute はアーカイブ済みのノートを更新日の新しい順で返す。
func (uc *ListArchivedNotesUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	return uc.notes.FindArchived(ctx, userID, page, limit)
}

// CountArchivedNotesUseCase はアーカイブ済みノート数を取得する。
type CountArchivedNotesUseCase struct {
	notes repository.NoteRepository
}

// NewCountArchivedNotesUseCase は CountArchivedNotesUseCase を生成する。
func NewCountArchivedNotesUseCase(notes repository.NoteRepository) *CountArchivedNotesUseCase {
	return &CountArchivedNotesUseCase{notes: notes}
}

// Execute はアーカイブ済みのノート総数を返す。
func (uc *CountArchivedNotesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.notes.CountArchivedByUserID(ctx, userID)
}

// ListNoteTagsUseCase はノートで使われているタグを重複なく取得する。
type ListNoteTagsUseCase struct {
	notes repository.NoteRepository
}

// NewListNoteTagsUseCase は ListNoteTagsUseCase を生成する。
func NewListNoteTagsUseCase(notes repository.NoteRepository) *ListNoteTagsUseCase {
	return &ListNoteTagsUseCase{notes: notes}
}

// noteTagScanLimit はタグ抽出のために読むノートの上限（移行前からの値）。
const noteTagScanLimit = 1000

// Execute は先頭 1000 件のノートからタグを重複なく抽出して返す。
func (uc *ListNoteTagsUseCase) Execute(ctx context.Context, userID uint) ([]string, error) {
	notes, err := uc.notes.FindByUserID(ctx, userID, 1, noteTagScanLimit)
	if err != nil {
		return nil, err
	}
	return ExtractUniqueNoteTags(notes), nil
}

// ExtractUniqueNoteTags はノート一覧からカンマ区切りのタグを重複なく抽出する純粋関数。
func ExtractUniqueNoteTags(notes []model.Note) []string {
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

// ExportNoteMarkdownUseCase はノートを Markdown として書き出す。
type ExportNoteMarkdownUseCase struct {
	notes repository.NoteRepository
}

// NewExportNoteMarkdownUseCase は ExportNoteMarkdownUseCase を生成する。
func NewExportNoteMarkdownUseCase(notes repository.NoteRepository) *ExportNoteMarkdownUseCase {
	return &ExportNoteMarkdownUseCase{notes: notes}
}

// Execute はノートを Markdown に整形して内容とタイトルを返す。所有者のみ。
func (uc *ExportNoteMarkdownUseCase) Execute(ctx context.Context, id, userID uint) ([]byte, string, error) {
	note, err := findOwnedNote(ctx, uc.notes, id, userID)
	if err != nil {
		return nil, "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", note.Title)
	if note.Tags != "" {
		fmt.Fprintf(&b, "**Tags:** %s\n", note.Tags)
	}
	fmt.Fprintf(&b, "**Created:** %s\n", note.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Fprintf(&b, "**Updated:** %s\n", note.UpdatedAt.Format("2006-01-02 15:04"))
	b.WriteString("\n---\n\n")
	b.WriteString(note.Content)
	b.WriteString("\n")

	return []byte(b.String()), note.Title, nil
}

// DuplicateNoteUseCase はノートを複製する。
type DuplicateNoteUseCase struct {
	notes repository.NoteRepository
}

// NewDuplicateNoteUseCase は DuplicateNoteUseCase を生成する。
func NewDuplicateNoteUseCase(notes repository.NoteRepository) *DuplicateNoteUseCase {
	return &DuplicateNoteUseCase{notes: notes}
}

// Execute はノートを複製する。所有者のみ。
// タイトルに「 (コピー)」を付け、お気に入り・アーカイブ状態はリセットする。
func (uc *DuplicateNoteUseCase) Execute(ctx context.Context, id, userID uint) (*model.Note, error) {
	original, err := findOwnedNote(ctx, uc.notes, id, userID)
	if err != nil {
		return nil, err
	}

	duplicate := &model.Note{
		UserID:     userID,
		Title:      original.Title + " (コピー)",
		Content:    original.Content,
		Tags:       original.Tags,
		FolderID:   original.FolderID,
		IsFavorite: false,
		IsArchived: false,
	}

	v := validator.NewNoteValidator()
	if err := v.ValidateCreateNote(duplicate.Title, duplicate.Content, duplicate.Tags); err != nil {
		return nil, err
	}

	if err := uc.notes.Create(ctx, duplicate); err != nil {
		return nil, err
	}
	return duplicate, nil
}
