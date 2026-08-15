import { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Trophy, Award, Zap, Users } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useAsyncData } from '../hooks/useAsyncData';
import { getUserBadges } from '../api/badges';
import { PageLoader } from '../components/common';

type BadgeCategory = 'all' | 'learning' | 'streak' | 'community';

const CATEGORY_CONFIG = {
  all: { label: 'すべて', Icon: Trophy },
  learning: { label: '学習', Icon: Award },
  streak: { label: 'ストリーク', Icon: Zap },
  community: { label: 'コミュニティ', Icon: Users },
};

export default function BadgeCollectionPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [selectedCategory, setSelectedCategory] = useState<BadgeCategory>('all');

  const { data: badgesData, loading } = useAsyncData(
    async () => {
      if (!user) return { badges: [] };
      const res = await getUserBadges(user.id);
      return res.data;
    },
    { initialData: { badges: [] }, deps: [user?.id], enabled: !!user }
  );

  const badges = badgesData.badges;

  const filteredBadges = useMemo(() => {
    if (selectedCategory === 'all') return badges;
    return badges.filter((badge) => badge.category === selectedCategory);
  }, [badges, selectedCategory]);

  const earnedBadges = filteredBadges.filter((b) => b.earned);
  const unearnedBadges = filteredBadges.filter((b) => !b.earned);
  const totalBadges = badges.length;
  const earnedCount = badges.filter((b) => b.earned).length;
  const progressPercentage = totalBadges > 0 ? Math.round((earnedCount / totalBadges) * 100) : 0;

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t('badges.title', 'バッジコレクション')}</h1>
          <p className="text-sm text-gray-400 mt-1">
            {t('badges.subtitle', '獲得したバッジと目標を確認しよう')}
          </p>
        </div>
      </div>

      {/* Progress */}
      <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Trophy className="w-5 h-5 text-yellow-400" />
            <span className="text-sm font-medium">
              {t('badges.progress', 'コレクション進捗')}
            </span>
          </div>
          <span className="text-2xl font-bold text-yellow-400">{progressPercentage}%</span>
        </div>
        <div className="h-3 bg-gray-700 rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-yellow-500 to-orange-500 transition-all"
            style={{ width: `${progressPercentage}%` }}
          />
        </div>
        <p className="text-xs text-gray-500 mt-2">
          {earnedCount} / {totalBadges} {t('badges.earned', 'バッジ獲得')}
        </p>
      </div>

      {/* Category Filters */}
      <div className="flex gap-2 flex-wrap">
        {(Object.keys(CATEGORY_CONFIG) as BadgeCategory[]).map((category) => {
          const { label, Icon } = CATEGORY_CONFIG[category];
          const isSelected = selectedCategory === category;

          return (
            <button
              key={category}
              onClick={() => setSelectedCategory(category)}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                isSelected
                  ? 'bg-blue-500/20 text-blue-400 border border-blue-400/30'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white border border-gray-700'
              }`}
            >
              <Icon className="w-4 h-4" />
              {label}
            </button>
          );
        })}
      </div>

      {/* Earned Badges */}
      {earnedBadges.length > 0 && (
        <div>
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <Trophy className="w-5 h-5 text-yellow-400" />
            {t('badges.earned_section', '獲得済みバッジ')} ({earnedBadges.length})
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {earnedBadges.map((badge) => (
              <div
                key={badge.id}
                className="bg-gray-900 border border-yellow-400/30 rounded-md p-4 hover:border-yellow-400/50 transition-colors"
              >
                <div className="flex items-start gap-3">
                  <div className="w-12 h-12 bg-gradient-to-br from-yellow-400 to-orange-500 rounded-full flex items-center justify-center flex-shrink-0">
                    <Trophy className="w-6 h-6 text-white" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-white">{badge.name}</h3>
                    <p className="text-xs text-gray-400 mt-1">{badge.description}</p>
                    <span className="inline-block mt-2 px-2 py-0.5 text-xs bg-yellow-400/10 text-yellow-400 rounded">
                      {t('badges.category_' + badge.category, badge.category)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Unearned Badges */}
      {unearnedBadges.length > 0 && (
        <div>
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <Award className="w-5 h-5 text-gray-500" />
            {t('badges.unearned_section', '未獲得バッジ')} ({unearnedBadges.length})
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {unearnedBadges.map((badge) => (
              <div
                key={badge.id}
                className="bg-gray-900 border border-gray-800 rounded-md p-4 opacity-60 hover:opacity-80 transition-opacity"
              >
                <div className="flex items-start gap-3">
                  <div className="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center flex-shrink-0">
                    <Award className="w-6 h-6 text-gray-500" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-400">{badge.name}</h3>
                    <p className="text-xs text-gray-500 mt-1">{badge.description}</p>
                    <span className="inline-block mt-2 px-2 py-0.5 text-xs bg-gray-700 text-gray-400 rounded">
                      {t('badges.category_' + badge.category, badge.category)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {filteredBadges.length === 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
          <Trophy className="w-16 h-16 mx-auto mb-4 text-gray-700" />
          <p className="text-gray-400">
            {t('badges.no_badges', 'このカテゴリにバッジはありません')}
          </p>
        </div>
      )}
    </div>
  );
}
