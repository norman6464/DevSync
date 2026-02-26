package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// BookmarkCollectionServiceInterface はBookmarkCollectionHandlerが依存するサービスのインターフェース。
type BookmarkCollectionServiceInterface interface {
	Create(collection *model.BookmarkCollection) error
	GetByUserID(userID uint) ([]model.BookmarkCollection, error)
	Update(id, userID uint, updates *model.BookmarkCollection) (*model.BookmarkCollection, error)
	Delete(id, userID uint) error
	AddPost(collectionID, postID, userID uint) error
	RemovePost(collectionID, postID, userID uint) error
	GetPosts(collectionID uint, limit, offset int) ([]model.Post, int64, error)
}

// BookmarkCollectionHandler はブックマークコレクション関連のHTTPハンドラ。
type BookmarkCollectionHandler struct {
	service BookmarkCollectionServiceInterface
}

// NewBookmarkCollectionHandler は新しいBookmarkCollectionHandlerインスタンスを生成する。
func NewBookmarkCollectionHandler(s BookmarkCollectionServiceInterface) *BookmarkCollectionHandler {
	return &BookmarkCollectionHandler{service: s}
}

// Create は新しいブックマークコレクションを作成する。
func (h *BookmarkCollectionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.BookmarkCollectionRequest](c)
	if input == nil {
		return
	}

	collection := &model.BookmarkCollection{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	}

	if err := h.service.Create(collection); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, collection)
}

// GetMyCollections はユーザー自身のコレクション一覧を返す。
func (h *BookmarkCollectionHandler) GetMyCollections(c *gin.Context) {
	userID := c.GetUint("userID")

	collections, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(collections))
}

// Update はコレクションを更新する。
func (h *BookmarkCollectionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.BookmarkCollectionRequest](c)
	if input == nil {
		return
	}

	updated, err := h.service.Update(id, userID, &model.BookmarkCollection{
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, updated)
}

// Delete はコレクションを削除する。
func (h *BookmarkCollectionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// AddPost はコレクションにブックマークを追加する。
func (h *BookmarkCollectionHandler) AddPost(c *gin.Context) {
	userID := c.GetUint("userID")
	collectionID, ok := parseID(c, "id")
	if !ok {
		return
	}
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	if err := h.service.AddPost(collectionID, postID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("コレクションに追加しました"))
}

// RemovePost はコレクションからブックマークを削除する。
func (h *BookmarkCollectionHandler) RemovePost(c *gin.Context) {
	userID := c.GetUint("userID")
	collectionID, ok := parseID(c, "id")
	if !ok {
		return
	}
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	if err := h.service.RemovePost(collectionID, postID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// GetPosts はコレクション内の投稿一覧を返す。
func (h *BookmarkCollectionHandler) GetPosts(c *gin.Context) {
	collectionID, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	posts, total, err := h.service.GetPosts(collectionID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.BookmarkCollectionPostsResponse{
		Posts: ensureSlice(posts),
		Total: total,
	})
}
