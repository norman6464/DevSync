package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostSearchRepository は投稿検索のリポジトリインターフェース
type PostSearchRepository interface {
	Search(query string, limit, offset int) (interface{}, int64, error)
}

// StudyCircleSearchRepository はスタディサークル検索のリポジトリインターフェース
type StudyCircleSearchRepository interface {
	Search(query string, limit, offset int) (interface{}, int64, error)
}

// SearchHandler is the handler for search operations
type SearchHandler struct {
	postRepo   PostSearchRepository
	circleRepo StudyCircleSearchRepository
}

// NewSearchHandler creates a new search handler instance
func NewSearchHandler(postRepo PostSearchRepository, circleRepo StudyCircleSearchRepository) *SearchHandler {
	return &SearchHandler{
		postRepo:   postRepo,
		circleRepo: circleRepo,
	}
}

// SearchPosts handles post search requests
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	results, _, err := h.postRepo.Search(query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, results)
}

// SearchCircles handles study circle search requests
func (h *SearchHandler) SearchCircles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "query parameter 'q' is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	results, _, err := h.circleRepo.Search(query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, results)
}
