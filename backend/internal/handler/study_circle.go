package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StudyCircleServiceInterface はStudyCircleHandlerが依存するサービスメソッドを定義する。
type StudyCircleServiceInterface interface {
	Create(circle *model.StudyCircle, memberIDs []uint) error
	GetMyCircles(userID uint) ([]model.StudyCircle, error)
	GetByID(id, userID uint) (*model.StudyCircle, error)
	Update(id, userID uint, name, topic, description *string) (*model.StudyCircle, error)
	Delete(id, userID uint) error
	GetMembers(circleID, userID uint) ([]model.StudyCircleMember, error)
	AddMember(circleID, userID, targetUserID uint) error
	RemoveMember(circleID, userID, targetUserID uint) error
	CreateStep(circleID, userID uint, step *model.StudyCircleStep) error
	UpdateStep(circleID, userID, stepID uint, title, description *string) (*model.StudyCircleStep, error)
	DeleteStep(circleID, userID, stepID uint) error
	ReorderSteps(circleID, userID uint, orders []repository.StepOrder) error
	UpdateProgress(circleID, userID, stepID uint, isCompleted bool) error
	GetProgress(circleID, userID uint) ([]model.StudyCircleMemberProgress, error)
	CreateCheckin(circleID, userID uint, content string) (*model.StudyCircleCheckin, error)
	GetCheckins(circleID, userID uint) ([]model.StudyCircleCheckin, error)
	GetStreakRanking(circleID, userID uint) ([]model.CircleMemberStreak, error)
}

// StudyCircleHandler はスタディサークル関連のHTTPリクエストを処理する。
type StudyCircleHandler struct {
	service StudyCircleServiceInterface
}

// NewStudyCircleHandler は新しいStudyCircleHandlerインスタンスを生成する。
func NewStudyCircleHandler(svc StudyCircleServiceInterface) *StudyCircleHandler {
	return &StudyCircleHandler{service: svc}
}

// Create はサークルを作成する。
func (h *StudyCircleHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Topic       string `json:"topic" binding:"required"`
		Description string `json:"description"`
		MaxMembers  int    `json:"max_members"`
		MemberIDs   []uint `json:"member_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")
	circle := &model.StudyCircle{
		Name:        req.Name,
		Topic:       req.Topic,
		Description: req.Description,
		OwnerID:     userID,
		MaxMembers:  req.MaxMembers,
	}

	if err := h.service.Create(circle, req.MemberIDs); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, circle)
}

// GetMyCircles は参加サークル一覧を返す。
func (h *StudyCircleHandler) GetMyCircles(c *gin.Context) {
	userID := c.GetUint("userID")
	circles, err := h.service.GetMyCircles(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, circles)
}

// GetByID はサークル詳細を返す。
func (h *StudyCircleHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	circle, err := h.service.GetByID(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, circle)
}

// Update はサークル情報を更新する。
func (h *StudyCircleHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		Name        *string `json:"name"`
		Topic       *string `json:"topic"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	circle, err := h.service.Update(id, userID, req.Name, req.Topic, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, circle)
}

// Delete はサークルを削除する。
func (h *StudyCircleHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// GetMembers はメンバー一覧を返す。
func (h *StudyCircleHandler) GetMembers(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	members, err := h.service.GetMembers(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, members)
}

// AddMember はメンバーを追加する。
func (h *StudyCircleHandler) AddMember(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.AddMember(id, userID, req.UserID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("member added"))
}

// RemoveMember はメンバーを除外/退出する。
func (h *StudyCircleHandler) RemoveMember(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	targetID, ok := parseID(c, "userId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.RemoveMember(id, userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("member removed"))
}

// CreateStep はステップを追加する。
func (h *StudyCircleHandler) CreateStep(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ResourceURL string `json:"resource_url"`
		OrderIndex  int    `json:"order_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")

	step := &model.StudyCircleStep{
		Title:       req.Title,
		Description: req.Description,
		ResourceURL: req.ResourceURL,
		OrderIndex:  req.OrderIndex,
	}
	if err := h.service.CreateStep(id, userID, step); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, step)
}

// UpdateStep はステップを更新する。
func (h *StudyCircleHandler) UpdateStep(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}
	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	step, err := h.service.UpdateStep(id, userID, stepID, req.Title, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, step)
}

// DeleteStep はステップを削除する。
func (h *StudyCircleHandler) DeleteStep(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.DeleteStep(id, userID, stepID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// ReorderSteps はステップの順序を変更する。
func (h *StudyCircleHandler) ReorderSteps(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		Orders []repository.StepOrder `json:"orders" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.ReorderSteps(id, userID, req.Orders); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("reordered"))
}

// UpdateProgress は自分のステップ進捗を更新する。
func (h *StudyCircleHandler) UpdateProgress(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}
	var req struct {
		IsCompleted bool `json:"is_completed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.UpdateProgress(id, userID, stepID, req.IsCompleted); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("progress updated"))
}

// GetProgress は全メンバーの進捗を返す。
func (h *StudyCircleHandler) GetProgress(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	progress, err := h.service.GetProgress(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, progress)
}

// CreateCheckin は日次チェックインを作成する。
func (h *StudyCircleHandler) CreateCheckin(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("userID")
	checkin, err := h.service.CreateCheckin(id, userID, req.Content)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, checkin)
}

// GetCheckins はチェックイン履歴を返す。
func (h *StudyCircleHandler) GetCheckins(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	checkins, err := h.service.GetCheckins(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, checkins)
}

// GetStreakRanking はストリークランキングを返す。
func (h *StudyCircleHandler) GetStreakRanking(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	ranking, err := h.service.GetStreakRanking(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ranking)
}
