package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteFolderHandler はノートフォルダ関連のHTTPハンドラ。
// フォルダのCRUD・階層管理を処理する。
type NoteFolderHandler struct {
	create       *usecase.CreateNoteFolderUseCase
	get          *usecase.GetNoteFolderUseCase
	list         *usecase.ListNoteFoldersUseCase
	listChild    *usecase.ListChildNoteFoldersUseCase
	listRoot     *usecase.ListRootNoteFoldersUseCase
	update       *usecase.UpdateNoteFolderUseCase
	count        *usecase.CountNoteFoldersUseCase
	deleteFolder *usecase.DeleteNoteFolderUseCase
}

// NewNoteFolderHandler は新しいNoteFolderHandlerインスタンスを生成する。
func NewNoteFolderHandler(
	create *usecase.CreateNoteFolderUseCase,
	get *usecase.GetNoteFolderUseCase,
	list *usecase.ListNoteFoldersUseCase,
	listChild *usecase.ListChildNoteFoldersUseCase,
	listRoot *usecase.ListRootNoteFoldersUseCase,
	update *usecase.UpdateNoteFolderUseCase,
	count *usecase.CountNoteFoldersUseCase,
	deleteFolder *usecase.DeleteNoteFolderUseCase,
) *NoteFolderHandler {
	return &NoteFolderHandler{
		create:       create,
		get:          get,
		list:         list,
		listChild:    listChild,
		listRoot:     listRoot,
		update:       update,
		count:        count,
		deleteFolder: deleteFolder,
	}
}

// Create は新しいノートフォルダを作成する。
func (h *NoteFolderHandler) Create(c *gin.Context) {
	input := bindJSON[dto.CreateNoteFolderRequest](c)
	if input == nil {
		return
	}

	folder, err := h.create.Execute(c.Request.Context(), usecase.CreateNoteFolderInput{
		UserID:   c.GetUint("userID"),
		Name:     input.Name,
		ParentID: input.ParentID,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, folder)
}

// GetByID は指定IDのフォルダを取得する。所有者本人のみ取得できる。
func (h *NoteFolderHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	handleGetByIDPublic(c, func(id uint) (*model.NoteFolder, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// GetByUserID は現在のユーザーのフォルダ一覧をページネーション付きで取得する。
func (h *NoteFolderHandler) GetByUserID(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	folders, total, err := h.list.Execute(c.Request.Context(), c.GetUint("userID"), limit, offset)
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

	children, err := h.listChild.Execute(c.Request.Context(), parentID, c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(children))
}

// GetRootFolders は現在のユーザーのルートフォルダ（親なし）一覧を取得する。
func (h *NoteFolderHandler) GetRootFolders(c *gin.Context) {
	folders, err := h.listRoot.Execute(c.Request.Context(), c.GetUint("userID"))
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

	input := bindJSON[dto.UpdateNoteFolderRequest](c)
	if input == nil {
		return
	}

	folder, err := h.update.Execute(c.Request.Context(), usecase.UpdateNoteFolderInput{
		ID:       id,
		UserID:   c.GetUint("userID"),
		Name:     input.Name,
		ParentID: input.ParentID,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, folder)
}

// GetMyCount は認証ユーザー自身のフォルダ総数を返す。
func (h *NoteFolderHandler) GetMyCount(c *gin.Context) {
	count, err := h.count.Execute(c.Request.Context(), c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// Delete はフォルダを削除する。
func (h *NoteFolderHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.deleteFolder.Execute(c.Request.Context(), id, userID)
	})
}
