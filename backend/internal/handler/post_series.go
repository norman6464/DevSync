package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostSeriesServiceInterface はPostSeriesHandlerが依存するサービスメソッドを定義する。
type PostSeriesServiceInterface interface {
	Create(series *model.PostSeries) error
	GetByID(id uint) (*model.PostSeries, error)
	GetByUserID(userID uint, page, limit int) ([]model.PostSeries, error)
	CountByUser(userID uint) (int64, error)
	Update(id, userID uint, updates *model.PostSeries) (*model.PostSeries, error)
	Delete(id, userID uint) error
	AddPost(seriesID, postID uint, orderIndex int, userID uint) error
	RemovePost(seriesID, postID, userID uint) error
	GetPosts(seriesID uint) ([]model.PostSeriesItem, error)
}

// PostSeriesHandler は投稿シリーズ関連のHTTPハンドラ。
type PostSeriesHandler struct {
	service PostSeriesServiceInterface
}

// NewPostSeriesHandler は新しいPostSeriesHandlerインスタンスを生成する。
func NewPostSeriesHandler(s PostSeriesServiceInterface) *PostSeriesHandler {
	return &PostSeriesHandler{service: s}
}

// Create は新しい投稿シリーズを作成する。
func (h *PostSeriesHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreatePostSeriesRequest](c)
	if req == nil {
		return
	}

	series := &model.PostSeries{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.service.Create(series); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, series)
}

// GetByID は指定IDのシリーズを取得する。
func (h *PostSeriesHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	series, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, series)
}

// GetByUserID は指定ユーザーのシリーズ一覧をページネーション付きで取得する。
func (h *PostSeriesHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	series, err := h.service.GetByUserID(userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.service.CountByUser(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPaginated(c, series, total, page, limit)
}

// Update は指定IDのシリーズを更新する。
func (h *PostSeriesHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdatePostSeriesRequest](c)
	if req == nil {
		return
	}

	updates := &model.PostSeries{
		Title:       req.Title,
		Description: req.Description,
	}

	series, err := h.service.Update(id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, series)
}

// Delete は指定IDのシリーズを削除する。
func (h *PostSeriesHandler) Delete(c *gin.Context) {
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

// GetPosts はシリーズ内の投稿一覧を取得する。
func (h *PostSeriesHandler) GetPosts(c *gin.Context) {
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

// AddPost はシリーズに投稿を追加する。
func (h *PostSeriesHandler) AddPost(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.AddPostToSeriesRequest](c)
	if req == nil {
		return
	}

	if err := h.service.AddPost(id, req.PostID, req.OrderIndex, userID); err != nil {
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

	if err := h.service.RemovePost(seriesID, postID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}
