package dto

// UploadResponse は単一ファイルアップロードレスポンス。
type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// UploadMultipleResponse は複数ファイルアップロードレスポンス。
type UploadMultipleResponse struct {
	URLs []string `json:"urls"`
}
