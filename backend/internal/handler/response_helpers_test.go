package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestRespondSuccessGeneric(t *testing.T) {
	t.Run("文字列データで200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		testData := "success message"

		RespondSuccess(c, testData)

		assert.Equal(t, http.StatusOK, w.Code)
		var response string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testData, response)
	})

	t.Run("構造体データで200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		testData := TestStruct{ID: 1, Name: "test"}

		RespondSuccess(c, testData)

		assert.Equal(t, http.StatusOK, w.Code)
		var response TestStruct
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testData, response)
	})
}

func TestRespondCreatedGeneric(t *testing.T) {
	t.Run("構造体データで201を返す", func(t *testing.T) {
		c, w := setupTestContext()
		type CreatedResource struct {
			ID        int    `json:"id"`
			CreatedAt string `json:"created_at"`
		}
		testData := CreatedResource{ID: 42, CreatedAt: "2024-01-01T00:00:00Z"}

		RespondCreated(c, testData)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response CreatedResource
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testData, response)
	})
}

func TestRespondMessageHelper(t *testing.T) {
	t.Run("メッセージレスポンスで200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		message := "操作が成功しました"

		RespondMessage(c, message)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, message, response.Message)
	})
}

func TestRespondDataHelper(t *testing.T) {
	t.Run("DataResponse形式で200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		type User struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}
		testUser := User{ID: 1, Email: "test@example.com"}

		RespondData(c, testUser)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.DataResponse[User]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testUser, response.Data)
	})

	t.Run("プリミティブ型でDataResponse形式で200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		testCount := 42

		RespondData(c, testCount)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.DataResponse[int]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCount, response.Data)
	})
}

func TestRespondListHelper(t *testing.T) {
	t.Run("ListResponse形式でtotal付きで200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		type Item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		testItems := []Item{
			{ID: 1, Name: "item1"},
			{ID: 2, Name: "item2"},
		}
		testTotal := 100

		RespondList(c, testItems, testTotal)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListResponse[Item]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testItems, response.Data)
		assert.Equal(t, testTotal, response.Total)
		assert.Equal(t, 0, response.Page)  // ページネーション無しなのでデフォルト値
		assert.Equal(t, 0, response.Limit) // ページネーション無しなのでデフォルト値
	})

	t.Run("空のリストでも正常に動作する", func(t *testing.T) {
		c, w := setupTestContext()
		type Item struct {
			ID int `json:"id"`
		}
		var testItems []Item

		RespondList(c, testItems, 0)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListResponse[Item]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Data))
		assert.Equal(t, 0, response.Total)
	})
}

func TestRespondListWithPaginationHelper(t *testing.T) {
	t.Run("ページネーション付きListResponse形式で200を返す", func(t *testing.T) {
		c, w := setupTestContext()
		type Post struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		}
		testPosts := []Post{
			{ID: 11, Title: "post11"},
			{ID: 12, Title: "post12"},
		}
		testTotal := 50
		testPage := 2
		testLimit := 10

		RespondListWithPagination(c, testPosts, testTotal, testPage, testLimit)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListResponse[Post]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testPosts, response.Data)
		assert.Equal(t, testTotal, response.Total)
		assert.Equal(t, testPage, response.Page)
		assert.Equal(t, testLimit, response.Limit)
	})

	t.Run("ページ1の場合も正常に動作する", func(t *testing.T) {
		c, w := setupTestContext()
		type Item struct {
			Value string `json:"value"`
		}
		testItems := []Item{{Value: "first"}}

		RespondListWithPagination(c, testItems, 1, 1, 20)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListResponse[Item]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, 1, response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 20, response.Limit)
	})
}
