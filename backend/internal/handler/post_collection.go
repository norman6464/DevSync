package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostCollectionServiceInterface はPostCollectionHandlerが依存するサービスメソッドを定義する。
type PostCollectionServiceInterface interface {
	Create(collection *model.PostCollection) (*model.PostCollection, error)
	GetByID(id uint) (*model.PostCollection, error)
	GetByUserID(userID uint) ([]model.PostCollection, error)
	GetPublicByUserID(userID uint) ([]model.PostCollection, error)
	Update(id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error)
	Delete(id, userID uint) error
	AddPost(collectionID, userID, postID uint, note string) error
	RemovePost(collectionID, userID, postID uint) error
	GetPosts(collectionID uint) ([]model.PostCollectionItem, error)
}

// PostCollectionHandler は投稿コレクション関連のHTTPハンドラ。
type PostCollectionHandler struct {
	service PostCollectionServiceInterface
}

// NewPostCollectionHandler は新しいPostCollectionHandlerインスタンスを生成する。
func NewPostCollectionHandler(s PostCollectionServiceInterface) *PostCollectionHandler {
	return &PostCollectionHandler{service: s}
}

// Create は新しい投稿コレクションを作成する。
func (h *PostCollectionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreatePostCollectionRequest](c)
	if req == nil {
		return
	}

	collection := &model.PostCollection{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}

	result, err := h.service.Create(collection)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, result)
}

// GetByID は指定IDのコレクションを取得する。
func (h *PostCollectionHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	collection, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, collection)
}

// GetByUserID は指定ユーザーのコレクション一覧を取得する。
// 自分のコレクションは全件、他人のコレクションは公開のみ返す。
func (h *PostCollectionHandler) GetByUserID(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	targetUserID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	var collections []model.PostCollection
	var err error

	if currentUserID == targetUserID {
		collections, err = h.service.GetByUserID(targetUserID)
	} else {
		collections, err = h.service.GetPublicByUserID(targetUserID)
	}

	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, collections)
}

// Update は指定IDのコレクションを更新する。
func (h *PostCollectionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdatePostCollectionRequest](c)
	if req == nil {
		return
	}

	collection, err := h.service.Update(id, userID, req.Title, req.Description, req.IsPublic)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, collection)
}

// Delete は指定IDのコレクションを削除する。
func (h *PostCollectionHandler) Delete(c *gin.Context) {
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

// GetPosts はコレクション内の投稿一覧を取得する。
func (h *PostCollectionHandler) GetPosts(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	items, err := h.service.GetPosts(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, items)
}

// AddPost はコレクションに投稿を追加する。
func (h *PostCollectionHandler) AddPost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.AddPostToCollectionRequest](c)
	if req == nil {
		return
	}

	if err := h.service.AddPost(id, userID, req.PostID, req.Note); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, nil)
}

// RemovePost はコレクションから投稿を削除する。
func (h *PostCollectionHandler) RemovePost(c *gin.Context) {
	userID := c.GetUint("userID")
	collectionID, ok := parseID(c, "id")
	if !ok {
		return
	}
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	if err := h.service.RemovePost(collectionID, userID, postID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}
