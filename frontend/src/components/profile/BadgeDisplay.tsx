import type { GitHubContribution } from '../../types/github';
import type { Post } from '../../types/post';

interface BadgeDisplayProps {
  contributions: GitHubContribution[];
  posts: Post[];
  followerCount: number;
  followingCount: number;
}

interface Badge {
  id: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  bgColor: string;
  earned: boolean;
}

export default function BadgeDisplay({
  contributions,
  posts,
  followerCount,
  followingCount,
}: BadgeDisplayProps) {
  // Calculate stats
  const totalContributions = contributions.reduce((sum, c) => sum + c.count, 0);
  const totalPosts = posts.length;
  const totalLikes = posts.reduce((sum, p) => sum + (p.like_count || 0), 0);

  // Calculate streak
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
      const diffDays = Math.floor(
        (currentDate.getTime() - contribDate.getTime()) / (24 * 60 * 60 * 1000)
      );
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

  // Define badges
  const badges: Badge[] = [
    // Contribution badges
    {
      id: 'first-commit',
      name: 'First Commit',
      description: 'Made your first contribution',
      icon: '🌱',
      color: 'text-green-400',
      bgColor: 'bg-green-500/10 border-green-500/30',
      earned: totalContributions >= 1,
    },
    {
      id: 'contributor',
      name: 'Contributor',
      description: '50+ contributions',
      icon: '💻',
      color: 'text-blue-400',
      bgColor: 'bg-blue-500/10 border-blue-500/30',
      earned: totalContributions >= 50,
    },
    {
      id: 'code-warrior',
      name: 'Code Warrior',
      description: '200+ contributions',
      icon: '⚔️',
      color: 'text-purple-400',
      bgColor: 'bg-purple-500/10 border-purple-500/30',
      earned: totalContributions >= 200,
    },
    {
      id: 'commit-master',
      name: 'Commit Master',
      description: '500+ contributions',
      icon: '👑',
      color: 'text-yellow-400',
      bgColor: 'bg-yellow-500/10 border-yellow-500/30',
      earned: totalContributions >= 500,
    },
    {
      id: 'legend',
      name: 'Legend',
      description: '1000+ contributions',
      icon: '🏆',
      color: 'text-orange-400',
      bgColor: 'bg-orange-500/10 border-orange-500/30',
      earned: totalContributions >= 1000,
    },

    // Streak badges
    {
      id: 'week-streak',
      name: 'Week Warrior',
      description: '7-day streak',
      icon: '🔥',
      color: 'text-orange-400',
      bgColor: 'bg-orange-500/10 border-orange-500/30',
      earned: streak >= 7,
    },
    {
      id: 'month-streak',
      name: 'Consistency King',
      description: '30-day streak',
      icon: '⚡',
      color: 'text-yellow-400',
      bgColor: 'bg-yellow-500/10 border-yellow-500/30',
      earned: streak >= 30,
    },

    // Post badges
    {
      id: 'first-post',
      name: 'Writer',
      description: 'Published first post',
      icon: '✍️',
      color: 'text-cyan-400',
      bgColor: 'bg-cyan-500/10 border-cyan-500/30',
      earned: totalPosts >= 1,
    },
    {
      id: 'blogger',
      name: 'Blogger',
      description: '10+ posts',
      icon: '📝',
      color: 'text-indigo-400',
      bgColor: 'bg-indigo-500/10 border-indigo-500/30',
      earned: totalPosts >= 10,
    },

    // Engagement badges
    {
      id: 'liked',
      name: 'Liked',
      description: 'Received 10+ likes',
      icon: '❤️',
      color: 'text-pink-400',
      bgColor: 'bg-pink-500/10 border-pink-500/30',
      earned: totalLikes >= 10,
    },
    {
      id: 'popular',
      name: 'Popular',
      description: 'Received 50+ likes',
      icon: '🌟',
      color: 'text-rose-400',
      bgColor: 'bg-rose-500/10 border-rose-500/30',
      earned: totalLikes >= 50,
    },

    // Social badges
    {
      id: 'friendly',
      name: 'Friendly',
      description: 'Following 5+ developers',
      icon: '🤝',
      color: 'text-teal-400',
      bgColor: 'bg-teal-500/10 border-teal-500/30',
      earned: followingCount >= 5,
    },
    {
      id: 'influencer',
      name: 'Influencer',
      description: '10+ followers',
      icon: '📢',
      color: 'text-violet-400',
      bgColor: 'bg-violet-500/10 border-violet-500/30',
      earned: followerCount >= 10,
    },
    {
      id: 'star',
      name: 'Rising Star',
      description: '50+ followers',
      icon: '⭐',
      color: 'text-amber-400',
      bgColor: 'bg-amber-500/10 border-amber-500/30',
      earned: followerCount >= 50,
    },
  ];

  const earnedBadges = badges.filter((b) => b.earned);
  const lockedBadges = badges.filter((b) => !b.earned);

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-800 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <svg
            className="w-4 h-4 text-yellow-400"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-2.748 0"
            />
          </svg>
          <h2 className="text-sm font-semibold">Achievements</h2>
        </div>
        <span className="text-xs text-gray-500">
          {earnedBadges.length}/{badges.length} unlocked
        </span>
      </div>

      <div className="p-6">
        {/* Earned Badges */}
        {earnedBadges.length > 0 && (
          <div className="mb-6">
            <h3 className="text-xs text-gray-400 uppercase tracking-wide mb-3">Earned</h3>
            <div className="flex flex-wrap gap-2">
              {earnedBadges.map((badge) => (
                <div
                  key={badge.id}
                  className={`group relative px-3 py-2 rounded-lg border ${badge.bgColor} cursor-pointer transition-all hover:scale-105`}
                >
                  <div className="flex items-center gap-2">
                    <span className="text-lg">{badge.icon}</span>
                    <span className={`text-sm font-medium ${badge.color}`}>{badge.name}</span>
                  </div>
                  {/* Tooltip */}
                  <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-800 rounded-lg text-xs text-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all whitespace-nowrap z-10 shadow-lg">
                    <div className="font-medium text-white">{badge.name}</div>
                    <div className="text-gray-400">{badge.description}</div>
                    <div className="absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-800" />
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Locked Badges */}
        {lockedBadges.length > 0 && (
          <div>
            <h3 className="text-xs text-gray-500 uppercase tracking-wide mb-3">Locked</h3>
            <div className="flex flex-wrap gap-2">
              {lockedBadges.map((badge) => (
                <div
                  key={badge.id}
                  className="group relative px-3 py-2 rounded-lg border border-gray-700/50 bg-gray-800/30 cursor-pointer transition-all hover:border-gray-600"
                >
                  <div className="flex items-center gap-2 opacity-40">
                    <span className="text-lg grayscale">{badge.icon}</span>
                    <span className="text-sm font-medium text-gray-500">{badge.name}</span>
                  </div>
                  {/* Tooltip */}
                  <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-800 rounded-lg text-xs text-center opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all whitespace-nowrap z-10 shadow-lg">
                    <div className="font-medium text-white">{badge.name}</div>
                    <div className="text-gray-400">{badge.description}</div>
                    <div className="absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-gray-800" />
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {earnedBadges.length === 0 && (
          <div className="text-center text-gray-500 py-4">
            <div className="text-3xl mb-2">🏅</div>
            <p className="text-sm">Start contributing to earn badges!</p>
          </div>
        )}
      </div>
    </div>
  );
}
