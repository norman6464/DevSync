package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler はファイルアップロード関連のHTTPハンドラ。
// 画像の単体・複数アップロードを処理する。
type UploadHandler struct {
	uploadDir string
}

// NewUploadHandler は新しいUploadHandlerインスタンスを生成する。
// アップロードディレクトリが存在しない場合は自動で作成する。
func NewUploadHandler() *UploadHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// アップロードディレクトリが存在しない場合は作成する
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &UploadHandler{uploadDir: uploadDir}
}

// UploadImage は単一の画像ファイルをアップロードする。
// 最大5MBまでのjpg/jpeg/png/gif/webp形式に対応する。
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// ファイルサイズのバリデーション（最大5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// ファイル形式のバリデーション
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// ユニークなファイル名を生成する
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// 日付ベースのサブディレクトリを作成する
	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// ファイルを保存する
	filePath := filepath.Join(fullDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// URLパスを返す
	urlPath := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)
	c.JSON(http.StatusOK, gin.H{
		"url":      urlPath,
		"filename": filename,
	})
}

// UploadMultipleImages は複数の画像ファイルを一括アップロードする。
// 最大10ファイルまで、各ファイル最大5MBまで対応する。
func (h *UploadHandler) UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image files provided"})
		return
	}

	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
		return
	}

	var urls []string
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	for _, file := range files {
		// ファイルサイズのバリデーション
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s exceeds 5MB limit", file.Filename)})
			return
		}

		// ファイル形式のバリデーション
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file type for %s", file.Filename)})
			return
		}

		// ユニークなファイル名を生成する
		filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
		filePath := filepath.Join(fullDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		urls = append(urls, fmt.Sprintf("/uploads/%s/%s", dateDir, filename))
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
