package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// BookmarkCollectionHandler はブックマークコレクション関連の HTTP ハンドラ。
type BookmarkCollectionHandler struct {
	create     *usecase.CreateBookmarkCollectionUseCase
	list       *usecase.ListBookmarkCollectionsUseCase
	update     *usecase.UpdateBookmarkCollectionUseCase
	delete     *usecase.DeleteBookmarkCollectionUseCase
	addPost    *usecase.AddPostToBookmarkCollectionUseCase
	removePost *usecase.RemovePostFromBookmarkCollectionUseCase
	listPosts  *usecase.ListBookmarkCollectionPostsUseCase
	count      *usecase.CountBookmarkCollectionsUseCase
}

// NewBookmarkCollectionHandler は BookmarkCollectionHandler を生成する。
func NewBookmarkCollectionHandler(
	create *usecase.CreateBookmarkCollectionUseCase,
	list *usecase.ListBookmarkCollectionsUseCase,
	update *usecase.UpdateBookmarkCollectionUseCase,
	deleteUC *usecase.DeleteBookmarkCollectionUseCase,
	addPost *usecase.AddPostToBookmarkCollectionUseCase,
	removePost *usecase.RemovePostFromBookmarkCollectionUseCase,
	listPosts *usecase.ListBookmarkCollectionPostsUseCase,
	count *usecase.CountBookmarkCollectionsUseCase,
) *BookmarkCollectionHandler {
	return &BookmarkCollectionHandler{
		create:     create,
		list:       list,
		update:     update,
		delete:     deleteUC,
		addPost:    addPost,
		removePost: removePost,
		listPosts:  listPosts,
		count:      count,
	}
}

// bookmarkCollectionRequest はブックマークコレクション作成・更新リクエスト。
type bookmarkCollectionRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Color       string `json:"color" binding:"omitempty,max=50"`
}

// Create は新しいブックマークコレクションを作成する。
func (h *BookmarkCollectionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[bookmarkCollectionRequest](c)
	if input == nil {
		return
	}

	collection := &model.BookmarkCollection{
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	}

	if err := h.create.Execute(c.Request.Context(), collection); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, collection)
}

// GetMyCollections はユーザー自身のコレクション一覧を返す。
func (h *BookmarkCollectionHandler) GetMyCollections(c *gin.Context) {
	userID := c.GetUint("userID")

	collections, err := h.list.Execute(c.Request.Context(), userID)
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

	input := bindJSON[bookmarkCollectionRequest](c)
	if input == nil {
		return
	}

	updated, err := h.update.Execute(c.Request.Context(), id, userID, &model.BookmarkCollection{
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

	if err := h.delete.Execute(c.Request.Context(), id, userID); err != nil {
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

	if err := h.addPost.Execute(c.Request.Context(), collectionID, postID, userID); err != nil {
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

	if err := h.removePost.Execute(c.Request.Context(), collectionID, postID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// bookmarkCollectionPostsResponse はコレクション内投稿一覧レスポンス。
type bookmarkCollectionPostsResponse struct {
	Posts []model.Post `json:"posts"`
	Total int64        `json:"total"`
}

// GetPosts はコレクション内の投稿一覧を返す。
func (h *BookmarkCollectionHandler) GetPosts(c *gin.Context) {
	collectionID, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	posts, total, err := h.listPosts.Execute(c.Request.Context(), collectionID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, bookmarkCollectionPostsResponse{
		Posts: ensureSlice(posts),
		Total: total,
	})
}

// GetMyCount は認証ユーザーのコレクション総数を返す。
func (h *BookmarkCollectionHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
