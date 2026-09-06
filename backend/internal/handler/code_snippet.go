package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// CodeSnippetHandler はコードスニペット関連のHTTPハンドラ。
type CodeSnippetHandler struct {
	create        *usecase.CreateCodeSnippetUseCase
	listByPost    *usecase.ListCodeSnippetsByPostUseCase
	byLanguage    *usecase.ListCodeSnippetsByLanguageUseCase
	update        *usecase.UpdateCodeSnippetUseCase
	remove        *usecase.DeleteCodeSnippetUseCase
	listComments  *usecase.ListSnippetCommentsUseCase
	createComment *usecase.CreateSnippetCommentUseCase
	deleteComment *usecase.DeleteSnippetCommentUseCase
	search        *usecase.SearchCodeSnippetsUseCase
	fork          *usecase.ForkCodeSnippetUseCase
	favorite      *usecase.FavoriteCodeSnippetUseCase
	unfavorite    *usecase.UnfavoriteCodeSnippetUseCase
	listFavorited *usecase.ListFavoritedCodeSnippetsUseCase
	count         *usecase.CountCodeSnippetsUseCase
}

// NewCodeSnippetHandler は新しいCodeSnippetHandlerインスタンスを生成する。
func NewCodeSnippetHandler(
	create *usecase.CreateCodeSnippetUseCase,
	listByPost *usecase.ListCodeSnippetsByPostUseCase,
	byLanguage *usecase.ListCodeSnippetsByLanguageUseCase,
	update *usecase.UpdateCodeSnippetUseCase,
	remove *usecase.DeleteCodeSnippetUseCase,
	listComments *usecase.ListSnippetCommentsUseCase,
	createComment *usecase.CreateSnippetCommentUseCase,
	deleteComment *usecase.DeleteSnippetCommentUseCase,
	search *usecase.SearchCodeSnippetsUseCase,
	fork *usecase.ForkCodeSnippetUseCase,
	favorite *usecase.FavoriteCodeSnippetUseCase,
	unfavorite *usecase.UnfavoriteCodeSnippetUseCase,
	listFavorited *usecase.ListFavoritedCodeSnippetsUseCase,
	count *usecase.CountCodeSnippetsUseCase,
) *CodeSnippetHandler {
	return &CodeSnippetHandler{
		create: create, listByPost: listByPost, byLanguage: byLanguage,
		update: update, remove: remove,
		listComments: listComments, createComment: createComment, deleteComment: deleteComment,
		search: search, fork: fork,
		favorite: favorite, unfavorite: unfavorite, listFavorited: listFavorited,
		count: count,
	}
}

// createCodeSnippetRequest はコードスニペット作成リクエスト。
type createCodeSnippetRequest struct {
	Language string `json:"language" binding:"required,max=100" validate:"required,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"required,max=50000" validate:"required,max=50000"`
}

// Create は投稿にコードスニペットを追加する。
func (h *CodeSnippetHandler) Create(c *gin.Context) {
	postID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[createCodeSnippetRequest](c)
	if input == nil {
		return
	}

	snippet := &model.CodeSnippet{
		PostID:   postID,
		UserID:   userID,
		Language: input.Language,
		FileName: input.FileName,
		Code:     input.Code,
	}
	created, err := h.create.Execute(c.Request.Context(), snippet)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, created)
}

// GetByPostID は投稿に紐づくスニペット一覧を返す。
func (h *CodeSnippetHandler) GetByPostID(c *gin.Context) {
	postID, ok := parseID(c, "id")
	if !ok {
		return
	}

	snippets, err := h.listByPost.Execute(c.Request.Context(), postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(snippets))
}

// GetByID は指定IDのスニペットを返す。
func (h *CodeSnippetHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	snippets, err := h.listByPost.Execute(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(snippets))
}

// updateCodeSnippetRequest はコードスニペット更新リクエスト。
type updateCodeSnippetRequest struct {
	Language string `json:"language" binding:"omitempty,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"omitempty,max=50000"`
}

// Update はスニペットを更新する。所有者のみ更新可能。
func (h *CodeSnippetHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[updateCodeSnippetRequest](c)
	if input == nil {
		return
	}

	snippet, err := h.update.Execute(c.Request.Context(), usecase.UpdateCodeSnippetInput{
		ID: id, UserID: userID,
		Language: input.Language, FileName: input.FileName, Code: input.Code,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, snippet)
}

// Delete はスニペットを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// GetComments はスニペットのインラインコメント一覧を返す。
func (h *CodeSnippetHandler) GetComments(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}

	comments, err := h.listComments.Execute(c.Request.Context(), snippetID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(comments))
}

// createSnippetCommentRequest はスニペットコメント作成リクエスト。
type createSnippetCommentRequest struct {
	LineNumber int    `json:"line_number" binding:"required" validate:"required"`
	Content    string `json:"content" binding:"required,max=5000" validate:"required,max=5000"`
}

// CreateComment はスニペットにインラインコメントを作成する。
func (h *CodeSnippetHandler) CreateComment(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[createSnippetCommentRequest](c)
	if input == nil {
		return
	}

	comment := &model.SnippetComment{
		SnippetID:  snippetID,
		UserID:     userID,
		LineNumber: input.LineNumber,
		Content:    input.Content,
	}
	if err := h.createComment.Execute(c.Request.Context(), comment); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, comment)
}

// DeleteComment はインラインコメントを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.deleteComment.Execute(c.Request.Context(), commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// forkSnippetRequest はスニペットフォークリクエスト。
type forkSnippetRequest struct {
	TargetPostID uint `json:"target_post_id" binding:"required"`
}

// Fork はスニペットをフォークして指定投稿にコピーする。
func (h *CodeSnippetHandler) Fork(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	req := bindJSON[forkSnippetRequest](c)
	if req == nil {
		return
	}

	forked, err := h.fork.Execute(c.Request.Context(), userID, snippetID, req.TargetPostID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, forked)
}

// GetByUserLanguage は認証ユーザーのスニペットを言語でフィルタリングして取得する。
func (h *CodeSnippetHandler) GetByUserLanguage(c *gin.Context) {
	userID := c.GetUint("userID")
	language := c.Param("language")

	snippets, err := h.byLanguage.Execute(c.Request.Context(), userID, language)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(snippets))
}

// codeSnippetListResponse はコードスニペット検索結果レスポンス。
type codeSnippetListResponse struct {
	Snippets []model.CodeSnippet `json:"snippets"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

// Search はコードスニペットをキーワード検索する。
func (h *CodeSnippetHandler) Search(c *gin.Context) {
	q, ok := parseSearchQuery(c)
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	snippets, total, err := h.search.Execute(c.Request.Context(), q, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, codeSnippetListResponse{
		Snippets: ensureSlice(snippets),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// Favorite はスニペットをお気に入りに追加する。
func (h *CodeSnippetHandler) Favorite(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.favorite.Execute(c.Request.Context(), userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("お気に入りに追加しました"))
}

// Unfavorite はスニペットのお気に入りを解除する。
func (h *CodeSnippetHandler) Unfavorite(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.unfavorite.Execute(c.Request.Context(), userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("お気に入りを解除しました"))
}

// codeSnippetFavoritesResponse はお気に入りスニペット一覧レスポンス。
type codeSnippetFavoritesResponse struct {
	Snippets []model.CodeSnippet `json:"snippets"`
	Total    int64               `json:"total"`
}

// GetFavorites はお気に入りスニペット一覧を取得する。
func (h *CodeSnippetHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	snippets, total, err := h.listFavorited.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, codeSnippetFavoritesResponse{
		Snippets: ensureSlice(snippets),
		Total:    total,
	})
}

// GetMyCount は認証ユーザーのコードスニペット総数を返す。
func (h *CodeSnippetHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
