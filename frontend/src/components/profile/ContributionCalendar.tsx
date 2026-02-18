import { useTranslation } from 'react-i18next';
import type { GitHubContribution } from '../../types/github';
import { useThemeStore } from '../../store/themeStore';

interface ContributionCalendarProps {
  contributions: GitHubContribution[];
}

function getColor(count: number, resolvedTheme: 'dark' | 'light'): string {
  if (count === 0) return resolvedTheme === 'dark' ? '#1e293b' : '#ebedf0';
  if (count <= 3) return '#166534';
  if (count <= 6) return '#16a34a';
  if (count <= 9) return '#22c55e';
  return '#4ade80';
}

export default function ContributionCalendar({ contributions }: ContributionCalendarProps) {
  const { t } = useTranslation();
  const resolvedTheme = useThemeStore((state) => state.resolvedTheme);
  const contributionMap = new Map(
    contributions.map((c) => [c.date.split('T')[0], c.count])
  );

  const today = new Date();
  const weeks: { date: string; count: number }[][] = [];

  // Build 52 weeks of data
  const startDate = new Date(today);
  startDate.setDate(startDate.getDate() - 364);
  startDate.setDate(startDate.getDate() - startDate.getDay()); // align to Sunday

  let currentWeek: { date: string; count: number }[] = [];
  const cursor = new Date(startDate);

  while (cursor <= today) {
    const dateStr = cursor.toISOString().split('T')[0];
    currentWeek.push({
      date: dateStr,
      count: contributionMap.get(dateStr) || 0,
    });

    if (currentWeek.length === 7) {
      weeks.push(currentWeek);
      currentWeek = [];
    }
    cursor.setDate(cursor.getDate() + 1);
  }
  if (currentWeek.length > 0) {
    weeks.push(currentWeek);
  }

  const totalContributions = contributions.reduce((sum, c) => sum + c.count, 0);

  // Calculate month labels for header
  const monthLabels: { label: string; weekIndex: number }[] = [];
  const months = t('common.monthsShort', { returnObjects: true }) as string[];
  let lastMonth = -1;

  weeks.forEach((week, weekIndex) => {
    if (week.length > 0) {
      const firstDayOfWeek = new Date(week[0].date);
      const month = firstDayOfWeek.getMonth();
      if (month !== lastMonth) {
        monthLabels.push({ label: months[month], weekIndex });
        lastMonth = month;
      }
    }
  });

  const dayLabels = t('common.dayLabelsShort', { returnObjects: true }) as string[];

  return (
    <div>
      <div className="overflow-x-auto">
        <div className="flex">
          {/* Day labels column */}
          <div className="flex flex-col gap-[3px] mr-2 text-xs text-gray-500">
            <div className="h-4" /> {/* Spacer for month header */}
            {dayLabels.map((label, i) => (
              <div key={i} className="h-3 flex items-center justify-end pr-1 text-[10px]">
                {label}
              </div>
            ))}
          </div>

          {/* Calendar grid */}
          <div>
            {/* Month labels row */}
            <div className="flex h-4 text-xs text-gray-500 mb-[3px]">
              {weeks.map((_, wi) => {
                const monthLabel = monthLabels.find(m => m.weekIndex === wi);
                return (
                  <div key={wi} className="w-3 mr-[3px] text-[10px]">
                    {monthLabel?.label || ''}
                  </div>
                );
              })}
            </div>

            {/* Contribution squares */}
            <div className="inline-flex gap-[3px]">
              {weeks.map((week, wi) => (
                <div key={wi} className="flex flex-col gap-[3px]">
                  {week.map((day, di) => (
                    <div
                      key={di}
                      title={t('common.dateContributions', { date: day.date, count: day.count })}
                      className="w-3 h-3 rounded-sm"
                      style={{ backgroundColor: getColor(day.count, resolvedTheme) }}
                    />
                  ))}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
      <p className="text-sm text-gray-500 mt-3">
        {t('common.contributionsInLastYear', { count: totalContributions })}
      </p>
    </div>
  );
}
