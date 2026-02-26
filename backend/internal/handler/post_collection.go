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
	GetCollectionsForViewer(viewerID, targetUserID uint, limit, offset int) ([]model.PostCollection, int64, error)
	Update(id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error)
	Delete(id, userID uint) error
	AddPost(collectionID, userID, postID uint, note string) error
	RemovePost(collectionID, userID, postID uint) error
	GetPosts(collectionID uint) ([]model.PostCollectionItem, error)
	CountByUserID(userID uint) (int64, error)
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
// 表示権限の判定（自分=全件/他人=公開のみ）はService層に委譲する。
func (h *PostCollectionHandler) GetByUserID(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	targetUserID, ok := parseID(c, "userId")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	collections, total, err := h.service.GetCollectionsForViewer(currentUserID, targetUserID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.PostCollectionListResponse{
		Collections: ensureSlice(collections),
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	})
}

// GetMyCollections は認証ユーザーのコレクション一覧を取得する。
func (h *PostCollectionHandler) GetMyCollections(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	collections, total, err := h.service.GetCollectionsForViewer(userID, userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.PostCollectionListResponse{
		Collections: ensureSlice(collections),
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	})
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

	respondOK(c, ensureSlice(items))
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

// GetMyCount は認証ユーザーのコレクション総数を返す。
func (h *PostCollectionHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
