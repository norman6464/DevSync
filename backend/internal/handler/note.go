package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteHandler は学習ノート関連のHTTPハンドラ。
// ノートのCRUD・検索・お気に入り管理を処理する。
type NoteHandler struct {
	create         *usecase.CreateNoteUseCase
	get            *usecase.GetNoteUseCase
	list           *usecase.ListNotesUseCase
	listByFolder   *usecase.ListNotesByFolderUseCase
	update         *usecase.UpdateNoteUseCase
	remove         *usecase.DeleteNoteUseCase
	search         *usecase.SearchNotesUseCase
	count          *usecase.CountNotesUseCase
	toggleFavorite *usecase.ToggleNoteFavoriteUseCase
	listFavorites  *usecase.ListFavoriteNotesUseCase
	countFavorites *usecase.CountFavoriteNotesUseCase
	archive        *usecase.ArchiveNoteUseCase
	unarchive      *usecase.UnarchiveNoteUseCase
	listArchived   *usecase.ListArchivedNotesUseCase
	countArchived  *usecase.CountArchivedNotesUseCase
	listTags       *usecase.ListNoteTagsUseCase
	export         *usecase.ExportNoteMarkdownUseCase
	duplicate      *usecase.DuplicateNoteUseCase
}

// NewNoteHandler は新しいNoteHandlerインスタンスを生成する。
func NewNoteHandler(
	create *usecase.CreateNoteUseCase,
	get *usecase.GetNoteUseCase,
	list *usecase.ListNotesUseCase,
	listByFolder *usecase.ListNotesByFolderUseCase,
	update *usecase.UpdateNoteUseCase,
	remove *usecase.DeleteNoteUseCase,
	search *usecase.SearchNotesUseCase,
	count *usecase.CountNotesUseCase,
	toggleFavorite *usecase.ToggleNoteFavoriteUseCase,
	listFavorites *usecase.ListFavoriteNotesUseCase,
	countFavorites *usecase.CountFavoriteNotesUseCase,
	archive *usecase.ArchiveNoteUseCase,
	unarchive *usecase.UnarchiveNoteUseCase,
	listArchived *usecase.ListArchivedNotesUseCase,
	countArchived *usecase.CountArchivedNotesUseCase,
	listTags *usecase.ListNoteTagsUseCase,
	export *usecase.ExportNoteMarkdownUseCase,
	duplicate *usecase.DuplicateNoteUseCase,
) *NoteHandler {
	return &NoteHandler{
		create: create, get: get, list: list, listByFolder: listByFolder,
		update: update, remove: remove, search: search, count: count,
		toggleFavorite: toggleFavorite, listFavorites: listFavorites, countFavorites: countFavorites,
		archive: archive, unarchive: unarchive, listArchived: listArchived, countArchived: countArchived,
		listTags: listTags, export: export, duplicate: duplicate,
	}
}

// Create は新しいノートを作成する。
func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[dto.CreateNoteRequest](c)
	if input == nil {
		return
	}

	note := &model.Note{
		UserID:   userID,
		Title:    input.Title,
		Content:  input.Content,
		Tags:     input.Tags,
		FolderID: input.FolderID,
	}

	if err := h.create.Execute(c.Request.Context(), note); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, note)
}

// GetByID は指定IDのノートを所有権検証付きで取得する。
func (h *NoteHandler) GetByID(c *gin.Context) {
	handleGetByID(c, func(id, userID uint) (*model.Note, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// GetByUserID は現在のユーザーのノート一覧をページネーション付きで取得する。
func (h *NoteHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	notes, err := h.list.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, notes, total, page, limit)
}

// GetMyCount は認証ユーザー自身のノート総数を返す。
func (h *NoteHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// GetByFolderID は指定フォルダ内のノート一覧を所有権検証付きで取得する。
func (h *NoteHandler) GetByFolderID(c *gin.Context) {
	folderID, ok := parseID(c, "folderId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	notes, err := h.listByFolder.Execute(c.Request.Context(), folderID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(notes))
}

// Update はノートを更新する。
func (h *NoteHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateNoteRequest](c)
	if input == nil {
		return
	}

	note, err := h.update.Execute(c.Request.Context(), id, userID, input.Title, input.Content, input.Tags, input.FolderID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, note)
}

// Delete はノートを削除する。
func (h *NoteHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// Search はキーワードでノートを検索する（ページネーション付き）。
func (h *NoteHandler) Search(c *gin.Context) {
	userID := c.GetUint("userID")
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "検索キーワードが必要です")
		return
	}

	page, limit := parsePagination(c)

	notes, total, err := h.search.Execute(c.Request.Context(), userID, query, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, notes, total, page, limit)
}

// ToggleFavorite はノートのお気に入り状態を切り替える。
func (h *NoteHandler) ToggleFavorite(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.toggleFavorite.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("お気に入り状態を更新しました"))
}

// Archive はノートをアーカイブする。
func (h *NoteHandler) Archive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.archive.Execute(c.Request.Context(), id, userID)
	}, "ノートをアーカイブしました")
}

// Unarchive はノートのアーカイブを解除する。
func (h *NoteHandler) Unarchive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.unarchive.Execute(c.Request.Context(), id, userID)
	}, "ノートのアーカイブを解除しました")
}

// GetFavorites は現在のユーザーのお気に入りノート一覧をページネーション付きで取得する。
func (h *NoteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	notes, err := h.listFavorites.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.countFavorites.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, notes, total, page, limit)
}

// GetArchived は現在のユーザーのアーカイブ済みノート一覧をページネーション付きで取得する。
func (h *NoteHandler) GetArchived(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	notes, err := h.listArchived.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.countArchived.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, notes, total, page, limit)
}

// Duplicate は既存のノートを複製する。
func (h *NoteHandler) Duplicate(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	duplicate, err := h.duplicate.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, duplicate)
}

// GetTags は認証ユーザーのノートで使用されているタグ一覧を返す。
func (h *NoteHandler) GetTags(c *gin.Context) {
	userID := c.GetUint("userID")

	tags, err := h.listTags.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(tags))
}

// Export はノートをMarkdownファイルとしてエクスポートする。
func (h *NoteHandler) Export(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	data, title, err := h.export.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	filename := fmt.Sprintf("%s.md", title)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/markdown; charset=utf-8", data)
}
