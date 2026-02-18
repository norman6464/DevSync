package dto

// UpdateUserRequest はユーザープロフィール更新リクエスト。
type UpdateUserRequest struct {
	Name                string  `json:"name" binding:"omitempty,max=200"`
	Bio                 string  `json:"bio" binding:"omitempty,max=500"`
	AvatarURL           string  `json:"avatar_url" binding:"omitempty,max=2000"`
	SkillsLanguages     *string `json:"skills_languages" binding:"omitempty,max=1000"`
	SkillsFrameworks    *string `json:"skills_frameworks" binding:"omitempty,max=1000"`
	AtCoderUsername     *string `json:"atcoder_username" binding:"omitempty,max=100"`
	PaizaRank           *string `json:"paiza_rank" binding:"omitempty,max=20"`
	OnboardingCompleted *bool   `json:"onboarding_completed"`
}
