import { useTranslation } from 'react-i18next';
import { Trophy, Flame, TrendingUp } from 'lucide-react';
import { Post, Badge } from '../../types';
import { Link } from 'react-router-dom';

interface ProfileHighlightCardProps {
  posts: Post[];
  badges: Badge[];
  streakDays: number;
}

export default function ProfileHighlightCard({ posts, badges, streakDays }: ProfileHighlightCardProps) {
  const { t } = useTranslation();

  // 最も反応が多い投稿を取得（いいね+コメント数）
  const topPost = posts.reduce<Post | null>((best, post) => {
    const score = (post.likes_count || 0) + (post.comments_count || 0);
    const bestScore = best ? (best.likes_count || 0) + (best.comments_count || 0) : 0;
    return score > bestScore ? post : best;
  }, null);

  // 最新バッジを取得
  const latestBadge = badges[0]; // バッジは新しい順にソートされていると仮定

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
                {topPost.likes_count || 0} {t('common.likes')} · {topPost.comments_count || 0} {t('common.comments')}
              </p>
            </div>
          </div>
        )}

        {/* 最新バッジ */}
        {latestBadge && (
          <div className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
            <span className="text-2xl shrink-0" aria-hidden="true">{latestBadge.icon}</span>
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
