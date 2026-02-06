import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getUser, getFollowers, getFollowing } from '../api/users';
import { getUserPosts } from '../api/posts';
import { getContributions, getLanguages, getRepos } from '../api/github';
import { useAuthStore } from '../store/authStore';
import type { User } from '../types/user';
import type { Post } from '../types/post';
import type { GitHubContribution, GitHubLanguageStat, GitHubRepository } from '../types/github';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';
import FollowButton from '../components/profile/FollowButton';
import ContributionCalendar from '../components/profile/ContributionCalendar';
import LanguageChart from '../components/profile/LanguageChart';
import BadgeDisplay from '../components/profile/BadgeDisplay';
import PostCard from '../components/posts/PostCard';

export default function ProfilePage() {
  const { t } = useTranslation();
  const { id } = useParams<{ id: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const [user, setUser] = useState<User | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [contributions, setContributions] = useState<GitHubContribution[]>([]);
  const [languages, setLanguages] = useState<GitHubLanguageStat[]>([]);
  const [repos, setRepos] = useState<GitHubRepository[]>([]);
  const [followerCount, setFollowerCount] = useState(0);
  const [followingCount, setFollowingCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    const userId = parseInt(id);

    const fetchData = async () => {
      setLoading(true);
      try {
        const [userRes, postsRes, followersRes, followingRes] = await Promise.all([
          getUser(userId),
          getUserPosts(userId),
          getFollowers(userId),
          getFollowing(userId),
        ]);
        setUser(userRes.data);
        setPosts(postsRes.data || []);
        setFollowerCount((followersRes.data || []).length);
        setFollowingCount((followingRes.data || []).length);

        if (userRes.data.github_connected) {
          const [contribRes, langRes, reposRes] = await Promise.all([
            getContributions(userId),
            getLanguages(userId),
            getRepos(userId),
          ]);
          setContributions(contribRes.data || []);
          setLanguages(langRes.data || []);
          setRepos(reposRes.data || []);
        }
      } catch {
        // handle error
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading) return <div className="py-12"><LoadingSpinner /></div>;
  if (!user) return <div className="text-center text-gray-400 py-12">{t('errors.notFound')}</div>;

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Profile Header */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
        <div className="flex items-start gap-5">
          <Avatar name={user.name} avatarUrl={user.avatar_url} size="lg" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-3 flex-wrap">
              <h1 className="text-2xl font-bold">{user.name}</h1>
              {!isOwnProfile && <FollowButton userId={user.id} />}
            </div>
            {user.bio && <p className="text-gray-400 mt-1 text-sm">{user.bio}</p>}
            {user.github_username && (
              <a
                href={`https://github.com/${user.github_username}`}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 mt-2 text-sm text-gray-500 hover:text-blue-400 transition-colors"
              >
                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                @{user.github_username}
                <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                </svg>
              </a>
            )}
            <div className="flex gap-4 mt-3 text-sm">
              <Link to={`/profile/${user.id}/followers`} className="text-gray-400 hover:text-blue-400 transition-colors">
                <strong className="text-white">{followerCount}</strong> {t('profile.followers')}
              </Link>
              <Link to={`/profile/${user.id}/following`} className="text-gray-400 hover:text-blue-400 transition-colors">
                <strong className="text-white">{followingCount}</strong> {t('profile.following')}
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Skills */}
      {(user.skills_languages || user.skills_frameworks) && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6 space-y-4">
          {user.skills_languages && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2">
                <span>✨</span> {t('profile.languages')}
              </h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer">
                <img
                  src={`https://skillicons.dev/icons?i=${user.skills_languages}&theme=dark`}
                  alt="Languages"
                  className="h-12"
                />
              </a>
            </div>
          )}
          {user.skills_frameworks && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2">
                <span>🚀</span> {t('profile.frameworks')}
              </h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer">
                <img
                  src={`https://skillicons.dev/icons?i=${user.skills_frameworks}&theme=dark`}
                  alt="Frameworks"
                  className="h-12"
                />
              </a>
            </div>
          )}
        </div>
      )}

      {/* GitHub Data */}
      {user.github_connected && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-800">
              <h2 className="text-sm font-semibold">{t('profile.contributions')}</h2>
            </div>
            <div className="p-6">
              <ContributionCalendar contributions={contributions} />
            </div>
          </div>

          {languages.length > 0 && (
            <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-800">
                <h2 className="text-sm font-semibold">{t('profile.languages')}</h2>
              </div>
              <div className="p-6">
                <LanguageChart languages={languages} />
              </div>
            </div>
          )}
        </div>
      )}

      {/* Badges / Achievements */}
      <BadgeDisplay
        contributions={contributions}
        posts={posts}
        followerCount={followerCount}
        followingCount={followingCount}
      />

      {/* Repositories */}
      {user.github_connected && repos.length > 0 && (
        <div>
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide">
              {t('profile.repositories')}
            </h2>
            {user.github_username && (
              <a
                href={`https://github.com/${user.github_username}?tab=repositories`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-xs text-gray-500 hover:text-blue-400 transition-colors flex items-center gap-1"
              >
                {t('profile.viewAllOnGitHub')}
                <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                </svg>
              </a>
            )}
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {repos.slice(0, 6).map((repo) => (
              <a
                key={repo.id}
                href={`https://github.com/${repo.full_name}`}
                target="_blank"
                rel="noopener noreferrer"
                className="bg-gray-900 border border-gray-800 rounded-xl p-4 hover:border-gray-600 transition-colors group"
              >
                <div className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-gray-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" />
                  </svg>
                  <div className="min-w-0 flex-1">
                    <div className="font-medium text-sm text-blue-400 group-hover:text-blue-300 truncate">
                      {repo.name}
                    </div>
                    {repo.description && (
                      <p className="text-xs text-gray-500 mt-1 line-clamp-2">{repo.description}</p>
                    )}
                    <div className="flex items-center gap-3 mt-2">
                      {repo.language && (
                        <span className="flex items-center gap-1 text-xs text-gray-400">
                          <span className="w-2.5 h-2.5 rounded-full bg-yellow-400" />
                          {repo.language}
                        </span>
                      )}
                      {repo.stars > 0 && (
                        <span className="flex items-center gap-1 text-xs text-gray-400">
                          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M11.48 3.499a.562.562 0 0 1 1.04 0l2.125 5.111a.563.563 0 0 0 .475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 0 0-.182.557l1.285 5.385a.562.562 0 0 1-.84.61l-4.725-2.885a.562.562 0 0 0-.586 0L6.982 20.54a.562.562 0 0 1-.84-.61l1.285-5.386a.562.562 0 0 0-.182-.557l-4.204-3.602a.562.562 0 0 1 .321-.988l5.518-.442a.563.563 0 0 0 .475-.345L11.48 3.5Z" />
                          </svg>
                          {repo.stars}
                        </span>
                      )}
                      {repo.forks > 0 && (
                        <span className="flex items-center gap-1 text-xs text-gray-400">
                          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M7.217 10.907a2.25 2.25 0 1 0 0 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186 9.566-5.314m-9.566 7.5 9.566 5.314m0 0a2.25 2.25 0 1 0 3.935 2.186 2.25 2.25 0 0 0-3.935-2.186Zm0-12.814a2.25 2.25 0 1 0 3.933-2.185 2.25 2.25 0 0 0-3.933 2.185Z" />
                          </svg>
                          {repo.forks}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </a>
            ))}
          </div>
        </div>
      )}

      {/* Posts */}
      <div>
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">{t('profile.posts')}</h2>
        {posts.length === 0 ? (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center text-gray-400 text-sm">
            {t('profile.noPosts')}
          </div>
        ) : (
          <div className="space-y-3">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
