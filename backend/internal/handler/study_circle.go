package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// StudyCircleHandler はスタディサークル関連のHTTPリクエストを処理する。
type StudyCircleHandler struct {
	create           *usecase.CreateStudyCircleUseCase
	listMine         *usecase.ListMyStudyCirclesUseCase
	listByStatus     *usecase.ListStudyCirclesByStatusUseCase
	get              *usecase.GetStudyCircleUseCase
	update           *usecase.UpdateStudyCircleUseCase
	remove           *usecase.DeleteStudyCircleUseCase
	listMembers      *usecase.ListStudyCircleMembersUseCase
	addMember        *usecase.AddStudyCircleMemberUseCase
	updateMemberRole *usecase.UpdateStudyCircleMemberRoleUseCase
	removeMember     *usecase.RemoveStudyCircleMemberUseCase
	createStep       *usecase.CreateStudyCircleStepUseCase
	updateStep       *usecase.UpdateStudyCircleStepUseCase
	deleteStep       *usecase.DeleteStudyCircleStepUseCase
	reorderSteps     *usecase.ReorderStudyCircleStepsUseCase
	updateProgress   *usecase.UpdateStudyCircleProgressUseCase
	listProgress     *usecase.ListStudyCircleProgressUseCase
	createCheckin    *usecase.CreateStudyCircleCheckinUseCase
	listCheckins     *usecase.ListStudyCircleCheckinsUseCase
	streakRanking    *usecase.GetStudyCircleStreakRankingUseCase
	search           *usecase.SearchStudyCirclesUseCase
	count            *usecase.CountStudyCirclesUseCase
}

// NewStudyCircleHandler は新しいStudyCircleHandlerインスタンスを生成する。
func NewStudyCircleHandler(
	create *usecase.CreateStudyCircleUseCase,
	listMine *usecase.ListMyStudyCirclesUseCase,
	listByStatus *usecase.ListStudyCirclesByStatusUseCase,
	get *usecase.GetStudyCircleUseCase,
	update *usecase.UpdateStudyCircleUseCase,
	remove *usecase.DeleteStudyCircleUseCase,
	listMembers *usecase.ListStudyCircleMembersUseCase,
	addMember *usecase.AddStudyCircleMemberUseCase,
	updateMemberRole *usecase.UpdateStudyCircleMemberRoleUseCase,
	removeMember *usecase.RemoveStudyCircleMemberUseCase,
	createStep *usecase.CreateStudyCircleStepUseCase,
	updateStep *usecase.UpdateStudyCircleStepUseCase,
	deleteStep *usecase.DeleteStudyCircleStepUseCase,
	reorderSteps *usecase.ReorderStudyCircleStepsUseCase,
	updateProgress *usecase.UpdateStudyCircleProgressUseCase,
	listProgress *usecase.ListStudyCircleProgressUseCase,
	createCheckin *usecase.CreateStudyCircleCheckinUseCase,
	listCheckins *usecase.ListStudyCircleCheckinsUseCase,
	streakRanking *usecase.GetStudyCircleStreakRankingUseCase,
	search *usecase.SearchStudyCirclesUseCase,
	count *usecase.CountStudyCirclesUseCase,
) *StudyCircleHandler {
	return &StudyCircleHandler{
		create: create, listMine: listMine, listByStatus: listByStatus, get: get,
		update: update, remove: remove, listMembers: listMembers, addMember: addMember,
		updateMemberRole: updateMemberRole, removeMember: removeMember,
		createStep: createStep, updateStep: updateStep, deleteStep: deleteStep,
		reorderSteps: reorderSteps, updateProgress: updateProgress, listProgress: listProgress,
		createCheckin: createCheckin, listCheckins: listCheckins, streakRanking: streakRanking,
		search: search, count: count,
	}
}

// createStudyCircleRequest はスタディサークル作成リクエストのDTO。
type createStudyCircleRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Topic       string `json:"topic" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	MaxMembers  int    `json:"max_members" binding:"omitempty,min=0,max=1000"`
	MemberIDs   []uint `json:"member_ids" binding:"omitempty,max=100"`
}

// Create はサークルを作成する。
func (h *StudyCircleHandler) Create(c *gin.Context) {
	input := bindJSON[createStudyCircleRequest](c)
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

	if err := h.create.Execute(c.Request.Context(), circle, input.MemberIDs); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, circle)
}

// studyCircleListResponse はスタディサークル一覧レスポンス（ページネーション付き）。
type studyCircleListResponse struct {
	Circles []model.StudyCircle `json:"circles"`
	Total   int64               `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
}

// GetMyCircles は参加サークル一覧をページネーション付きで返す。
func (h *StudyCircleHandler) GetMyCircles(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)
	circles, total, err := h.listMine.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, studyCircleListResponse{
		Circles: ensureSlice(circles),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetByID はサークル詳細を返す。
func (h *StudyCircleHandler) GetByID(c *gin.Context) {
	handleGetByID(c, func(id, userID uint) (*model.StudyCircle, error) {
		return h.get.Execute(c.Request.Context(), id, userID)
	})
}

// updateStudyCircleRequest はスタディサークル更新リクエストのDTO。
type updateStudyCircleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Topic       *string `json:"topic" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

// Update はサークル情報を更新する。
func (h *StudyCircleHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[updateStudyCircleRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	circle, err := h.update.Execute(c.Request.Context(), id, userID, input.Name, input.Topic, input.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, circle)
}

// Delete はサークルを削除する。
func (h *StudyCircleHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// GetMembers はメンバー一覧を返す。
func (h *StudyCircleHandler) GetMembers(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")
	members, err := h.listMembers.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(members))
}

// addStudyCircleMemberRequest はスタディサークルメンバー追加リクエストのDTO。
type addStudyCircleMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// AddMember はメンバーを追加する。
func (h *StudyCircleHandler) AddMember(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[addStudyCircleMemberRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.addMember.Execute(c.Request.Context(), id, userID, input.UserID); err != nil {
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
	if err := h.removeMember.Execute(c.Request.Context(), id, userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("member removed"))
}

// createStudyCircleStepRequest はスタディサークルステップ作成リクエストのDTO。
type createStudyCircleStepRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	ResourceURL string `json:"resource_url" binding:"omitempty,max=2000"`
	OrderIndex  int    `json:"order_index" binding:"omitempty,min=0"`
}

// CreateStep はステップを追加する。
func (h *StudyCircleHandler) CreateStep(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[createStudyCircleStepRequest](c)
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
	if err := h.createStep.Execute(c.Request.Context(), id, userID, step); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, step)
}

// updateStudyCircleStepRequest はスタディサークルステップ更新リクエストのDTO。
type updateStudyCircleStepRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
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
	input := bindJSON[updateStudyCircleStepRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	step, err := h.updateStep.Execute(c.Request.Context(), id, userID, stepID, input.Title, input.Description)
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
	if err := h.deleteStep.Execute(c.Request.Context(), id, userID, stepID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// reorderStudyCircleStepsRequest はスタディサークルステップ並べ替えリクエストのDTO。
type reorderStudyCircleStepsRequest struct {
	Orders []model.StepOrder `json:"orders" binding:"required"`
}

// ReorderSteps はステップの順序を変更する。
func (h *StudyCircleHandler) ReorderSteps(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[reorderStudyCircleStepsRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.reorderSteps.Execute(c.Request.Context(), id, userID, input.Orders); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("reordered"))
}

// updateStudyCircleProgressRequest はスタディサークル進捗更新リクエストのDTO。
type updateStudyCircleProgressRequest struct {
	IsCompleted bool `json:"is_completed"`
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
	input := bindJSON[updateStudyCircleProgressRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	if err := h.updateProgress.Execute(c.Request.Context(), id, userID, stepID, input.IsCompleted); err != nil {
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
	progress, err := h.listProgress.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, progress)
}

// createStudyCircleCheckinRequest はスタディサークルチェックイン作成リクエストのDTO。
type createStudyCircleCheckinRequest struct {
	Content string `json:"content" binding:"required,max=5000"`
}

// CreateCheckin は日次チェックインを作成する。
func (h *StudyCircleHandler) CreateCheckin(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	input := bindJSON[createStudyCircleCheckinRequest](c)
	if input == nil {
		return
	}
	userID := c.GetUint("userID")
	checkin, err := h.createCheckin.Execute(c.Request.Context(), id, userID, input.Content)
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
	checkins, err := h.listCheckins.Execute(c.Request.Context(), id, userID)
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
	ranking, err := h.streakRanking.Execute(c.Request.Context(), id, userID)
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

	circles, err := h.listByStatus.Execute(c.Request.Context(), userID, status)
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

	if err := h.updateMemberRole.Execute(c.Request.Context(), circleID, userID, targetUserID, req.Role); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("メンバー役割を更新しました"))
}

// Search はスタディサークルをキーワード検索する。
func (h *StudyCircleHandler) Search(c *gin.Context) {
	query, ok := parseSearchQuery(c)
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	circles, total, err := h.search.Execute(c.Request.Context(), query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, studyCircleListResponse{
		Circles: circles,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetMyCount は認証ユーザーが参加しているスタディサークル総数を返す。
func (h *StudyCircleHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
