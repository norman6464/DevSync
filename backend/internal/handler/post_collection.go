package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostCollectionHandler は投稿コレクション関連の HTTP ハンドラ。
type PostCollectionHandler struct {
	create        *usecase.CreatePostCollectionUseCase
	get           *usecase.GetPostCollectionUseCase
	listForViewer *usecase.ListPostCollectionsForViewerUseCase
	count         *usecase.CountPostCollectionsUseCase
	update        *usecase.UpdatePostCollectionUseCase
	delete        *usecase.DeletePostCollectionUseCase
	addPost       *usecase.AddPostToCollectionUseCase
	removePost    *usecase.RemovePostFromCollectionUseCase
	listPosts     *usecase.ListPostCollectionPostsUseCase
}

// NewPostCollectionHandler は PostCollectionHandler を生成する。
func NewPostCollectionHandler(
	create *usecase.CreatePostCollectionUseCase,
	get *usecase.GetPostCollectionUseCase,
	listForViewer *usecase.ListPostCollectionsForViewerUseCase,
	count *usecase.CountPostCollectionsUseCase,
	update *usecase.UpdatePostCollectionUseCase,
	deleteUC *usecase.DeletePostCollectionUseCase,
	addPost *usecase.AddPostToCollectionUseCase,
	removePost *usecase.RemovePostFromCollectionUseCase,
	listPosts *usecase.ListPostCollectionPostsUseCase,
) *PostCollectionHandler {
	return &PostCollectionHandler{
		create:        create,
		get:           get,
		listForViewer: listForViewer,
		count:         count,
		update:        update,
		delete:        deleteUC,
		addPost:       addPost,
		removePost:    removePost,
		listPosts:     listPosts,
	}
}

// createPostCollectionRequest はコレクション作成のリクエストボディ。
type createPostCollectionRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	IsPublic    bool   `json:"is_public"`
}

// Create は新しい投稿コレクションを作成する。
func (h *PostCollectionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[createPostCollectionRequest](c)
	if req == nil {
		return
	}

	collection := &model.PostCollection{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}

	result, err := h.create.Execute(c.Request.Context(), collection)
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

	collection, err := h.get.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, collection)
}

// postCollectionListResponse はコレクション一覧レスポンス（ページネーション付き）。
type postCollectionListResponse struct {
	Collections []model.PostCollection `json:"collections"`
	Total       int64                  `json:"total"`
	Limit       int                    `json:"limit"`
	Offset      int                    `json:"offset"`
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

	collections, total, err := h.listForViewer.Execute(c.Request.Context(), currentUserID, targetUserID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, postCollectionListResponse{
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

	collections, total, err := h.listForViewer.Execute(c.Request.Context(), userID, userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, postCollectionListResponse{
		Collections: ensureSlice(collections),
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	})
}

// updatePostCollectionRequest はコレクション更新のリクエストボディ。
type updatePostCollectionRequest struct {
	Title       string `json:"title" binding:"omitempty,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	IsPublic    bool   `json:"is_public"`
}

// Update は指定IDのコレクションを更新する。
func (h *PostCollectionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updatePostCollectionRequest](c)
	if req == nil {
		return
	}

	collection, err := h.update.Execute(c.Request.Context(), id, userID, req.Title, req.Description, req.IsPublic)
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

	if err := h.delete.Execute(c.Request.Context(), id, userID); err != nil {
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

	items, err := h.listPosts.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(items))
}

// addPostToCollectionRequest はコレクションへの投稿追加リクエスト。
type addPostToCollectionRequest struct {
	PostID uint   `json:"post_id" binding:"required"`
	Note   string `json:"note" binding:"omitempty,max=1000"`
}

// AddPost はコレクションに投稿を追加する。
func (h *PostCollectionHandler) AddPost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[addPostToCollectionRequest](c)
	if req == nil {
		return
	}

	if err := h.addPost.Execute(c.Request.Context(), id, userID, req.PostID, req.Note); err != nil {
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

	if err := h.removePost.Execute(c.Request.Context(), collectionID, userID, postID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// GetMyCount は認証ユーザーのコレクション総数を返す。
func (h *PostCollectionHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
