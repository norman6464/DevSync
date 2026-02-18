package dto

// UpdateNotificationSettingsRequest は通知設定更新のリクエストボディ。
type UpdateNotificationSettingsRequest struct {
	EnableLikes    bool `json:"enable_likes"`
	EnableComments bool `json:"enable_comments"`
	EnableFollows  bool `json:"enable_follows"`
	EnableMessages bool `json:"enable_messages"`
	EnableMentions bool `json:"enable_mentions"`
	EnableWebPush  bool `json:"enable_web_push"`
	EnableEmail    bool `json:"enable_email"`
	EnableSound    bool `json:"enable_sound"`
}
