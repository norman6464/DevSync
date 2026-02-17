package dto

// UpdateUserRequest はユーザープロフィール更新リクエスト。
type UpdateUserRequest struct {
	Name                string  `json:"name"`
	Bio                 string  `json:"bio"`
	AvatarURL           string  `json:"avatar_url"`
	SkillsLanguages     *string `json:"skills_languages"`
	SkillsFrameworks    *string `json:"skills_frameworks"`
	AtCoderUsername     *string `json:"atcoder_username"`
	PaizaRank           *string `json:"paiza_rank"`
	OnboardingCompleted *bool   `json:"onboarding_completed"`
}
