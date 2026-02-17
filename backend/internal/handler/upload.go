package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/norman6464/devsync/backend/internal/dto"
)

// allowedMIMETypes はアップロードを許可するMIMEタイプの一覧。
// ファイルのマジックバイト（先頭512バイト）から検出したMIMEタイプで検証する。
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// detectMIMEType はファイルの先頭512バイトからMIMEタイプを検出する。
func detectMIMEType(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	// 読み取り位置をリセット
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

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
		respondBadRequest(c, "No image file provided")
		return
	}

	// ファイルサイズのバリデーション（最大5MB）
	if file.Size > 5*1024*1024 {
		respondBadRequest(c, "File size exceeds 5MB limit")
		return
	}

	// ファイル拡張子のバリデーション
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		respondBadRequest(c, "Invalid file type. Allowed: jpg, jpeg, png, gif, webp")
		return
	}

	// マジックバイトによるMIMEタイプ検証（拡張子偽装防止）
	src, err := file.Open()
	if err != nil {
		respondInternalError(c, "Failed to read file")
		return
	}
	defer src.Close()

	mimeType, err := detectMIMEType(src)
	if err != nil {
		respondInternalError(c, "Failed to detect file type")
		return
	}
	if !allowedMIMETypes[mimeType] {
		respondBadRequest(c, fmt.Sprintf("Invalid file content type: %s. Only image files are allowed", mimeType))
		return
	}

	// ユニークなファイル名を生成する
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// 日付ベースのサブディレクトリを作成する
	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		respondInternalError(c, "Failed to create directory")
		return
	}

	// ファイルを保存する
	filePath := filepath.Join(fullDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		respondInternalError(c, "Failed to save file")
		return
	}

	// URLパスを返す
	urlPath := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)
	respondOK(c, dto.UploadResponse{
		URL:      urlPath,
		Filename: filename,
	})
}

// UploadMultipleImages は複数の画像ファイルを一括アップロードする。
// 最大10ファイルまで、各ファイル最大5MBまで対応する。
func (h *UploadHandler) UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		respondBadRequest(c, "Failed to parse form")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		respondBadRequest(c, "No image files provided")
		return
	}

	if len(files) > 10 {
		respondBadRequest(c, "Maximum 10 images allowed")
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
		respondInternalError(c, "Failed to create directory")
		return
	}

	for _, file := range files {
		// ファイルサイズのバリデーション
		if file.Size > 5*1024*1024 {
			respondBadRequest(c, fmt.Sprintf("File %s exceeds 5MB limit", file.Filename))
			return
		}

		// ファイル拡張子のバリデーション
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			respondBadRequest(c, fmt.Sprintf("Invalid file type for %s", file.Filename))
			return
		}

		// マジックバイトによるMIMEタイプ検証（拡張子偽装防止）
		src, err := file.Open()
		if err != nil {
			respondInternalError(c, "Failed to read file")
			return
		}
		mimeType, err := detectMIMEType(src)
		src.Close()
		if err != nil {
			respondInternalError(c, "Failed to detect file type")
			return
		}
		if !allowedMIMETypes[mimeType] {
			respondBadRequest(c, fmt.Sprintf("Invalid content type for %s: %s", file.Filename, mimeType))
			return
		}

		// ユニークなファイル名を生成する
		filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
		filePath := filepath.Join(fullDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			respondInternalError(c, "Failed to save file")
			return
		}

		urls = append(urls, fmt.Sprintf("/uploads/%s/%s", dateDir, filename))
	}

	respondOK(c, dto.UploadMultipleResponse{URLs: urls})
}
