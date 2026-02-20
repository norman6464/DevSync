package service

import (
	"regexp"
	"strings"
)

// maxHashtags は1コンテンツから抽出するハッシュタグの最大数。
const maxHashtags = 10

// hashtagRegex はテキスト中の #hashtag パターンを抽出する正規表現。
// 英数字・アンダースコア・CJK文字に対応。URL中のフラグメントを除外するため、
// #の前が英数字・スラッシュ・ドットでないことを確認する。
var hashtagRegex = regexp.MustCompile(`(?:^|[^\w/.])(#([\p{L}\p{N}_][\p{L}\p{N}_]*))`)

// codeBlockRegex はMarkdownコードブロック（```...```）にマッチする正規表現。
var codeBlockRegex = regexp.MustCompile("(?s)```.*?```")

// inlineCodeRegex はMarkdownインラインコード（`...`）にマッチする正規表現。
var inlineCodeRegex = regexp.MustCompile("`[^`]+`")

// ExtractHashtags はテキストからハッシュタグ（#tag）を重複なしで抽出する。
// コードブロック・インラインコード内のタグ、URLフラグメント、1文字タグは除外する。
// 大文字小文字を区別せずに重複排除し、最初の出現形を保持する。最大10個まで。
func ExtractHashtags(content string) []string {
	if content == "" {
		return nil
	}

	// コードブロックとインラインコードを除去
	cleaned := codeBlockRegex.ReplaceAllString(content, "")
	cleaned = inlineCodeRegex.ReplaceAllString(cleaned, "")

	matches := hashtagRegex.FindAllStringSubmatch(cleaned, -1)
	seen := make(map[string]bool)
	var tags []string

	for _, match := range matches {
		tag := match[2]
		// 1文字のタグは無視
		if len([]rune(tag)) < 2 {
			continue
		}
		lower := strings.ToLower(tag)
		if !seen[lower] {
			seen[lower] = true
			tags = append(tags, tag)
			if len(tags) >= maxHashtags {
				break
			}
		}
	}

	return tags
}
