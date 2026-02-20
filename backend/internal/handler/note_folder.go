package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteFolderServiceInterface は NoteFolderService のインターフェース。
// テスト時のモック化を容易にするため、インターフェースとして定義する。
type NoteFolderServiceInterface interface {
	Create(folder *model.NoteFolder) error
	GetByID(id uint) (*model.NoteFolder, error)
	GetByUserID(userID uint, limit, offset int) ([]model.NoteFolder, int64, error)
	GetChildren(parentID uint) ([]model.NoteFolder, error)
	GetRootFolders(userID uint) ([]model.NoteFolder, error)
	Update(id, userID uint, name string, parentID *uint) (*model.NoteFolder, error)
	Delete(id, userID uint) error
}

// NoteFolderHandler はノートフォルダ関連のHTTPハンドラ。
// フォルダのCRUD・階層管理を処理する。
type NoteFolderHandler struct {
	service NoteFolderServiceInterface
}

// NewNoteFolderHandler は新しいNoteFolderHandlerインスタンスを生成する。
func NewNoteFolderHandler(s NoteFolderServiceInterface) *NoteFolderHandler {
	return &NoteFolderHandler{service: s}
}

// Create は新しいノートフォルダを作成する。
func (h *NoteFolderHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[dto.CreateNoteFolderRequest](c)
	if input == nil {
		return
	}

	folder := &model.NoteFolder{
		UserID:   userID,
		Name:     input.Name,
		ParentID: input.ParentID,
	}

	if err := h.service.Create(folder); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, folder)
}

// GetByID は指定IDのフォルダを取得する。
func (h *NoteFolderHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	folder, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, folder)
}

// GetByUserID は現在のユーザーのフォルダ一覧をページネーション付きで取得する。
func (h *NoteFolderHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	folders, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.NoteFolderListResponse{
		Folders: ensureSlice(folders),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetChildren は指定フォルダの子フォルダ一覧を取得する。
func (h *NoteFolderHandler) GetChildren(c *gin.Context) {
	parentID, ok := parseID(c, "id")
	if !ok {
		return
	}

	children, err := h.service.GetChildren(parentID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(children))
}

// GetRootFolders は現在のユーザーのルートフォルダ（親なし）一覧を取得する。
func (h *NoteFolderHandler) GetRootFolders(c *gin.Context) {
	userID := c.GetUint("userID")

	folders, err := h.service.GetRootFolders(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(folders))
}

// Update はフォルダを更新する。
func (h *NoteFolderHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateNoteFolderRequest](c)
	if input == nil {
		return
	}

	folder, err := h.service.Update(id, userID, input.Name, input.ParentID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, folder)
}

// Delete はフォルダを削除する。
func (h *NoteFolderHandler) Delete(c *gin.Context) {
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
