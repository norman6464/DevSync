package dto

// UpdateEmailPreferencesRequest はメール配信設定更新リクエスト。
type UpdateEmailPreferencesRequest struct {
	EmailWeeklyReport *bool   `json:"email_weekly_report"`
	EmailLanguage     *string `json:"email_language"`
}

// EmailPreferencesResponse はメール配信設定レスポンス。
type EmailPreferencesResponse struct {
	EmailWeeklyReport bool   `json:"email_weekly_report"`
	EmailLanguage     string `json:"email_language"`
}
