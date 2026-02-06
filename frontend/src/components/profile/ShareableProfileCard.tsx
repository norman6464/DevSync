import { forwardRef } from 'react';
import type { User } from '../../types/user';
import type { GitHubLanguageStat } from '../../types/github';

interface ShareableProfileCardProps {
  user: User;
  followerCount: number;
  followingCount: number;
  totalContributions: number;
  languages: GitHubLanguageStat[];
  postCount: number;
}

const ShareableProfileCard = forwardRef<HTMLDivElement, ShareableProfileCardProps>(
  ({ user, followerCount, followingCount, totalContributions, languages, postCount }, ref) => {
    const topLanguages = languages.slice(0, 5);
    const languageColors: Record<string, string> = {
      TypeScript: '#3178c6',
      JavaScript: '#f7df1e',
      Python: '#3776ab',
      Go: '#00add8',
      Rust: '#dea584',
      Java: '#b07219',
      'C++': '#f34b7d',
      C: '#555555',
      Ruby: '#701516',
      PHP: '#4f5d95',
      Swift: '#ffac45',
      Kotlin: '#a97bff',
      Scala: '#c22d40',
      HTML: '#e34c26',
      CSS: '#563d7c',
      Vue: '#41b883',
      Svelte: '#ff3e00',
      Dart: '#00b4ab',
    };

    return (
      <div
        ref={ref}
        className="w-[600px] h-[315px] bg-gradient-to-br from-gray-900 via-gray-800 to-purple-900 p-6 rounded-xl relative overflow-hidden"
        style={{ fontFamily: 'system-ui, sans-serif' }}
      >
        {/* Background decoration */}
        <div className="absolute top-0 right-0 w-64 h-64 bg-purple-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-blue-500/10 rounded-full blur-3xl" />

        <div className="relative z-10 flex h-full">
          {/* Left side - Avatar and basic info */}
          <div className="flex flex-col items-center justify-center w-48">
            <img
              src={user.avatar_url || `https://ui-avatars.com/api/?name=${encodeURIComponent(user.name || 'User')}&background=7c3aed&color=fff&size=128`}
              alt={user.name}
              className="w-24 h-24 rounded-full border-4 border-purple-500/50"
            />
            <h2 className="mt-3 text-xl font-bold text-white truncate max-w-full">
              {user.name || 'Developer'}
            </h2>
            {user.github_username && (
              <p className="text-gray-400 text-sm">@{user.github_username}</p>
            )}
          </div>

          {/* Right side - Stats */}
          <div className="flex-1 pl-6 flex flex-col justify-center">
            {/* Stats row */}
            <div className="grid grid-cols-2 gap-4 mb-4">
              <div className="bg-gray-800/50 rounded-lg p-3 text-center">
                <p className="text-2xl font-bold text-purple-400">{totalContributions}</p>
                <p className="text-xs text-gray-400">Contributions</p>
              </div>
              <div className="bg-gray-800/50 rounded-lg p-3 text-center">
                <p className="text-2xl font-bold text-blue-400">{postCount}</p>
                <p className="text-xs text-gray-400">Posts</p>
              </div>
              <div className="bg-gray-800/50 rounded-lg p-3 text-center">
                <p className="text-2xl font-bold text-green-400">{followerCount}</p>
                <p className="text-xs text-gray-400">Followers</p>
              </div>
              <div className="bg-gray-800/50 rounded-lg p-3 text-center">
                <p className="text-2xl font-bold text-pink-400">{followingCount}</p>
                <p className="text-xs text-gray-400">Following</p>
              </div>
            </div>

            {/* Languages */}
            {topLanguages.length > 0 && (
              <div className="bg-gray-800/30 rounded-lg p-3">
                <p className="text-xs text-gray-400 mb-2">Top Languages</p>
                <div className="flex flex-wrap gap-2">
                  {topLanguages.map((lang) => (
                    <span
                      key={lang.language}
                      className="px-2 py-1 text-xs rounded-full text-white"
                      style={{
                        backgroundColor: languageColors[lang.language] || '#6b7280',
                      }}
                    >
                      {lang.language}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="absolute bottom-3 right-3 flex items-center gap-2">
          <span className="text-gray-500 text-xs">DevSync</span>
          <div className="w-4 h-4 bg-gradient-to-r from-purple-500 to-blue-500 rounded" />
        </div>
      </div>
    );
  }
);

ShareableProfileCard.displayName = 'ShareableProfileCard';

export default ShareableProfileCard;
