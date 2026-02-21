package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAllowedReactionEmoji(t *testing.T) {
	t.Run("許可された絵文字の場合trueを返す", func(t *testing.T) {
		allowedEmojis := []string{"👍", "🎉", "❤️", "🔥", "👀"}
		for _, emoji := range allowedEmojis {
			assert.True(t, IsAllowedReactionEmoji(emoji), "絵文字 %s は許可されているべき", emoji)
		}
	})

	t.Run("許可されていない絵文字の場合falseを返す", func(t *testing.T) {
		notAllowedEmojis := []string{"😀", "💩", "🐶", "🌸", "🍕"}
		for _, emoji := range notAllowedEmojis {
			assert.False(t, IsAllowedReactionEmoji(emoji), "絵文字 %s は許可されていないべき", emoji)
		}
	})

	t.Run("空文字列の場合falseを返す", func(t *testing.T) {
		assert.False(t, IsAllowedReactionEmoji(""))
	})

	t.Run("通常の文字列の場合falseを返す", func(t *testing.T) {
		assert.False(t, IsAllowedReactionEmoji("thumbsup"))
		assert.False(t, IsAllowedReactionEmoji("like"))
		assert.False(t, IsAllowedReactionEmoji("react"))
	})
}

func TestAllowedReactionEmojis(t *testing.T) {
	t.Run("AllowedReactionEmojisに5つの絵文字が含まれている", func(t *testing.T) {
		assert.Len(t, AllowedReactionEmojis, 5)
	})

	t.Run("AllowedReactionEmojisに重複がない", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, emoji := range AllowedReactionEmojis {
			assert.False(t, seen[emoji], "絵文字 %s が重複しています", emoji)
			seen[emoji] = true
		}
	})

	t.Run("AllowedReactionEmojisに空文字列が含まれていない", func(t *testing.T) {
		for _, emoji := range AllowedReactionEmojis {
			assert.NotEmpty(t, emoji, "空文字列が含まれています")
		}
	})

	t.Run("AllowedReactionEmojisの順序が正しい", func(t *testing.T) {
		expected := []string{"👍", "🎉", "❤️", "🔥", "👀"}
		assert.Equal(t, expected, AllowedReactionEmojis)
	})
}
