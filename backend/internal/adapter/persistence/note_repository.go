package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteRepository は [repository.NoteRepository] の sqlc(pgx) 実装。
type noteRepository struct {
	q *sqlcgen.Queries
}

// NewNoteRepository は NoteRepository の sqlc(pgx) 実装を返す。
func NewNoteRepository(q *sqlcgen.Queries) repository.NoteRepository {
	return &noteRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteRepository = (*noteRepository)(nil)

// toModelUser は sqlc の生成行を model.User へ変換する（Note.User の Preload 相当分のみ）。
func toModelUser(u sqlcgen.User) model.User {
	return model.User{
		ID:                  uint(u.ID),
		Username:            u.Username,
		Name:                u.Name,
		Email:               u.Email,
		Password:            fromStringPtr(u.Password),
		AvatarURL:           fromStringPtr(u.AvatarUrl),
		Bio:                 fromStringPtr(u.Bio),
		GitHubID:            fromInt64PtrValue(u.GitHubID),
		GitHubUsername:      fromStringPtr(u.GitHubUsername),
		GitHubToken:         fromStringPtr(u.GitHubToken),
		GitHubConnected:     fromBoolPtr(u.GitHubConnected),
		SpotifyConnected:    fromBoolPtr(u.SpotifyConnected),
		SpotifyToken:        fromStringPtr(u.SpotifyToken),
		SpotifyRefreshToken: fromStringPtr(u.SpotifyRefreshToken),
		SpotifyTokenExpiry:  timeValue(fromTimestamptz(u.SpotifyTokenExpiry)),
		ZennUsername:        fromStringPtr(u.ZennUsername),
		QiitaUsername:       fromStringPtr(u.QiitaUsername),
		AtCoderUsername:     fromStringPtr(u.AtCoderUsername),
		PaizaRank:           fromStringPtr(u.PaizaRank),
		SkillsLanguages:     fromStringPtr(u.SkillsLanguages),
		SkillsFrameworks:    fromStringPtr(u.SkillsFrameworks),
		OnboardingCompleted: fromBoolPtr(u.OnboardingCompleted),
		EmailWeeklyReport:   fromBoolPtr(u.EmailWeeklyReport),
		EmailLanguage:       fromStringPtr(u.EmailLanguage),
		CreatedAt:           timeValue(fromTimestamptz(u.CreatedAt)),
		UpdatedAt:           timeValue(fromTimestamptz(u.UpdatedAt)),
	}
}

// fromInt64PtrValue は NULL 許容の int64 を、不在なら 0 として返す。
func fromInt64PtrValue(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// timeValue は *time.Time を、不在ならゼロ値として返す。
func timeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// attachFolder は LEFT JOIN で取得した note_folders の個別カラムを model.Note へ Folder として付与する。
// folderID が nil の場合はフォルダなし（GORMの Preload("Folder") で Folder == nil の状態）として何もしない。
func attachFolder(note *model.Note, folderID, folderUserID, folderParentID *int64, folderName *string, folderCreatedAt, folderUpdatedAt pgtype.Timestamptz) {
	if folderID == nil {
		return
	}
	note.Folder = &model.NoteFolder{
		ID:        uint(*folderID),
		UserID:    uint(*folderUserID),
		ParentID:  fromInt64PtrToUintPtr(folderParentID),
		Name:      fromStringPtr(folderName),
		CreatedAt: timeValue(fromTimestamptz(folderCreatedAt)),
		UpdatedAt: timeValue(fromTimestamptz(folderUpdatedAt)),
	}
}

// Create は新しいノートを作成する。
func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
	row, err := r.q.CreateNote(ctx, sqlcgen.CreateNoteParams{
		UserID:     int64(note.UserID),
		FolderID:   toInt64PtrFromUintPtr(note.FolderID),
		Title:      note.Title,
		Content:    &note.Content,
		Tags:       &note.Tags,
		IsFavorite: &note.IsFavorite,
		IsArchived: &note.IsArchived,
	})
	if err != nil {
		return err
	}
	*note = toModelNote(row)
	return nil
}

// Update は既存のノートを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	row, err := r.q.UpdateNote(ctx, sqlcgen.UpdateNoteParams{
		ID:         int64(note.ID),
		Title:      note.Title,
		Content:    &note.Content,
		Tags:       &note.Tags,
		FolderID:   toInt64PtrFromUintPtr(note.FolderID),
		IsFavorite: &note.IsFavorite,
		IsArchived: &note.IsArchived,
	})
	if err != nil {
		return err
	}
	*note = toModelNote(row)
	return nil
}

// Delete はノートを削除する。
func (r *noteRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteNote(ctx, int64(id))
}

// FindByID は指定 ID のノートをユーザー・フォルダ付きで取得する。不在の場合は (nil, nil) を返す。
func (r *noteRepository) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	row, err := r.q.GetNoteByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	note := toModelNote(row.Note)
	note.User = toModelUser(row.User)
	attachFolder(&note, row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	return &note, nil
}

// FindByUserID はアーカイブ済みを除いたノートを更新日の新しい順で取得する。
func (r *noteRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListNotesByUser(ctx, sqlcgen.ListNotesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	notes := make([]model.Note, len(rows))
	for i, row := range rows {
		notes[i] = toModelNote(row.Note)
		attachFolder(&notes[i], row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	}
	return notes, nil
}

// FindFavorites はお気に入りのノートを更新日の新しい順で取得する。
func (r *noteRepository) FindFavorites(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListFavoriteNotesByUser(ctx, sqlcgen.ListFavoriteNotesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	notes := make([]model.Note, len(rows))
	for i, row := range rows {
		notes[i] = toModelNote(row.Note)
		attachFolder(&notes[i], row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	}
	return notes, nil
}

// FindArchived はアーカイブ済みのノートを更新日の新しい順で取得する。
func (r *noteRepository) FindArchived(ctx context.Context, userID uint, page, limit int) ([]model.Note, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListArchivedNotesByUser(ctx, sqlcgen.ListArchivedNotesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	notes := make([]model.Note, len(rows))
	for i, row := range rows {
		notes[i] = toModelNote(row.Note)
		attachFolder(&notes[i], row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	}
	return notes, nil
}

// FindByFolderID は指定フォルダ内のノートを更新日の新しい順で取得する。
func (r *noteRepository) FindByFolderID(ctx context.Context, folderID, userID uint) ([]model.Note, error) {
	rows, err := r.q.ListNotesByFolder(ctx, sqlcgen.ListNotesByFolderParams{
		FolderID: toInt64PtrFromUintPtr(&folderID),
		UserID:   int64(userID),
	})
	if err != nil {
		return nil, err
	}
	notes := make([]model.Note, len(rows))
	for i, row := range rows {
		notes[i] = toModelNote(row)
	}
	return notes, nil
}

// Search はアーカイブ済みを除き、タイトルまたは本文への部分一致で検索する。
func (r *noteRepository) Search(ctx context.Context, userID uint, query string, limit, offset int) ([]model.Note, int64, error) {
	pattern := escapeLikePattern(query)

	total, err := r.q.CountSearchNotes(ctx, sqlcgen.CountSearchNotesParams{
		UserID: int64(userID),
		Title:  pattern,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchNotes(ctx, sqlcgen.SearchNotesParams{
		UserID: int64(userID),
		Title:  pattern,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	notes := make([]model.Note, len(rows))
	for i, row := range rows {
		notes[i] = toModelNote(row.Note)
		attachFolder(&notes[i], row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	}
	return notes, total, nil
}

// CountByUserID はアーカイブ済みを除いたノート総数を返す。
func (r *noteRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountActiveNotesByUser(ctx, int64(userID))
}

// CountFavoritesByUserID はお気に入りのノート総数を返す。
func (r *noteRepository) CountFavoritesByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountFavoriteNotesByUser(ctx, int64(userID))
}

// CountArchivedByUserID はアーカイブ済みのノート総数を返す。
func (r *noteRepository) CountArchivedByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountArchivedNotesByUser(ctx, int64(userID))
}

// ToggleFavorite はお気に入り状態を反転させる。
func (r *noteRepository) ToggleFavorite(ctx context.Context, id uint) error {
	return r.q.ToggleNoteFavorite(ctx, int64(id))
}

// Archive はノートをアーカイブする。
func (r *noteRepository) Archive(ctx context.Context, id uint) error {
	return r.q.ArchiveNote(ctx, int64(id))
}

// Unarchive はノートのアーカイブを解除する。
func (r *noteRepository) Unarchive(ctx context.Context, id uint) error {
	return r.q.UnarchiveNote(ctx, int64(id))
}
