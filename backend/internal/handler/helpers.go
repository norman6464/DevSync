package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
)

// parseID はURLパラメータからuint型のIDを取得する。
// パース失敗時は400レスポンスを返しfalseを返す。
func parseID(c *gin.Context, param string) (uint, bool) {
	raw := c.Param(param)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		respondBadRequest(c, "invalid "+param)
		return 0, false
	}
	return uint(id), true
}

// parsePagination はクエリパラメータからpage/limitを取得する。
// デフォルト: page=1, limit=20。limitの上限は100。
// domain.ValidatePaginationを使用してlimitの上限を正規化する。
func parsePagination(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// domain.ValidatePaginationでlimitの上限（100）を正規化（エラーは無視）
	offset := (page - 1) * limit
	limit, _, _ = domain.ValidatePagination(limit, offset)

	return page, limit
}

// parseLimitOffset はクエリパラメータからlimit/offsetを取得し正規化する。
// デフォルト: limit=20, offset=0。limitの上限は100（domain.ValidatePagination準拠）。
func parseLimitOffset(c *gin.Context) (limit, offset int) {
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, offset, _ = domain.ValidatePagination(limit, offset)
	return limit, offset
}

// maxSearchQueryLen は検索クエリの最大文字数（ルーン数）。
const maxSearchQueryLen = 500

// maxUserSearchQueryLen はユーザー検索クエリの最大文字数（ルーン数）。
const maxUserSearchQueryLen = 100

// parseSearchQuery はクエリパラメータ "q" を取得・検証する共通ヘルパー。
// 未指定・空文字列・500文字超の場合は400を返しfalseを返す。
// 前後の空白はTrimSpaceで除去して返す。
func parseSearchQuery(c *gin.Context) (string, bool) {
	query := c.Query("q")
	if query == "" {
		respondBadRequest(c, "query parameter 'q' is required")
		return "", false
	}
	if len([]rune(query)) > maxSearchQueryLen {
		respondBadRequest(c, "検索クエリは500文字以下である必要があります")
		return "", false
	}
	return strings.TrimSpace(query), true
}

// parseOptionalSearchQuery はクエリパラメータ "q" をオプションとして取得・検証する。
// 未指定・空文字列の場合は空文字列を返す（エラーなし）。
// maxLen文字超の場合は400を返しfalseを返す。
// 前後の空白はTrimSpaceで除去して返す。
func parseOptionalSearchQuery(c *gin.Context, maxLen int) (string, bool) {
	query := c.Query("q")
	if query == "" {
		return "", true
	}
	if len([]rune(query)) > maxLen {
		respondBadRequest(c, "検索クエリは"+strconv.Itoa(maxLen)+"文字以下である必要があります")
		return "", false
	}
	return strings.TrimSpace(query), true
}

// parseQueryInt はクエリパラメータを整数としてパースする。
// 未指定時はdefaultValueを使用する。パース失敗時は400レスポンスを返しfalseを返す。
func parseQueryInt(c *gin.Context, param, defaultValue string) (int, bool) {
	raw := c.DefaultQuery(param, defaultValue)
	val, err := strconv.Atoi(raw)
	if err != nil {
		respondBadRequest(c, param+"は数値で指定してください")
		return 0, false
	}
	return val, true
}

// parseQueryIntSilent はクエリパラメータを整数としてパースする。
// 未指定時やパース失敗時はデフォルト値を返す。
func parseQueryIntSilent(c *gin.Context, param string, defaultValue int) int {
	raw := c.Query(param)
	if raw == "" {
		return defaultValue
	}
	if val, err := strconv.Atoi(raw); err == nil && val > 0 {
		return val
	}
	return defaultValue
}

// parseExportPeriod はエクスポート期間パラメータをパースする。
// "all"の場合は0を返す。未指定時はデフォルト30日。不正な値は400を返しfalseを返す。
func parseExportPeriod(c *gin.Context) (int, bool) {
	p := c.Query("period")
	if p == "" {
		return 30, true
	}
	if p == "all" {
		return 0, true
	}
	n, err := strconv.Atoi(p)
	if err != nil || n < 0 {
		respondBadRequest(c, "periodは7/30/90/allのいずれかを指定してください")
		return 0, false
	}
	return n, true
}

// ensureSlice はnilスライスを空スライスに変換する。
// JSONレスポンスでnullの代わりに[]を返すために使用する。
func ensureSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

// bindJSON はリクエストボディをJSON構造体にバインドする。
// バインド失敗時は400レスポンスを返しnilを返す。
func bindJSON[T any](c *gin.Context) *T {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストの形式が不正です")
		return nil
	}
	return &req
}

// respondError はサービス層のエラーを適切なHTTPステータスコードに変換してレスポンスを返す。
// DomainError を使用して統一的なエラーハンドリングを実現する。
func respondError(c *gin.Context, err error) {
	// DomainError の場合
	if domainErr := domain.GetDomainError(err); domainErr != nil {
		response := domain.NewErrorResponse(domainErr.Message, string(domainErr.Code), nil)
		c.JSON(domainErr.HTTPStatus(), response)
		return
	}

	// DomainError でない場合は内部エラーとして扱う（内部エラーの詳細は隠蔽）
	response := domain.NewErrorResponse("内部エラーが発生しました", string(domain.ErrCodeInternal), nil)
	c.JSON(http.StatusInternalServerError, response)
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
	response := domain.NewMessageResponse("deleted")
	c.JSON(http.StatusOK, response)
}

// respondBadRequest は400 Bad Requestでエラーメッセージを返す。
func respondBadRequest(c *gin.Context, message string) {
	response := domain.NewErrorResponse(message, string(domain.ErrCodeValidation), nil)
	c.JSON(http.StatusBadRequest, response)
}

// respondForbidden は403 Forbiddenでエラーメッセージを返す。
func respondForbidden(c *gin.Context, message string) {
	response := domain.NewErrorResponse(message, string(domain.ErrCodeForbidden), nil)
	c.JSON(http.StatusForbidden, response)
}

// respondNotFound は404 Not Foundでエラーメッセージを返す。
func respondNotFound(c *gin.Context, message string) {
	response := domain.NewErrorResponse(message, string(domain.ErrCodeNotFound), nil)
	c.JSON(http.StatusNotFound, response)
}

// respondUnauthorized は401 Unauthorizedでエラーメッセージを返す。
func respondUnauthorized(c *gin.Context, message string) {
	response := domain.NewErrorResponse(message, string(domain.ErrCodeUnauthorized), nil)
	c.JSON(http.StatusUnauthorized, response)
}

// respondInternalError は500 Internal Server Errorでエラーメッセージを返す。
func respondInternalError(c *gin.Context, message string) {
	response := domain.NewErrorResponse(message, string(domain.ErrCodeInternal), nil)
	c.JSON(http.StatusInternalServerError, response)
}

// respondPaginated はページネーション付きレスポンスを返す。
func respondPaginated(c *gin.Context, data interface{}, total int64, page, limit int) {
	response := domain.NewPaginatedResponse(data, total, page, limit)
	c.JSON(http.StatusOK, response)
}

// handleDelete はリソース削除の共通パターンを実装する。
// parseID → GetUint("userID") → deleteFn → respondDeleted
func handleDelete(c *gin.Context, deleteFn func(id, userID uint) error) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := deleteFn(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// handleGetByID はID+userIDによるリソース取得の共通パターンを実装する。
// parseID → GetUint("userID") → getter → respondOK
func handleGetByID[T any](c *gin.Context, getter func(id, userID uint) (*T, error)) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	result, err := getter(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, result)
}

// handleGetByIDPublic は公開リソース取得の共通パターンを実装する。
// parseID → getter → respondOK（userIDチェックなし）
func handleGetByIDPublic[T any](c *gin.Context, getter func(id uint) (*T, error)) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := getter(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, result)
}

// handleToggleAction はLike/Unlikeなどのトグル操作を共通化するヘルパー。
// parseID → サービス呼び出し → レスポンス返却のパターンを統一する。
func handleToggleAction(c *gin.Context, action func(userID, id uint) error, message string) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := action(userID, id); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse(message))
}
