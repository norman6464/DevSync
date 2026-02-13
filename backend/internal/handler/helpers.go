package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
)

// parseID はURLパラメータからuint型のIDを取得する。
// パース失敗時は400レスポンスを返しfalseを返す。
func parseID(c *gin.Context, param string) (uint, bool) {
	raw := c.Param(param)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + param})
		return 0, false
	}
	return uint(id), true
}

// parsePagination はクエリパラメータからpage/limitを取得する。
// デフォルト: page=1, limit=20。limitの上限は100。
func parsePagination(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// bindJSON はリクエストボディをJSON構造体にバインドする。
// バインド失敗時は400レスポンスを返しnilを返す。
func bindJSON[T any](c *gin.Context) *T {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return &req
}

// respondError はサービス層のエラーを適切なHTTPステータスコードに変換してレスポンスを返す。
// DomainError を使用して統一的なエラーハンドリングを実現する。
func respondError(c *gin.Context, err error) {
	// DomainError の場合
	if domainErr := domain.GetDomainError(err); domainErr != nil {
		c.JSON(domainErr.HTTPStatus(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
		return
	}

	// DomainError でない場合は内部エラーとして扱う
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// respondOK は200 OKでデータを返す。
func respondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// respondCreated は201 Createdでデータを返す。
func respondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// respondDeleted は200 OKで削除メッセージを返す。
func respondDeleted(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
