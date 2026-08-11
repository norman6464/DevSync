package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// QuestionHandler は質問関連のHTTPハンドラ。
// 質問のCRUD・検索・投票を処理する。
type QuestionHandler struct {
	create         *usecase.CreateQuestionUseCase
	list           *usecase.ListQuestionsUseCase
	search         *usecase.SearchQuestionsUseCase
	get            *usecase.GetQuestionUseCase
	listByUser     *usecase.ListQuestionsByUserUseCase
	userVote       *usecase.GetQuestionUserVoteUseCase
	update         *usecase.UpdateQuestionUseCase
	remove         *usecase.DeleteQuestionUseCase
	vote           *usecase.VoteQuestionUseCase
	removeVote     *usecase.RemoveQuestionVoteUseCase
	listSolved     *usecase.ListSolvedQuestionsUseCase
	listUnanswered *usecase.ListUnansweredQuestionsUseCase
	bookmark       *usecase.BookmarkQuestionUseCase
	unbookmark     *usecase.UnbookmarkQuestionUseCase
	listBookmarked *usecase.ListBookmarkedQuestionsUseCase
	count          *usecase.CountQuestionsUseCase
}

// NewQuestionHandler は新しいQuestionHandlerインスタンスを生成する。
func NewQuestionHandler(
	create *usecase.CreateQuestionUseCase,
	list *usecase.ListQuestionsUseCase,
	search *usecase.SearchQuestionsUseCase,
	get *usecase.GetQuestionUseCase,
	listByUser *usecase.ListQuestionsByUserUseCase,
	userVote *usecase.GetQuestionUserVoteUseCase,
	update *usecase.UpdateQuestionUseCase,
	remove *usecase.DeleteQuestionUseCase,
	vote *usecase.VoteQuestionUseCase,
	removeVote *usecase.RemoveQuestionVoteUseCase,
	listSolved *usecase.ListSolvedQuestionsUseCase,
	listUnanswered *usecase.ListUnansweredQuestionsUseCase,
	bookmark *usecase.BookmarkQuestionUseCase,
	unbookmark *usecase.UnbookmarkQuestionUseCase,
	listBookmarked *usecase.ListBookmarkedQuestionsUseCase,
	count *usecase.CountQuestionsUseCase,
) *QuestionHandler {
	return &QuestionHandler{
		create: create, list: list, search: search, get: get,
		listByUser: listByUser, userVote: userVote, update: update, remove: remove,
		vote: vote, removeVote: removeVote, listSolved: listSolved,
		listUnanswered: listUnanswered, bookmark: bookmark, unbookmark: unbookmark,
		listBookmarked: listBookmarked, count: count,
	}
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

	if err := h.create.Execute(c.Request.Context(), question); err != nil {
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

	questions, total, err := h.list.Execute(c.Request.Context(), limit, offset, tag, sort)
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

	questions, total, err := h.search.Execute(c.Request.Context(), q, limit, offset)
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

	question, err := h.get.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	userVote, _ := h.userVote.Execute(c.Request.Context(), userID, id)

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

	questions, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
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

	questions, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
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

	question, err := h.update.Execute(c.Request.Context(), id, userID, req.Title, req.Body, req.Tags)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, question)
}

// Delete は指定された質問を削除する。
func (h *QuestionHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
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

	if err := h.vote.Execute(c.Request.Context(), userID, id, req.Value); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("投票しました"))
}

// GetSolved は解決済みの質問一覧を取得する。
func (h *QuestionHandler) GetSolved(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.listSolved.Execute(c.Request.Context(), limit, offset)
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

	questions, total, err := h.listUnanswered.Execute(c.Request.Context(), limit, offset)
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

	if err := h.removeVote.Execute(c.Request.Context(), userID, id); err != nil {
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

	if err := h.bookmark.Execute(c.Request.Context(), userID, id); err != nil {
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

	if err := h.unbookmark.Execute(c.Request.Context(), userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("ブックマークを解除しました"))
}

// GetBookmarks はブックマーク済み質問一覧を取得する。
func (h *QuestionHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	questions, total, err := h.listBookmarked.Execute(c.Request.Context(), userID, limit, offset)
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

// GetMyCount は認証ユーザーの質問総数を取得する。
func (h *QuestionHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}
