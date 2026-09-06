package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uploadDirPerm  os.FileMode = 0o750 // rwxr-x---
	uploadFilePerm os.FileMode = 0o640 // rw-r-----
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
	if err != nil && !errors.Is(err, io.EOF) {
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
func NewUploadHandler() (*UploadHandler, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// アップロードディレクトリが存在しない場合は作成する
	if err := os.MkdirAll(uploadDir, uploadDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &UploadHandler{uploadDir: uploadDir}, nil
}

// uploadResponse は単一ファイルアップロードレスポンス。
type uploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// UploadImage は単一の画像ファイルをアップロードする。
// 最大5MBまでのjpg/jpeg/png/gif/webp形式に対応する。
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		respondBadRequest(c, "画像ファイルが必要です")
		return
	}

	// ファイルサイズのバリデーション（最大5MB）
	if file.Size > 5*1024*1024 {
		respondBadRequest(c, "ファイルサイズは5MB以下にしてください")
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
		respondBadRequest(c, "無効なファイル形式です。許可: jpg, jpeg, png, gif, webp")
		return
	}

	// マジックバイトによるMIMEタイプ検証（拡張子偽装防止）
	src, err := file.Open()
	if err != nil {
		respondInternalError(c, "ファイルの読み込みに失敗しました")
		return
	}
	defer src.Close()

	mimeType, err := detectMIMEType(src)
	if err != nil {
		respondInternalError(c, "ファイルタイプの検出に失敗しました")
		return
	}
	if !allowedMIMETypes[mimeType] {
		respondBadRequest(c, fmt.Sprintf("無効なファイルコンテンツタイプ: %s。画像ファイルのみ許可されています", mimeType))
		return
	}

	// ユニークなファイル名を生成する
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// 日付ベースのサブディレクトリを作成する
	dateDir := time.Now().Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(fullDir, uploadDirPerm); err != nil {
		respondInternalError(c, "ディレクトリの作成に失敗しました")
		return
	}

	// ファイルを保存する
	filePath := filepath.Join(fullDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		respondInternalError(c, "ファイルの保存に失敗しました")
		return
	}

	// ファイルパーミッションを制限する
	if err := os.Chmod(filePath, uploadFilePerm); err != nil {
		respondInternalError(c, "ファイル権限の設定に失敗しました")
		return
	}

	// URLパスを返す
	urlPath := fmt.Sprintf("/uploads/%s/%s", dateDir, filename)
	respondOK(c, uploadResponse{
		URL:      urlPath,
		Filename: filename,
	})
}

// uploadMultipleResponse は複数ファイルアップロードレスポンス。
type uploadMultipleResponse struct {
	URLs []string `json:"urls"`
}

// UploadMultipleImages は複数の画像ファイルを一括アップロードする。
// 最大10ファイルまで、各ファイル最大5MBまで対応する。
func (h *UploadHandler) UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		respondBadRequest(c, "フォームの解析に失敗しました")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		respondBadRequest(c, "画像ファイルが必要です")
		return
	}

	if len(files) > 10 {
		respondBadRequest(c, "画像は最大10枚までです")
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
	if err := os.MkdirAll(fullDir, uploadDirPerm); err != nil {
		respondInternalError(c, "ディレクトリの作成に失敗しました")
		return
	}

	for _, file := range files {
		// ファイルサイズのバリデーション
		if file.Size > 5*1024*1024 {
			respondBadRequest(c, fmt.Sprintf("ファイル %s は5MB以下にしてください", file.Filename))
			return
		}

		// ファイル拡張子のバリデーション
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			respondBadRequest(c, fmt.Sprintf("ファイル %s の形式が無効です", file.Filename))
			return
		}

		// マジックバイトによるMIMEタイプ検証（拡張子偽装防止）
		src, err := file.Open()
		if err != nil {
			respondInternalError(c, "ファイルの読み込みに失敗しました")
			return
		}
		mimeType, err := detectMIMEType(src)
		_ = src.Close()
		if err != nil {
			respondInternalError(c, "ファイルタイプの検出に失敗しました")
			return
		}
		if !allowedMIMETypes[mimeType] {
			respondBadRequest(c, fmt.Sprintf("ファイル %s のコンテンツタイプが無効です: %s", file.Filename, mimeType))
			return
		}

		// ユニークなファイル名を生成する
		filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
		filePath := filepath.Join(fullDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			respondInternalError(c, "ファイルの保存に失敗しました")
			return
		}

		// ファイルパーミッションを制限する
		if err := os.Chmod(filePath, uploadFilePerm); err != nil {
			respondInternalError(c, "ファイル権限の設定に失敗しました")
			return
		}

		urls = append(urls, fmt.Sprintf("/uploads/%s/%s", dateDir, filename))
	}

	respondOK(c, uploadMultipleResponse{URLs: urls})
}
