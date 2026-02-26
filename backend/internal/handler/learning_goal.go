package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningGoalServiceInterface はLearningGoalServiceが実装すべきインターフェース。
type LearningGoalServiceInterface interface {
	Create(goal *model.LearningGoal) error
	GetByID(id, userID uint) (*model.LearningGoal, error)
	GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetByCategory(userID uint, category string) ([]model.LearningGoal, error)
	GetByStatus(userID uint, status string) ([]model.LearningGoal, error)
	GetStats(userID uint) (*model.LearningGoalStats, error)
	Update(id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error)
	Delete(id, userID uint) error
	GetDeadlineAlerts(userID uint) ([]model.GoalDeadlineAlert, error)
	Duplicate(id, userID uint) (*model.LearningGoal, error)
	ToggleShare(id, userID uint) (*model.LearningGoal, error)
	GetPublicGoals(limit, offset int) ([]model.LearningGoal, int64, error)
	GetPublicByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetActiveByUserID(userID uint) ([]model.LearningGoal, error)
	GetForecast(userID uint) ([]model.GoalForecast, error)
	BatchUpdateProgress(userID uint, updates []struct {
		GoalID   uint
		Progress int
	}) ([]model.LearningGoal, error)
}

// LearningGoalHandler は学習目標関連のHTTPハンドラ。
// 学習目標のCRUD・統計情報の取得を処理する。
type LearningGoalHandler struct {
	service LearningGoalServiceInterface
}

// NewLearningGoalHandler は新しいLearningGoalHandlerインスタンスを生成する。
func NewLearningGoalHandler(s LearningGoalServiceInterface) *LearningGoalHandler {
	return &LearningGoalHandler{service: s}
}

// Create は新しい学習目標を作成する。
func (h *LearningGoalHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreateGoalRequest](c)
	if req == nil {
		return
	}

	goal := &model.LearningGoal{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    model.GoalCategory(req.Category),
		Status:      model.GoalStatusActive,
		Progress:    0,
	}

	// 目標日が指定されている場合はパースして設定
	if targetDate, ok := parseDateParam(req.TargetDate); ok && !targetDate.IsZero() {
		goal.TargetDate = &targetDate
	}

	if err := h.service.Create(goal); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, goal)
}

// Update は指定された学習目標を更新する。
func (h *LearningGoalHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	goalID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateGoalRequest](c)
	if req == nil {
		return
	}

	updates := &model.LearningGoal{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.Category != nil {
		updates.Category = model.GoalCategory(*req.Category)
	}
	if req.TargetDate != nil {
		if *req.TargetDate == "" {
			updates.TargetDate = nil
		} else if targetDate, ok := parseDateParam(*req.TargetDate); ok {
			updates.TargetDate = &targetDate
		}
	}
	if req.Progress != nil {
		updates.Progress = *req.Progress
	}
	if req.Status != nil {
		updates.Status = model.GoalStatus(*req.Status)
	}

	goal, err := h.service.Update(goalID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// Delete は指定された学習目標を削除する。
func (h *LearningGoalHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// GetByID は指定されたIDの学習目標を取得する。所有者のみ取得可能。
func (h *LearningGoalHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// GetByUserID は指定されたユーザーの学習目標一覧を取得する。
func (h *LearningGoalHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)
	goals, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.GoalListResponse{
		Goals:  goals,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetMyGoals は認証ユーザー自身の学習目標一覧を取得する。
func (h *LearningGoalHandler) GetMyGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	limit, offset := parseLimitOffset(c)
	goals, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.GoalListResponse{
		Goals:  goals,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetDeadlineAlerts は認証ユーザーのデッドラインアラートを取得する。
func (h *LearningGoalHandler) GetDeadlineAlerts(c *gin.Context) {
	userID := c.GetUint("userID")

	alerts, err := h.service.GetDeadlineAlerts(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(alerts))
}

// GetByCategory は認証ユーザーの学習目標をカテゴリでフィルタリングして取得する。
func (h *LearningGoalHandler) GetByCategory(c *gin.Context) {
	userID := c.GetUint("userID")
	category := c.Param("category")

	goals, err := h.service.GetByCategory(userID, category)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(goals))
}

// GetByStatus は認証ユーザーの学習目標をステータスでフィルタリングして取得する。
func (h *LearningGoalHandler) GetByStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Param("status")

	goals, err := h.service.GetByStatus(userID, status)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(goals))
}

// Duplicate は学習目標を複製する。所有者のみ複製可能。
func (h *LearningGoalHandler) Duplicate(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	goal, err := h.service.Duplicate(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, goal)
}

// GetStats は指定されたユーザーの学習目標統計情報を取得する。
func (h *LearningGoalHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	stats, err := h.service.GetStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}

// ToggleShare は学習目標の公開/非公開を切り替える。所有者のみ操作可能。
func (h *LearningGoalHandler) ToggleShare(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	goal, err := h.service.ToggleShare(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// GetPublicGoals は全ユーザーの公開済み学習目標一覧を返す。
func (h *LearningGoalHandler) GetPublicGoals(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	goals, total, err := h.service.GetPublicGoals(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.GoalListResponse{
		Goals:  goals,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetPublicByUserID は指定ユーザーの公開済み学習目標一覧を返す。
func (h *LearningGoalHandler) GetPublicByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)
	goals, total, err := h.service.GetPublicByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.GoalListResponse{
		Goals:  goals,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// BatchUpdateProgress は複数の学習目標の進捗を一括更新する。
func (h *LearningGoalHandler) BatchUpdateProgress(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.BatchUpdateProgressRequest](c)
	if req == nil {
		return
	}

	updates := make([]struct {
		GoalID   uint
		Progress int
	}, len(req.Updates))
	for i, u := range req.Updates {
		updates[i].GoalID = u.GoalID
		updates[i].Progress = u.Progress
	}

	results, err := h.service.BatchUpdateProgress(userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(results))
}

// GetForecast は認証ユーザーのアクティブ目標の達成予測一覧を返す。
func (h *LearningGoalHandler) GetForecast(c *gin.Context) {
	userID := c.GetUint("userID")

	forecasts, err := h.service.GetForecast(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(forecasts))
}
