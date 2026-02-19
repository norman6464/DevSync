package model

// QAStats はユーザーのQ&A活動集計統計を表す。
type QAStats struct {
	TotalQuestions     int64 `json:"total_questions"`
	TotalAnswers       int64 `json:"total_answers"`
	BestAnswerCount    int64 `json:"best_answer_count"`
	TotalVotesReceived int64 `json:"total_votes_received"`
}
