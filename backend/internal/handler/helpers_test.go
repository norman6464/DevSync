package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- parseID ---

func TestParseID_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items/:id", func(c *gin.Context) {
		id, ok := parseID(c, "id")
		if ok {
			c.JSON(http.StatusOK, gin.H{"id": id})
		}
	})
	c.Request = httptest.NewRequest("GET", "/items/42", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(42), resp["id"])
}

func TestParseID_Invalid(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items/:id", func(c *gin.Context) {
		_, ok := parseID(c, "id")
		assert.False(t, ok)
	})
	c.Request = httptest.NewRequest("GET", "/items/abc", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "invalid id")
}

func TestParseID_Zero(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items/:id", func(c *gin.Context) {
		id, ok := parseID(c, "id")
		if ok {
			c.JSON(http.StatusOK, gin.H{"id": id})
		}
	})
	c.Request = httptest.NewRequest("GET", "/items/0", nil)
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
}

// --- parsePagination ---

func TestParsePagination_Defaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items", func(c *gin.Context) {
		page, limit := parsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit})
	})
	c.Request = httptest.NewRequest("GET", "/items", nil)
	r.ServeHTTP(w, c.Request)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, float64(20), resp["limit"])
}

func TestParsePagination_Custom(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items", func(c *gin.Context) {
		page, limit := parsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit})
	})
	c.Request = httptest.NewRequest("GET", "/items?page=3&limit=50", nil)
	r.ServeHTTP(w, c.Request)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(3), resp["page"])
	assert.Equal(t, float64(50), resp["limit"])
}

func TestParsePagination_ClampNegative(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items", func(c *gin.Context) {
		page, limit := parsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit})
	})
	c.Request = httptest.NewRequest("GET", "/items?page=-1&limit=-5", nil)
	r.ServeHTTP(w, c.Request)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, float64(20), resp["limit"])
}

func TestParsePagination_ClampMax(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/items", func(c *gin.Context) {
		page, limit := parsePagination(c)
		c.JSON(http.StatusOK, gin.H{"page": page, "limit": limit})
	})
	c.Request = httptest.NewRequest("GET", "/items?limit=999", nil)
	r.ServeHTTP(w, c.Request)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(100), resp["limit"])
}

// --- bindJSON ---

type testRequest struct {
	Name  string `json:"name" binding:"required"`
	Value int    `json:"value"`
}

func TestBindJSON_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/test", func(c *gin.Context) {
		req := bindJSON[testRequest](c)
		if req != nil {
			c.JSON(http.StatusOK, req)
		}
	})
	body, _ := json.Marshal(testRequest{Name: "hello", Value: 42})
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp testRequest
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "hello", resp.Name)
	assert.Equal(t, 42, resp.Value)
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/test", func(c *gin.Context) {
		req := bindJSON[testRequest](c)
		assert.Nil(t, req)
	})
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBindJSON_MissingRequired(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/test", func(c *gin.Context) {
		req := bindJSON[testRequest](c)
		assert.Nil(t, req)
	})
	body, _ := json.Marshal(map[string]int{"value": 10})
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- respondError ---

func TestRespondError_MapsServiceErrors(t *testing.T) {
	tests := []struct {
		err      error
		wantCode int
	}{
		{service.ErrNotFound, http.StatusNotFound},
		{service.ErrForbidden, http.StatusForbidden},
		{service.ErrBadRequest, http.StatusBadRequest},
		{service.ErrUnauthorized, http.StatusUnauthorized},
		{service.ErrConflict, http.StatusConflict},
		{service.ErrRateLimitExceeded, http.StatusTooManyRequests},
		{service.ErrLLMNotConfigured, http.StatusServiceUnavailable},
		{fmt.Errorf("unknown error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.GET("/test", func(c *gin.Context) {
				respondError(c, tt.err)
			})
			req := httptest.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Contains(t, resp, "error")
		})
	}
}

func TestRespondError_WrappedError(t *testing.T) {
	wrapped := fmt.Errorf("user 123: %w", service.ErrNotFound)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondError(c, wrapped)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRespondError_DomainError(t *testing.T) {
	tests := []struct {
		name     string
		err      *domain.DomainError
		wantCode int
		wantErrCode domain.ErrorCode
	}{
		{"NotFound", domain.ErrNotFound, http.StatusNotFound, domain.ErrCodeNotFound},
		{"Forbidden", domain.ErrForbidden, http.StatusForbidden, domain.ErrCodeForbidden},
		{"BadRequest", domain.ErrBadRequest, http.StatusBadRequest, domain.ErrCodeBadRequest},
		{"Unauthorized", domain.ErrUnauthorized, http.StatusUnauthorized, domain.ErrCodeUnauthorized},
		{"Conflict", domain.ErrConflict, http.StatusConflict, domain.ErrCodeConflict},
		{"RateLimitExceeded", domain.ErrRateLimitExceeded, http.StatusTooManyRequests, domain.ErrCodeRateLimitExceeded},
		{"ServiceUnavailable", domain.ErrServiceUnavailable, http.StatusServiceUnavailable, domain.ErrCodeServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			r.GET("/test", func(c *gin.Context) {
				respondError(c, tt.err)
			})
			req := httptest.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Contains(t, resp, "error")
			assert.Contains(t, resp, "code")
			assert.Equal(t, string(tt.wantErrCode), resp["code"])
		})
	}
}

func TestRespondError_CustomDomainError(t *testing.T) {
	customErr := domain.NewError(domain.ErrCodeValidation, "カスタムバリデーションエラー", errors.New("field is invalid"))
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondError(c, customErr)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "カスタムバリデーションエラー", resp["error"])
	assert.Equal(t, string(domain.ErrCodeValidation), resp["code"])
}

func TestRespondError_NonDomainError(t *testing.T) {
	stdErr := errors.New("標準エラー")
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondError(c, stdErr)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "error")
}

// --- respondOK / respondCreated / respondDeleted ---

func TestRespondOK(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondOK(c, gin.H{"key": "value"})
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "value", resp["key"])
}

func TestRespondCreated(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/test", func(c *gin.Context) {
		respondCreated(c, gin.H{"id": 1})
	})
	req := httptest.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRespondDeleted(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.DELETE("/test", func(c *gin.Context) {
		respondDeleted(c)
	})
	req := httptest.NewRequest("DELETE", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "deleted", resp["message"])
}

// --- respondBadRequest / respondForbidden / respondNotFound / respondUnauthorized / respondInternalError ---

func TestRespondBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondBadRequest(c, "invalid input")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid input", resp["error"])
	assert.Equal(t, string(domain.ErrCodeValidation), resp["code"])
}

func TestRespondForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondForbidden(c, "access denied")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "access denied", resp["error"])
	assert.Equal(t, string(domain.ErrCodeForbidden), resp["code"])
}

func TestRespondNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondNotFound(c, "resource not found")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "resource not found", resp["error"])
	assert.Equal(t, string(domain.ErrCodeNotFound), resp["code"])
}

func TestRespondUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondUnauthorized(c, "authentication required")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "authentication required", resp["error"])
	assert.Equal(t, string(domain.ErrCodeUnauthorized), resp["code"])
}

func TestRespondInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondInternalError(c, "internal server error")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "internal server error", resp["error"])
	assert.Equal(t, string(domain.ErrCodeInternal), resp["code"])
}

// --- respondPaginated ---

func TestRespondPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		items := []string{"a", "b", "c"}
		respondPaginated(c, items, 25, 2, 10)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(25), resp["total"])
	assert.Equal(t, float64(2), resp["page"])
	assert.Equal(t, float64(10), resp["limit"])
	assert.Equal(t, float64(3), resp["total_pages"])
	data := resp["data"].([]interface{})
	assert.Len(t, data, 3)
}

func TestRespondPaginated_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.GET("/test", func(c *gin.Context) {
		respondPaginated(c, []string{}, 0, 1, 20)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["total"])
	// total_pagesはomitemptyのため0の場合はJSONに含まれない
	assert.Nil(t, resp["total_pages"])
	data := resp["data"].([]interface{})
	assert.Len(t, data, 0)
}
