package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostSeriesHandler は投稿シリーズ関連の HTTP ハンドラ。
type PostSeriesHandler struct {
	create     *usecase.CreatePostSeriesUseCase
	get        *usecase.GetPostSeriesUseCase
	list       *usecase.ListPostSeriesUseCase
	count      *usecase.CountPostSeriesUseCase
	update     *usecase.UpdatePostSeriesUseCase
	delete     *usecase.DeletePostSeriesUseCase
	addPost    *usecase.AddPostToSeriesUseCase
	removePost *usecase.RemovePostFromSeriesUseCase
	listPosts  *usecase.ListPostSeriesPostsUseCase
}

// NewPostSeriesHandler は PostSeriesHandler を生成する。
func NewPostSeriesHandler(
	create *usecase.CreatePostSeriesUseCase,
	get *usecase.GetPostSeriesUseCase,
	list *usecase.ListPostSeriesUseCase,
	count *usecase.CountPostSeriesUseCase,
	update *usecase.UpdatePostSeriesUseCase,
	deleteUC *usecase.DeletePostSeriesUseCase,
	addPost *usecase.AddPostToSeriesUseCase,
	removePost *usecase.RemovePostFromSeriesUseCase,
	listPosts *usecase.ListPostSeriesPostsUseCase,
) *PostSeriesHandler {
	return &PostSeriesHandler{
		create:     create,
		get:        get,
		list:       list,
		count:      count,
		update:     update,
		delete:     deleteUC,
		addPost:    addPost,
		removePost: removePost,
		listPosts:  listPosts,
	}
}

// createPostSeriesRequest はシリーズ作成のリクエストボディ。
type createPostSeriesRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
}

// Create は新しい投稿シリーズを作成する。
func (h *PostSeriesHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[createPostSeriesRequest](c)
	if req == nil {
		return
	}

	series := &model.PostSeries{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.create.Execute(c.Request.Context(), series); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, series)
}

// GetByID は指定IDのシリーズを取得する。
func (h *PostSeriesHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	handleGetByIDPublic(c, func(id uint) (*model.PostSeries, error) {
		return h.get.Execute(ctx, id)
	})
}

// GetByUserID は指定ユーザーのシリーズ一覧をページネーション付きで取得する。
func (h *PostSeriesHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	series, err := h.list.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, series, total, page, limit)
}

// GetMySeries は認証ユーザーの投稿シリーズ一覧を取得する。
func (h *PostSeriesHandler) GetMySeries(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)

	series, err := h.list.Execute(c.Request.Context(), userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, series, total, page, limit)
}

// GetMySeriesCount は認証ユーザーの投稿シリーズ数を取得する。
func (h *PostSeriesHandler) GetMySeriesCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

// updatePostSeriesRequest はシリーズ更新のリクエストボディ。
type updatePostSeriesRequest struct {
	Title       string `json:"title" binding:"omitempty,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
}

// Update は指定IDのシリーズを更新する。
func (h *PostSeriesHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updatePostSeriesRequest](c)
	if req == nil {
		return
	}

	updates := &model.PostSeries{
		Title:       req.Title,
		Description: req.Description,
	}

	series, err := h.update.Execute(c.Request.Context(), id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, series)
}

// Delete は指定IDのシリーズを削除する。
func (h *PostSeriesHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	handleDelete(c, func(id, userID uint) error {
		return h.delete.Execute(ctx, id, userID)
	})
}

// GetPosts はシリーズ内の投稿一覧を取得する。
func (h *PostSeriesHandler) GetPosts(c *gin.Context) {
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

// addPostToSeriesRequest はシリーズへの投稿追加リクエスト。
type addPostToSeriesRequest struct {
	PostID     uint `json:"post_id" binding:"required"`
	OrderIndex int  `json:"order_index" binding:"omitempty,min=0"`
}

// AddPost はシリーズに投稿を追加する。
func (h *PostSeriesHandler) AddPost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[addPostToSeriesRequest](c)
	if req == nil {
		return
	}

	if err := h.addPost.Execute(c.Request.Context(), id, req.PostID, req.OrderIndex, userID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, nil)
}

// RemovePost はシリーズから投稿を削除する。
func (h *PostSeriesHandler) RemovePost(c *gin.Context) {
	userID := c.GetUint("userID")
	seriesID, ok := parseID(c, "id")
	if !ok {
		return
	}
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	if err := h.removePost.Execute(c.Request.Context(), seriesID, postID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}
