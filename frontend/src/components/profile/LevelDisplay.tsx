import { useTranslation } from 'react-i18next';
import { Star } from 'lucide-react';
import { useLevel, useLevelBreakdown } from '../../hooks';
import type { XPBreakdown } from '../../types/level';

function getLevelTitleKey(level: number): { key: string; params?: Record<string, unknown> } {
  if (level === 0) return { key: 'level.titleNewcomer' };
  if (level <= 5) return { key: 'level.titleBeginner' };
  if (level <= 10) return { key: 'level.titleIntermediate' };
  if (level <= 20) return { key: 'level.titleAdvanced' };
  if (level <= 30) return { key: 'level.titleExpert' };
  if (level <= 40) return { key: 'level.titleMaster' };
  return { key: 'level.titleLegend', params: { level } };
}

interface LevelDisplayProps {
  userId: number;
}

/** XP内訳の各項目 */
function BreakdownItem({ label, xp, maxXP }: { label: string; xp: number; maxXP: number }) {
  const { t } = useTranslation();
  const percent = maxXP > 0 ? (xp / maxXP) * 100 : 0;

  if (xp === 0) return null;

  return (
    <div className="flex items-center gap-3">
      <span className="text-xs text-gray-400 w-28 shrink-0 truncate">{label}</span>
      <div className="flex-1 h-2 bg-gray-800 rounded-full overflow-hidden">
        <div
          className="h-full bg-gradient-to-r from-yellow-500/80 to-orange-500/80 rounded-full transition-all"
          style={{ width: `${Math.min(percent, 100)}%` }}
        />
      </div>
      <span className="text-xs text-gray-300 w-16 text-right shrink-0">
        {xp.toLocaleString()} {t('level.xpUnit')}
      </span>
    </div>
  );
}

export default function LevelDisplay({ userId }: LevelDisplayProps) {
  const { t } = useTranslation();
  const { levelInfo } = useLevel(userId);
  const { breakdown } = useLevelBreakdown(userId);

  if (!levelInfo) return null;

  const { level, total_xp, next_level_xp, progress_percent } = levelInfo;

  // XP内訳の最大値を取得（バーの相対幅計算用）
  const breakdownItems: { label: string; xp: number }[] = breakdown
    ? [
        { label: t('level.breakdownLearningLogs'), xp: breakdown.learning_logs },
        { label: t('level.breakdownPosts'), xp: breakdown.posts },
        { label: t('level.breakdownGitHub'), xp: breakdown.github },
        { label: t('level.breakdownGoals'), xp: breakdown.goals },
        { label: t('level.breakdownComments'), xp: breakdown.comments },
        { label: t('level.breakdownLikes'), xp: breakdown.likes },
        { label: t('level.breakdownStreakBonus'), xp: breakdown.streak_bonus },
      ]
    : [];
  const maxXP = Math.max(...breakdownItems.map((i) => i.xp), 1);

  if (total_xp === 0) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
      <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide flex items-center gap-2 mb-4">
        <Star className="w-5 h-5 text-yellow-400" />
        {t('level.title')}
      </h2>

      {/* Level + Progress */}
      <div className="flex items-center gap-4 mb-4">
        <div className="text-center">
          <span className="text-xs text-gray-500">{t('level.level')}</span>
          <div className="text-3xl font-bold text-yellow-400">{level}</div>
          <p className="text-xs text-yellow-400/70 mt-0.5">
            {t(getLevelTitleKey(level).key, getLevelTitleKey(level).params)}
          </p>
        </div>
        <div className="flex-1">
          <div className="flex justify-between text-xs text-gray-400 mb-1">
            <span>{t('level.totalXP')}: {total_xp.toLocaleString()}</span>
            <span>{Math.round(progress_percent)}%</span>
          </div>
          <div className="h-2.5 bg-gray-800 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-yellow-500 to-orange-500 rounded-full transition-all"
              style={{ width: `${Math.min(progress_percent, 100)}%` }}
            />
          </div>
          <div className="text-xs text-gray-500 mt-1">
            {t('level.xpNeeded', { xp: (next_level_xp - total_xp).toLocaleString() })}
          </div>
        </div>
      </div>

      {/* XP Breakdown */}
      {breakdown && breakdown.total > 0 && (
        <div>
          <h3 className="text-xs font-medium text-gray-400 mb-3">{t('level.breakdown')}</h3>
          <div className="space-y-2">
            {breakdownItems.map((item) => (
              <BreakdownItem
                key={item.label}
                label={item.label}
                xp={item.xp}
                maxXP={maxXP}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
