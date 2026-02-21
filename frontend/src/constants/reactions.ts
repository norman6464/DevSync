// 利用可能なリアクション絵文字
// バックエンドの AllowedReactionEmojis と一致させる必要がある
export const AVAILABLE_REACTION_EMOJIS = ['👍', '🎉', '❤️', '🔥', '👀'] as const;

export type ReactionEmoji = typeof AVAILABLE_REACTION_EMOJIS[number];
