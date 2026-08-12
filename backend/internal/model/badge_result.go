package model

// BadgeResult は個別バッジの獲得状況を表す。
// DB テーブルには対応せず、判定結果を返すための型。
type BadgeResult struct {
	ID          string `json:"id"`          // バッジ識別子
	Name        string `json:"name"`        // バッジ名（i18nキー）
	Description string `json:"description"` // バッジ説明（i18nキー）
	Category    string `json:"category"`    // バッジカテゴリ
	Earned      bool   `json:"earned"`      // 獲得済みフラグ
}
