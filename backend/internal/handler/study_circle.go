package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// StudyCircleServiceInterface はStudyCircleHandlerが依存するサービスメソッドを定義する。
type StudyCircleServiceInterface interface {
	Create(circle *model.StudyCircle, memberIDs []uint) error
	GetMyCircles(userID uint, limit, offset int) ([]model.StudyCircle, int64, error)
	GetByID(id, userID uint) (*model.StudyCircle, error)
	Update(id, userID uint, name, topic, description *string) (*model.StudyCircle, error)
	Delete(id, userID uint) error
	GetMembers(circleID, userID uint) ([]model.StudyCircleMember, error)
	AddMember(circleID, userID, targetUserID uint) error
	RemoveMember(circleID, userID, targetUserID uint) error
	CreateStep(circleID, userID uint, step *model.StudyCircleStep) error
	UpdateStep(circleID, userID, stepID uint, title, description *string) (*model.StudyCircleStep, error)
	DeleteStep(circleID, userID, stepID uint) error
	ReorderSteps(circleID, userID uint, orders []model.StepOrder) error
	UpdateProgress(circleID, userID, stepID uint, isCompleted bool) error
	GetProgress(circleID, userID uint) ([]model.StudyCircleMemberProgress, error)
	CreateCheckin(circleID, userID uint, content string) (*model.StudyCircleCheckin, error)
	GetCheckins(circleID, userID uint) ([]model.StudyCircleCheckin, error)
	GetStreakRanking(circleID, userID uint) ([]model.CircleMemberStreak, error)
	GetByStatus(userID uint, status string) ([]model.StudyCircle, error)
	UpdateMemberRole(circleID, userID, targetUserID uint, role string) error
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
	input := bindJSON[dto.CreateStudyCircleRequest](c)
	if input == nil {
		return
	}

	userID := c.GetUint("userID")
	circle := &model.StudyCircle{
		Name:        input.Name,
		Topic:       input.Topic,
		Description: input.Description,
		OwnerID:     userID,
		MaxMembers:  input.MaxMembers,
	}

	if err := h.service.Create(circle, input.MemberIDs); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, circle)
}

// GetMyCircles は参加サークル一覧をページネーション付きで返す。
func (h *StudyCircleHandler) GetMyCircles(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)
	circles, total, err := h.service.GetMyCircles(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.StudyCircleListResponse{
		Circles: ensureSlice(circles),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetByID はサークル詳細を返す。
func (h *StudyCircleHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// Update はサークル情報を更新する。
func (h *StudyCircleHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[dto.UpdateStudyCircleRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	circle, err := h.service.Update(id, userID, input.Name, input.Topic, input.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, circle)
}

// Delete はサークルを削除する。
func (h *StudyCircleHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
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
	respondOK(c, ensureSlice(members))
}

// AddMember はメンバーを追加する。
func (h *StudyCircleHandler) AddMember(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[dto.AddStudyCircleMemberRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.AddMember(id, userID, input.UserID); err != nil {
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
	input := bindJSON[dto.CreateStudyCircleStepRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")

	step := &model.StudyCircleStep{
		Title:       input.Title,
		Description: input.Description,
		ResourceURL: input.ResourceURL,
		OrderIndex:  input.OrderIndex,
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
	input := bindJSON[dto.UpdateStudyCircleStepRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	step, err := h.service.UpdateStep(id, userID, stepID, input.Title, input.Description)
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
	input := bindJSON[dto.ReorderStudyCircleStepsRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.ReorderSteps(id, userID, input.Orders); err != nil {
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
	input := bindJSON[dto.UpdateStudyCircleProgressRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.service.UpdateProgress(id, userID, stepID, input.IsCompleted); err != nil {
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
	input := bindJSON[dto.CreateStudyCircleCheckinRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	checkin, err := h.service.CreateCheckin(id, userID, input.Content)
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

// GetByStatus はステータスでサークルをフィルタリングして取得する。
func (h *StudyCircleHandler) GetByStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	status := c.Param("status")

	circles, err := h.service.GetByStatus(userID, status)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(circles))
}

// UpdateMemberRole はメンバーの役割を更新する。オーナーのみ操作可能。
func (h *StudyCircleHandler) UpdateMemberRole(c *gin.Context) {
	userID := c.GetUint("userID")
	circleID, ok := parseID(c, "id")
	if !ok {
		return
	}
	targetUserID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, domain.NewError(domain.ErrCodeBadRequest, "役割の指定が必要です", err))
		return
	}

	if err := h.service.UpdateMemberRole(circleID, userID, targetUserID, req.Role); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "メンバー役割を更新しました"})
}
