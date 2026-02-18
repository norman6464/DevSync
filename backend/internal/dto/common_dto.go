package dto

// VoteRequest は投票のリクエストボディ。
// 質問・回答の投票で共通使用する。
type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=1 -1" validate:"required,oneof=1 -1"`
}

// ConnectUsernameRequest は外部サービスのユーザー名接続リクエスト。
// AtCoder・Zenn・Qiitaなど共通で使用する。
type ConnectUsernameRequest struct {
	Username string `json:"username" binding:"required,max=100" validate:"required,max=100"`
}

// ConnectSyncResponse は外部サービス接続・同期のレスポンス。
// Qiita・Zennで共通使用する。
type ConnectSyncResponse struct {
	Message       string `json:"message"`
	ArticlesCount int    `json:"articles_count"`
}
