import type { GitHubContribution } from '../../types/github';
import type { Post } from '../../types/post';

interface ActivityReportProps {
  contributions: GitHubContribution[];
  posts: Post[];
  followerCount: number;
}

export default function ActivityReport({ contributions, posts, followerCount }: ActivityReportProps) {
  const now = new Date();
  const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
  const twoWeeksAgo = new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000);
  const oneMonthAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);

  // Calculate contributions for different periods
  const thisWeekContributions = contributions.filter((c) => new Date(c.date) >= oneWeekAgo);
  const lastWeekContributions = contributions.filter(
    (c) => new Date(c.date) >= twoWeeksAgo && new Date(c.date) < oneWeekAgo
  );
  const thisMonthContributions = contributions.filter((c) => new Date(c.date) >= oneMonthAgo);

  const thisWeekTotal = thisWeekContributions.reduce((sum, c) => sum + c.count, 0);
  const lastWeekTotal = lastWeekContributions.reduce((sum, c) => sum + c.count, 0);
  const thisMonthTotal = thisMonthContributions.reduce((sum, c) => sum + c.count, 0);
  const totalContributions = contributions.reduce((sum, c) => sum + c.count, 0);

  // Calculate posts stats
  const thisWeekPosts = posts.filter((p) => new Date(p.created_at) >= oneWeekAgo).length;
  const thisMonthPosts = posts.filter((p) => new Date(p.created_at) >= oneMonthAgo).length;
  const totalLikes = posts.reduce((sum, p) => sum + (p.like_count || 0), 0);

  // Calculate percentage change
  const getPercentChange = (current: number, previous: number) => {
    if (previous === 0) return current > 0 ? 100 : 0;
    return Math.round(((current - previous) / previous) * 100);
  };

  const contributionChange = getPercentChange(thisWeekTotal, lastWeekTotal);

  // Find streak (consecutive days with contributions)
  const getStreak = () => {
    if (contributions.length === 0) return 0;

    const sortedContribs = [...contributions]
      .filter((c) => c.count > 0)
      .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

    if (sortedContribs.length === 0) return 0;

    let streak = 0;
    let currentDate = new Date();
    currentDate.setHours(0, 0, 0, 0);

    for (const contrib of sortedContribs) {
      const contribDate = new Date(contrib.date);
      contribDate.setHours(0, 0, 0, 0);

      const diffDays = Math.floor((currentDate.getTime() - contribDate.getTime()) / (24 * 60 * 60 * 1000));

      if (diffDays <= 1) {
        streak++;
        currentDate = contribDate;
      } else {
        break;
      }
    }

    return streak;
  };

  const streak = getStreak();

  const stats = [
    {
      label: 'This Week',
      value: thisWeekTotal,
      icon: (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5" />
        </svg>
      ),
      suffix: 'contributions',
      change: contributionChange,
    },
    {
      label: 'Current Streak',
      value: streak,
      icon: (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M15.362 5.214A8.252 8.252 0 0 1 12 21 8.25 8.25 0 0 1 6.038 7.047 8.287 8.287 0 0 0 9 9.601a8.983 8.983 0 0 1 3.361-6.867 8.21 8.21 0 0 0 3 2.48Z" />
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 18a3.75 3.75 0 0 0 .495-7.468 5.99 5.99 0 0 0-1.925 3.547 5.975 5.975 0 0 1-2.133-1.001A3.75 3.75 0 0 0 12 18Z" />
        </svg>
      ),
      suffix: 'days',
      color: streak >= 7 ? 'text-orange-400' : undefined,
    },
    {
      label: 'This Month',
      value: thisMonthTotal,
      icon: (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
        </svg>
      ),
      suffix: 'contributions',
    },
    {
      label: 'Total',
      value: totalContributions,
      icon: (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z" />
        </svg>
      ),
      suffix: 'contributions',
    },
  ];

  const socialStats = [
    { label: 'Posts', value: posts.length, thisWeek: thisWeekPosts, thisMonth: thisMonthPosts },
    { label: 'Total Likes', value: totalLikes },
    { label: 'Followers', value: followerCount },
  ];

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-800 flex items-center gap-2">
        <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
        </svg>
        <h2 className="text-sm font-semibold">Activity Report</h2>
      </div>

      <div className="p-6">
        {/* Contribution Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          {stats.map((stat) => (
            <div key={stat.label} className="bg-gray-800/50 rounded-lg p-4">
              <div className="flex items-center gap-2 text-gray-400 mb-2">
                {stat.icon}
                <span className="text-xs">{stat.label}</span>
              </div>
              <div className="flex items-baseline gap-1">
                <span className={`text-2xl font-bold ${stat.color || 'text-white'}`}>
                  {stat.value.toLocaleString()}
                </span>
                <span className="text-xs text-gray-500">{stat.suffix}</span>
              </div>
              {stat.change !== undefined && (
                <div className={`text-xs mt-1 ${stat.change >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                  {stat.change >= 0 ? '↑' : '↓'} {Math.abs(stat.change)}% vs last week
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Motivation Message */}
        {contributionChange > 0 && (
          <div className="bg-green-500/10 border border-green-500/20 rounded-lg p-4 mb-6">
            <div className="flex items-center gap-2">
              <span className="text-lg">🎉</span>
              <p className="text-sm text-green-300">
                Great job! You're {contributionChange}% more active than last week. Keep it up!
              </p>
            </div>
          </div>
        )}

        {streak >= 7 && (
          <div className="bg-orange-500/10 border border-orange-500/20 rounded-lg p-4 mb-6">
            <div className="flex items-center gap-2">
              <span className="text-lg">🔥</span>
              <p className="text-sm text-orange-300">
                You're on fire! {streak} day streak! Don't break the chain!
              </p>
            </div>
          </div>
        )}

        {/* Social Stats */}
        <div className="border-t border-gray-800 pt-4">
          <h3 className="text-xs text-gray-500 uppercase tracking-wide mb-3">Social Activity</h3>
          <div className="flex gap-6">
            {socialStats.map((stat) => (
              <div key={stat.label} className="text-center">
                <div className="text-xl font-bold text-white">{stat.value}</div>
                <div className="text-xs text-gray-500">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
