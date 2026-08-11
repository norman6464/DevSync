package model

// AtCoderRatingEntry はAtCoderのレーティング履歴エントリを表す。
// AtCoderのAPIレスポンスをそのまま受け取るため、フィールド名は外部APIのキーに合わせている。
type AtCoderRatingEntry struct {
	IsRated     bool   `json:"IsRated"`
	Place       int    `json:"Place"`
	OldRating   int    `json:"OldRating"`
	NewRating   int    `json:"NewRating"`
	Performance int    `json:"Performance"`
	ContestName string `json:"ContestName"`
	EndTime     string `json:"EndTime"`
}

// AtCoderRatingInfo はAtCoderのレーティング情報を表す。
type AtCoderRatingInfo struct {
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Color    string `json:"color"`
	Rank     string `json:"rank"`
}
