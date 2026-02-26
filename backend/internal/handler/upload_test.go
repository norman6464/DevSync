package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// minimalPNG は1x1の透過PNGバイト列（テスト用最小サイズ）
var minimalPNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, // IDAT chunk
	0x54, 0x78, 0x9C, 0x62, 0x00, 0x00, 0x00, 0x02,
	0x00, 0x01, 0xE5, 0x27, 0xDE, 0xFC, 0x00, 0x00,
	0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, // IEND chunk
	0x60, 0x82,
}

func setupUploadHandler(t *testing.T) *UploadHandler {
	t.Helper()
	dir := t.TempDir()
	return &UploadHandler{uploadDir: dir}
}

// ---------- UploadImage テスト ----------

func TestUploadImage_Success(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write(minimalPNG)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadImage(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "url")
	assert.Contains(t, w.Body.String(), "filename")
}

func TestUploadImage_NoFile(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "No image file provided")
}

func TestUploadImage_InvalidExtension(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.txt")
	part.Write([]byte("not an image"))
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid file type")
}

func TestUploadImage_InvalidMIMEType(t *testing.T) {
	h := setupUploadHandler(t)

	// .png拡張子だが中身はテキスト（拡張子偽装）
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "fake.png")
	part.Write([]byte("this is not a real png file content"))
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid file content type")
}

// ---------- UploadMultipleImages テスト ----------

func TestUploadMultipleImages_Success(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := 0; i < 2; i++ {
		part, _ := writer.CreateFormFile("images", fmt.Sprintf("img%d.png", i))
		part.Write(minimalPNG)
	}
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "urls")
}

func TestUploadMultipleImages_NoFiles(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadMultipleImages_TooManyFiles(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := 0; i < 11; i++ {
		part, _ := writer.CreateFormFile("images", fmt.Sprintf("img%d.png", i))
		part.Write(minimalPNG)
	}
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Maximum 10 images allowed")
}

func TestUploadMultipleImages_InvalidExtension(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part1, _ := writer.CreateFormFile("images", "valid.png")
	part1.Write(minimalPNG)
	part2, _ := writer.CreateFormFile("images", "invalid.txt")
	part2.Write([]byte("not an image"))
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid file type")
}

func TestUploadMultipleImages_InvalidMIMEType(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("images", "fake.png")
	part.Write([]byte("this is not a real png file content"))
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid content type")
}

// ---------- detectMIMEType テスト ----------

func TestDetectMIMEType_PNG(t *testing.T) {
	reader := bytes.NewReader(minimalPNG)
	mime, err := detectMIMEType(reader)
	assert.NoError(t, err)
	assert.Equal(t, "image/png", mime)
}

func TestDetectMIMEType_TextFile(t *testing.T) {
	reader := bytes.NewReader([]byte("hello world this is plain text"))
	mime, err := detectMIMEType(reader)
	assert.NoError(t, err)
	assert.NotEqual(t, "image/png", mime)
}

// ---------- NewUploadHandler テスト ----------

func TestNewUploadHandler_DefaultDir(t *testing.T) {
	orig := os.Getenv("UPLOAD_DIR")
	defer os.Setenv("UPLOAD_DIR", orig)

	tmpDir := t.TempDir()
	os.Setenv("UPLOAD_DIR", tmpDir+"/test_uploads")

	h, err := NewUploadHandler()
	assert.NoError(t, err)
	assert.NotNil(t, h)

	// ディレクトリが作成されたことを確認
	_, err = os.Stat(tmpDir + "/test_uploads")
	assert.NoError(t, err)
}

// ---------- ファイルパーミッションテスト ----------

func TestNewUploadHandler_DirPermissions(t *testing.T) {
	orig := os.Getenv("UPLOAD_DIR")
	defer os.Setenv("UPLOAD_DIR", orig)

	tmpDir := t.TempDir()
	os.Setenv("UPLOAD_DIR", tmpDir+"/secure_uploads")

	_, _ = NewUploadHandler()

	info, err := os.Stat(tmpDir + "/secure_uploads")
	assert.NoError(t, err)
	// ディレクトリは0750（rwxr-x---）であること
	assert.Equal(t, os.FileMode(0750), info.Mode().Perm())
}

func TestUploadImage_FilePermissions(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write(minimalPNG)
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadImage(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// アップロードされたファイルのパーミッションを検証
	files, _ := filepath.Glob(filepath.Join(h.uploadDir, "*", "*", "*.png"))
	assert.NotEmpty(t, files)
	info, err := os.Stat(files[0])
	assert.NoError(t, err)
	// ファイルは0640（rw-r-----）であること
	assert.Equal(t, os.FileMode(0640), info.Mode().Perm())
}

func TestUploadMultipleImages_FilePermissions(t *testing.T) {
	h := setupUploadHandler(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for i := 0; i < 2; i++ {
		part, _ := writer.CreateFormFile("images", fmt.Sprintf("img%d.png", i))
		part.Write(minimalPNG)
	}
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/upload/multiple", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	h.UploadMultipleImages(c)

	assert.Equal(t, http.StatusOK, w.Code)

	// アップロードされた全ファイルのパーミッションを検証
	files, _ := filepath.Glob(filepath.Join(h.uploadDir, "*", "*", "*.png"))
	assert.Len(t, files, 2)
	for _, f := range files {
		info, err := os.Stat(f)
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0640), info.Mode().Perm())
	}
}
