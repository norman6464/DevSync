import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flame, Trophy, Calendar, TrendingUp, Shield } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useAsyncData } from '../../hooks/useAsyncData';
import { getStreakInfo, getCalendarData } from '../../api/learningLogs';
import { useStreakFreeze } from '../../hooks/useStreakFreeze';
import type { CalendarEntry } from '../../types/learningLog';

export default function LearningStreakCard() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const { data: streakInfo, loading: streakLoading } = useAsyncData(
    async () => {
      if (!user) return null;
      const res = await getStreakInfo(user.id);
      return res.data;
    },
    { initialData: null, deps: [user?.id], enabled: !!user }
  );

  const { data: calendarData, loading: calendarLoading } = useAsyncData(
    async () => {
      if (!user) return [];
      const res = await getCalendarData(user.id);
      return res.data;
    },
    { initialData: [] as CalendarEntry[], deps: [user?.id], enabled: !!user }
  );

  const { freezeStatus, applyFreeze } = useStreakFreeze();

  const loading = streakLoading || calendarLoading;

  const currentStreak = streakInfo?.current_streak ?? 0;
  const longestStreak = streakInfo?.longest_streak ?? 0;
  const totalDays = streakInfo?.total_days ?? 0;

  // ストリークレベル判定
  const streakLevel = useMemo(() => {
    if (currentStreak === 0) return { label: '-', color: 'text-gray-500' };
    if (currentStreak < 7) return { label: t('streak.levelBeginner'), color: 'text-green-400' };
    if (currentStreak < 30) return { label: t('streak.levelIntermediate'), color: 'text-blue-400' };
    if (currentStreak < 100) return { label: t('streak.levelAdvanced'), color: 'text-purple-400' };
    return { label: t('streak.levelMaster'), color: 'text-yellow-400' };
  }, [currentStreak, t]);

  // マイルストーンメッセージ
  const milestoneMessage = useMemo(() => {
    if (currentStreak === 0) return t('streak.noStreak');
    if (currentStreak === 1) return t('streak.milestoneFirstDay');
    if (currentStreak === 7) return t('streak.milestone7');
    if (currentStreak === 30) return t('streak.milestone30');
    if (currentStreak === 100) return t('streak.milestone100');
    if (currentStreak % 7 === 0) return t('streak.milestoneWeekly', { count: currentStreak });
    const nextMilestone = currentStreak < 7 ? 7 : currentStreak < 30 ? 30 : currentStreak < 100 ? 100 : (Math.floor(currentStreak / 100) + 1) * 100;
    return t('streak.milestoneNext', { count: nextMilestone });
  }, [currentStreak, t]);

  // 直近30日のカレンダーデータ生成
  const last30Days = useMemo(() => {
    const days: { date: string; count: number; dayOfWeek: number }[] = [];
    const today = new Date();

    for (let i = 29; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);
      const dateStr = date.toISOString().split('T')[0];
      const dayOfWeek = date.getDay();

      const entry = calendarData.find((e) => e.date.split('T')[0] === dateStr);
      days.push({
        date: dateStr,
        count: entry?.count ?? 0,
        dayOfWeek,
      });
    }

    return days;
  }, [calendarData]);

  // カウント数に応じた色クラス
  const getCountColor = (count: number) => {
    if (count === 0) return 'bg-gray-800';
    if (count === 1) return 'bg-green-700';
    if (count === 2) return 'bg-green-600';
    if (count === 3) return 'bg-blue-500';
    return 'bg-orange-500';
  };

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
        <div className="h-6 bg-gray-800 rounded animate-pulse w-1/2 mb-4" />
        <div className="h-20 bg-gray-800 rounded animate-pulse mb-4" />
        <div className="grid grid-cols-3 gap-3 mb-4">
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
        </div>
        <div className="h-24 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-2 text-lg font-semibold text-white">
          <Flame className={`w-5 h-5 ${currentStreak > 0 ? 'text-orange-400' : 'text-gray-500'}`} />
          {t('streak.title')}
        </h3>
        <div className={`px-3 py-1 rounded-full text-xs font-medium ${streakLevel.color} bg-gray-800`}>
          {streakLevel.label}
        </div>
      </div>

      {/* Current Streak */}
      <div className="text-center py-4 bg-gradient-to-br from-gray-800/50 to-gray-900 rounded-lg border border-gray-700">
        <div className="flex items-center justify-center gap-3 mb-2">
          <span className="text-5xl font-bold text-white">{currentStreak}</span>
          <div className="text-left">
            <div className="text-sm text-gray-400">{t('streak.daysConsecutive')}</div>
            <div className="text-xs text-gray-500">{milestoneMessage}</div>
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-3 gap-3">
        <div className="bg-gray-800/50 rounded-lg p-3 text-center">
          <Trophy className="w-4 h-4 text-yellow-400 mx-auto mb-1" />
          <div className="text-xl font-bold text-white">{longestStreak}</div>
          <div className="text-xs text-gray-500">{t('streak.longestStreak')}</div>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-3 text-center">
          <Calendar className="w-4 h-4 text-blue-400 mx-auto mb-1" />
          <div className="text-xl font-bold text-white">{totalDays}</div>
          <div className="text-xs text-gray-500">{t('streak.totalDays')}</div>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-3 text-center">
          <TrendingUp className="w-4 h-4 text-green-400 mx-auto mb-1" />
          <div className="text-xl font-bold text-white">
            {last30Days.filter((d) => d.count > 0).length}
          </div>
          <div className="text-xs text-gray-500">{t('streak.monthlyDays')}</div>
        </div>
      </div>

      {/* Streak Freeze */}
      {freezeStatus && (
        <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Shield className="w-4 h-4 text-cyan-400" />
              <div>
                <div className="text-sm font-medium text-white">{t('streak.freezeTitle')}</div>
                <div className="text-xs text-gray-400">
                  {t('streak.freezeCount', { used: freezeStatus.used_freezes, max: freezeStatus.max_freezes })}
                </div>
              </div>
            </div>
            <button
              onClick={() => {
                if (confirm(t('streak.freezeConfirm'))) {
                  applyFreeze();
                }
              }}
              disabled={!freezeStatus.can_use_today}
              className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors ${
                freezeStatus.can_use_today
                  ? 'bg-cyan-600 hover:bg-cyan-500 text-white'
                  : 'bg-gray-700 text-gray-500 cursor-not-allowed'
              }`}
            >
              {freezeStatus.today_used
                ? t('streak.freezeAlreadyUsed')
                : freezeStatus.remaining === 0
                ? t('streak.freezeUnavailable')
                : t('streak.freezeButton')}
            </button>
          </div>
        </div>
      )}

      {/* Calendar Grid */}
      <div>
        <h4 className="text-xs font-medium text-gray-400 mb-2">{t('streak.last30Days')}</h4>
        <div
          data-testid="streak-calendar"
          className="grid grid-cols-10 gap-1.5"
        >
          {last30Days.map((day) => (
            <div
              key={day.date}
              data-date={day.date}
              data-count={day.count}
              className={`aspect-square rounded ${getCountColor(day.count)} transition-colors hover:ring-2 hover:ring-blue-400`}
              title={t('streak.calendarTooltip', { date: day.date, count: day.count })}
            />
          ))}
        </div>
        <div className="flex items-center justify-end gap-2 mt-2 text-xs text-gray-500">
          <span>{t('streak.legendLess')}</span>
          <div className="flex gap-1">
            <div className="w-3 h-3 rounded bg-gray-800" />
            <div className="w-3 h-3 rounded bg-green-700" />
            <div className="w-3 h-3 rounded bg-green-600" />
            <div className="w-3 h-3 rounded bg-blue-500" />
            <div className="w-3 h-3 rounded bg-orange-500" />
          </div>
          <span>{t('streak.legendMore')}</span>
        </div>
      </div>
    </div>
  );
}
