package model

// StudyCircleStats はスタディサークルの集計統計を表す。
type StudyCircleStats struct {
	MemberCount    int64 `json:"member_count"`
	CheckinCount   int64 `json:"checkin_count"`
	TotalSteps     int64 `json:"total_steps"`
	CompletedSteps int64 `json:"completed_steps"`
}
