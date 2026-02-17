package dto

// ConnectUsernameRequest は外部サービスのユーザー名接続リクエスト。
// AtCoder・Zenn・Qiitaなど共通で使用する。
type ConnectUsernameRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
}

// ConnectSyncResponse は外部サービス接続・同期のレスポンス。
// Qiita・Zennで共通使用する。
type ConnectSyncResponse struct {
	Message       string `json:"message"`
	ArticlesCount int    `json:"articles_count"`
}
