package model

// RecommendedUser はスキルマッチングに基づくおすすめユーザーを表す。
type RecommendedUser struct {
	User         User     `json:"user"`
	CommonSkills []string `json:"common_skills"`
	MatchScore   int      `json:"match_score"`
}
