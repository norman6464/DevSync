package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteServiceInterface はNoteServiceのインターフェース。
// テスト時のモック化を容易にするため、インターフェースとして定義する。
type NoteServiceInterface interface {
	Create(note *model.Note) error
	GetByID(id, userID uint) (*model.Note, error)
	GetByUserID(userID uint, page, limit int) ([]model.Note, error)
	GetByFolderID(folderID, userID uint) ([]model.Note, error)
	Update(id, userID uint, title, content, tags string, folderID *uint) (*model.Note, error)
	Delete(id, userID uint) error
	Search(userID uint, query string, page, limit int) ([]model.Note, int64, error)
	CountByUserID(userID uint) (int64, error)
	ToggleFavorite(id, userID uint) error
	GetFavorites(userID uint, page, limit int) ([]model.Note, error)
	CountFavoritesByUserID(userID uint) (int64, error)
	Archive(id, userID uint) error
	Unarchive(id, userID uint) error
	GetArchived(userID uint, page, limit int) ([]model.Note, error)
	CountArchivedByUserID(userID uint) (int64, error)
	Duplicate(id uint, userID uint) (*model.Note, error)
	ExportMarkdown(id, userID uint) ([]byte, string, error)
	GetTags(userID uint) ([]string, error)
}

// NoteHandler は学習ノート関連のHTTPハンドラ。
// ノートのCRUD・検索・お気に入り管理を処理する。
type NoteHandler struct {
	service NoteServiceInterface
}

// NewNoteHandler は新しいNoteHandlerインスタンスを生成する。
func NewNoteHandler(s NoteServiceInterface) *NoteHandler {
	return &NoteHandler{service: s}
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

	if err := h.service.Create(note); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, note)
}

// GetByID は指定IDのノートを所有権検証付きで取得する。
func (h *NoteHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// GetByUserID は現在のユーザーのノート一覧をページネーション付きで取得する。
func (h *NoteHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	notes, err := h.service.GetByUserID(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, notes, total, page, limit)
}

// GetByFolderID は指定フォルダ内のノート一覧を所有権検証付きで取得する。
func (h *NoteHandler) GetByFolderID(c *gin.Context) {
	folderID, ok := parseID(c, "folderId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	notes, err := h.service.GetByFolderID(folderID, userID)
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

	note, err := h.service.Update(id, userID, input.Title, input.Content, input.Tags, input.FolderID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, note)
}

// Delete はノートを削除する。
func (h *NoteHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
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

	notes, total, err := h.service.Search(userID, query, page, limit)
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

	if err := h.service.ToggleFavorite(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("お気に入り状態を更新しました"))
}

// Archive はノートをアーカイブする。
func (h *NoteHandler) Archive(c *gin.Context) {
	handleAction(c, h.service.Archive, "ノートをアーカイブしました")
}

// Unarchive はノートのアーカイブを解除する。
func (h *NoteHandler) Unarchive(c *gin.Context) {
	handleAction(c, h.service.Unarchive, "ノートのアーカイブを解除しました")
}

// GetFavorites は現在のユーザーのお気に入りノート一覧をページネーション付きで取得する。
func (h *NoteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	notes, err := h.service.GetFavorites(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.service.CountFavoritesByUserID(userID)
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

	notes, err := h.service.GetArchived(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.service.CountArchivedByUserID(userID)
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

	duplicate, err := h.service.Duplicate(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, duplicate)
}

// GetTags は認証ユーザーのノートで使用されているタグ一覧を返す。
func (h *NoteHandler) GetTags(c *gin.Context) {
	userID := c.GetUint("userID")

	tags, err := h.service.GetTags(userID)
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

	data, title, err := h.service.ExportMarkdown(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	filename := fmt.Sprintf("%s.md", title)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/markdown; charset=utf-8", data)
}
