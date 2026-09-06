package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningResourceHandler は学習リソース関連のHTTPハンドラ。
// 学習リソースのCRUD・検索・いいね・保存を処理する。
type LearningResourceHandler struct {
	create        *usecase.CreateLearningResourceUseCase
	get           *usecase.GetLearningResourceUseCase
	listByUser    *usecase.ListLearningResourcesByUserUseCase
	listPublic    *usecase.ListPublicLearningResourcesUseCase
	listByDiff    *usecase.ListLearningResourcesByDifficultyUseCase
	search        *usecase.SearchLearningResourcesUseCase
	update        *usecase.UpdateLearningResourceUseCase
	updateVisible *usecase.UpdateLearningResourceVisibilityUseCase
	remove        *usecase.DeleteLearningResourceUseCase
	like          *usecase.LikeLearningResourceUseCase
	unlike        *usecase.UnlikeLearningResourceUseCase
	hasLiked      *usecase.HasLikedLearningResourceUseCase
	save          *usecase.SaveLearningResourceUseCase
	unsave        *usecase.UnsaveLearningResourceUseCase
	hasSaved      *usecase.HasSavedLearningResourceUseCase
	listSaved     *usecase.ListSavedLearningResourcesUseCase
	count         *usecase.CountLearningResourcesUseCase
}

// NewLearningResourceHandler は新しいLearningResourceHandlerインスタンスを生成する。
func NewLearningResourceHandler(
	create *usecase.CreateLearningResourceUseCase,
	get *usecase.GetLearningResourceUseCase,
	listByUser *usecase.ListLearningResourcesByUserUseCase,
	listPublic *usecase.ListPublicLearningResourcesUseCase,
	listByDiff *usecase.ListLearningResourcesByDifficultyUseCase,
	search *usecase.SearchLearningResourcesUseCase,
	update *usecase.UpdateLearningResourceUseCase,
	updateVisible *usecase.UpdateLearningResourceVisibilityUseCase,
	remove *usecase.DeleteLearningResourceUseCase,
	like *usecase.LikeLearningResourceUseCase,
	unlike *usecase.UnlikeLearningResourceUseCase,
	hasLiked *usecase.HasLikedLearningResourceUseCase,
	save *usecase.SaveLearningResourceUseCase,
	unsave *usecase.UnsaveLearningResourceUseCase,
	hasSaved *usecase.HasSavedLearningResourceUseCase,
	listSaved *usecase.ListSavedLearningResourcesUseCase,
	count *usecase.CountLearningResourcesUseCase,
) *LearningResourceHandler {
	return &LearningResourceHandler{
		create: create, get: get, listByUser: listByUser, listPublic: listPublic,
		listByDiff: listByDiff, search: search, update: update, updateVisible: updateVisible,
		remove: remove, like: like, unlike: unlike, hasLiked: hasLiked,
		save: save, unsave: unsave, hasSaved: hasSaved, listSaved: listSaved, count: count,
	}
}

// createResourceRequest は学習リソース作成のリクエストボディ。
type createResourceRequest struct {
	Title       string `json:"title" binding:"required,max=300"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"omitempty,http_url,max=2000"`
	Category    string `json:"category" binding:"required"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url" binding:"omitempty,http_url,max=2000"`
	IsPublic    *bool  `json:"is_public"`
}

// Create は新しい学習リソースを作成する。
func (h *LearningResourceHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[createResourceRequest](c)
	if req == nil {
		return
	}

	// 公開設定のデフォルト値はtrue
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	resource := &model.LearningResource{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    model.ResourceCategory(req.Category),
		Difficulty:  model.ResourceDifficulty(req.Difficulty),
		Tags:        req.Tags,
		ImageURL:    req.ImageURL,
		IsPublic:    isPublic,
	}

	if err := h.create.Execute(c.Request.Context(), resource); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, resource)
}

// resourceDetailResponse は学習リソース詳細レスポンス（いいね・保存状態付き）。
type resourceDetailResponse struct {
	Resource model.LearningResource `json:"resource"`
	HasLiked bool                   `json:"has_liked"`
	HasSaved bool                   `json:"has_saved"`
}

// GetByID は指定されたIDの学習リソースを取得する。
func (h *LearningResourceHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")

	resource, err := h.get.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	// 現在のユーザーがいいね・保存済みかを確認
	hasLiked, _ := h.hasLiked.Execute(c.Request.Context(), userID, id)
	hasSaved, _ := h.hasSaved.Execute(c.Request.Context(), userID, id)

	respondOK(c, resourceDetailResponse{
		Resource: *resource,
		HasLiked: hasLiked,
		HasSaved: hasSaved,
	})
}

// resourceListResponse は学習リソース一覧レスポンス。
type resourceListResponse struct {
	Resources []model.LearningResource `json:"resources"`
	Total     int64                    `json:"total"`
	Limit     int                      `json:"limit"`
	Offset    int                      `json:"offset"`
}

// GetByUserID は指定されたユーザーの学習リソース一覧をページネーション付きで取得する。
func (h *LearningResourceHandler) GetByUserID(c *gin.Context) {
	targetUserID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	currentUserID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.listByUser.Execute(c.Request.Context(), targetUserID, currentUserID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetMyResources は認証ユーザーの学習リソース一覧を取得する。
func (h *LearningResourceHandler) GetMyResources(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.listByUser.Execute(c.Request.Context(), userID, userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetPublic は公開学習リソース一覧をページネーション・フィルター付きで取得する。
func (h *LearningResourceHandler) GetPublic(c *gin.Context) {
	limit, offset := parseLimitOffset(c)
	category := c.Query("category")
	difficulty := c.Query("difficulty")

	resources, total, err := h.listPublic.Execute(c.Request.Context(), limit, offset, category, difficulty)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// Search はキーワードで学習リソースを検索する。
func (h *LearningResourceHandler) Search(c *gin.Context) {
	query, ok := parseSearchQuery(c)
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.search.Execute(c.Request.Context(), query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// updateResourceRequest は学習リソース更新のリクエストボディ。
type updateResourceRequest struct {
	Title       string `json:"title" binding:"max=300"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"omitempty,http_url,max=2000"`
	Category    string `json:"category"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url" binding:"omitempty,http_url,max=2000"`
	IsPublic    *bool  `json:"is_public"`
}

// Update は指定された学習リソースを更新する。
func (h *LearningResourceHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updateResourceRequest](c)
	if req == nil {
		return
	}

	updates := &model.LearningResource{}
	if req.Title != "" {
		updates.Title = req.Title
	}
	if req.Description != "" {
		updates.Description = req.Description
	}
	if req.URL != "" {
		updates.URL = req.URL
	}
	if req.Category != "" {
		updates.Category = model.ResourceCategory(req.Category)
	}
	if req.Difficulty != "" {
		updates.Difficulty = model.ResourceDifficulty(req.Difficulty)
	}
	if req.Tags != "" {
		updates.Tags = req.Tags
	}
	if req.ImageURL != "" {
		updates.ImageURL = req.ImageURL
	}

	resource, err := h.update.Execute(c.Request.Context(), id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// 公開設定が指定されている場合は別途更新
	if req.IsPublic != nil {
		resource, err = h.updateVisible.Execute(c.Request.Context(), id, userID, *req.IsPublic)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, resource)
}

// Delete は指定された学習リソースを削除する。
func (h *LearningResourceHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// Like は学習リソースにいいねする。
func (h *LearningResourceHandler) Like(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.like.Execute(c.Request.Context(), userID, id)
	}, "リソースにいいねしました")
}

// Unlike は学習リソースのいいねを取り消す。
func (h *LearningResourceHandler) Unlike(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.unlike.Execute(c.Request.Context(), userID, id)
	}, "いいねを取り消しました")
}

// SaveResource は学習リソースを保存する。
func (h *LearningResourceHandler) SaveResource(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.save.Execute(c.Request.Context(), userID, id)
	}, "リソースを保存しました")
}

// UnsaveResource は学習リソースの保存を取り消す。
func (h *LearningResourceHandler) UnsaveResource(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.unsave.Execute(c.Request.Context(), userID, id)
	}, "保存を解除しました")
}

// GetByDifficulty は難易度別の公開学習リソースを取得する。
func (h *LearningResourceHandler) GetByDifficulty(c *gin.Context) {
	difficulty := c.Param("difficulty")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.listByDiff.Execute(c.Request.Context(), difficulty, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetSaved は認証ユーザーの保存済み学習リソース一覧を取得する。
func (h *LearningResourceHandler) GetSaved(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.listSaved.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetMyCount は認証ユーザーの学習リソース総数を返す。
func (h *LearningResourceHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
