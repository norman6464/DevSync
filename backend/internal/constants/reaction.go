package constants

// AllowedReactionEmojis はリアクションに使用可能な絵文字一覧。
// フロントエンドとバックエンドで共通の定数として一元管理する。
var AllowedReactionEmojis = []string{
	"👍",
	"🎉",
	"❤️",
	"🔥",
	"👀",
}

// IsAllowedReactionEmoji は指定された絵文字がリアクションに使用可能かを判定する。
func IsAllowedReactionEmoji(emoji string) bool {
	for _, allowed := range AllowedReactionEmojis {
		if allowed == emoji {
			return true
		}
	}
	return false
}
