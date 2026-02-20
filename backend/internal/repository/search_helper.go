package repository

import "strings"

// EscapeLikeChars はLIKE/ILIKEクエリのワイルドカード文字（%, _, \）をエスケープする。
// ラッピングなしのエスケープのみ。タグ検索など独自パターンを組み立てる場合に使用する。
func EscapeLikeChars(query string) string {
	query = strings.ReplaceAll(query, "\\", "\\\\")
	query = strings.ReplaceAll(query, "%", "\\%")
	query = strings.ReplaceAll(query, "_", "\\_")
	return query
}

// EscapeLikePattern はLIKE/ILIKEクエリのワイルドカード文字（%, _, \）をエスケープし、
// 部分一致検索パターンを生成する。ユーザー入力をそのままLIKEに渡すと
// ワイルドカード注入のリスクがあるため、この関数でエスケープする。
func EscapeLikePattern(query string) string {
	return "%" + EscapeLikeChars(query) + "%"
}
