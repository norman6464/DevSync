import { useTranslation } from 'react-i18next';
import { useThemeStore } from '../../store/themeStore';
import type { CalendarEntry } from '../../types/learningLog';

interface LogCalendarProps {
  entries: CalendarEntry[];
  onDateClick?: (date: string) => void;
}

function getColor(count: number, resolvedTheme: 'dark' | 'light'): string {
  if (count === 0) return resolvedTheme === 'dark' ? '#1e293b' : '#ebedf0';
  if (count === 1) return '#7c3aed';
  if (count === 2) return '#8b5cf6';
  return '#a78bfa';
}

export default function LogCalendar({ entries, onDateClick }: LogCalendarProps) {
  const { t } = useTranslation();
  const resolvedTheme = useThemeStore((state) => state.resolvedTheme);

  const entryMap = new Map(entries.map((e) => [e.date, e.count]));

  const today = new Date();
  const weeks: { date: string; count: number }[][] = [];

  const startDate = new Date(today);
  startDate.setDate(startDate.getDate() - 364);
  startDate.setDate(startDate.getDate() - startDate.getDay());

  let currentWeek: { date: string; count: number }[] = [];
  const cursor = new Date(startDate);

  while (cursor <= today) {
    const dateStr = cursor.toISOString().split('T')[0];
    currentWeek.push({
      date: dateStr,
      count: entryMap.get(dateStr) || 0,
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

  const totalLogs = entries.reduce((sum, e) => sum + e.count, 0);

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
          <div className="flex flex-col gap-[3px] mr-2 text-xs text-gray-500">
            <div className="h-4" />
            {dayLabels.map((label, i) => (
              <div key={i} className="h-3 flex items-center justify-end pr-1 text-[10px]">
                {label}
              </div>
            ))}
          </div>

          <div>
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

            <div className="inline-flex gap-[3px]">
              {weeks.map((week, wi) => (
                <div key={wi} className="flex flex-col gap-[3px]">
                  {week.map((day, di) => (
                    <div
                      key={di}
                      title={t('common.dateLogs', { date: day.date, count: day.count })}
                      className={`w-3 h-3 rounded-sm ${onDateClick && day.count > 0 ? 'cursor-pointer hover:ring-1 hover:ring-purple-400' : ''}`}
                      style={{ backgroundColor: getColor(day.count, resolvedTheme) }}
                      onClick={() => onDateClick && day.count > 0 && onDateClick(day.date)}
                    />
                  ))}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
      <p className="text-sm text-gray-500 mt-3">
        {t('learningLogs.totalLogs', { count: totalLogs })}
      </p>
    </div>
  );
}
