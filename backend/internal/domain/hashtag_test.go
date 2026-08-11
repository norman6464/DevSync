package domain_test

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestExtractHashtags_SingleTag(t *testing.T) {
	result := domain.ExtractHashtags("今日は #golang を学んだ")
	assert.Equal(t, []string{"golang"}, result)
}

func TestExtractHashtags_MultipleTags(t *testing.T) {
	result := domain.ExtractHashtags("#React と #TypeScript でアプリを作成")
	assert.Equal(t, []string{"React", "TypeScript"}, result)
}

func TestExtractHashtags_JapaneseTags(t *testing.T) {
	result := domain.ExtractHashtags("#Go言語 を使って #Web開発 を学習中")
	assert.Equal(t, []string{"Go言語", "Web開発"}, result)
}

func TestExtractHashtags_Deduplication(t *testing.T) {
	result := domain.ExtractHashtags("#golang は素晴らしい。#golang が好き")
	assert.Equal(t, []string{"golang"}, result)
}

func TestExtractHashtags_CasePreserved(t *testing.T) {
	result := domain.ExtractHashtags("#TypeScript と #typescript")
	// 大文字小文字を区別せず重複排除（最初の出現を保持）
	assert.Len(t, result, 1)
	assert.Equal(t, "TypeScript", result[0])
}

func TestExtractHashtags_IgnoreInCodeBlock(t *testing.T) {
	content := "本文 #validTag\n```\n#notATag\nfmt.Println(\"hello\")\n```\n#anotherTag"
	result := domain.ExtractHashtags(content)
	assert.Equal(t, []string{"validTag", "anotherTag"}, result)
}

func TestExtractHashtags_IgnoreInInlineCode(t *testing.T) {
	content := "これは `#notATag` ではなく #realTag です"
	result := domain.ExtractHashtags(content)
	assert.Equal(t, []string{"realTag"}, result)
}

func TestExtractHashtags_IgnoreURLFragments(t *testing.T) {
	content := "https://example.com/page#section は参照です #validTag"
	result := domain.ExtractHashtags(content)
	assert.Equal(t, []string{"validTag"}, result)
}

func TestExtractHashtags_EmptyContent(t *testing.T) {
	result := domain.ExtractHashtags("")
	assert.Empty(t, result)
}

func TestExtractHashtags_NoHashtags(t *testing.T) {
	result := domain.ExtractHashtags("ハッシュタグなしの通常テキスト")
	assert.Empty(t, result)
}

func TestExtractHashtags_MaxLimit(t *testing.T) {
	// 最大10個まで抽出
	content := "#a1 #a2 #a3 #a4 #a5 #a6 #a7 #a8 #a9 #a10 #a11 #a12"
	result := domain.ExtractHashtags(content)
	assert.Len(t, result, 10)
}

func TestExtractHashtags_MinLength(t *testing.T) {
	// 1文字のタグは無視
	content := "#a #ab #abc"
	result := domain.ExtractHashtags(content)
	assert.Equal(t, []string{"ab", "abc"}, result)
}

func TestExtractHashtags_HashOnly(t *testing.T) {
	result := domain.ExtractHashtags("# だけ")
	assert.Empty(t, result)
}

func TestExtractHashtags_TagAtStartOfLine(t *testing.T) {
	result := domain.ExtractHashtags("#開発日記\n今日の成果")
	assert.Equal(t, []string{"開発日記"}, result)
}

func TestExtractHashtags_TagWithNumbers(t *testing.T) {
	result := domain.ExtractHashtags("#100DaysOfCode チャレンジ中")
	assert.Equal(t, []string{"100DaysOfCode"}, result)
}
