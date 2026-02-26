package handler

import (

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// QuestionServiceInterface はQuestionServiceが実装すべきインターフェース。
type QuestionServiceInterface interface {
	Create(question *model.Question) error
	GetAll(limit, offset int, tag, sort string) ([]model.Question, int64, error)
	Search(q string, limit, offset int) ([]model.Question, int64, error)
	GetByID(id uint) (*model.Question, error)
	GetByUserID(userID uint, limit, offset int) ([]model.Question, int64, error)
	GetUserVote(userID, questionID uint) (int, error)
	Update(id, userID uint, title, body, tags string) (*model.Question, error)
	Delete(id, userID uint) error
	Vote(userID, questionID uint, value int) error
	RemoveVote(userID, questionID uint) error
	GetSolved(limit, offset int) ([]model.Question, int64, error)
	Bookmark(userID, questionID uint) error
	Unbookmark(userID, questionID uint) error
	GetBookmarkedByUserID(userID uint, limit, offset int) ([]model.Question, int64, error)
	GetUnanswered(limit, offset int) ([]model.Question, int64, error)
}

// QuestionHandler は質問関連のHTTPハンドラ。
// 質問のCRUD・検索・投票を処理する。
type QuestionHandler struct {
	service QuestionServiceInterface
}

// NewQuestionHandler は新しいQuestionHandlerインスタンスを生成する。
func NewQuestionHandler(s QuestionServiceInterface) *QuestionHandler {
	return &QuestionHandler{service: s}
}

// Create は新しい質問を作成する。
func (h *QuestionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreateQuestionRequest](c)
	if req == nil {
		return
	}

	question := &model.Question{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
		Tags:   req.Tags,
	}

	if err := h.service.Create(question); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, question)
}

// GetAll は質問一覧をページネーション・タグ・ソート付きで取得する。
func (h *QuestionHandler) GetAll(c *gin.Context) {
	limit, offset := parseLimitOffset(c)
	tag := c.Query("tag")
	sort := c.DefaultQuery("sort", "newest")

	questions, total, err := h.service.GetAll(limit, offset, tag, sort)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// Search はキーワードで質問を検索する。
func (h *QuestionHandler) Search(c *gin.Context) {
	q, ok := parseSearchQuery(c)
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.Search(q, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetByID は指定されたIDの質問を取得する。
func (h *QuestionHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	question, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	userVote, _ := h.service.GetUserVote(userID, id)

	respondOK(c, dto.QuestionDetailResponse{
		Question: *question,
		UserVote: userVote,
	})
}

// GetByUserID は指定されたユーザーの質問一覧をページネーション付きで取得する。
func (h *QuestionHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetMyQuestions は認証ユーザーの質問一覧を取得する。
func (h *QuestionHandler) GetMyQuestions(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// Update は指定された質問を更新する。
func (h *QuestionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateQuestionRequest](c)
	if req == nil {
		return
	}

	question, err := h.service.Update(id, userID, req.Title, req.Body, req.Tags)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, question)
}

// Delete は指定された質問を削除する。
func (h *QuestionHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// Vote は質問に投票する。
func (h *QuestionHandler) Vote(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.VoteRequest](c)
	if req == nil {
		return
	}

	if err := h.service.Vote(userID, id, req.Value); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("投票しました"))
}

// GetSolved は解決済みの質問一覧を取得する。
func (h *QuestionHandler) GetSolved(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.GetSolved(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetUnanswered は未回答の質問一覧を取得する。
func (h *QuestionHandler) GetUnanswered(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.GetUnanswered(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// RemoveVote は質問への投票を取り消す。
func (h *QuestionHandler) RemoveVote(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.RemoveVote(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("投票を取り消しました"))
}

// Bookmark は質問をブックマークする。
func (h *QuestionHandler) Bookmark(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Bookmark(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ブックマークしました"))
}

// Unbookmark は質問のブックマークを解除する。
func (h *QuestionHandler) Unbookmark(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Unbookmark(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ブックマークを解除しました"))
}

// GetBookmarks はブックマーク済み質問一覧を取得する。
func (h *QuestionHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.service.GetBookmarkedByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.QuestionListResponse{
		Questions: questions,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}
