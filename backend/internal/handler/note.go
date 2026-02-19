package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteServiceInterface はNoteServiceのインターフェース。
// テスト時のモック化を容易にするため、インターフェースとして定義する。
type NoteServiceInterface interface {
	Create(note *model.Note) error
	GetByID(id uint) (*model.Note, error)
	GetByUserID(userID uint, page, limit int) ([]model.Note, error)
	GetByFolderID(folderID uint) ([]model.Note, error)
	Update(id, userID uint, title, content, tags string, folderID *uint) (*model.Note, error)
	Delete(id, userID uint) error
	Search(userID uint, query string, page, limit int) ([]model.Note, int64, error)
	CountByUserID(userID uint) (int64, error)
	ToggleFavorite(id, userID uint) error
	Archive(id, userID uint) error
	Unarchive(id, userID uint) error
	GetArchived(userID uint, page, limit int) ([]model.Note, error)
	CountArchivedByUserID(userID uint) (int64, error)
	Duplicate(id uint, userID uint) (*model.Note, error)
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

// CreateNoteInput はノート作成のリクエストボディ。
type CreateNoteInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Tags     string `json:"tags"`
	FolderID *uint  `json:"folder_id"`
}

// Create は新しいノートを作成する。
func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[CreateNoteInput](c)
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

// GetByID は指定IDのノートを取得する。
func (h *NoteHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	note, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, note)
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

// GetByFolderID は指定フォルダ内のノート一覧を取得する。
func (h *NoteHandler) GetByFolderID(c *gin.Context) {
	folderID, ok := parseID(c, "folderId")
	if !ok {
		return
	}

	notes, err := h.service.GetByFolderID(folderID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, notes)
}

// UpdateNoteInput はノート更新のリクエストボディ。
type UpdateNoteInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Tags     string `json:"tags"`
	FolderID *uint  `json:"folder_id"`
}

// Update はノートを更新する。
func (h *NoteHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[UpdateNoteInput](c)
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
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
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
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Archive(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ノートをアーカイブしました"))
}

// Unarchive はノートのアーカイブを解除する。
func (h *NoteHandler) Unarchive(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unarchive(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ノートのアーカイブを解除しました"))
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
