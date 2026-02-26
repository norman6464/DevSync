package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// AnswerServiceInterface はAnswerServiceが実装すべきインターフェース。
type AnswerServiceInterface interface {
	GetByQuestionID(questionID uint) ([]model.Answer, error)
	Create(answer *model.Answer) error
	Update(answerID, userID uint, body string) (*model.Answer, error)
	Delete(answerID, userID uint) error
	SetBestAnswer(questionID, answerID, userID uint) error
	Vote(userID, answerID uint, value int) error
	RemoveVote(userID, answerID uint) error
	GetByVoteRange(questionID uint, minVote, maxVote int) ([]model.Answer, error)
}

// AnswerHandler は回答関連のHTTPハンドラ。
// 回答のCRUD・ベストアンサー選定・投票を処理する。
type AnswerHandler struct {
	service AnswerServiceInterface
}

// NewAnswerHandler は新しいAnswerHandlerインスタンスを生成する。
func NewAnswerHandler(s AnswerServiceInterface) *AnswerHandler {
	return &AnswerHandler{service: s}
}

// GetByQuestionID は指定された質問の回答一覧を取得する。
func (h *AnswerHandler) GetByQuestionID(c *gin.Context) {
	questionID, ok := parseID(c, "id")
	if !ok {
		return
	}

	answers, err := h.service.GetByQuestionID(questionID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(answers))
}

// Create は質問に対して新しい回答を作成する。
func (h *AnswerHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	questionID, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.CreateAnswerRequest](c)
	if req == nil {
		return
	}

	answer := &model.Answer{
		UserID:     userID,
		QuestionID: questionID,
		Body:       req.Body,
	}

	if err := h.service.Create(answer); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, answer)
}

// Update は指定された回答を更新する。
func (h *AnswerHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, ok := parseID(c, "answerId")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateAnswerRequest](c)
	if req == nil {
		return
	}

	answer, err := h.service.Update(answerID, userID, req.Body)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, answer)
}

// Delete は指定された回答を削除する。
func (h *AnswerHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, ok := parseID(c, "answerId")
	if !ok {
		return
	}

	if err := h.service.Delete(answerID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// SetBestAnswer は質問のベストアンサーを選定する。
// 質問の投稿者のみがベストアンサーを選定できる。
func (h *AnswerHandler) SetBestAnswer(c *gin.Context) {
	userID := c.GetUint("userID")
	questionID, ok := parseID(c, "id")
	if !ok {
		return
	}

	answerID, ok := parseID(c, "answerId")
	if !ok {
		return
	}

	if err := h.service.SetBestAnswer(questionID, answerID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ベストアンサーを設定しました"))
}

// Vote は回答に投票する。
func (h *AnswerHandler) Vote(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, ok := parseID(c, "answerId")
	if !ok {
		return
	}

	req := bindJSON[dto.VoteRequest](c)
	if req == nil {
		return
	}

	if err := h.service.Vote(userID, answerID, req.Value); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("投票しました"))
}

// GetByVoteRange は指定質問の回答を投票スコア範囲でフィルタリングして取得する。
func (h *AnswerHandler) GetByVoteRange(c *gin.Context) {
	questionID, ok := parseID(c, "id")
	if !ok {
		return
	}

	minVote, ok := parseQueryInt(c, "min_vote", "0")
	if !ok {
		return
	}
	maxVote, ok := parseQueryInt(c, "max_vote", "100")
	if !ok {
		return
	}

	answers, err := h.service.GetByVoteRange(questionID, minVote, maxVote)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(answers))
}

// RemoveVote は回答への投票を取り消す。
func (h *AnswerHandler) RemoveVote(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, ok := parseID(c, "answerId")
	if !ok {
		return
	}

	if err := h.service.RemoveVote(userID, answerID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("投票を取り消しました"))
}
