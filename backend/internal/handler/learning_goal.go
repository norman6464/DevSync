package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningGoalHandler は学習目標関連のHTTPハンドラ。
// 学習目標のCRUD・統計情報の取得を処理する。
type LearningGoalHandler struct {
	create        *usecase.CreateLearningGoalUseCase
	get           *usecase.GetLearningGoalUseCase
	list          *usecase.ListLearningGoalsUseCase
	listActive    *usecase.ListActiveLearningGoalsUseCase
	byCategory    *usecase.ListLearningGoalsByCategoryUseCase
	byStatus      *usecase.ListLearningGoalsByStatusUseCase
	stats         *usecase.GetLearningGoalStatsUseCase
	update        *usecase.UpdateLearningGoalUseCase
	deadlineAlert *usecase.GetGoalDeadlineAlertsUseCase
	duplicate     *usecase.DuplicateLearningGoalUseCase
	toggleShare   *usecase.ToggleLearningGoalShareUseCase
	listPublic    *usecase.ListPublicLearningGoalsUseCase
	listPublicBy  *usecase.ListPublicLearningGoalsByUserUseCase
	count         *usecase.CountLearningGoalsUseCase
	remove        *usecase.DeleteLearningGoalUseCase
	batchProgress *usecase.BatchUpdateGoalProgressUseCase
	forecast      *usecase.GetGoalForecastUseCase
}

// NewLearningGoalHandler は新しいLearningGoalHandlerインスタンスを生成する。
func NewLearningGoalHandler(
	create *usecase.CreateLearningGoalUseCase,
	get *usecase.GetLearningGoalUseCase,
	list *usecase.ListLearningGoalsUseCase,
	listActive *usecase.ListActiveLearningGoalsUseCase,
	byCategory *usecase.ListLearningGoalsByCategoryUseCase,
	byStatus *usecase.ListLearningGoalsByStatusUseCase,
	stats *usecase.GetLearningGoalStatsUseCase,
	update *usecase.UpdateLearningGoalUseCase,
	deadlineAlert *usecase.GetGoalDeadlineAlertsUseCase,
	duplicate *usecase.DuplicateLearningGoalUseCase,
	toggleShare *usecase.ToggleLearningGoalShareUseCase,
	listPublic *usecase.ListPublicLearningGoalsUseCase,
	listPublicBy *usecase.ListPublicLearningGoalsByUserUseCase,
	count *usecase.CountLearningGoalsUseCase,
	remove *usecase.DeleteLearningGoalUseCase,
	batchProgress *usecase.BatchUpdateGoalProgressUseCase,
	forecast *usecase.GetGoalForecastUseCase,
) *LearningGoalHandler {
	return &LearningGoalHandler{
		create: create, get: get, list: list, listActive: listActive,
		byCategory: byCategory, byStatus: byStatus, stats: stats, update: update,
		deadlineAlert: deadlineAlert, duplicate: duplicate, toggleShare: toggleShare,
		listPublic: listPublic, listPublicBy: listPublicBy, count: count,
		remove: remove, batchProgress: batchProgress, forecast: forecast,
	}
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

	if err := h.create.Execute(c.Request.Context(), goal); err != nil {
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

	goal, err := h.update.Execute(c.Request.Context(), goalID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// Delete は指定された学習目標を削除する。
func (h *LearningGoalHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// GetByID は指定されたIDの学習目標を取得する。所有者のみ取得可能。
func (h *LearningGoalHandler) GetByID(c *gin.Context) {
	handleGetByID(c, func(id, userID uint) (*model.LearningGoal, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// GetByUserID は指定されたユーザーの学習目標一覧を取得する。
func (h *LearningGoalHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)
	goals, total, err := h.list.Execute(c.Request.Context(), userID, limit, offset)
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
	goals, total, err := h.list.Execute(c.Request.Context(), userID, limit, offset)
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

	alerts, err := h.deadlineAlert.Execute(c.Request.Context(), userID)
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

	goals, err := h.byCategory.Execute(c.Request.Context(), userID, category)
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

	goals, err := h.byStatus.Execute(c.Request.Context(), userID, status)
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

	goal, err := h.duplicate.Execute(c.Request.Context(), id, userID)
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

	stats, err := h.stats.Execute(c.Request.Context(), userID)
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

	goal, err := h.toggleShare.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, goal)
}

// GetPublicGoals は全ユーザーの公開済み学習目標一覧を返す。
func (h *LearningGoalHandler) GetPublicGoals(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	goals, total, err := h.listPublic.Execute(c.Request.Context(), limit, offset)
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
	goals, total, err := h.listPublicBy.Execute(c.Request.Context(), userID, limit, offset)
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

	updates := make([]usecase.GoalProgressUpdate, len(req.Updates))
	for i, u := range req.Updates {
		updates[i] = usecase.GoalProgressUpdate{GoalID: u.GoalID, Progress: u.Progress}
	}

	results, err := h.batchProgress.Execute(c.Request.Context(), userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(results))
}

// GetMyStats は認証ユーザー自身の学習目標統計情報を取得する。
func (h *LearningGoalHandler) GetMyStats(c *gin.Context) {
	userID := c.GetUint("userID")

	stats, err := h.stats.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}

// GetActiveGoals は認証ユーザーのアクティブな学習目標のみを取得する。
func (h *LearningGoalHandler) GetActiveGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	goals, err := h.listActive.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(goals))
}

// GetMyCount は認証ユーザー自身の学習目標総数を返す。
func (h *LearningGoalHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// GetForecast は認証ユーザーのアクティブ目標の達成予測一覧を返す。
func (h *LearningGoalHandler) GetForecast(c *gin.Context) {
	userID := c.GetUint("userID")

	forecasts, err := h.forecast.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(forecasts))
}
