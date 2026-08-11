package persistence

import "strings"

// escapeLikeChars は LIKE/ILIKE クエリのワイルドカード文字（%, _, \）をエスケープする。
func escapeLikeChars(query string) string {
	query = strings.ReplaceAll(query, "\\", "\\\\")
	query = strings.ReplaceAll(query, "%", "\\%")
	query = strings.ReplaceAll(query, "_", "\\_")
	return query
}

// escapeLikePattern は LIKE/ILIKE クエリのワイルドカード文字をエスケープし、部分一致パターンを作る。
// ユーザー入力をそのまま LIKE に渡すとワイルドカード注入のリスクがあるため必ずこれを通す。
//
// 旧 repository パッケージにも同じ関数があるが、そちらは未移行の repository が使っており、
// adapter から撤去予定のパッケージへ依存させないためにこちらへ持たせている。
func escapeLikePattern(query string) string {
	return "%" + escapeLikeChars(query) + "%"
}
