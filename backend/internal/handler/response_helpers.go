package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
)

// RespondSuccess は成功レスポンスを返す（200 OK）
func RespondSuccess[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, data)
}

// RespondCreated は作成成功レスポンスを返す（201 Created）
func RespondCreated[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, data)
}

// RespondMessage はメッセージレスポンスを返す（200 OK）
func RespondMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, dto.MessageResponse{Message: message})
}

// RespondData はデータレスポンスを返す（200 OK）
func RespondData[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, dto.DataResponse[T]{Data: data})
}

// RespondList はリストレスポンスを返す（200 OK）
func RespondList[T any](c *gin.Context, data []T, total int) {
	c.JSON(http.StatusOK, dto.ListResponse[T]{
		Data:  data,
		Total: total,
	})
}

// RespondListWithPagination はページネーション付きリストレスポンスを返す（200 OK）
func RespondListWithPagination[T any](c *gin.Context, data []T, total, page, limit int) {
	c.JSON(http.StatusOK, dto.ListResponse[T]{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
