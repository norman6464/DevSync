import { useTranslation } from 'react-i18next';
import { Trophy, Flame, TrendingUp, Medal } from 'lucide-react';
import type { Post } from '../../types/post';
import type { BadgeResult } from '../../types/badge';
import { Link } from 'react-router-dom';

interface ProfileHighlightCardProps {
  posts: Post[];
  badges: BadgeResult[];
  streakDays: number;
}

/** 投稿の反応の多さ（いいね + コメント）。 */
function reactionScore(post: Post): number {
  return (post.like_count || 0) + (post.comment_count || 0);
}

export default function ProfileHighlightCard({ posts, badges, streakDays }: ProfileHighlightCardProps) {
  const { t } = useTranslation();

  // 最も反応が多い投稿を取得（いいね+コメント数）
  const topPost = posts.reduce<Post | null>((best, post) => {
    if (!best) return post;
    return reactionScore(post) > reactionScore(best) ? post : best;
  }, null);

  // 最新バッジを取得（獲得済みのみ。一覧は新しい順に並んでいる前提）
  const latestBadge = badges.find((badge) => badge.earned);

  // ハイライトが何もない場合は表示しない
  if (!topPost && !latestBadge && streakDays === 0) {
    return null;
  }

  return (
    <div className="bg-gradient-to-br from-gray-900 to-gray-800 border border-gray-700 rounded-md p-5">
      <h3 className="text-sm font-semibold text-gray-200 mb-4 flex items-center gap-2">
        <Trophy className="w-4 h-4 text-yellow-400" aria-hidden="true" />
        {t('profile.highlights.title')}
      </h3>

      <div className="space-y-3">
        {/* 人気の投稿 */}
        {topPost && (
          <div className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
            <TrendingUp className="w-4 h-4 text-blue-400 mt-1 shrink-0" aria-hidden="true" />
            <div className="flex-1 min-w-0">
              <p className="text-xs text-gray-400 mb-1">{t('profile.highlights.topPost')}</p>
              <Link
                to={`/posts/${topPost.id}`}
                className="text-sm text-gray-200 hover:text-white transition-colors line-clamp-2"
              >
                {topPost.title}
              </Link>
              <p className="text-xs text-gray-500 mt-1">
                {topPost.like_count || 0} {t('common.likes')} · {topPost.comment_count || 0} {t('common.comments')}
              </p>
            </div>
          </div>
        )}

        {/* 最新バッジ */}
        {latestBadge && (
          <div className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
            <Medal className="w-4 h-4 text-yellow-400 mt-1 shrink-0" aria-hidden="true" />
            <div className="flex-1 min-w-0">
              <p className="text-xs text-gray-400 mb-1">{t('profile.highlights.latestBadge')}</p>
              <p className="text-sm text-gray-200 font-medium">{latestBadge.name}</p>
              <p className="text-xs text-gray-500 mt-1">{latestBadge.description}</p>
            </div>
          </div>
        )}

        {/* 連続学習日数 */}
        {streakDays > 0 && (
          <div className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
            <Flame className="w-4 h-4 text-orange-400 mt-1 shrink-0" aria-hidden="true" />
            <div className="flex-1 min-w-0">
              <p className="text-xs text-gray-400 mb-1">{t('profile.highlights.streak')}</p>
              <p className="text-sm text-gray-200 font-medium">
                {streakDays} {t('profile.highlights.streakDays')}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
