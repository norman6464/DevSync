package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteFolderServiceInterface は NoteFolderService のインターフェース。
// テスト時のモック化を容易にするため、インターフェースとして定義する。
type NoteFolderServiceInterface interface {
	Create(folder *model.NoteFolder) error
	GetByID(id uint) (*model.NoteFolder, error)
	GetByUserID(userID uint) ([]model.NoteFolder, error)
	GetChildren(parentID uint) ([]model.NoteFolder, error)
	GetRootFolders(userID uint) ([]model.NoteFolder, error)
	Update(folder *model.NoteFolder) error
	Delete(id uint) error
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

// CreateNoteFolderInput はフォルダ作成のリクエストボディ。
type CreateNoteFolderInput struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// Create は新しいノートフォルダを作成する。
func (h *NoteFolderHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[CreateNoteFolderInput](c)
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

// GetByUserID は現在のユーザーのフォルダ一覧を取得する。
func (h *NoteFolderHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	folders, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, folders)
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

	respondOK(c, children)
}

// GetRootFolders は現在のユーザーのルートフォルダ（親なし）一覧を取得する。
func (h *NoteFolderHandler) GetRootFolders(c *gin.Context) {
	userID := c.GetUint("userID")

	folders, err := h.service.GetRootFolders(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, folders)
}

// UpdateNoteFolderInput はフォルダ更新のリクエストボディ。
type UpdateNoteFolderInput struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

// Update はフォルダを更新する。
func (h *NoteFolderHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[UpdateNoteFolderInput](c)
	if input == nil {
		return
	}

	// 既存のフォルダを取得して所有者確認
	folder, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	// 所有者チェック
	if folder.UserID != userID {
		respondForbidden(c, "この操作を行う権限がありません")
		return
	}

	// 更新フィールドを適用（空でない場合のみ）
	if input.Name != "" {
		folder.Name = input.Name
	}
	if input.ParentID != nil {
		folder.ParentID = input.ParentID
	}

	if err := h.service.Update(folder); err != nil {
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

	if err := h.service.Delete(id); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}
