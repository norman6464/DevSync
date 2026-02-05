import type { GitHubContribution } from '../../types/github';

interface ContributionCalendarProps {
  contributions: GitHubContribution[];
}

function getColor(count: number): string {
  if (count === 0) return '#1e293b';
  if (count <= 3) return '#166534';
  if (count <= 6) return '#16a34a';
  if (count <= 9) return '#22c55e';
  return '#4ade80';
}

export default function ContributionCalendar({ contributions }: ContributionCalendarProps) {
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

  return (
    <div>
      <div className="overflow-x-auto">
        <div className="inline-flex gap-[3px]">
          {weeks.map((week, wi) => (
            <div key={wi} className="flex flex-col gap-[3px]">
              {week.map((day, di) => (
                <div
                  key={di}
                  title={`${day.date}: ${day.count} contributions`}
                  className="w-3 h-3 rounded-sm"
                  style={{ backgroundColor: getColor(day.count) }}
                />
              ))}
            </div>
          ))}
        </div>
      </div>
      <p className="text-sm text-gray-400 mt-3">
        {totalContributions.toLocaleString()} contributions in the last year
      </p>
    </div>
  );
}
