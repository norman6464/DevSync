package dto

// ConnectUsernameRequest は外部サービスのユーザー名接続リクエスト。
// AtCoder・Zenn・Qiitaなど共通で使用する。
type ConnectUsernameRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
}
